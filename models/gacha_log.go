package models

import (
	"time"

	"gorm.io/gorm"
)

type GachaLog struct {
	gorm.Model
	Uid       uint
	GachaType string
	Time      time.Time
	ItemID    uint
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
type LogDB struct {
	db *gorm.DB
}

// 添加条目到数据库
func (d *LogDB) Add(items []GachaLog) error {
	return d.db.Create(items).Error
}
func fixGachaType(query *gorm.DB, gachaType string) *gorm.DB {
	if gachaType == "301" {
		query.Or("gacha_type=?", "400")
	}
	return query
}
