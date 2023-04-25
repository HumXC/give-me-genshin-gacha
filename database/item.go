package database

import (
	"give-me-genshin-gacha/models"
)

type ItemMan interface {
	BaseMan[models.Item]
	Put(...models.Item) ([]models.Item, error)
}

type itemMan struct {
	baseMan[models.Item]
}

func (m *itemMan) Put(items ...models.Item) ([]models.Item, error) {
	add := make([]models.Item, 0, 5)
	for i := 0; i < len(items); i++ {
		var item models.Item
		err := m.db.Where(models.Item{ItemID: items[i].ItemID}).
			FirstOrInit(&item).Error
		if err != nil {
			return nil, err
		}
		if item.ID == 0 {
			add = append(add, items[i])
			err = m.db.Create(&items[i]).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return add, nil
}
