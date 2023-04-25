package database

import (
	"give-me-genshin-gacha/models"
)

type LogMan interface {
	BaseMan[models.Log]
}

type logMan struct {
	baseMan[models.Log]
}
