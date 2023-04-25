package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64 `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Log struct {
	Model
	// OriginGachaType 是米哈游自带的 gacha_type, 会有 400 的值
	// GachaType 是 uigf_gacha_type
	// 见 https://github.com/DGP-Studio/Snap.Genshin/wiki/StandardFormat
	OriginGachaType string
	GachaType       string
	UserID          uint64
	Time            time.Time
	ItemID          uint64
	Count           int32
	LogID           uint64
}
