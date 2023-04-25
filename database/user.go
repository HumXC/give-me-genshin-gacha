package database

import (
	"give-me-genshin-gacha/models"
)

type UserMan interface {
	BaseMan[models.User]
}

type userMan struct {
	baseMan[models.User]
}
