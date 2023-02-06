package assets

import (
	"give-me-genshin-gacha/gacha"
	"give-me-genshin-gacha/models"
)

type ItemStore struct {
	itemDB      *models.ItemDB
	loadedLang  map[string]struct{}
	assetsStore gacha.AssetsStore
}

// 加载一种语言的资源
func (i *ItemStore) Load(lang string) error {
	if _, ok := i.loadedLang[lang]; ok {
		return nil
	}
	obcItems, err := i.assetsStore.Get(lang)
	if err != nil {
		return err
	}
	for _, v := range obcItems {
		item, err := i.itemDB.GetWithID(v.ItemID)
		if err != nil {
			return err
		}
		item.ID = v.ItemID
		item.ItemType = v.ItemType
		item.RankType = v.RankType
		newItem := setItemName(item, lang, v.Name)
		err = i.itemDB.Update(newItem)
		if err != nil {
			return err
		}
	}
	i.loadedLang[lang] = struct{}{}
	return nil
}

// 卸载语言资源，如果已有的语言资源过期了则调用此方法
// 例如有原神发布新角色了，但是语言资源还没更新，就会存在无法获取物品名称的问题
func (i *ItemStore) UnLoad(lang string) {
	delete(i.loadedLang, lang)
}
func setItemName(item models.Item, lang, name string) models.Item {
	switch lang {
	case "zh-cn":
		item.ZhCN = name
	case "zh-tw":
		item.ZhTW = name
	case "ja-jp":
		item.JaJP = name
	case "en-us":
		item.EnUS = name
	}
	return item
}

func NewItemStore(itemDB *models.ItemDB) (*ItemStore, error) {
	is := ItemStore{
		itemDB:     itemDB,
		loadedLang: make(map[string]struct{}),
	}
	loadedLang, err := is.itemDB.LoadedLang()
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(loadedLang); i++ {
		is.loadedLang[loadedLang[i]] = struct{}{}
	}
	return &is, nil
}
