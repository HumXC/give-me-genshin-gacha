package models

import (
	"time"

	"gorm.io/gorm"
)

type GachaInfo struct {
	GachaType string `json:"gachaType"`
	AllCount  int    `json:"allCount"`
	Avatar5   int    `json:"avatar5"`
	Avatar4   int    `json:"avatar4"`
	Weapon5   int    `json:"weapon5"`
	Weapon4   int    `json:"weapon4"`
	Weapon3   int    `json:"weapon3"`
}
type GachaLog struct {
	gorm.Model
	// OriginGachaType 是米哈游自带的 gacha_type, 会有 400 的值
	// GachaType 是 uigf_gacha_type
	// 见 https://github.com/DGP-Studio/Snap.Genshin/wiki/StandardFormat
	OriginGachaType string    `json:"origin_gacha_type"`
	GachaType       string    `json:"gacha_type"`
	Uid             uint64    `json:"uid"`
	Time            time.Time `json:"time"`
	ItemID          int       `json:"item_id"`
	Count           int       `json:"count"`
	LogID           uint64    `json:"log_id"`
}
type LogDB struct {
	db *gorm.DB
}

func (d *LogDB) GetInfo() ([]GachaInfo, error) {
	result := make([]struct {
		Count     int
		GachaType string
		RankType  int
		ItemType  int
	}, 0)
	err := d.db.Model(&GachaLog{}).
		Select("COUNT(*) as count", "gacha_type", "rank_type", "item_type").
		Joins("join items on gacha_logs.item_id=items.id").
		Group("gacha_type").
		Group("rank_type").
		Group("item_type").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	infosMap := make(map[string]GachaInfo, 0)
	for _, v := range result {
		info := infosMap[v.GachaType]
		info.GachaType = v.GachaType
		info.AllCount += v.Count
		if v.ItemType == ItemAvatar {
			switch v.RankType {
			case 4:
				info.Avatar4 += v.Count
			case 5:
				info.Avatar5 += v.Count
			}
			continue
		}
		switch v.RankType {
		case 3:
			info.Weapon3 += v.Count
		case 4:
			info.Weapon4 += v.Count
		case 5:
			info.Weapon5 += v.Count
		}
		infosMap[v.GachaType] = info
	}
	infos := make([]GachaInfo, 0)
	for _, v := range infosMap {
		infos = append(infos, v)
	}
	return infos, nil
}

// 添加条目到数据库
func (d *LogDB) Add(items []GachaLog) error {
	if len(items) == 0 {
		return nil
	}
	return d.db.Create(items).Error
}

// TODO: 测试多用户时此方法的正确性
// 获取每个池子每个 uid 最新的一次祈愿记录
func (d *LogDB) EndLogIDs() (map[string]map[uint64]uint64, error) {
	result := make(map[string]map[uint64]uint64, 0)
	col := make([]GachaLog, 0)
	subQuery := d.db.Model(&GachaLog{}).Order("id DESC")
	err := d.db.Table("(?) as u", subQuery).Select("uid", "log_id", "gacha_type").Group("gacha_type").Find(&col).Error
	for _, log := range col {
		if result[log.GachaType] == nil {
			result[log.GachaType] = make(map[uint64]uint64, 0)
		}
		result[log.GachaType][log.Uid] = log.LogID
	}
	return result, err
}
