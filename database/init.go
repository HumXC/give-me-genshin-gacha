package database

import (
	"give-me-genshin-gacha/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Models interface {
	models.Log | models.Item | models.User | models.ItemName
}
type BaseMan[T Models] interface {
	Model() *gorm.DB
	Create(...T) error
	Remove(id ...uint64) error
	Update(T) error
}
type baseMan[T Models] struct {
	db    gorm.DB
	model T
}

func (m *baseMan[T]) Model() *gorm.DB {
	return m.db.Model(m.model)
}

func (m *baseMan[T]) Create(values ...T) error {
	return m.Model().Create(values).Error
}

func (m *baseMan[T]) Remove(ids ...uint64) error {
	return m.db.Delete(m.model, ids).Error
}

func (m *baseMan[T]) Update(value T) error {
	return m.db.Save(&value).Error
}

type DB struct {
	Logs      LogMan
	Items     ItemMan
	ItemNames ItemNameMan
	Users     UserMan
}

func NewDB(file string) (*DB, error) {
	gormDB, err := gorm.Open(sqlite.Open(file))
	if err != nil {
		return nil, err
	}
	err = gormDB.AutoMigrate(
		models.Log{},
		models.Item{},
		models.ItemName{},
		models.User{},
	)
	if err != nil {
		return nil, err
	}
	db := DB{
		Logs:      &logMan{baseMan: baseMan[models.Log]{db: *gormDB}},
		Users:     &userMan{baseMan: baseMan[models.User]{db: *gormDB}},
		Items:     &itemMan{baseMan: baseMan[models.Item]{db: *gormDB}},
		ItemNames: &itemNameMan{baseMan: baseMan[models.ItemName]{db: *gormDB}},
	}

	return &db, nil
}
