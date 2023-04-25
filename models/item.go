package models

type ItemType = uint8
type RankType = uint8

const (
	ItemTypeWeapon ItemType = 0
	ItemTypeAvatar ItemType = 1

	RankType3 RankType = 3
	RankType4 RankType = 4
	RankType5 RankType = 5
)

type ItemName struct {
	Model
	ItemID uint64
	Lang   string
	Name   string
}
type Item struct {
	Model
	ItemID   uint64
	RankType RankType // 星数
	ItemType ItemType // 武器还是角色
}
