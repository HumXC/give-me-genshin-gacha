package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type ItemType = int

const (
	ItemWeapon ItemType = 0
	ItemAvatar ItemType = 1
)

type Name struct {
	gorm.Model
	ItemID uint
	Lang   string
	Value  string
}
type Item struct {
	gorm.Model
	RankType int      // 星数
	ItemType ItemType // 武器还是角色
	// 一个物品有多个 Name
	Names []Name
}

type ItemDB struct {
	db *gorm.DB
}

func (i *ItemDB) GetWithID(id uint) (item Item, err error) {
	err = i.db.Where("id = ?", id).Find(&item).Error
	return
}

func (i *ItemDB) GetWithName(lang, name string) (Item, error) {
	item := Item{}
	var id uint = 0
	err := i.db.Model(&Name{}).Debug().
		Select("item_id").
		Where("lang=? AND value=?", lang, name).
		Find(&id).Error
	if err != nil {
		return item, err
	}
	return i.GetWithID(id)
}

// 获取已经加载的语言种类
// FIXME 似乎不准确
func (i *ItemDB) LoadedLang() ([]string, error) {
	result := make([]string, 0, 10)
	err := i.db.Model(&Item{}).Select("lang").Group("lang").Association("Names").Find(&result)
	return result, err
}

// Update 会保存 item 里的所有字段，即使里面包含零值
// 所以 Update 前先从 Get() 获取数据库里已有的 item
func (i *ItemDB) Update(item Item) error {
	return i.db.Save(&item).Error
}

// 设置名称
func (i *ItemDB) SetName(id uint, lang, name string) error {
	var _name sql.NullString
	err := i.db.Model(&Name{}).Select("value").Where("item_id=? AND lang=?", id, lang).Find(&_name).Error
	if err != nil {
		return err
	}
	if _name.Valid {
		return nil
	}
	err = i.db.Model(&Item{
		Model: gorm.Model{ID: id},
	}).
		Association("Names").
		Append(&Name{Lang: lang, ItemID: id, Value: name})
	return err
}
