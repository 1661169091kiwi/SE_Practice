<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '../utils/http'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const studentId = ref('')
const password = ref('')
const showPassword = ref(false)
const loginError = ref('')

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const handleLogin = async () => {
  const studentIdRegex = /^\w{5,20}$/
  if (!studentIdRegex.test(studentId.value)) {
    loginError.value = '请输入正确的学号'
    return
  }
  if (!password.value) {
    loginError.value = '请输入学号和密码'
    return
  }
  try {
    const data = await post('/login', { student_id: studentId.value, password: password.value }, { auth: false })
    if (data.code === 200 && data.data && data.data.token) {
      authStore.setToken(data.data.token)
      const userRole = data.data.user?.role
      const userStudentId = data.data.user?.student_id
      const userAvatar = data.data.user?.avatar_url
      
      if (userRole) {
        authStore.setRole(userRole)
      }
      if (data.data.user?.roles) {
        authStore.setRoles(data.data.user.roles)
      }
      if (userStudentId) {
        authStore.setStudentId(userStudentId)
      }
      if (userAvatar) {
        authStore.setAvatar(userAvatar)
      }
      
      loginError.value = ''
      
      // 根据角色跳转
      const roles = authStore.roles
      if (roles.includes('admin') || roles.includes('super_admin')) {
        router.push('/admin/dashboard')
      } else if (roles.includes('collector')) {
        router.push('/events')
      } else {
        router.push('/student/events')
      }
    } else {
      loginError.value = data.message || '登录失败'
    }
  } catch {
    loginError.value = '登录失败，请稍后重试'
  }
}

const handleForgotPassword = () => {
  // 跳转到密码找回页面（暂时简单提示）
  alert('密码找回功能暂未实现')
}
</script>

<template>
  <div class="login-container">
    <!-- 动态背景装饰 -->
    <div class="bg-shape shape-1"></div>
    <div class="bg-shape shape-2"></div>
    
    <div class="login-content fade-in">
      <div class="login-header">
        <h1>中大校园体育赛事平台</h1>
        <p>统一认证入口</p>
      </div>

      <div class="login-form">
        <div class="form-group">
          <label for="studentId">学号</label>
          <input
            id="studentId"
            v-model="studentId"
            type="text"
            placeholder="请输入学号"
            class="form-input"
          />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <div class="password-input-container">
            <input
              id="password"
              v-model="password"
              :type="showPassword ? 'text' : 'password'"
              placeholder="请输入密码"
              class="form-input"
            />
            <button type="button" class="password-toggle" @click="togglePasswordVisibility">
              {{ showPassword ? '隐藏' : '显示' }}
            </button>
          </div>
        </div>

        <div v-if="loginError" class="error-message error-shake">
          {{ loginError }}
        </div>

        <button type="button" class="login-button" @click="handleLogin">
          <span>登录</span>
        </button>

        <div class="login-footer">
          <button type="button" class="forgot-password" @click="handleForgotPassword">
            忘记密码？
          </button>
          <button type="button" class="forgot-password" @click="router.push('/register')">
            还没有账号？去注册
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: var(--spacing-3xl) var(--spacing-xl);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

/* 动态背景图形 */
.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  z-index: 0;
  opacity: 0.6;
}

.shape-1 {
  width: 300px;
  height: 300px;
  background: var(--primary-gradient-start);
  top: -50px;
  left: -50px;
  animation: float 8s infinite ease-in-out;
}

.shape-2 {
  width: 400px;
  height: 400px;
  background: var(--accent-gradient-end);
  bottom: -100px;
  right: -100px;
  animation: float 10s infinite ease-in-out reverse;
}

@keyframes float {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, 40px); }
}

.login-content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-2xl);
  text-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.login-header h1 {
  font-size: 26px;
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
  background: linear-gradient(to right, var(--primary-dark), var(--primary-active));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.login-header p {
  font-size: 15px;
  color: var(--text-secondary);
  font-weight: 500;
}

.login-form {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: var(--border-radius-xl);
  padding: 32px 28px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.4);
  transition: transform 0.3s ease;
}

.login-form:hover {
  transform: translateY(-2px);
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 600;
}

.form-input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: 2px solid transparent;
  background: rgba(255, 255, 255, 0.9);
  border-radius: var(--border-radius-lg);
  font-size: 15px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 6px rgba(0,0,0,0.02);
}

.form-input:focus {
  border-color: var(--primary-color);
  background: #fff;
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.2);
}

.password-input-container {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--primary-color);
  font-weight: 600;
  font-size: 13px;
  padding: 4px 8px;
  border-radius: 4px;
}

.password-toggle:hover {
  background: rgba(79, 172, 254, 0.1);
}

.login-button {
  width: 100%;
  height: 50px;
  background: linear-gradient(135deg, var(--primary-gradient-start) 0%, var(--primary-gradient-end) 100%);
  color: white;
  border-radius: var(--border-radius-lg);
  font-size: 16px;
  font-weight: 600;
  margin-top: 12px;
  box-shadow: 0 4px 15px rgba(79, 172, 254, 0.4);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 172, 254, 0.6);
}

.login-button:active {
  transform: translateY(1px);
}

.error-message {
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
  padding: 10px;
  border-radius: 8px;
  font-size: 13px;
  margin-bottom: 16px;
  text-align: center;
  border: 1px solid rgba(255, 77, 79, 0.2);
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.forgot-password {
  color: var(--text-tertiary);
  font-size: 14px;
  transition: color 0.3s;
}

.forgot-password:hover {
  color: var(--primary-color);
  text-decoration: underline;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .login-container {
    padding: 20px;
    align-items: center; /* 保持居中 */
  }
  
  .login-content {
    padding: 0;
    width: 100%;
  }

  .login-form {
    padding: 24px 20px;
  }

  .login-header h1 {
    font-size: 22px;
  }
  
  .shape-1 {
    width: 150px;
    height: 150px;
    top: -40px;
    left: -40px;
  }
  
  .shape-2 {
    width: 180px;
    height: 180px;
    bottom: -50px;
    right: -50px;
  }
}
</style>
