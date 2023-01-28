package models

import (
	"gorm.io/gorm"
)

const (
	ItemWeapon = 0
	ItemAvatar = 1
)

type Item struct {
	gorm.Model
	Level int // 星数
	Type  int // 武器还是角色
	// 以下是各语言的物品名称
	ZhCN string
	ZhTW string
	JaJP string
	EnUS string
}

type ItemDB struct{}

func (i *ItemDB) Get(id uint) (item Item, err error) {
	err = db.Model(&Item{}).Where("id = ?", id).Find(&item).Error
	return
}

// Update 会保存 item 里的所有字段，即使里面包含零值
// 所以 Update 前先从 Get() 获取数据库里已有的 item
func (i *ItemDB) Update(item Item) error {
	return db.Model(&User{}).Save(&item).Error
}
