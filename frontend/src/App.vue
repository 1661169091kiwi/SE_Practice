<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { get } from '@/utils/http'

const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  // 检查会话有效性：如果有 token 但对应用户不存在（例如数据库重置），则登出
  if (authStore.token && authStore.studentId) {
    try {
      // 尝试获取用户信息以验证 token 和用户状态
      // http.js 会自动处理 401 Unauthorized
      const res = await get(`/user/profile/${authStore.studentId}`)
      if (res && res.code === 200 && res.data) {
        if (res.data.avatar_url) {
          authStore.setAvatar(res.data.avatar_url)
        }
      }
    } catch (err) {
      console.warn('Session check failed:', err)
      // 如果是 404 (用户不存在) 或其他非网络错误，清除无效会话
      if (err.status === 404) {
        authStore.clearAuth()
        router.push('/login')
      }
    }
  } else {
     // 如果没有登录，确保我们在登录页
     if (router.currentRoute.value.path !== '/login' && router.currentRoute.value.path !== '/register') {
         router.push('/login')
     }
  }
})
</script>

<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<style>
.app-container {
  min-height: 100vh;
  background-color: var(--background-secondary);
  color: var(--text-primary);
}

/* 页面切换的淡入效果，让整体过渡更柔和 */
.app-container > * {
  animation: app-fade-in 0.25s var(--transition-bezier);
}

@keyframes app-fade-in {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
