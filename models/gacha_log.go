package models

import (
	"time"

	"gorm.io/gorm"
)

type FullGachaLog struct {
	// 在同一个 GachaType 中，此物品与上一个同 RankType 的物品之间所花费的抽数
	// 例如 出了 5 星之后的第 3 发又出了 5 星，那后者的 Cost 就是 3
	Cost            int       `json:"cost"`
	GachaType       string    `json:"gachaType"`
	OriginGachaType string    `json:"originGachaType"`
	RankType        int       `json:"rankType"`
	ItemType        ItemType  `json:"itemType"`
	Name            string    `json:"name"`
	Time            time.Time `json:"time"`
	ItemID          uint      `json:"itemId"`
	ID              int       // 此处的 id 用于计算 Cost
}
type GachaInfo struct {
	GachaType string `json:"gachaType"`
	AllCount  int    `json:"allCount"`
	Avatar5   int    `json:"avatar5"`
	Avatar4   int    `json:"avatar4"`
	Weapon5   int    `json:"weapon5"`
	Weapon4   int    `json:"weapon4"`
	Weapon3   int    `json:"weapon3"`
}
type GachaLog struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	// OriginGachaType 是米哈游自带的 gacha_type, 会有 400 的值
	// GachaType 是 uigf_gacha_type
	// 见 https://github.com/DGP-Studio/Snap.Genshin/wiki/StandardFormat
	OriginGachaType string
	GachaType       string
	UserID          uint
	Time            time.Time
	ItemID          uint
	Count           int
	LogID           uint
}
type LogDB struct {
	db *gorm.DB
}

// FIXME 删除祈愿的语言信息（Name）前端获取对应的语言名称获取物品名称
// 获取祈愿记录，最后的参数用与筛选，a和w表示角色和武器，后面的数字代表rank_type
func (d *LogDB) GetFullGacha(page int, uid uint, lang, gachaType string, a4, a5, w3, w4, w5, desc bool) ([]FullGachaLog, error) {
	offset := page * 100
	result := make([]FullGachaLog, 0)
	tx := d.db.Model(&GachaLog{}).Debug().
		Where("user_id = ? AND gacha_type = ?", uid, gachaType).Offset(offset)
	var tx2 *gorm.DB
	or := func(it, rt int) {
		if tx2 == nil {
			tx2 = d.db.Or("item_type = ? AND rank_type = ?", it, rt)
		} else {
			tx2.Or("item_type = ? AND rank_type = ?", it, rt)
		}
	}
	if a4 {
		or(1, 4)
	}
	if a5 {
		or(1, 5)
	}
	if w3 {
		or(0, 3)
	}
	if w4 {
		or(0, 4)
	}
	if w5 {
		or(0, 5)
	}
	tx.Where(tx2)
	tx.Joins("join items on gacha_logs.item_id=items.id").
		Select(
			"gacha_type",
			"origin_gacha_type",
			"rank_type",
			"item_type",
			"time",
			"gacha_logs.id as id",
			"item_id",
		)
	if !desc {
		tx.Order("gacha_logs.id DESC")
	}
	err := tx.
		Limit(100).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(result); i++ {
		// 计算 cost
		c, err := d.cost(uid, result[i].ID, result[i].RankType, result[i].GachaType)
		if err != nil {
			return nil, err
		}
		result[i].Cost = c
		// 获取对应语言的名称
		n, err := d.getName(result[i].ItemID, lang)
		if err != nil {
			return nil, err
		}
		result[i].Name = n
	}
	return result, nil
}
func (d *LogDB) getName(id uint, lang string) (n string, err error) {
	err = d.db.Model(&Item{
		Model: gorm.Model{ID: id},
	}).Select("value").Where("lang=?", lang).Association("Names").Find(&n)
	return
}

// 计算 Cost，Cost 定义见 FullGachaLog 结构体定义
func (d *LogDB) cost(uid uint, id int, rankType int, gachaType string) (int, error) {
	if rankType == 3 {
		return 1, nil
	}
	// 查询相同 GachaType 相同 RankType 的上一个物品
	tmp := make([]int, 2)
	err := d.db.Model(&GachaLog{}).
		Select("gacha_logs.id").
		Joins("join items on gacha_logs.item_id=items.id").
		Where("user_id = ? AND gacha_type = ? AND rank_type >= ? AND gacha_logs.id BETWEEN 0 and ? ", uid, gachaType, rankType, id).
		Order("gacha_logs.id DESC").
		Limit(2).
		Scan(&tmp).Error
	if err != nil {
		return 0, err
	}
	var result int64
	if len(tmp) < 2 {
		// 说明在之前的记录里已经没有更多星的物品了
		err := d.db.Model(&GachaLog{}).
			Select("id").
			Where("user_id = ? AND gacha_type = ? AND id <= ?", uid, gachaType, tmp[0]).
			Count(&result).Error
		return int(result), err
	}

	err = d.db.Model(&GachaLog{}).
		Select("id").
		Where("user_id = ? AND gacha_type = ? AND id BETWEEN ? and ? ", uid, gachaType, tmp[1], tmp[0]).
		Count(&result).Error
	return int(result) - 1, err
}
func (d *LogDB) GetInfo(uid int) ([]GachaInfo, error) {
	result := make([]struct {
		Count     int
		GachaType string
		RankType  int
		ItemType  int
	}, 0)
	err := d.db.Model(&GachaLog{}).
		Where("user_id = ?", uid).
		Select("COUNT(*) as count", "gacha_type", "rank_type", "item_type").
		Joins("join items on gacha_logs.item_id=items.id").
		Group("gacha_type").
		Group("rank_type").
		Group("item_type").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	infosMap := make(map[string]GachaInfo, 0)
	for _, v := range result {
		info := infosMap[v.GachaType]
		info.GachaType = v.GachaType
		info.AllCount += v.Count
		if v.ItemType == ItemAvatar {
			switch v.RankType {
			case 4:
				info.Avatar4 += v.Count
			case 5:
				info.Avatar5 += v.Count
			}
			infosMap[v.GachaType] = info
			continue
		}
		switch v.RankType {
		case 3:
			info.Weapon3 += v.Count
		case 4:
			info.Weapon4 += v.Count
		case 5:
			info.Weapon5 += v.Count
		}
		infosMap[v.GachaType] = info
	}
	infos := make([]GachaInfo, 0)
	for _, v := range infosMap {
		infos = append(infos, v)
	}
	return infos, nil
}

// 添加条目到数据库
func (d *LogDB) Add(items []GachaLog) error {
	if len(items) == 0 {
		return nil
	}
	return d.db.Create(items).Error
}

// TODO: 测试多用户时此方法的正确性
// 获取每个池子每个 uid 最新的一次祈愿记录
func (d *LogDB) EndLogIDs() (map[string]map[uint]uint, error) {
	result := make(map[string]map[uint]uint, 0)
	col := make([]GachaLog, 0)
	subQuery := d.db.Model(&GachaLog{}).Order("id DESC")
	err := d.db.Table("(?) as u", subQuery).Select("user_id", "log_id", "gacha_type").Group("gacha_type").Find(&col).Error
	for _, log := range col {
		if result[log.GachaType] == nil {
			result[log.GachaType] = make(map[uint]uint, 0)
		}
		result[log.GachaType][log.UserID] = log.LogID
	}
	return result, err
}
