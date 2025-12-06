package model

import "time"

// Match 比赛模型
type Match struct {
	ID           int64     `json:"match_id"`
	EventID      int64     `json:"event_id"`
	Name         string    `json:"match_name"`
	Round        string    `json:"round"`
	Time         time.Time `json:"match_time"`
	TeamAID      int64     `json:"team_a_id"`
	TeamBID      int64     `json:"team_b_id"`
	ScoreA       int       `json:"score_team_a"`
	ScoreB       int       `json:"score_team_b"`
	HalfScoreA   int       `json:"half_score_team_a"`
	HalfScoreB   int       `json:"half_score_team_b"`
	Status       string    `json:"status"` // not_started/ongoing/finished/cancelled
	Collector1ID int64     `json:"collector1_id"`
	Collector2ID int64     `json:"collector2_id"`
}

// CreateMatchRequest 创建比赛请求
type CreateMatchRequest struct {
	EventID   int64     `json:"event_id"`
	MatchName string    `json:"match_name"`
	Round     string    `json:"round"`
	MatchTime time.Time `json:"match_time"`
	TeamAID   int64     `json:"team_a_id"`
	TeamBID   int64     `json:"team_b_id"`
}

// MatchDetailResponse 比赛详情响应
type MatchDetailResponse struct {
	MatchID    int64       `json:"match_id"`
	EventID    int64       `json:"event_id"`
	MatchName  string      `json:"match_name"`
	Round      string      `json:"round"`
	MatchTime  time.Time   `json:"match_time"`
	TeamA      TeamBrief   `json:"team_a"`
	TeamB      TeamBrief   `json:"team_b"`
	ScoreA     int         `json:"score_team_a"`
	ScoreB     int         `json:"score_team_b"`
	Status     string      `json:"status"`
	Collectors []UserBrief `json:"collectors"`
}

// TeamBrief 队伍简要信息
type TeamBrief struct {
	ID   int64  `json:"team_id"`
	Name string `json:"team_name"`
}

// UserBrief 用户简要信息
type UserBrief struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// MatchDataResponse 赛事数据统一响应
type MatchDataResponse struct {
	DataList []ResultItem `json:"dataList"`
}

// ResultItem 赛事数据单项（用于积分榜/赛程/榜单）
type ResultItem struct {
	Rank   int    `json:"rank,omitempty"`   // 排名
	ImgURL string `json:"imgUrl,omitempty"` // 图片链接
	Name   string `json:"name"`             // 名称
	Played int    `json:"赛,omitempty"`      // 比赛场数
	Won    int    `json:"胜,omitempty"`      // 胜场数
	Drawn  int    `json:"平,omitempty"`      // 平场数
	Lost   int    `json:"负,omitempty"`      // 负场数
	Goals  string `json:"进/失,omitempty"`    // 进/失球数
	Points int    `json:"积分,omitempty"`     // 积分
	// Extra fields for Schedule/History
	MatchTime string `json:"matchTime,omitempty"` // 比赛时间
	Status    string `json:"status,omitempty"`    // 状态
}

// SubscribedMatchResponse 已订阅比赛列表响应
type SubscribedMatchResponse struct {
	SubscribedMatches []SubscribedMatchItem `json:"subscribedMatches"`
}

// SubscribedMatchItem 已订阅比赛单项
type SubscribedMatchItem struct {
	MatchID     string      `json:"matchId"`     // 比赛唯一编号
	MatchTime   string      `json:"matchTime"`   // 比赛时间
	MatchVenue  string      `json:"matchVenue"`  // 比赛地点
	HomeTeam    TeamWrapper `json:"homeTeam"`    // 主队信息
	AwayTeam    TeamWrapper `json:"awayTeam"`    // 客队信息
	MatchStatus string      `json:"matchStatus"` // 比赛赛况 (1-0 / VS)
	MatchState  string      `json:"matchState"`  // 比赛状态 (未开始/进行中/已结束)
}

// TeamWrapper 队伍信息包装 (for SubscribedMatchItem)
type TeamWrapper struct {
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
}

// SubscribeRequest 订阅/取消订阅请求
type SubscribeRequest struct {
	StudentID   string `json:"studentId"`   // 用户学号
	MatchID     string `json:"matchId"`     // 比赛唯一编号
	OperateType int    `json:"operateType"` // 0=新增订阅，1=取消订阅
}

// SubscribeResponse 订阅/取消订阅响应
type SubscribeResponse struct {
	StudentID   string `json:"studentId"`
	MatchID     string `json:"matchId"`
	OperateType int    `json:"operateType"`
}
