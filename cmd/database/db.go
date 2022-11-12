package database

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"give-me-genshin-gacha/network"
	"strconv"
	"time"
)

const INIT_TABLE = `
	CREATE TABLE IF NOT EXISTS gacha (
		uid TEXT NOT NULL,
        gacha_type TEXT,
        item_id TEXT,
        count TEXT,
        time INT,
        name TEXT,
        lang TEXT,
        item_type TEXT,
        rank_type TEXT,
        id TEXT PRIMARY KEY
	)`
const INSERT_ITEM = `
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

const GET_LAST_ID = `
	SELECT * FROM ( 
		SELECT uid, gacha_type, id 
		FROM gacha ORDER BY time DESC
	) GROUP BY uid,gacha_type;
`

type DB struct {
	db *sql.DB
}

func genError(err error) error {
	return errors.New("Database Error: \n\t" + err.Error())
}

// 获取最新的物品 id, 用于 Fetcher.Get()
func (d *DB) GetLastIDs() (map[string]map[string]string, error) {
	m := make(map[string]map[string]string)
	rows, err := d.db.Query(GET_LAST_ID)
	if err != nil {
		return nil, genError(err)
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
			return nil, genError(err)
		}
		if _, ok := m[uid]; !ok {
			m[uid] = make(map[string]string)
		}
		m[uid][gachaType] = id
	}
	return m, nil
}

func (d *DB) Add(items []network.RespDataListItem) error {
	if len(items) == 0 {
		return nil
	}
	var b bytes.Buffer
	for i, v := range items {
		t, _ := time.Parse("2006-01-02 15:04:05", v.Time)
		str := fmt.Sprintf("(\"%s\",\"%s\",\"%s\",\"%s\",%s,\"%s\",\"%s\",\"%s\",\"%s\",\"%s\")",
			v.Uid,
			v.GachaType,
			v.ItemId,
			v.Count,
			strconv.FormatInt(t.Unix(), 10),
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
	_, err := d.db.Exec(INSERT_ITEM + s)
	if err != nil {
		return genError(err)
	}
	return nil
}

func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, genError(err)
	}
	_, err = db.Exec(INIT_TABLE)
	if err != nil {
		return nil, genError(err)
	}
	return &DB{
		db: db,
	}, nil
}
