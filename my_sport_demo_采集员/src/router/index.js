import { createRouter, createWebHistory } from 'vue-router'

// 导入所有页面组件
const LoginView = () => import('../views/LoginView.vue')
const EventSelectionView = () => import('../views/EventSelectionView.vue')
const BaseDataCollection = () => import('../views/BaseDataCollection.vue')
const FootballDataCollection = () => import('../views/FootballDataCollection.vue')
const BasketballDataCollection = () => import('../views/BasketballDataCollection.vue')
const BadmintonDataCollection = () => import('../views/BadmintonDataCollection.vue')
const VolleyballDataCollection = () => import('../views/VolleyballDataCollection.vue')
const WaterSportsDataCollection = () => import('../views/WaterSportsDataCollection.vue')
const DataPreview = () => import('../views/DataPreview.vue')
const HistoryView = () => import('../views/HistoryView.vue')
const ProfileView = () => import('../views/ProfileView.vue')

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
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
      path: '/events',
      name: 'events',
      component: EventSelectionView,
      meta: { requiresAuth: true },
    },
    {
      path: '/data-collection/:sportType/:eventId',
      name: 'data-collection',
      component: BaseDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/football-data/:eventId',
      name: 'football-data',
      component: FootballDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/basketball-data/:eventId',
      name: 'basketball-data',
      component: BasketballDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/badminton-data/:eventId',
      name: 'badminton-data',
      component: BadmintonDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/volleyball-data/:eventId',
      name: 'volleyball-data',
      component: VolleyballDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/water-sports-data/:eventId',
      name: 'water-sports-data',
      component: WaterSportsDataCollection,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/data-preview/:sportType/:eventId',
      name: 'data-preview',
      component: DataPreview,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/history',
      name: 'history',
      component: HistoryView,
      meta: { requiresAuth: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true },
    },
    // 404页面
    {
      path: '/:pathMatch(.*)*',
      redirect: '/login',
    },
  ],
})

// 路由守卫：处理身份验证
// router.beforeEach((to, from, next) => {
//   // 检查路由是否需要身份验证
//   const requiresAuth = to.meta.requiresAuth !== false

//   // 检查用户是否已登录（通过localStorage中的token）
//   const isAuthenticated = !!localStorage.getItem('authToken')

//   // 如果路由需要身份验证且用户未登录，则重定向到登录页
//   if (requiresAuth && !isAuthenticated) {
//     next('/login')
//   }
//   // 如果用户已登录且尝试访问登录页，则重定向到赛事列表页
//   else if (isAuthenticated && to.name === 'login') {
//     next('/events')
//   }
//   // 其他情况正常导航
//   else {
//     next()
//   }
// })

export default router
