package database

import (
	"give-me-genshin-gacha/models"
)

type ItemMan interface {
	BaseMan[models.Item]
	Put(...models.Item) error
}

type itemMan struct {
	baseMan[models.Item]
}

func (m *itemMan) Put(items ...models.Item) error {
	for i := 0; i < len(items); i++ {
		var item models.Item
		err := m.db.Where(models.Item{ItemID: items[i].ItemID}).
			FirstOrInit(&item).Error
		if err != nil {
			return err
		}
		if item.ID == 0 {
			err = m.db.Create(&items[i]).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
