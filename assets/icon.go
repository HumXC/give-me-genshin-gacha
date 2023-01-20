package assets

const API_ICON string = "https://sg-public-api-static.hoyolab.com/common/map_user/ys_obc/v1/map/game_item?map_id=2&app_sn=ys_obc&lang=zh-cn"

type Icon struct {
	Url  string `json:"icon"`
	Name string `json:"name"`
}

// 为 Server 提供处理 icon 请求的能力
type iconHandler struct {
	IconAvatar []Icon
	IconWeapon []Icon
	IconDir    string
}
