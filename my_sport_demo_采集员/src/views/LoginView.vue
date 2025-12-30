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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 50%, #f093fb 100%);
  background-size: 200% 200%;
  animation: gradientBG 15s ease infinite;
}

@keyframes gradientBG {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
  animation: slideDown 0.6s var(--transition-bezier);
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header h1 {
  font-size: 28px;
  color: white;
  margin-bottom: var(--spacing-sm);
  letter-spacing: 2px;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
  font-weight: 600;
}

.login-header p {
  font-size: 15px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 300;
  letter-spacing: 1px;
}

.login-form {
  width: 100%;
  max-width: 420px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 32px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.3);
  animation: fadeInCard 0.5s var(--transition-bezier);
}

@keyframes fadeInCard {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  color: #374151;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.form-input {
  width: 100%;
  height: 48px;
  padding: 0 16px;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  font-size: 15px;
  box-sizing: border-box;
  background-color: #ffffff;
  transition: all 0.3s var(--transition-bezier);
}

.form-input:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.15);
  transform: translateY(-1px);
}

.form-input:hover {
  border-color: #d1d5db;
}

.password-input-container {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #6366f1;
  font-size: 13px;
  padding: 6px 10px;
  cursor: pointer;
  border-radius: 6px;
  transition: all 0.2s;
  font-weight: 500;
}

.password-toggle:hover {
  background-color: #eef2ff;
}

.error-message {
  color: #ef4444;
  font-size: 14px;
  margin-bottom: 20px;
  padding: 12px 16px;
  background-color: #fef2f2;
  border-radius: 8px;
  border-left: 4px solid #ef4444;
  animation: shake 0.4s var(--transition-bezier);
}

@keyframes shake {
  0%, 100% {
    transform: translateX(0);
  }
  25% {
    transform: translateX(-5px);
  }
  75% {
    transform: translateX(5px);
  }
}

.login-button {
  width: 100%;
  height: 50px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s var(--transition-bezier);
  box-shadow: 0 4px 14px 0 rgba(99, 102, 241, 0.39);
  letter-spacing: 0.5px;
}

.login-button:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  box-shadow: 0 6px 20px 0 rgba(99, 102, 241, 0.5);
  transform: translateY(-2px);
}

.login-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px 0 rgba(99, 102, 241, 0.3);
}

.login-footer {
  margin-top: 24px;
  text-align: center;
}

.forgot-password {
  background: none;
  border: none;
  color: #6366f1;
  font-size: 14px;
  cursor: pointer;
  margin-bottom: 12px;
  font-weight: 500;
  transition: all 0.2s;
  padding: 4px 8px;
  border-radius: 6px;
}

.forgot-password:hover {
  background-color: #eef2ff;
}

.auth-hint {
  font-size: 12px;
  color: #9ca3af;
  margin: 0;
  font-weight: 400;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .login-container {
    justify-content: flex-start;
    padding: 24px 20px;
    padding-top: 80px;
  }

  .login-form {
    padding: 24px;
    margin-top: 20px;
    border-radius: 16px;
  }

  .login-header h1 {
    font-size: 24px;
  }

  .form-input,
  .login-button {
    height: 46px;
  }
}
</style>
