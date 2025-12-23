package router

import (
	"net/http"

	"se_practice/backend/internal/handler"
	"se_practice/backend/internal/middleware"
)

// New registers routes and returns an http.Handler.
func New() http.Handler {
	mux := http.NewServeMux()

	// health
	mux.HandleFunc("/api/health", handler.Health)

	// auth
	mux.HandleFunc("/api/register", handler.Register)
	mux.HandleFunc("/api/login", handler.Login)
	mux.HandleFunc("/api/user/change-password", handler.ChangePassword) // 修改密码
	// 用户资料
	mux.HandleFunc("/api/user/profile/", handler.GetUserProfile)                                                        // expects /api/user/profile/{student_id}
	mux.HandleFunc("/api/user/avatar/", handler.UpdateAvatar)                                                           // 更新用户头像 /api/user/avatar/{student_id}
	mux.HandleFunc("/api/user/college", handler.UpdateCollege)                                                          // 更新用户学院信息
	mux.HandleFunc("/api/user/apply-collector", handler.ApplyCollector)                                                 // 申请成为采集员
	mux.HandleFunc("/api/user/apply-athlete", handler.ApplyAthlete)                                                     // 申请成为运动员
	mux.HandleFunc("/api/admin/collectors/pending", handler.ListPendingCollectors)                                      // 待审核列表
	mux.HandleFunc("/api/admin/collectors/approve", handler.ApproveCollector)                                           // 审核通过
	mux.HandleFunc("/api/admin/collectors/list", handler.ListCollectors)                                                // 所有采集员列表
	mux.Handle("/api/admin/collectors/delete", middleware.RoleAuth("admin")(http.HandlerFunc(handler.DeleteCollector))) // 删除采集员
	mux.HandleFunc("/api/admin/athletes/pending", handler.GetPendingAthletes)                                           // 待审核运动员列表
	mux.HandleFunc("/api/admin/athletes/approve", handler.ApproveAthlete)                                               // 审核通过运动员

	// events
	mux.HandleFunc("/api/events", handler.Events) // GET list
	mux.Handle("/api/events/create", middleware.RoleAuth("admin")(http.HandlerFunc(handler.CreateEvent)))
	mux.HandleFunc("/api/events/", handler.GetEventDetail)
	mux.HandleFunc("/api/events/list", handler.ListEvents)
	mux.Handle("/api/events/delete", middleware.RoleAuth("admin")(http.HandlerFunc(handler.DeleteEvent)))

	// teams
	mux.Handle("/api/teams/create", middleware.RoleAuth("admin")(http.HandlerFunc(handler.CreateTeam)))
	mux.Handle("/api/teams/apply-create", middleware.Auth(http.HandlerFunc(handler.ApplyCreateTeam))) // 申请创建队伍
	mux.HandleFunc("/api/teams/list", handler.ListTeams)
	mux.HandleFunc("/api/teams/my", handler.GetMyTeams)                 // 获取我的队伍
	mux.HandleFunc("/api/admin/teams/pending", handler.GetPendingTeams) // 待审核队伍
	mux.HandleFunc("/api/admin/teams/approve", handler.ApproveTeam)     // 审核通过队伍
	mux.Handle("/api/teams/delete", middleware.RoleAuth("admin")(http.HandlerFunc(handler.DeleteTeam)))
	mux.HandleFunc("/api/teams/members", handler.GetTeamMembersHandler)
	mux.Handle("/api/teams/members/pending", middleware.Auth(http.HandlerFunc(handler.GetPendingTeamMembersHandler)))
	mux.Handle("/api/teams/members/approve", middleware.Auth(http.HandlerFunc(handler.ApproveTeamMemberHandler)))
	// RemoveTeamMember now supports both admin and captain (handled inside via claims)
	mux.Handle("/api/teams/remove-member", middleware.Auth(http.HandlerFunc(handler.RemoveTeamMember)))
	mux.Handle("/api/teams/update-name", middleware.Auth(http.HandlerFunc(handler.UpdateTeamName)))

	// matches
	mux.HandleFunc("/api/matches", handler.Matches)
	mux.HandleFunc("/api/matches/", handler.MatchesEntry)
	mux.Handle("/api/matches/create", middleware.RoleAuth("admin")(http.HandlerFunc(handler.CreateMatch)))
	mux.Handle("/api/matches/delete", middleware.RoleAuth("admin")(http.HandlerFunc(handler.DeleteMatch)))
	// 保护更新比分接口（需采集员角色）
	mux.Handle("/api/matches/update-score/", middleware.RoleAuth("collector")(http.HandlerFunc(handler.UpdateMatchScore)))
	// 运动员参加的比赛列表
	mux.HandleFunc("/api/athlete/matches", handler.AthleteMatches)
	// 运动员可报名的比赛列表
	mux.HandleFunc("/api/athlete/available-matches", handler.AvailableMatchesForAthlete)
	// 运动员报名参加比赛
	mux.HandleFunc("/api/athlete/join-match", handler.JoinMatch)
	// 运动员相关功能
	mux.Handle("/api/athlete/info", middleware.Auth(http.HandlerFunc(handler.GetMyAthleteInfo)))                    // 获取我的运动员信息
	mux.Handle("/api/athlete/info/team", middleware.Auth(http.HandlerFunc(handler.GetMyAthleteInfoByTeam)))        // 获取我在指定队伍中的信息
	mux.Handle("/api/athlete/update", middleware.Auth(http.HandlerFunc(handler.UpdateMyAthleteInfo)))             // 更新我的运动员信息
	mux.Handle("/api/athlete/leave-team", middleware.Auth(http.HandlerFunc(handler.LeaveTeam)))                     // 退出队伍
	mux.Handle("/api/athlete/team-members", middleware.Auth(http.HandlerFunc(handler.GetTeamMembersForAthlete)))   // 查看队伍成员
	mux.Handle("/api/athlete/check-captain", middleware.Auth(http.HandlerFunc(handler.CheckIfCaptain)))             // 检查是否是队长

	// match data (积分榜等)
	mux.HandleFunc("/api/match/data", handler.GetMatchData)

	// user subscribed matches
	mux.HandleFunc("/api/user/subscribed-matches", handler.GetUserSubscribedMatches)

	// user subscribe operation
	// 统一订阅接口，支持通过 eventId 直接订阅或 matchId 查找订阅
	mux.HandleFunc("/api/user/subscribe", handler.SubscribeMatch)
	mux.HandleFunc("/api/events/subscribe", handler.SubscribeMatch)

	// collector 命名空间入口：默认鉴权；POST 在 handler 内进行角色校验
	mux.Handle("/api/collector/", middleware.Auth(http.HandlerFunc(handler.CollectorEntry)))

	// 队内聊天相关功能
	mux.Handle("/api/team/chat/messages", middleware.Auth(http.HandlerFunc(handler.CreateMessage)))           // 创建消息
	mux.Handle("/api/team/chat/messages/list", middleware.Auth(http.HandlerFunc(handler.GetTeamMessages)))     // 获取消息列表
	mux.Handle("/api/team/chat/messages/read", middleware.Auth(http.HandlerFunc(handler.MarkMessageAsRead)))  // 标记已读
	mux.Handle("/api/team/chat/votes", middleware.Auth(http.HandlerFunc(handler.CreateVote)))                  // 创建投票
	mux.Handle("/api/team/chat/votes/list", middleware.Auth(http.HandlerFunc(handler.GetTeamVotes)))          // 获取投票列表
	mux.Handle("/api/team/chat/votes/vote", middleware.Auth(http.HandlerFunc(handler.Vote)))                 // 投票
	mux.Handle("/api/team/chat/notifications", middleware.Auth(http.HandlerFunc(handler.CreateNotification))) // 创建通知
	mux.Handle("/api/team/chat/notifications/list", middleware.Auth(http.HandlerFunc(handler.GetTeamNotifications))) // 获取通知列表
	mux.Handle("/api/team/chat/leave-requests", middleware.Auth(http.HandlerFunc(handler.CreateLeaveRequest))) // 创建请假申请
	mux.Handle("/api/team/chat/leave-requests/list", middleware.Auth(http.HandlerFunc(handler.GetLeaveRequests))) // 获取请假申请列表
	mux.Handle("/api/team/chat/leave-requests/review", middleware.Auth(http.HandlerFunc(handler.ReviewLeaveRequest))) // 审核请假申请

	return middleware.CORS(mux)
}
