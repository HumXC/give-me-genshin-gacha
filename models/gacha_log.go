package models

import (
	"time"

	"gorm.io/gorm"
)

type GachaLog struct {
	// OriginGachaType 是米哈游自带的 gacha_type, 会有 400 的值
	// GachaType 是 uigf_gacha_type
	// 见 https://github.com/DGP-Studio/Snap.Genshin/wiki/StandardFormat
	OriginGachaType string    `json:"origin_gacha_type"`
	GachaType       string    `json:"gacha_type"`
	Uid             uint64    `json:"uid"`
	Time            time.Time `json:"time"`
	ItemID          uint      `json:"item_id"`
	Count           int       `json:"count"`
	LogID           uint64    `json:"log_id"`
	gorm.Model
}
type LogDB struct {
	db *gorm.DB
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
