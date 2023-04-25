package models

import (
	"time"
)

type User struct {
	Model
	UserID   uint64
	RawURL   string
	SyncTime time.Time
}
