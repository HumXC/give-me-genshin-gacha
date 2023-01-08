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
