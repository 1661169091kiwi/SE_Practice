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
    // 登录成功后跳转到管理员页面
    router.push('/arrange')
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
/* 使用vw单位适配2880×1920屏幕 */
/* 1vw = 28.8px (2880px屏幕宽度的1%) */
.login-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 0.694vw; /* 20px → 20/2880*100 = 0.694vw */
  background-color: #f5f5f5;
  box-sizing: border-box;
}

.login-header {
  text-align: center;
  margin-bottom: 1.389vw; /* 40px → 1.389vw */
}

.login-header h1 {
  font-size: 0.833vw; /* 24px → 0.833vw */
  color: #333;
  margin-bottom: 0.278vw; /* 8px → 0.278vw */
}

.login-header p {
  font-size: 0.556vw; /* 16px → 0.556vw */
  color: #666;
}

.login-form {
  width: 100%;
  max-width: 13.889vw; /* 400px → 13.889vw */
  background-color: #fff;
  border-radius: 0.417vw; /* 12px → 0.417vw */
  padding: 1.042vw; /* 30px → 1.042vw */
  box-shadow: 0 0.069vw 0.347vw rgba(0, 0, 0, 0.1); /* 2px 10px → 0.069vw 0.347vw */
  box-sizing: border-box;
}

.form-group {
  margin-bottom: 0.694vw; /* 20px → 0.694vw */
}

.form-group label {
  display: block;
  margin-bottom: 0.278vw; /* 8px → 0.278vw */
  font-size: 0.486vw; /* 14px → 0.486vw */
  color: #333;
  font-weight: 500;
}

.form-input {
  width: 100%;
  height: 1.528vw; /* 44px → 1.528vw */
  padding: 0 0.417vw; /* 12px → 0.417vw */
  border: 0.035vw solid #ddd; /* 1px → 0.035vw */
  border-radius: 0.278vw; /* 8px → 0.278vw */
  font-size: 0.556vw; /* 16px → 0.556vw */
  box-sizing: border-box;
}

.form-input:focus {
  outline: none;
  border-color: #1890ff;
  box-shadow: 0 0 0 0.069vw rgba(24, 144, 255, 0.2); /* 2px → 0.069vw */
}

.password-input-container {
  position: relative;
}

.password-toggle {
  position: absolute;
  right: 0.417vw; /* 12px → 0.417vw */
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #1890ff;
  font-size: 0.486vw; /* 14px → 0.486vw */
  padding: 0.139vw 0.278vw; /* 4px 8px → 0.139vw 0.278vw */
  cursor: pointer;
}

.error-message {
  color: #f5222d;
  font-size: 0.486vw; /* 14px → 0.486vw */
  margin-bottom: 0.556vw; /* 16px → 0.556vw */
}

.login-button {
  width: 100%;
  height: 1.528vw; /* 44px → 1.528vw */
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 0.278vw; /* 8px → 0.278vw */
  font-size: 0.556vw; /* 16px → 0.556vw */
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-button:hover {
  background-color: #40a9ff;
}

.login-footer {
  margin-top: 0.694vw; /* 20px → 0.694vw */
  text-align: center;
}

.forgot-password {
  background: none;
  border: none;
  color: #1890ff;
  font-size: 0.486vw; /* 14px → 0.486vw */
  cursor: pointer;
  margin-bottom: 0.417vw; /* 12px → 0.417vw */
}

.auth-hint {
  font-size: 0.417vw; /* 12px → 0.417vw */
  color: #999;
  margin: 0;
}

/* 2880×1920 屏幕优化（可选） */
@media screen and (min-width: 2880px) {
  .login-form {
    max-width: 15vw; /* 微调最大宽度 */
    padding: 1.2vw;
  }

  .login-header h1 {
    font-size: 0.9vw;
  }

  .form-input {
    height: 1.6vw;
    font-size: 0.6vw;
  }

  .login-button {
    height: 1.6vw;
    font-size: 0.6vw;
  }
}

/* 移动端适配 - 限制最小尺寸 */
@media (max-width: 768px) {
  .login-container {
    padding: 20px;
  }

  .login-form {
    padding: 20px;
    max-width: 100%;
  }

  .login-header h1 {
    font-size: 20px;
  }

  .login-header p {
    font-size: 16px;
  }

  .form-input {
    height: 44px;
    font-size: 16px;
    padding: 0 12px;
  }

  .password-toggle {
    font-size: 14px;
    right: 12px;
  }

  .error-message {
    font-size: 14px;
  }

  .login-button {
    height: 44px;
    font-size: 16px;
  }

  .forgot-password {
    font-size: 14px;
  }

  .auth-hint {
    font-size: 12px;
  }
}
</style>
