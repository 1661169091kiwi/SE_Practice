<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const studentId = ref('')
const password = ref('')
const showPassword = ref(false)
const loginError = ref('')

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

const handleLogin = () => {
  // 这里模拟登录验证
  // 验证学号是否为8位纯数字
  const studentIdRegex = /^\d{8}$/
  if (!studentIdRegex.test(studentId.value)) {
    loginError.value = '请输入8位纯数字学号'
    return
  }
  
  if (studentId.value && password.value) {
    loginError.value = ''
    // 登录成功后跳转到赛事选择页面
    router.push('/events')
  } else {
    loginError.value = '请输入学号和密码'
  }
}

const handleForgotPassword = () => {
  // 跳转到密码找回页面（暂时简单提示）
  alert('密码找回功能暂未实现')
}
</script>

<template>
  <div class="login-container">
    <div class="login-header">
      <h1>赛事数据采集系统</h1>
      <p>采集员端</p>
    </div>

    <div class="login-form">
      <div class="form-group">
        <label for="studentId">学号</label>
        <input
          id="studentId"
          v-model="studentId"
          type="text"
          placeholder="请输入8位纯数字学号"
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

      <div v-if="loginError" class="error-message">
        {{ loginError }}
      </div>

      <button type="button" class="login-button" @click="handleLogin">登录</button>

      <div class="login-footer">
        <button type="button" class="forgot-password" @click="handleForgotPassword">
          忘记密码？
        </button>
        <p class="auth-hint">仅授权采集员可登录</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: var(--spacing-3xl) var(--spacing-xl);
  background: radial-gradient(circle at top, #e6f1ff 0, var(--background-secondary) 45%, #edf0f5 100%);
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.login-header h1 {
  font-size: 24px;
  color: var(--text-primary);
  margin-bottom: var(--spacing-sm);
  letter-spacing: 1px;
}

.login-header p {
  font-size: 14px;
  color: var(--text-secondary);
}

.login-form {
  width: 100%;
  max-width: 400px;
  background-color: var(--background-primary);
  border-radius: var(--border-radius-lg);
  padding: var(--spacing-2xl);
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
  animation: fadeInCard 0.3s var(--transition-bezier);
}

.form-group {
  margin-bottom: var(--spacing-xl);
}

.form-group label {
  display: block;
  margin-bottom: var(--spacing-xs);
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

.form-input {
  width: 100%;
  height: 44px;
  padding: 0 var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-size: 15px;
  box-sizing: border-box;
  background-color: var(--background-primary);
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(42, 122, 226, 0.18);
}

.password-input-container {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: var(--primary-color);
  font-size: 14px;
  padding: 4px 6px;
  cursor: pointer;
}

.error-message {
  color: var(--danger-color);
  font-size: 14px;
  margin-bottom: var(--spacing-md);
}

.login-button {
  width: 100%;
  height: 44px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  border: none;
  border-radius: var(--border-radius-md);
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-normal);
  box-shadow: var(--shadow-md);
}

.login-button:hover {
  background: linear-gradient(135deg, var(--primary-hover) 0%, var(--primary-active) 100%);
  box-shadow: var(--shadow-lg);
  transform: translateY(-1px);
}

.login-button:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

.login-footer {
  margin-top: var(--spacing-xl);
  text-align: center;
}

.forgot-password {
  background: none;
  border: none;
  color: var(--primary-color);
  font-size: 14px;
  cursor: pointer;
  margin-bottom: var(--spacing-sm);
}

.auth-hint {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0;
}

@keyframes fadeInCard {
  from {
    opacity: 0;
    transform: translateY(6px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .login-container {
    justify-content: flex-start;
    padding: var(--spacing-2xl) var(--spacing-lg);
    padding-top: var(--spacing-4xl);
  }

  .login-form {
    padding: var(--spacing-xl);
    margin-top: var(--spacing-lg);
  }

  .login-header h1 {
    font-size: 20px;
  }
}
</style>
