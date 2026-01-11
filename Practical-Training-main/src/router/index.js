import { createRouter, createWebHistory } from 'vue-router'
import LoginPassword from '@/views/login/login_password.vue'
import LoginPhone from '@/views/login/login_phone.vue'
import Register from '@/views/login/regis_index.vue'
import Recommend from '@/views/layout/recommend_index.vue'
import SearchIndex from '@/views/layout/search_index.vue'
import SearchResult from '@/views/layout/search_result.vue'
import MyDate from '@/views/data/data_index.vue'
import Event from '@/views/event/eventPage.vue'
import Menu from '@/views/event/eventMenu.vue'
import ToStart from '@/views/event/event_to_start.vue'
import Underway from '@/views/event/event_underway.vue'
import SporterScore from '@/views/event/sporter_score.vue'
import TeamScore from '@/views/event/team_score.vue'
import InfoIndex from '@/views/myinfo/info_index.vue'
import MyFollow from '@/views/myinfo/my_follow.vue'
import MyInformation from '@/views/myinfo/my_information.vue'
import MyRating from '@/views/myinfo/my_rating.vue'
import MyHistory from '@/views/myinfo/my_history.vue'
import PasswordChange from '@/views/myinfo/password_change.vue'
import Layout from '@/views/layout/layout_index.vue'
import ChangeMyName from '@/views/myinfo/change_myname.vue'
import Team from '@/views/event/team_score.vue'
import Admin from '@/views/admin/admin_main.vue'
import AdEvent from '@/views/admin/admin_event.vue'
import Rule from '@/views/admin/admin_rule.vue'
import Grade from '@/views/admin/admin_grade.vue'
import AdTeam from '@/views/admin/admin_team.vue'
import OngoingMatch from '@/views/event/OngoingMatch.vue'
import BaseMatch from '@/views/event/BaseMatch.vue'
import Start from '@/views/event/startMatch.vue'
import Finish from '@/views/event/finishMatch.vue'
import Football from '@/views/event/football_team.vue'
import TableTennis from '@/views/event/table_tennis_team.vue'
import Badminton from '@/views/event/badminton_team.vue'
import Water from '@/views/event/water_sports_team.vue'











const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/login_phone', component: LoginPhone },
    { path: '/login_password', component: LoginPassword },
    { path: '/register', component: Register },
    {
      path: '/',
      component: Layout,
      redirect: '/recommend',
      children: [
        {path: '/event', component:Event,
          children:[
            {path: 'menu', component: Menu, props: true },
          ], props: true
        },
        { path: '/recommend', component: Recommend },
        { path: '/data', component: MyDate },
        { path: '/InfoIndex', component: InfoIndex },
      ],
    },

    {path:'/event/tostart',component: ToStart},
    {path:'/event/underway',component:Underway},
    {path:'/event/finished',component:Underway},
    {path:'/event/teampage/football',component:Football},
    {path:'/event/teampage/pingpong',component:TableTennis},
    {path:'/event/teampage/badminton',component:Badminton},
    {path:'/event/teampage/water',component:Water},
    {path:'/event/teampage/basketball',component:Team},//队伍页面的路径可能需要优化
    { path: '/search', component: SearchIndex },
    { path: '/search_result', component: SearchResult },
    { path: '/sporter_score', component: SporterScore },
    { path: '/team_score', component: TeamScore },
    { path: '/my_follow', component: MyFollow },
    { path: '/my_information', component: MyInformation },
    { path: '/my_rating', component: MyRating },
    { path: '/my_history', component: MyHistory },
    { path: '/password_change', component: PasswordChange },
    { path: '/change_myname', component: ChangeMyName },
    { path: '/admin', component: Admin },
    { path: '/arrange', component: AdEvent },
    { path: '/rule', component: Rule },
    { path: '/grade', component: Grade },
    { path: '/team', component: AdTeam },
    { path: '/base/:matchId', component: BaseMatch, props: true },
    { path: '/ongoing/:sportType', component: OngoingMatch, name: 'OngoingMatch',
      props: route => ({
      sportType: route.params.sportType, // 从路径参数获取 sportType
      matchType: route.query.matchType // 从查询参数获取 matchType
      })
     },
    {
      path: '/event-to-start/:sportType', // 动态路由参数
      name: 'EventToStart',
      component: Start,
      props: route => ({
        sportType: route.params.sportType, // 从路径参数获取 sportType
        matchType: route.query.matchType // 从查询参数获取 matchType
      }),
      // 可选：路由参数验证
      beforeEnter: (to, from, next) => {
        const validSports = ['football', 'basketball', 'tabletennis', 'badminton', 'water']
        const validTypes = ['singles', 'doubles', 'team']

        // 验证运动类型
        if (!validSports.includes(to.params.sportType)) {
          return next('/404') // 跳转到404页面
        }

        // 验证比赛类型（羽毛球特有）
        if (to.params.matchType && !validTypes.includes(to.params.matchType)) {
          return next('/404')
        }

        next()
      }
    },
    { path: '/finished/:sportType', component: Finish, name: 'finishMatch',
      props: route => ({
      sportType: route.params.sportType, // 从路径参数获取 sportType
      matchType: route.query.matchType // 从查询参数获取 matchType
      }),
      beforeEnter: (to, from, next) => {
        const validSports = ['football', 'basketball', 'tabletennis', 'badminton', 'water']
        const validTypes = ['singles', 'doubles', 'team']

        // 验证运动类型
        if (!validSports.includes(to.params.sportType)) {
          return next('/404') // 跳转到404页面
        }

        // 验证比赛类型（羽毛球特有）
        if (to.params.matchType && !validTypes.includes(to.params.matchType)) {
          return next('/404')
        }

        next()
      }
     },
  ],
})

export default router
