package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type ItemType = int

const (
	ItemWeapon ItemType = 0
	ItemAvatar ItemType = 1
)

type Item struct {
	ID        int `gorm:"primarykey" josn:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RankType  ItemType       `json:"rank_type"` // 星数
	ItemType  ItemType       `json:"item_type"` // 武器还是角色
	// 以下是各语言的物品名称
	// 如果需要增加语言的支持，还需要更改 LoadedLang 方法和 assets/item.go 的代码
	ZhCN string `gorm:"column:zh-cn" json:"zh-cn"`
	ZhTW string `gorm:"column:zh-tw" json:"zh-tw"`
	JaJP string `gorm:"column:ja-jp" json:"ja-jp"`
	EnUS string `gorm:"column:en-us" json:"en-us"`
}

type ItemDB struct {
	db *gorm.DB
}

func (i *ItemDB) GetWithID(id int) (item Item, err error) {
	err = i.db.Where("id = ?", id).Find(&item).Error
	return
}

func (i *ItemDB) GetWithName(lang, name string) (Item, error) {
	item := Item{}
	err := i.db.Where("`"+lang+"`=?", name).Find(&item).Error
	return item, err
}

// 获取已经加载的语言种类
func (i *ItemDB) LoadedLang() ([]string, error) {
	result := make([]string, 0)
	r := Item{}
	err := i.db.Last(&r).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return result, nil
		}
		return nil, err
	}
	// 此处的判断 ID == 1 是因为在 NewDB 里增加了一条初始记录 ID 是 1
	if r.ID == 1 {
		return result, nil
	}
	// TODO: 是否有更好的办法，而不是 if
	if r.ZhCN != "" {
		result = append(result, "zh-cn")
	}
	if r.ZhTW != "" {
		result = append(result, "zh-tw")
	}
	if r.JaJP != "" {
		result = append(result, "ja-jp")
	}
	if r.EnUS != "" {
		result = append(result, "en-us")
	}
	return result, nil
}

// Update 会保存 item 里的所有字段，即使里面包含零值
// 所以 Update 前先从 Get() 获取数据库里已有的 item
func (i *ItemDB) Update(item Item) error {
	return i.db.Save(&item).Error
}
