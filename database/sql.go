package database

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type GachaLog struct {
	GachaType string `json:"gachaType"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	Itemtype  string `json:"itemType"`
	RankType  string `json:"rankType"`
	ID        string `json:"id"`
}

type GachaUsedCost struct {
	GachaType string `json:"gachaType"`
	Cost      int    `json:"cost"`
}

type GachaTotal struct {
	Total    int    `json:"total"`
	ItemType string `json:"itemType"`
	RankType string `json:"rankType"`
}
type GachaDB struct {
	db   *sql.DB
	Uids []string
}

// 获取物品
func (d *GachaDB) GetLogs(uid, gachaType string, num, page int) ([]GachaLog, error) {
	result := make([]GachaLog, 0)
	sql := `
	SELECT gacha_type,time,name,lang,item_type,rank_type,id FROM 'UID'
	WHERE gacha_type=? ORDER BY i DESC LIMIT ?,?
	`
	offset := num * page
	sql = strings.Replace(sql, "UID", uid, 1)
	sql = fixGachaType(sql, gachaType)
	rows, err := d.db.Query(sql, gachaType, offset, num)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		i := GachaLog{}
		err := rows.Scan(&i.GachaType, &i.Time, &i.Name, &i.Lang, &i.Itemtype, &i.RankType, &i.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

// 获取与某物品与其同品质的上一件物品之间的祈愿次数
func (d *GachaDB) GetNumWithLast(uid, gachaType, id string) (int, error) {
	sql := `
	SELECT rank_type FROM 'UID' WHERE gacha_type=? AND
	i<=(SELECT i FROM 'UID' WHERE id=?)
	ORDER BY i DESC
	LIMIT 100;
	`
	sql = strings.Replace(sql, "UID", uid, 2)
	sql = fixGachaType(sql, gachaType)
	rows, err := d.db.Query(sql, gachaType, id)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	count := 1
	rankType := ""
	for rows.Next() {
		rt := ""
		err := rows.Scan(&rt)
		if err != nil {
			return 0, err
		}
		if rankType == "" {
			rankType = rt
			continue
		}
		if rankType == rt || rt == "5" {
			return count, nil
		} else {
			count++
		}
	}
	return count, nil
}

// 获取上一次五星之后抽的次数
func (d *GachaDB) GetUsedCost(uid string) ([]GachaUsedCost, error) {
	sql := `
	SELECT COUNT ( * ) FROM
	(SELECT i,gacha_type,rank_type FROM 'UID' 
		WHERE gacha_type = ? order by i desc)
	WHERE i>
		(SELECT MAX( i ) FROM 'UID' WHERE gacha_type=? AND rank_type=5 order by i desc)`
	gachaTypes := []string{"301", "302", "200", "100"}
	result := make([]GachaUsedCost, 0)
	sql = strings.ReplaceAll(sql, "UID", uid)
	for _, gachaType := range gachaTypes {
		rows, err := d.db.Query(fixGachaType(sql, gachaType), gachaType, gachaType)
		if err != nil {
			return result, err
		}
		defer rows.Close()
		if rows.Next() {
			var cost int
			err := rows.Scan(&cost)
			if err != nil {
				return result, err
			}
			result = append(result, GachaUsedCost{
				GachaType: gachaType,
				Cost:      cost,
			})
		}
	}

	return result, nil
}

func (d *GachaDB) GetTotals(uid string) (map[string][]GachaTotal, error) {
	result := make(map[string][]GachaTotal, 0)
	sql := `SELECT 
	gacha_type, rank_type,
	item_type, COUNT(*) AS total 
	FROM "?" GROUP BY gacha_type, rank_type, item_type`
	rows, err := d.db.Query(strings.Replace(sql, "?", uid, 1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		t := GachaTotal{}
		gt := ""
		err := rows.Scan(&gt, &t.RankType, &t.ItemType, &t.Total)
		if err != nil {
			return nil, err
		}
		_, ok := result[gt]
		if !ok {
			result[gt] = make([]GachaTotal, 0)
		}
		result[gt] = append(result[gt], t)
	}
	for _, g400 := range result["400"] {
		if _, ok := result["301"]; !ok {
			result["301"] = make([]GachaTotal, 0)
		}
		if len(result["301"]) == 0 {
			result["301"] = result["400"]
			break
		}
		for i := 0; i < len(result["301"]); i++ {
			if g400.ItemType == result["301"][i].ItemType && g400.RankType == result["301"][i].RankType {
				result["301"][i].Total += g400.Total
			}
		}
	}
	delete(result, "400")
	return result, nil
}
func (d *GachaDB) GetLastIDs() (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	sql := `
	SELECT * FROM ( 
		SELECT i, gacha_type, id 
		FROM "?" ORDER BY i DESC
	) GROUP BY gacha_type`
	for _, u := range d.Uids {
		rows, err := d.db.Query(strings.Replace(sql, "?", u, 1))
		if err != nil {
			return nil, err
		}
		m := make(map[string]int, 0)
		for rows.Next() {
			var (
				i         int
				gachaType string
				id        string
			)
			err := rows.Scan(&i, &gachaType, &id)
			if err != nil {
				return nil, err
			}
			if _, ok := result[u]; !ok {
				result[u] = make(map[string]string)
			}
			m[gachaType] = i
			result[u][gachaType] = id
		}
		// 在 400 和 301 的 gacha_type 之间选一个最新的，另一个删除
		// 因为 400 和 301 是两个池子，但是请求的 gacha_type 都是301
		id301 := m["301"]
		id400 := m["400"]
		if id301 != 0 && id400 != 0 {
			if id301 > id400 {
				delete(result[u], "400")
			} else {
				delete(result[u], "301")
			}
		}
		rows.Close()
	}
	return result, nil
}
func (d *GachaDB) Add(uid string, items []GachaLog) error {
	// 字段 i 保证插入顺序，祈愿记录越新，i 越大
	createTable := `
	CREATE TABLE ? (
		i INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        gacha_type TEXT,
		name TEXT,
        time TEXT,
        item_type TEXT,
        rank_type TEXT,
        id TEXT,
		lang TEXT
	)`
	insert := `
	INSERT INTO ? (
		gacha_type, name, time, item_type, rank_type, id, lang
	) VALUES (
		?, ?, ?, ?, ?, ?, ?
	)
	`
	insert = strings.Replace(insert, "?", "'"+uid+"'", 1)
	isHasUid := false
	for _, u := range d.Uids {
		if uid == u {
			isHasUid = true
			break
		}
	}
	if !isHasUid {
		_, err := d.db.Exec(strings.Replace(createTable, "?", "'"+uid+"'", 1))
		if err != nil {
			return err
		}
		d.Uids = append(d.Uids, uid)
	}
	stmt, err := d.db.Prepare(insert)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		_, err := stmt.Exec(
			item.GachaType, item.Name,
			item.Time, item.Itemtype,
			item.RankType, item.ID, item.Lang)
		if err != nil {
			return err
		}
	}
	return nil
}
func fixGachaType(query, gachaType string) string {
	if gachaType == "301" {
		return strings.ReplaceAll(query, "gacha_type=?", "(gacha_type=? OR gacha_type='400')")
	}
	if gachaType == "400" {
		return strings.ReplaceAll(query, "gacha_type=?", "(gacha_type=? OR gacha_type='301')")
	}
	return query
}
func NewDB(name string) (*GachaDB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	gachaDB := &GachaDB{
		db:   db,
		Uids: make([]string, 0),
	}
	sql := "SELECT name FROM sqlite_master WHERE type='table' AND name!='sqlite_sequence'"
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		gachaDB.Uids = append(gachaDB.Uids, s)
	}

	return gachaDB, nil
}
