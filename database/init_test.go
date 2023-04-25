package database_test

import (
	"give-me-genshin-gacha/database"
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
	"os"
	"testing"
)

var TestDir = "../test/database"

func TestMain(m *testing.M) {
	// _ = os.RemoveAll(TestDir)
	os.MkdirAll(TestDir, 0775)
	m.Run()
}

func TestXxx(t *testing.T) {
	lang := "zh-tw"
	db, err := database.NewDB(TestDir + "/data.db")
	if err != nil {
		panic(err)
	}
	avatar, weapon, err := gacha.GetGameItem(lang)
	if err != nil {
		panic(err)
	}
	items := make([]models.Item, 0, len(avatar)+len(weapon))
	names := make([]models.ItemName, 0, len(avatar)+len(weapon))
	for _, a := range avatar {
		items = append(items, models.Item{
			ItemID:   a.ItemID,
			RankType: a.RankType,
			ItemType: models.ItemTypeAvatar,
		})
		names = append(names, models.ItemName{
			ItemID: a.ItemID,
			Lang:   lang,
			Name:   a.Name,
		})
	}
	for _, w := range weapon {
		items = append(items, models.Item{
			ItemID:   w.ItemID,
			RankType: w.RankType,
			ItemType: models.ItemTypeWeapon,
		})
		names = append(names, models.ItemName{
			ItemID: w.ItemID,
			Lang:   lang,
			Name:   w.Name,
		})
	}
	db.Items.Put(items...)
	db.ItemNames.Put(names...)
}
