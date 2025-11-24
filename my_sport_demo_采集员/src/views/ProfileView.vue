<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MatchHeader from '../components/MatchHeader.vue'

const router = useRouter()

// 赛事信息（仅用于头部显示）
const eventInfo = ref({
  id: '',
  name: '个人中心',
  teamA: '',
  teamB: '',
  status: '',
  collectorName: '采集员姓名'
})

// 用户信息
const userInfo = reactive({
  name: '张三',
  studentId: '20201234', // 8位纯数字学号
  avatar: '', // 可以替换为实际头像URL
  totalCollected: 42,
  completedTasks: 36,
  pendingTasks: 6,
  joinDate: '2023-09-01',
  department: '体育部'
})

// 密码修改表单
const passwordForm = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// UI状态
const uiState = reactive({
  showPasswordForm: false,
  isLoading: false,
  isSubmitting: false,
  error: '',
  success: '',
  showCurrentPassword: false,
  showNewPassword: false,
  showConfirmPassword: false
})

// 加载状态
const isLoading = ref(true)

// 数据统计
const stats = computed(() => [
  { label: '已采集场次', value: userInfo.totalCollected, icon: '📊' },
  { label: '已完成任务', value: userInfo.completedTasks, icon: '✅' },
  { label: '待处理任务', value: userInfo.pendingTasks, icon: '⏳' }
])

// 菜单项
const menuItems = computed(() => [
  {
    id: 'account',
    title: '账号信息',
    icon: '👤',
    items: [
      { label: '学号', value: userInfo.studentId },
      { label: '所属部门', value: userInfo.department },
      { label: '注册日期', value: userInfo.joinDate }
    ]
  },
  {
    id: 'settings',
    title: '设置',
    icon: '⚙️',
    actions: [
      { label: '修改密码', action: () => uiState.showPasswordForm = true },
      { label: '通知设置', action: () => showToast('通知设置功能开发中') },
      { label: '隐私设置', action: () => showToast('隐私设置功能开发中') }
    ]
  },
  {
    id: 'help',
    title: '帮助与支持',
    icon: '❓',
    actions: [
      { label: '使用指南', action: () => showToast('使用指南功能开发中') },
      { label: '常见问题', action: () => showToast('常见问题功能开发中') },
      { label: '联系客服', action: () => showToast('联系客服功能开发中') }
    ]
  }
])

// 显示提示信息
const showToast = (message, type = 'info') => {
  if (type === 'error') {
    uiState.error = message
    setTimeout(() => { uiState.error = '' }, 3000)
  } else {
    uiState.success = message
    setTimeout(() => { uiState.success = '' }, 3000)
  }
}

// 加载用户信息
const loadUserInfo = () => {
  isLoading.value = true
  
  // 模拟API请求
  setTimeout(() => {
    // 这里可以替换为真实的API请求逻辑
    // 模拟加载成功
    isLoading.value = false
  }, 1500)
}

// 验证密码格式
const validatePassword = (password) => {
  if (password.length < 6) {
    return '密码长度不能少于6位'
  }
  return null
}

// 提交密码修改
const handlePasswordSubmit = () => {
  // 验证表单
  if (!passwordForm.currentPassword) {
    showToast('请输入当前密码', 'error')
    return
  }
  
  const newPasswordError = validatePassword(passwordForm.newPassword)
  if (newPasswordError) {
    showToast(newPasswordError, 'error')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    showToast('两次输入的新密码不一致', 'error')
    return
  }
  
  uiState.isSubmitting = true
  
  // 模拟API请求
  setTimeout(() => {
    uiState.isSubmitting = false
    
    // 模拟修改成功
    showToast('密码修改成功')
    
    // 重置表单
    passwordForm.currentPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    uiState.showPasswordForm = false
  }, 1500)
}

// 取消密码修改
const handlePasswordCancel = () => {
  passwordForm.currentPassword = ''
  passwordForm.newPassword = ''
  passwordForm.confirmPassword = ''
  uiState.showPasswordForm = false
}

// 退出登录
const handleLogout = () => {
  if (confirm('确定要退出登录吗？')) {
    // 清除本地存储的认证信息
    localStorage.removeItem('authToken')
    localStorage.removeItem('userInfo')
    
    // 重定向到登录页面
    router.replace('/login')
  }
}

// 生命周期钩子
onMounted(() => {
  loadUserInfo()
})
</script>

<template>
  <div class="profile-view">
    <!-- 赛事信息头部 -->
    <MatchHeader :event="eventInfo" />
    
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">个人中心</h1>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 主内容区域 -->
    <div v-else class="profile-content">
      <!-- 用户信息卡片 -->
      <div class="user-profile-card">
        <div class="avatar-container">
          <div v-if="!userInfo.avatar" class="avatar-placeholder">
            {{ userInfo.name.substring(0, 1) }}
          </div>
          <img v-else :src="userInfo.avatar" alt="用户头像" class="avatar">
        </div>
        
        <div class="user-info">
          <h2 class="user-name">{{ userInfo.name }}</h2>
          <p class="user-username">{{ userInfo.studentId }}</p>
        </div>
      </div>
      
      <!-- 数据统计 -->
      <div class="stats-container">
        <div 
          v-for="stat in stats" 
          :key="stat.label"
          class="stat-card"
        >
          <div class="stat-icon">{{ stat.icon }}</div>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div>
      
      <!-- 菜单列表 -->
      <div class="menu-container">
        <div 
          v-for="menu in menuItems" 
          :key="menu.id"
          class="menu-section"
        >
          <div class="menu-title">
            <span class="menu-icon">{{ menu.icon }}</span>
            <span class="menu-text">{{ menu.title }}</span>
          </div>
          
          <div class="menu-content">
            <!-- 信息项 -->
            <div v-if="menu.items" class="info-items">
              <div 
                v-for="item in menu.items" 
                :key="item.label"
                class="info-item"
              >
                <span class="info-label">{{ item.label }}</span>
                <span class="info-value">{{ item.value }}</span>
              </div>
            </div>
            
            <!-- 操作项 -->
            <div v-if="menu.actions" class="action-items">
              <div 
                v-for="action in menu.actions" 
                :key="action.label"
                class="action-item"
                @click="action.action"
              >
                <span class="action-label">{{ action.label }}</span>
                <span class="action-arrow">›</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 退出登录按钮 -->
        <button 
          @click="handleLogout"
          class="logout-button"
        >
          退出登录
        </button>
      </div>
    </div>
    
    <!-- 密码修改弹窗 -->
    <div v-if="uiState.showPasswordForm" class="modal-overlay" @click="handlePasswordCancel">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3 class="modal-title">修改密码</h3>
          <button @click="handlePasswordCancel" class="modal-close">×</button>
        </div>
        
        <div class="modal-body">
          <div class="form-group">
            <label class="form-label">当前密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.currentPassword"
                :type="uiState.showCurrentPassword ? 'text' : 'password'"
                placeholder="请输入当前密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showCurrentPassword = !uiState.showCurrentPassword"
                class="toggle-password"
              >
                {{ uiState.showCurrentPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
          </div>
          
          <div class="form-group">
            <label class="form-label">新密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.newPassword"
                :type="uiState.showNewPassword ? 'text' : 'password'"
                placeholder="请输入新密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showNewPassword = !uiState.showNewPassword"
                class="toggle-password"
              >
                {{ uiState.showNewPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
            <p class="form-hint">密码长度不少于6位</p>
          </div>
          
          <div class="form-group">
            <label class="form-label">确认新密码</label>
            <div class="password-input-container">
              <input 
                v-model="passwordForm.confirmPassword"
                :type="uiState.showConfirmPassword ? 'text' : 'password'"
                placeholder="请再次输入新密码"
                class="form-input"
              />
              <button 
                type="button"
                @click="uiState.showConfirmPassword = !uiState.showConfirmPassword"
                class="toggle-password"
              >
                {{ uiState.showConfirmPassword ? '👁️' : '👁️‍🗨️' }}
              </button>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button 
            @click="handlePasswordCancel"
            class="btn-secondary"
            :disabled="uiState.isSubmitting"
          >
            取消
          </button>
          <button 
            @click="handlePasswordSubmit"
            class="btn-primary"
            :disabled="uiState.isSubmitting"
          >
            <span v-if="!uiState.isSubmitting">确定修改</span>
            <span v-else>修改中...</span>
          </button>
        </div>
      </div>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="uiState.error" class="toast error-toast">
      <span class="toast-icon">❌</span>
      <span class="toast-message">{{ uiState.error }}</span>
    </div>
    
    <!-- 成功提示 -->
    <div v-if="uiState.success" class="toast success-toast">
      <span class="toast-icon">✅</span>
      <span class="toast-message">{{ uiState.success }}</span>
    </div>
  </div>
</template>

<style scoped>
.profile-view {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  padding-bottom: 40px;
}

.page-header {
  padding: 20px;
  background-color: white;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

/* 加载状态 */
.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 16px;
  color: #666;
  padding: 40px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #1890ff;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 主内容区域 */
.profile-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 0 20px;
}

/* 用户信息卡片 */
.user-profile-card {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.avatar-container {
  position: relative;
}

.avatar-placeholder {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 32px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.user-info {
  flex: 1;
}

.user-name {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0 0 8px 0;
}

.user-username {
  font-size: 14px;
  color: #999;
  margin: 0;
}

/* 数据统计 */
.stats-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  font-size: 32px;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #999;
}

/* 菜单区域 */
.menu-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.menu-section {
  background-color: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.menu-title {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  background-color: #fafafa;
  border-bottom: 1px solid #f0f0f0;
}

.menu-icon {
  font-size: 20px;
}

.menu-text {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.menu-content {
  padding: 8px 0;
}

.info-items,
.action-items {
  display: flex;
  flex-direction: column;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid #f5f5f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: #666;
}

.info-value {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.action-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.action-item:last-child {
  border-bottom: none;
}

.action-item:hover {
  background-color: #fafafa;
}

.action-label {
  font-size: 14px;
  color: #333;
}

.action-arrow {
  font-size: 18px;
  color: #999;
}

/* 退出登录按钮 */
.logout-button {
  padding: 16px;
  background-color: white;
  border: 1px solid #ff4d4f;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  color: #ff4d4f;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.logout-button:hover {
  background-color: #ff4d4f;
  color: white;
  box-shadow: 0 4px 16px rgba(255, 77, 79, 0.3);
}

/* 密码修改弹窗 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  width: 100%;
  max-width: 480px;
  max-height: 80vh;
  overflow-y: auto;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #f0f0f0;
}

.modal-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border: none;
  background: none;
  font-size: 24px;
  color: #999;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.3s;
}

.modal-close:hover {
  background-color: #f5f5f5;
  color: #333;
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 14px;
  font-weight: 500;
  color: #333;
  margin-bottom: 8px;
}

.password-input-container {
  position: relative;
}

.form-input {
  width: 100%;
  padding: 12px 40px 12px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #1890ff;
}

.toggle-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  border: none;
  background: none;
  font-size: 18px;
  cursor: pointer;
}

.form-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  margin-bottom: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid #f0f0f0;
}

.btn-secondary,
.btn-primary {
  padding: 10px 20px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-secondary {
  border: 1px solid #ddd;
  background-color: white;
  color: #333;
}

.btn-secondary:hover:not(:disabled) {
  border-color: #1890ff;
  color: #1890ff;
}

.btn-primary {
  border: 1px solid #1890ff;
  background-color: #1890ff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #40a9ff;
  border-color: #40a9ff;
}

.btn-secondary:disabled,
.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 提示信息 */
.toast {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 1001;
  animation: slideDown 0.3s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translate(-50%, -20px);
  }
  to {
    opacity: 1;
    transform: translate(-50%, 0);
  }
}

.error-toast {
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
  color: #f5222d;
}

.success-toast {
  background-color: #f6ffed;
  border: 1px solid #b7eb8f;
  color: #52c41a;
}

.toast-icon {
  font-size: 18px;
}

.toast-message {
  font-size: 14px;
  font-weight: 500;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    padding: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .profile-content {
    padding: 0 16px;
    gap: 16px;
  }
  
  .user-profile-card {
    flex-direction: column;
    text-align: center;
    padding: 20px;
  }
  
  .avatar-placeholder,
  .avatar {
    width: 64px;
    height: 64px;
  }
  
  .avatar-placeholder {
    font-size: 24px;
  }
  
  .user-name {
    font-size: 20px;
  }
  
  .stats-container {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .stat-card {
    flex-direction: column;
    text-align: center;
    padding: 16px;
    gap: 12px;
  }
  
  .stat-icon {
    font-size: 24px;
  }
  
  .stat-value {
    font-size: 20px;
  }
  
  .menu-section {
    border-radius: 8px;
  }
  
  .menu-title {
    padding: 14px 16px;
  }
  
  .info-item {
    padding: 12px 16px;
  }
  
  .action-item {
    padding: 14px 16px;
  }
  
  .logout-button {
    border-radius: 8px;
  }
  
  .modal-overlay {
    padding: 16px;
  }
  
  .modal-content {
    border-radius: 8px;
  }
  
  .modal-header,
  .modal-body,
  .modal-footer {
    padding: 16px;
  }
  
  .modal-title {
    font-size: 16px;
  }
  
  .modal-footer {
    flex-direction: column;
  }
  
  .btn-secondary,
  .btn-primary {
    width: 100%;
  }
  
  .toast {
    left: 16px;
    right: 16px;
    transform: none;
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
}
</style>