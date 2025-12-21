import { createRouter, createWebHistory } from 'vue-router'
import Event from '@/views/admin_event.vue'
import Grade from '@/views/admin_grade.vue'
import Main from '@/views/admin_main.vue'
import Rule from '@/views/admin_rule.vue'
import Team from '@/views/admin_team.vue'


// 导入所有页面组件
const LoginView = () => import('../views/LoginView.vue')

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
      path:'/arrange',
      name:'arrange',
      component:Event,
      meta:{requiresAuth:true},
    },
    {
      path:'/grade',
      name:'grade',
      component:Grade,
      meta:{requiresAuth:true},
    },
    {
      path:'/admin',
      name:'admin',
      component:Main,
      meta:{requiresAuth:true},
    },
    {
      path:'/rule',
      name:'rule',
      component:Rule,
      meta:{requiresAuth:true},
    },
    {
      path:'/team',
      name:'team',
      component:Team,
      meta:{requiresAuth:true},
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
