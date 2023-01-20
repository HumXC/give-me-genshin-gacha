package app

import (
	"give-me-genshin-gacha/models"
	"give-me-genshin-gacha/webview"
)

type PieData struct {
	UsedCost   int    `json:"usedCost"`
	Arms3Total int    `json:"arms3Total"`
	Arms4Total int    `json:"arms4Total"`
	Arms5Total int    `json:"arms5Total"`
	Role4Total int    `json:"role4Total"`
	Role5Total int    `json:"role5Total"`
	GachaType  string `json:"gachaType"`
}
type GachaPieTotals struct {
	T301 []models.GachaTotal `json:"t301"`
	T302 []models.GachaTotal `json:"t302"`
	T200 []models.GachaTotal `json:"t200"`
	T100 []models.GachaTotal `json:"t100"`
}
type GachaPieDate struct {
	UsedCosts []models.GachaUsedCost `json:"usedCosts"`
	Totals    GachaPieTotals         `json:"totals"`
}

type GachaMan struct {
	db      *models.GachaDB
	WebView *webview.WebView
}

func (g *GachaMan) GetNumWithLast(uid, gachaType, id string) int {
	result, err := g.db.GetNumWithLast(uid, gachaType, id)
	if err != nil {
		g.WebView.Alert.Error("从数据库获取计数时出现错误:" + err.Error())
		return result
	}
	return result
}
func (g *GachaMan) GetLogs(uid, gachaType string, num, page int) []models.GachaLog {
	result, err := g.db.GetLogs(uid, gachaType, num, page)
	if err != nil {
		g.WebView.Alert.Error("从数据库获取记录时出现错误:" + err.Error())
		return make([]models.GachaLog, 0)
	}
	return result
}

// 饼图数据
func (g *GachaMan) GetPieDatas(uid string) GachaPieDate {
	result := GachaPieDate{}
	r, err := g.db.GetTotals(uid)
	if err != nil {
		g.WebView.Alert.Error("从数据库获取记录时出现错误:" + err.Error())
		return result
	}
	c, err := g.db.GetUsedCost(uid)
	if err != nil {
		g.WebView.Alert.Error("从数据库获取记录时出现错误:" + err.Error())
		return result
	}
	result.UsedCosts = c
	result.Totals.T301 = r["301"]
	result.Totals.T302 = r["302"]
	result.Totals.T200 = r["200"]
	result.Totals.T100 = r["100"]
	return result
}
func (g *GachaMan) GetUids() []string {
	result, err := g.db.GetUids()
	if err != nil {
		g.WebView.Alert.Error("无法从数据库获取 uid:" + err.Error())
	}
	return result
}
