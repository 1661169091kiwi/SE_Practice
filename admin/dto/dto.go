package dto

// 通用表格单元格
type Cell struct {
	Value string `json:"value"`
}

// 1. 发布赛事请求
type PublishEventReq struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Time struct {
		Year  string `json:"year"`
		Month string `json:"month"`
		Day   string `json:"day"`
	} `json:"time"`
	Location string `json:"location"`
}

// 2. & 5. 通用表格数据请求 (编排、成绩)
type GridDataReq struct {
	Type string   `json:"type"`
	Data [][]Cell `json:"data"` // 二维数组，对应前端的表格行和列
}

// 3. 保存队伍请求
type SaveTeamReq struct {
	Type     string   `json:"type"`
	TeamName string   `json:"teamName"`
	Captain  string   `json:"captain"`
	Members  [][]Cell `json:"members"`
}

// 4. 保存规则请求
type SaveRuleReq struct {
	Type      string `json:"type"`
	MatchTime struct {
		Year  string `json:"year"`
		Month string `json:"month"`
		Day   string `json:"day"`
	} `json:"matchTime"`
	WinCondition string `json:"winCondition"`
}
