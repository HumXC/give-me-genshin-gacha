package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"give-me-genshin-gacha/gacha"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	INIT_TABLE = `
	CREATE TABLE IF NOT EXISTS gacha (
		uid TEXT NOT NULL,
        gacha_type TEXT,
        item_id TEXT,
        count TEXT,
        time DATETIME,
        name TEXT,
        lang TEXT,
        item_type TEXT,
        rank_type TEXT,
        id TEXT,
		i INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT
	)`
	INSERT_ITEM = `
	INSERT INTO gacha  (
		uid,
        gacha_type,
        item_id,
        count,
        time,
        name,
        lang,
        item_type,
        rank_type,
        id
	) VALUES`

	GET_LAST_ID = `
	SELECT * FROM ( 
		SELECT uid, gacha_type, id 
		FROM gacha ORDER BY time DESC
	) GROUP BY uid,gacha_type;
	`

	GET_UIDS = `
	SELECT uid FROM gacha GROUP BY uid
	`

	GET_GACHA_RANK = `
	SELECT 
		uid,
        gacha_type,
        item_id,
        count,
        time,
        name,
        lang,
        item_type,
        rank_type,
        id
	FROM gacha WHERE uid=? AND ( gacha_type=? ) AND rank_type=? ORDER BY time DESC
	`

	COUNT_GACHA_IN_ID = `
	SELECT COUNT ( * ) num FROM (
		SELECT id FROM gacha WHERE uid=? AND ( gacha_type=? ) ORDER BY time
	) WHERE id BETWEEN ? and ?
	`

	COUNT_GACHA = `
	SELECT COUNT ( * ) num FROM gacha WHERE uid=? AND ( gacha_type=? )
	`

	GET_OLDEST = `
	SELECT 
		uid,
        gacha_type,
        item_id,
        count,
        time,
        name,
        lang,
        item_type,
        rank_type,
        id
	FROM gacha WHERE uid=? ( AND gacha_type=? ) ORDER BY time LIMIT 1
	`
)

type GachaItem struct {
	Uid       string    `json:"uid"`
	GachaType string    `json:"gacha_type"`
	ItemId    string    `json:"item_id"`
	Count     string    `json:"count"`
	Time      time.Time `json:"time"`
	Name      string    `json:"name"`
	Lang      string    `json:"lang"`
	Itemtype  string    `json:"item_type"`
	RankType  string    `json:"rank_type"`
	ID        string    `json:"id"`
}

type GachaDB interface {
	Add(items []gacha.RespDataListItem) error
	// 获取最新的物品 id, 用于 Fetcher.Get()
	GetLastIDs() (map[string]map[string]string, error)
	// 获取最久远的记录
	GetOldest(uid, gachaType string) (GachaItem, error)
	// 根据物品星数来获取物品
	GetGachaRank(uid, gachaType, rankType string) ([]GachaItem, error)
	// 获取所有 uid
	GetUids() ([]string, error)
	// 获取两个物品之前的物品数量
	CountIn(uid, gachaType, id1, id2 string) (int, error)
	// 获取某个卡池的祈愿数量
	CountGacha(uid, gachaType string) (int, error)
	// 获取某个卡池某种品质的物品数量
	CountGachaRank(uid, gachaType, rank string) (int, error)
	Close() error
}
type DB struct {
	*sql.DB
}

func (d *DB) GetLastIDs() (map[string]map[string]string, error) {
	m := make(map[string]map[string]string)
	rows, err := d.Query(GET_LAST_ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			uid       string
			gachaType string
			id        string
		)
		err := rows.Scan(&uid, &gachaType, &id)
		if err != nil {
			return nil, err
		}
		if _, ok := m[uid]; !ok {
			m[uid] = make(map[string]string)
		}
		m[uid][gachaType] = id
	}
	return m, nil
}

func (d *DB) Add(items []gacha.RespDataListItem) error {
	if len(items) == 0 {
		return nil
	}
	var b bytes.Buffer

	for i, v := range items {
		t, _ := time.Parse("2006-01-02 15:04:05", v.Time)
		str := fmt.Sprintf("(\"%s\",\"%s\",\"%s\",\"%s\",%d,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")",
			v.Uid,
			v.GachaType,
			v.ItemId,
			v.Count,
			t.Unix(),
			v.Name,
			v.Lang,
			v.Itemtype,
			v.RankType,
			v.ID)
		b.WriteString(str)
		if i != len(items)-1 {
			b.WriteString(",")
		}
	}
	s := b.String()
	_, err := d.Exec(INSERT_ITEM + s)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) GetGachaRank(uid, gachaType, rankType string) ([]GachaItem, error) {

	var rows *sql.Rows
	// 因为活动祈愿有两个池子，所以要用这种方法
	if gachaType == "301" {
		s := strings.Replace(GET_GACHA_RANK, "gacha_type=?", " gacha_type=? OR gacha_type=? ", 1)
		rs, err := d.Query(s, uid, "301", "400", rankType)
		if err != nil {
			return nil, err
		}
		rows = rs
	} else {
		rs, err := d.Query(GET_GACHA_RANK, uid, gachaType, rankType)
		if err != nil {
			return nil, err
		}
		rows = rs
	}

	result := make([]GachaItem, 0)
	for rows.Next() {
		i := GachaItem{}
		err := rows.Scan(
			&i.Uid, &i.GachaType,
			&i.ItemId, &i.Count,
			&i.Time, &i.Name,
			&i.Lang, &i.Itemtype,
			&i.RankType, &i.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return result, nil
}

func (d *DB) GetUids() ([]string, error) {
	result := make([]string, 0)
	rows, err := d.Query(GET_UIDS)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var s string
		err := rows.Scan(&s)
		if err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (d *DB) CountIn(uid, gachaType, id1, id2 string) (int, error) {
	var rows *sql.Rows
	// 因为活动祈愿有两个池子，所以要用这种方法
	if gachaType == "301" {
		s := strings.Replace(COUNT_GACHA_IN_ID, "gacha_type=?", " gacha_type=? OR gacha_type=? ", 1)
		rs, err := d.Query(s, uid, "301", "400", id1, id2)
		if err != nil {
			return 0, err
		}
		rows = rs
	} else {
		rs, err := d.Query(COUNT_GACHA_IN_ID, uid, gachaType, id1, id2)
		if err != nil {
			return 0, err
		}
		rows = rs
	}
	if !rows.Next() {
		return 0, nil
	}
	num := 0
	err := rows.Scan(&num)
	if err != nil {
		return 0, err
	}
	return num - 1, nil
}

func (d *DB) CountGacha(uid, gachaType string) (int, error) {
	var rows *sql.Rows
	// 因为活动祈愿有两个池子，所以要用这种方法
	if gachaType == "301" {
		s := strings.Replace(COUNT_GACHA, "gacha_type=?", " gacha_type=? OR gacha_type=? ", 1)
		rs, err := d.Query(s, uid, "301", "400")
		if err != nil {
			return 0, err
		}
		rows = rs
	} else {
		rs, err := d.Query(COUNT_GACHA, uid, gachaType)
		if err != nil {
			return 0, err
		}
		rows = rs
	}
	if !rows.Next() {
		return 0, nil
	}
	num := 0
	err := rows.Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (d *DB) CountGachaRank(uid, gachaType, rankType string) (int, error) {
	var rows *sql.Rows
	// 因为活动祈愿有两个池子，所以要用这种方法
	if gachaType == "301" {
		s := strings.Replace(COUNT_GACHA, "gacha_type=?", " gacha_type=? OR gacha_type=? ", 1)
		rs, err := d.Query(s+" AND rank_type=? ", uid, "301", "400", rankType)
		if err != nil {
			return 0, err
		}
		rows = rs
	} else {
		rs, err := d.Query(COUNT_GACHA+" AND rank_type=? ", uid, gachaType, rankType)
		if err != nil {
			return 0, err
		}
		rows = rs
	}

	if !rows.Next() {
		return 0, nil
	}
	num := 0
	err := rows.Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func (d *DB) GetOldest(uid, gachaType string) (GachaItem, error) {
	i := GachaItem{}
	rows, err := d.Query(GET_OLDEST, uid, gachaType)
	if err != nil {
		return i, err
	}
	if rows.Next() {
		err := rows.Scan(
			&i.Uid, &i.GachaType,
			&i.ItemId, &i.Count,
			&i.Time, &i.Name,
			&i.Lang, &i.Itemtype,
			&i.RankType, &i.ID)
		if err != nil {
			return i, err
		}
	}
	return i, nil
}
func NewDB(name string) (GachaDB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(INIT_TABLE)
	if err != nil {
		return nil, err
	}
	return &DB{
		db,
	}, nil
}
