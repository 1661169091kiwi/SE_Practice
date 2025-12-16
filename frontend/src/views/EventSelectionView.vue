<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { get } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isAdmin = computed(() => authStore.role === 'admin')

// 筛选条件
const selectedSportType = ref('all')
const selectedStatus = ref('all')
const isLoading = ref(false)

// 运动类型选项
const sportTypes = [
  { value: 'all', label: '全部' },
  { value: 'football', label: '足球' },
  { value: 'basketball', label: '篮球' },
  { value: 'badminton', label: '羽毛球' },
  { value: 'volleyball', label: '排球' },
  { value: 'water', label: '水上运动' }
]

// 比赛状态选项
const matchStatuses = [
  { value: 'all', label: '全部' },
  { value: 'not_started', label: '未开始' },
  { value: 'in_progress', label: '进行中' },
  { value: 'finished', label: '已结束' }
]

// 赛事数据
const events = ref([])

// 过滤后的赛事列表
const filteredEvents = ref([])
// 已结束的时间筛选（按月 YYYY-MM）
const selectedMonth = ref('')

// 运动ID映射
const sportIdMap = {
  1: 'football',
  2: 'basketball',
  3: 'badminton',
  4: 'volleyball',
  5: 'water'
}

// 获取赛事列表
const fetchEvents = async () => {
  isLoading.value = true
  try {
    const res = await get('/matches', { params: { view: 'collector' } })
    if (res.code === 200 && res.data) {
      // 映射后端数据到前端结构
      events.value = res.data.map(match => ({
        id: match.match_id,
        name: match.match_name,
        teamA: match.team_a_name || '主队',
        teamB: match.team_b_name || '客队',
        time: formatTime(match.match_time),
        timestamp: match.match_time ? new Date(match.match_time).getTime() : 0,
        venue: '校体育馆', // 后端暂未返回场地信息，默认显示
        sportType: sportIdMap[match.sport_id] || 'football', // 默认足球
        // 归一化状态：ongoing -> in_progress
        status: match.status === 'ongoing' ? 'in_progress' : match.status
      }))
    } else {
      console.error('Failed to fetch events:', res.msg)
    }
  } catch (error) {
    console.error('Error fetching events:', error)
  } finally {
    isLoading.value = false
    filterEvents()
  }
}

// 格式化时间
const formatTime = (timeStr) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 刷新赛事列表
const refreshEvents = () => {
  fetchEvents()
}

// 过滤赛事
const filterEvents = () => {
  filteredEvents.value = events.value.filter(event => {
    const sportMatch = selectedSportType.value === 'all' || event.sportType === selectedSportType.value
    let statusMatch = true
    if (selectedStatus.value === 'finished') {
      statusMatch = event.status === 'finished'
    } else if (selectedStatus.value === 'not_started' || selectedStatus.value === 'in_progress') {
      statusMatch = event.status === selectedStatus.value
    } else {
      // 全部：不展示已结束
      statusMatch = event.status !== 'finished'
    }
    // 时间筛选：仅在已结束分类下生效，按月 YYYY-MM
    let timeMatch = true
    if (selectedStatus.value === 'finished' && selectedMonth.value) {
      const d = event.timestamp ? new Date(event.timestamp) : null
      const ym = d ? `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}` : ''
      timeMatch = ym === selectedMonth.value
    }
    return sportMatch && statusMatch && timeMatch
  })
}

// 进入数据采集
const goToDataCollection = (event) => {
  // 根据运动类型跳转到对应的数据录入页面
  router.push(`/data-collection/${event.sportType}/${event.id}`)
}

// 切换到学生版视图
const switchToStudentMode = () => {
  router.push('/student/events')
}

// 跳转到个人中心
const goToProfile = () => {
  router.push('/profile')
}

// 生命周期钩子
onMounted(() => {
  fetchEvents()
})
</script>

<template>
  <div class="event-selection-container">
    <!-- 页面头部 -->
    <header class="page-header glass-effect">
      <div class="header-left">
        <h1>待采集赛事</h1>
      </div>
      <div class="header-right">
        <button v-if="isAdmin" class="action-btn admin-btn" @click="router.push('/admin/dashboard')" title="管理后台">
          管理后台
        </button>
        <button class="action-btn student-btn" @click="switchToStudentMode" title="切换至学生版">
          学生版
        </button>
        <button class="action-btn secondary-btn" @click="refreshEvents" :disabled="isLoading" title="刷新">
          {{ isLoading ? '刷新中...' : '刷新' }}
        </button>
        <button class="action-btn secondary-btn" @click="goToProfile" title="个人中心">
          个人中心
        </button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar glass-effect">
      <div class="filter-group">
        <label class="filter-label">运动类型</label>
        <div class="filter-chips">
          <button 
            v-for="type in sportTypes" 
            :key="type.value" 
            class="filter-chip"
            :class="{ active: selectedSportType === type.value }"
            @click="selectedSportType = type.value; filterEvents()"
          >
            {{ type.label }}
          </button>
        </div>
      </div>
      
      <div class="filter-group">
        <label class="filter-label">状态</label>
        <div class="filter-chips">
          <button 
            v-for="status in matchStatuses" 
            :key="status.value" 
            class="filter-chip"
            :class="{ active: selectedStatus === status.value }"
            @click="selectedStatus = status.value; filterEvents()"
          >
            {{ status.value === 'in_progress' ? '采集中' : status.label }}
          </button>
        </div>
        <div v-if="selectedStatus === 'finished'" class="finished-time-filter">
          <label class="filter-label">按月份筛选</label>
          <input type="month" v-model="selectedMonth" @change="filterEvents" class="month-input" />
          <button class="filter-chip reset" @click="selectedMonth=''; filterEvents()">重置</button>
        </div>
      </div>
    </div>

    <!-- 赛事列表 -->
    <div class="event-list-container">
      <div v-if="isLoading" class="loading-state">
        <div class="spinner"></div>
        <p>加载中...</p>
      </div>
      <div v-else-if="filteredEvents.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <p>暂无符合条件的赛事</p>
      </div>
      <div v-else class="event-grid">
        <div v-for="event in filteredEvents" :key="event.id" class="event-card glass-card">
          <div class="card-header">
            <span class="sport-badge" :class="event.sportType">
              {{ sportTypes.find(t => t.value === event.sportType)?.label || '未知' }}
            </span>
            <span class="status-badge" :class="event.status">
              {{ event.status === 'finished' ? '已结束' : (event.status === 'not_started' ? '未开始' : '采集中') }}
            </span>
          </div>
          
          <h3 class="event-title">{{ event.name }}</h3>
          
          <div class="match-up">
            <div class="team-info">
              <span class="team-name">{{ event.teamA }}</span>
            </div>
            <div class="vs-badge">VS</div>
            <div class="team-info">
              <span class="team-name">{{ event.teamB }}</span>
            </div>
          </div>
          
          <div class="match-details">
            <div class="detail-row">
              <span class="icon">🕒</span>
              <span>{{ event.time }}</span>
            </div>
            <div class="detail-row">
              <span class="icon">📍</span>
              <span>{{ event.venue }}</span>
            </div>
          </div>
          
          <button class="collect-btn" @click="goToDataCollection(event)">
            进入采集
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.event-selection-container {
  min-height: 100vh;
  background: radial-gradient(circle at top left, #eef2f3 0%, #d9e4f5 100%);
  display: flex;
  flex-direction: column;
}

.glass-effect {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.5);
}

/* Header Styles */
.page-header {
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
}

.header-left h1 {
  font-size: 24px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(135deg, var(--primary-color) 0%, #1890ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-right {
  display: flex;
  gap: 12px;
  align-items: center;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 12px;
  border: none;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  height: 40px;
}

.action-btn .icon {
  font-size: 16px;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.action-btn:active {
  transform: scale(0.95);
}

.admin-btn {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(118, 75, 162, 0.3);
}

.student-btn {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.secondary-btn {
  background: white;
  color: var(--text-primary);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.danger-btn {
  background: rgba(255, 77, 79, 0.1);
  color: #ff4d4f;
}

.danger-btn:hover {
  background: #ff4d4f;
  color: white;
}

.spin {
  animation: spin 1s linear infinite;
}

/* Filter Bar */
.filter-bar {
  padding: 16px 24px;
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.filter-chips {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.finished-time-filter {
  display: flex;
  align-items: center;
  gap: 8px;
}
.month-input {
  padding: 6px 10px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  border-radius: 6px;
  font-size: 12px;
  background: rgba(255, 255, 255, 0.8);
}
.filter-chip.reset {
  background: #f5f5f5;
  color: #333;
}

.filter-chip {
  padding: 6px 14px;
  border-radius: 20px;
  border: 1px solid transparent;
  background: rgba(255, 255, 255, 0.5);
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.filter-chip:hover {
  background: white;
  transform: translateY(-1px);
}

.filter-chip.active {
  background: var(--primary-color);
  color: white;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
}

/* Event Grid */
.event-list-container {
  flex: 1;
  padding: 24px;
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}

.event-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.event-card {
  background: rgba(255, 255, 255, 0.8);
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.5);
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.event-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
  background: white;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sport-badge {
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  color: white;
}

.sport-badge.football { background: linear-gradient(135deg, #36d1dc, #5b86e5); }
.sport-badge.basketball { background: linear-gradient(135deg, #f2994a, #f2c94c); }
.sport-badge.badminton { background: linear-gradient(135deg, #FF512F, #DD2476); }
.sport-badge.volleyball { background: linear-gradient(135deg, #11998e, #38ef7d); }
.sport-badge.water { background: linear-gradient(135deg, #00c6ff, #0072ff); }

.status-badge {
  font-size: 12px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.05);
}

.status-badge.in_progress {
  color: #52c41a;
  background: rgba(82, 196, 26, 0.1);
}

.status-badge.not_started {
  color: #faad14;
  background: rgba(250, 173, 20, 0.1);
}

.status-badge.finished {
  color: #8c8c8c;
  background: rgba(140, 140, 140, 0.1);
}

.event-title {
  font-size: 18px;
  font-weight: 700;
  margin: 0;
  color: var(--text-primary);
}

.match-up {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(240, 242, 245, 0.5);
  padding: 16px;
  border-radius: 12px;
}

.team-info {
  flex: 1;
  text-align: center;
}

.team-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 15px;
}

.vs-badge {
  font-size: 12px;
  font-weight: 800;
  color: rgba(0, 0, 0, 0.2);
  margin: 0 12px;
}

.match-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary);
}

.collect-btn {
  width: 100%;
  padding: 12px;
  border-radius: 10px;
  border: none;
  background: var(--primary-color);
  color: white;
  font-weight: 600;
  font-size: 15px;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: auto;
}

.collect-btn:hover {
  background: var(--primary-hover);
  box-shadow: 0 4px 15px rgba(24, 144, 255, 0.3);
}

/* Loading & Empty States */
.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  color: var(--text-tertiary);
  width: 100%;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top-color: var(--primary-color);
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Mobile Responsiveness */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 16px;
  }
  
  .header-right {
    width: 100%;
    justify-content: space-between;
  }
  
  .action-btn {
    flex: 1;
    justify-content: center;
    padding: 8px;
  }
  
  .action-btn .text {
    display: none;
  }
  
  .filter-bar {
    padding: 16px;
    gap: 16px;
    overflow-x: auto;
    flex-wrap: nowrap;
  }
  .finished-time-filter {
    min-width: max-content;
  }
  
  .filter-group {
    min-width: fit-content;
  }
  
  .filter-chips {
    flex-wrap: nowrap;
  }
  
  .event-list-container {
    padding: 16px;
  }
  
  .event-grid {
    grid-template-columns: 1fr;
  }
}
</style>
