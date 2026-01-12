import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

// 导入所有页面组件
const LoginView = () => import('../views/LoginView.vue')
const EventSelectionView = () => import('../views/EventSelectionView.vue')
const BaseDataCollection = () => import('../views/BaseDataCollection.vue')
const FootballDataCollection = () => import('../views/FootballDataCollection.vue')
const BasketballDataCollection = () => import('../views/BasketballDataCollection.vue')
const BadmintonDataCollection = () => import('../views/BadmintonDataCollection.vue')
const VolleyballDataCollection = () => import('../views/VolleyballDataCollection.vue')
const DataPreview = () => import('../views/DataPreview.vue')
const HistoryView = () => import('../views/HistoryView.vue')
const ProfileView = () => import('../views/ProfileView.vue')
const AdminDashboard = () => import('../views/AdminDashboard.vue')

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { requiresAuth: false },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('../views/RegisterView.vue'),
      meta: { requiresAuth: false },
    },
    // 采集员端路由
    {
      path: '/events',
      name: 'events',
      component: EventSelectionView,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    // 学生端路由
    {
      path: '/student/events',
      name: 'student-events',
      component: () => import('../views/StudentEventList.vue'),
      meta: { roles: ['student', 'collector', 'admin', 'athlete'] }
    },
    {
      path: '/match/:id',
      name: 'match-detail',
      component: () => import('../views/MatchDetailView.vue'),
      meta: { roles: ['student', 'collector', 'admin', 'athlete'] }
    },
    {
      path: '/student/standings',
      name: 'student-standings-index',
      component: () => import('../views/StandingsView.vue'),
      meta: { roles: ['student', 'collector', 'admin', 'athlete'] }
    },
    {
      path: '/student/standings/:eventId',
      name: 'student-standings',
      component: () => import('../views/StandingsView.vue'),
      meta: { roles: ['student', 'collector', 'admin', 'athlete'] }
    },
    {
      path: '/data-collection/:sportType/:eventId',
      name: 'data-collection',
      component: BaseDataCollection,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/football-data/:eventId',
      name: 'football-data',
      component: FootballDataCollection,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/basketball-data/:eventId',
      name: 'basketball-data',
      component: BasketballDataCollection,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/badminton-data/:eventId',
      name: 'badminton-data',
      component: BadmintonDataCollection,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/volleyball-data/:eventId',
      name: 'volleyball-data',
      component: VolleyballDataCollection,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/data-preview/:sportType/:eventId',
      name: 'data-preview',
      component: DataPreview,
      props: true,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/history',
      name: 'history',
      component: HistoryView,
      meta: { requiresAuth: true, roles: ['collector', 'admin'] },
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true },
    },
    // 管理员端路由
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: AdminDashboard,
      meta: { requiresAuth: true, roles: ['admin'] },
    },
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login',
    },
  ],
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.meta.requiresAuth !== false
  const token = authStore.token
  const role = authStore.role
  let roles = authStore.roles || []
  if (roles.length === 0 && role) {
    roles = [role]
  }

  if (requiresAuth && !token) {
    next('/login')
  } else if (token && to.name === 'login') {
    // 即使已登录，也允许访问登录页（方便切换账号或重新登录）
    next()
  } else {
    // 检查角色权限
    if (to.meta.roles) {
      const hasPermission = to.meta.roles.some(r => roles.includes(r))
      if (hasPermission) {
        next()
      } else {
        // 角色不符
        if (roles.length > 0) {
          if (roles.includes('admin')) {
            next('/admin/dashboard')
          } else if (roles.includes('collector')) {
            next('/events')
          } else {
            next('/student/events')
          }
        } else {
          authStore.clearAuth()
          next('/login')
        }
      }
    } else {
      next()
    }
  }
})

export default router
