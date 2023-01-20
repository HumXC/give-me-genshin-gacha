package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Uid        string
	RawURL     string
	UpdateTime time.Time
}
