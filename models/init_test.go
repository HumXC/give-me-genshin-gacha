package models

import (
	"encoding/json"
	"path"
	"testing"
)

var testData = `
 		[{
            "uid": "111111111",
            "gacha_type": "200",
            "item_id": "",
            "count": "1",
            "time": "2022-11-10 12:12:19",
            "name": "以理服人",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1668053160000458186"
        },
        {
            "uid": "222222222",
            "gacha_type": "200",
            "item_id": "",
            "count": "1",
            "time": "2022-11-09 21:23:03",
            "name": "冷刃",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667999160001203086"
        },
        {
            "uid": "111111111",
            "gacha_type": "200",
            "item_id": "",
            "count": "1",
            "time": "2022-11-05 14:58:08",
            "name": "黎明神剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667628360004797286"
        },
        {
            "uid": "111111111",
            "gacha_type": "200",
            "item_id": "",
            "count": "1",
            "time": "2022-11-04 00:03:21",
            "name": "鸦羽弓",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667491560000316486"
        },
        {
            "uid": "111111111",
            "gacha_type": "200",
            "item_id": "",
            "count": "1",
            "time": "2022-11-02 20:37:19",
            "name": "黎明神剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667390760012370286"
        }, 
        {
            "uid": "111111111",
            "gacha_type": "302",
            "item_id": "",
            "count": "1",
            "time": "2022-11-02 11:54:30",
            "name": "匣里龙吟",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "4",
            "id": "1667358360213844286"
        },
        {
            "uid": "111111111",
            "gacha_type": "302",
            "item_id": "",
            "count": "1",
            "time": "2022-11-02 11:03:03",
            "name": "冷刃",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667358360021459086"
        },
        {
            "uid": "111111111",
            "gacha_type": "302",
            "item_id": "",
            "count": "1",
            "time": "2022-11-01 10:18:51",
            "name": "黎明神剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667268360001062186"
        },
        {
            "uid": "222222222",
            "gacha_type": "302",
            "item_id": "",
            "count": "1",
            "time": "2022-11-01 10:18:48",
            "name": "弹弓",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667268360001060986"
        },
        {
            "uid": "111111111",
            "gacha_type": "302",
            "item_id": "",
            "count": "1",
            "time": "2022-11-01 10:18:46",
            "name": "沐浴龙血的剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667268360001058586"
        },
        {
            "uid": "111111111",
            "gacha_type": "301",
            "item_id": "",
            "count": "1",
            "time": "2022-11-01 10:18:44",
            "name": "飞天御剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667268360001057286"
        },
        {
            "uid": "111111111",
            "gacha_type": "301",
            "item_id": "",
            "count": "1",
            "time": "2022-11-01 10:18:40",
            "name": "翡玉法球",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1667268360001054286"
        },
        {
            "uid": "222222222",
            "gacha_type": "301",
            "item_id": "",
            "count": "1",
            "time": "2022-10-21 01:18:24",
            "name": "黎明神剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "4",
            "id": "1666285560000425386"
        },
        {
            "uid": "111111111",
            "gacha_type": "301",
            "item_id": "",
            "count": "1",
            "time": "2022-10-21 01:17:53",
            "name": "铁影阔剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1666285560000415086"
        },
        {
            "uid": "111111111",
            "gacha_type": "301",
            "item_id": "",
            "count": "1",
            "time": "2022-10-21 00:51:58",
            "name": "铁影阔剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1666281960001662986"
        },
        {
            "uid": "111111111",
            "gacha_type": "100",
            "item_id": "",
            "count": "1",
            "time": "2022-10-14 20:11:55",
            "name": "班尼特",
            "lang": "zh-cn",
            "item_type": "角色",
            "rank_type": "4",
            "id": "1665749160003901986"
        },
        {
            "uid": "111111111",
            "gacha_type": "100",
            "item_id": "",
            "count": "1",
            "time": "2022-10-14 18:11:44",
            "name": "弹弓",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "5",
            "id": "1665741960039490486"
        },
        {
            "uid": "111111111",
            "gacha_type": "100",
            "item_id": "",
            "count": "1",
            "time": "2022-10-14 18:11:05",
            "name": "飞天御剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1665741960038216486"
        },
        {
            "uid": "111111111",
            "gacha_type": "100",
            "item_id": "",
            "count": "1",
            "time": "2022-10-14 18:11:02",
            "name": "飞天御剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1665741960038126786"
        },
        {
            "uid": "111111111",
            "gacha_type": "100",
            "item_id": "",
            "count": "1",
            "time": "2022-10-14 18:10:59",
            "name": "飞天御剑",
            "lang": "zh-cn",
            "item_type": "武器",
            "rank_type": "3",
            "id": "1665741960038012486"
        }]
`

func TestInit(t *testing.T) {
	db, err := NewDB(path.Join(TestDir, "data.db"))
	if err != nil {
		t.Fatal("can not init db", err)
	}
	g := GachaLog{}
	db.db.AutoMigrate(&g)
}

func TestAdd(t *testing.T) {
	db, err := NewDB(path.Join(TestDir, "data.db"))
	if err != nil {
		t.Fatal("can not init db", err)
	}
	g := GachaLog{}
	db.db.AutoMigrate(&g)
	data := make([]GachaLog, 0)
	err = json.Unmarshal([]byte(testData), &data)
	if err != nil {
		t.Fatal("can not unmarshal test data", err)
	}
	db.Add(data)
}
