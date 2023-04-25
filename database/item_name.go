package database

import (
	"give-me-genshin-gacha/models"
)

type ItemNameMan interface {
	BaseMan[models.ItemName]
	Put(...models.ItemName) ([]models.ItemName, error)
}

type itemNameMan struct {
	baseMan[models.ItemName]
}

func (m *itemNameMan) Put(names ...models.ItemName) ([]models.ItemName, error) {
	add := make([]models.ItemName, 0, 5)
	for i := 0; i < len(names); i++ {
		var name models.ItemName
		err := m.db.Where(models.ItemName{ItemID: names[i].ItemID, Lang: names[i].Lang}).
			FirstOrInit(&name).Error
		if err != nil {
			return nil, err
		}
		if name.ID == 0 {
			add = append(add, names[i])
			err = m.db.Create(&names[i]).Error
			if err != nil {
				return nil, err
			}
		}
	}
	return add, nil
}
