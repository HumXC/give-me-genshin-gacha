package models

import (
	"time"

	"gorm.io/gorm"
)

type GachaLog struct {
	I         uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Uid       string         `json:"uid"`
	GachaType string         `json:"gacha_type"`
	Time      string         `json:"time"`
	Name      string         `json:"name"`
	Lang      string         `json:"lang"`
	ItemType  string         `json:"item_type"`
	RankType  string         `json:"rank_type"`
	ID        string         `json:"id"`
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
type GachaDB struct{}

func (d *GachaDB) GetUids() ([]string, error) {
	result := make([]string, 1)
	err := db.Table("gacha_logs").Select("uid").Group("uid").Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

func (d *GachaDB) GetLogs(uid, gachaType string, num, page int) ([]GachaLog, error) {
	result := make([]GachaLog, num)
	offset := num * page
	query := db.Table("gacha_logs").
		Where("uid=?", uid).
		Where("gacha_type=?", gachaType).
		Limit(num).
		Offset(offset)
	err := fixGachaType(query, gachaType).Find(&result).Error
	if err != nil {
		return result, err
	}
	return result, nil
}

// 返回与某物品同品质的上一件物品 之间的祈愿次数
func (d *GachaDB) GetNumWithLast(uid, gachaType, id string) (int, error) {
	query := db.Table("gacha_logs").
		Select("rank_type").
		Where("uid=?", uid).
		Where("gacha_type=?", gachaType).
		Order("i DESC").
		Limit(100).
		Where("i<=(?)", db.Table("gacha_logs").Select("i").Where("uid=?", uid).Where("id=?", id))

	rows, err := fixGachaType(query, gachaType).Rows()
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
	gachaTypes := []string{"301", "302", "200", "100"}
	result := make([]GachaUsedCost, 1)

	for _, gachaType := range gachaTypes {
		query := db.Table("gacha_logs").
			Select("COUNT(*)").
			Where("gacha_type=?", gachaType).
			Order("i DESC").
			Where("uid=?", uid).
			Where("i>(?)",
				db.Table("gacha_logs").Where("uid=?", uid).Select("MAX(i)").Where("gacha_type=?", gachaType).Where("rank_type=5"))
		cost := 0
		err := fixGachaType(query, gachaType).Find(&cost).Error
		if err != nil {
			return result, nil
		}
		result = append(result, GachaUsedCost{
			GachaType: gachaType,
			Cost:      cost,
		})
	}
	return result, nil
}

// 获取各品质各物品种类的总数
func (d *GachaDB) GetTotals(uid string) (map[string][]GachaTotal, error) {
	result := make(map[string][]GachaTotal, 0)
	rows, err := db.Table("gacha_logs").
		Select("gacha_type", "rank_type", "item_type", "COUNT(*) AS total").
		Where("uid=?", uid).
		Group("gacha_type, rank_type, item_type").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		t_ := GachaTotal{}
		gt := ""
		err := rows.Scan(&gt, &t_.RankType, &t_.ItemType, &t_.Total)
		if err != nil {
			return result, err
		}
		_, ok := result[gt]
		if !ok {
			result[gt] = make([]GachaTotal, 0)
		}
		result[gt] = append(result[gt], t_)
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

// 获取账号最后同步的物品
func (d *GachaDB) GetLastIDs() (map[string]map[string]string, error) {
	result := make(map[string]map[string]string, 4)
	rows, err := db.Table("gacha_logs").Select("uid", "i", "gacha_type", "id").Order("i DESC").Group("uid").Group("gacha_type").Rows()
	if err != nil {
		return result, err
	}
	defer rows.Close()
	// 保存物品的 i
	itemsI := make(map[string]map[string]int)
	for rows.Next() {
		var (
			uid       string
			i         int
			gachaType string
			id        string
		)
		err := rows.Scan(&uid, &i, &gachaType, &id)
		if err != nil {
			return result, err
		}

		if _, ok := result[uid]; !ok {
			result[uid] = make(map[string]string)
		}
		if _, ok := itemsI[uid]; !ok {
			itemsI[uid] = make(map[string]int)
		}
		itemsI[uid][gachaType] = i
		result[uid][gachaType] = id
	}
	// 在 400 和 301 的 gacha_type 之间选一个最新的，另一个删除
	// 因为 400 和 301 是两个池子，但是通过米哈游的 API 请求的 gacha_type 都是301
	for uid := range itemsI {
		i301 := itemsI[uid]["301"]
		i400 := itemsI[uid]["400"]
		if i301 != 0 && i400 != 0 {
			if i301 > i400 {
				delete(result[uid], "400")
			} else {
				delete(result[uid], "301")
			}
		}
	}
	return result, nil
}

// 添加条目到数据库
func (d *GachaDB) Add(items []GachaLog) error {
	err := db.Create(items).Error
	return err
}
func fixGachaType(query *gorm.DB, gachaType string) *gorm.DB {
	if gachaType == "301" {
		query.Or("gacha_type=?", "400")
	}
	return query
}
