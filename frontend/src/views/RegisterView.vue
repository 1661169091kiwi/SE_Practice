<template>
  <div class="register-container">
    <div class="bg-shape shape-1"></div>
    <div class="bg-shape shape-2"></div>

    <div class="register-content fade-in">
      <div class="register-header">
        <h1>创建账号</h1>
        <p>注册并加入校园赛事平台</p>
      </div>

      <div class="register-form">
        <div class="form-group">
          <label for="studentId">学号</label>
          <input id="studentId" v-model="studentId" type="text" placeholder="请输入学号" class="form-input" />
        </div>

        <div class="form-group">
          <label for="name">姓名</label>
          <input id="name" v-model="name" type="text" placeholder="请输入姓名" class="form-input" />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input id="password" v-model="password" :type="showPassword ? 'text' : 'password'" placeholder="请设置密码" class="form-input" />
        </div>

        <div class="form-group two-cols">
          <div class="col">
            <label for="college">学院</label>
            <input id="college" v-model="college" type="text" placeholder="如：计算机学院" class="form-input" />
          </div>
          <div class="col">
            <label for="grade">年级</label>
            <input id="grade" v-model="grade" type="text" placeholder="如：2023级" class="form-input" />
          </div>
        </div>

        <div v-if="error" class="error-message">{{ error }}</div>
        <div v-if="success" class="success-message">{{ success }}</div>

        <button class="register-button" @click="handleRegister">
          <span>注册</span>
        </button>

        <div class="register-footer">
          <button class="to-login" @click="router.push('/login')">已有账号？去登录</button>
        </div>
      </div>
    </div>
  </div>
 </template>

 <script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '@/utils/http'

const router = useRouter()
const studentId = ref('')
const name = ref('')
const password = ref('')
const college = ref('')
const grade = ref('')
const showPassword = ref(false)
const error = ref('')
const success = ref('')

const validate = () => {
  if (!studentId.value || !name.value || !password.value) {
    error.value = '请填写学号、姓名和密码'
    return false
  }
  error.value = ''
  return true
}

const handleRegister = async () => {
  if (!validate()) return
  try {
    const res = await post('/register', {
      student_id: studentId.value,
      password: password.value,
      name: name.value,
      college: college.value,
      grade: grade.value
    }, { auth: false })
    if (res.code === 200) {
      success.value = '注册成功，正在跳转登录...'
      setTimeout(() => router.push('/login'), 800)
    } else {
      error.value = res.message || '注册失败'
    }
  } catch {
    error.value = '注册失败，请稍后重试'
  }
}
 </script>

 <style scoped>
.register-container {
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: var(--spacing-3xl) var(--spacing-xl);
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  position: relative;
  overflow: hidden;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  z-index: 0;
  opacity: 0.6;
}
.shape-1 { width: 280px; height: 280px; background: var(--primary-gradient-start); top: -50px; left: -60px; }
.shape-2 { width: 380px; height: 380px; background: var(--accent-gradient-end); bottom: -100px; right: -120px; }

.register-content {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
}

.register-header { text-align: center; margin-bottom: var(--spacing-2xl); }
.register-header h1 {
  font-size: 26px;
  background: linear-gradient(to right, var(--primary-dark), var(--primary-active));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
.register-header p { color: var(--text-secondary); }

.register-form {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: var(--border-radius-xl);
  padding: 28px 24px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.4);
}

.form-group { margin-bottom: 16px; }
.form-group.two-cols { display: flex; gap: 12px; }
.form-group.two-cols .col { flex: 1; }
.form-group label { display: block; margin-bottom: 8px; color: var(--text-secondary); font-weight: 600; font-size: 14px; }
.form-input {
  width: 100%; height: 46px; padding: 0 14px;
  border: 2px solid transparent; border-radius: var(--border-radius-lg);
  background: rgba(255,255,255,0.9); font-size: 15px; transition: all 0.3s ease;
  box-shadow: 0 2px 6px rgba(0,0,0,0.02);
}
.form-input:focus { border-color: var(--primary-color); background: #fff; box-shadow: 0 4px 12px rgba(79,172,254,0.18); }

.error-message {
  background: rgba(255, 77, 79, 0.1); color: #ff4d4f;
  padding: 10px; border-radius: 8px; font-size: 13px;
  margin-bottom: 12px; text-align: center; border: 1px solid rgba(255, 77, 79, 0.2);
}
.success-message {
  background: rgba(82, 196, 26, 0.1); color: #52c41a;
  padding: 10px; border-radius: 8px; font-size: 13px;
  margin-bottom: 12px; text-align: center; border: 1px solid rgba(82,196,26,0.2);
}

.register-button {
  width: 100%; height: 48px; border-radius: var(--border-radius-lg);
  background: linear-gradient(135deg, var(--primary-gradient-start) 0%, var(--primary-gradient-end) 100%);
  color: #fff; font-size: 16px; font-weight: 600;
  box-shadow: 0 4px 15px rgba(79, 172, 254, 0.4); transition: all 0.3s ease;
}
.register-button:hover { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(79, 172, 254, 0.6); }

.register-footer { margin-top: 18px; text-align: center; }
.to-login { color: var(--text-tertiary); font-size: 14px; }
.to-login:hover { color: var(--primary-color); text-decoration: underline; }

@media (max-width: 768px) {
  .register-container { align-items: flex-start; padding-top: 12vh; }
  .register-content { padding: 0 16px; }
}
 </style>
