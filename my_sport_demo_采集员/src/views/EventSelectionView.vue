<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

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
  { value: 'in_progress', label: '进行中' }
]

// 模拟赛事数据
const events = ref([
  {
    id: 1,
    name: '2024校级足球联赛半决赛',
    teamA: '计算机学院',
    teamB: '机械学院',
    time: '2024-11-19 14:30',
    venue: '主体育场',
    sportType: 'football',
    status: 'not_started'
  },
  {
    id: 2,
    name: '2024校级篮球联赛决赛',
    teamA: '商学院',
    teamB: '电子学院',
    time: '2024-11-19 16:00',
    venue: '篮球馆',
    sportType: 'basketball',
    status: 'in_progress'
  },
  {
    id: 3,
    name: '2024校级羽毛球团体赛',
    teamA: '文学院',
    teamB: '外国语学院',
    time: '2024-11-20 09:00',
    venue: '综合馆',
    sportType: 'badminton',
    status: 'not_started'
  },
  {
    id: 4,
    name: '2024校级排球联赛',
    teamA: '医学院',
    teamB: '化工学院',
    time: '2024-11-20 10:30',
    venue: '排球场',
    sportType: 'volleyball',
    status: 'not_started'
  },
  {
    id: 5,
    name: '2024校级龙舟比赛',
    teamA: '建筑学院',
    teamB: '环境学院',
    time: '2024-11-21 13:00',
    venue: '东湖',
    sportType: 'water',
    status: 'not_started'
  }
])

// 过滤后的赛事列表
const filteredEvents = ref([])

// 刷新赛事列表
const refreshEvents = () => {
  isLoading.value = true
  // 模拟加载延迟
  setTimeout(() => {
    filterEvents()
    isLoading.value = false
  }, 1000)
}

// 过滤赛事
const filterEvents = () => {
  filteredEvents.value = events.value.filter(event => {
    const sportMatch = selectedSportType.value === 'all' || event.sportType === selectedSportType.value
    const statusMatch = selectedStatus.value === 'all' || event.status === selectedStatus.value
    return sportMatch && statusMatch
  })
}

// 进入数据采集
const goToDataCollection = (event) => {
  // 根据运动类型跳转到对应的数据录入页面
  router.push(`/data-collection/${event.sportType}/${event.id}`)
}

// 登出
const logout = () => {
  router.push('/login')
}

// 生命周期钩子
onMounted(() => {
  filterEvents()
})
</script>

<template>
  <div class="event-selection-container">
    <!-- 页面头部 -->
    <header class="page-header">
      <div class="header-left">
        <h1>待采集赛事</h1>
      </div>
      <div class="header-right">
        <button class="refresh-button" @click="refreshEvents" :disabled="isLoading">
          <span v-if="isLoading">刷新中...</span>
          <span v-else>刷新</span>
        </button>
        <button class="logout-button" @click="logout">
          退出
        </button>
      </div>
    </header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <div class="filter-item">
        <label>运动类型：</label>
        <select v-model="selectedSportType" @change="filterEvents" class="filter-select">
          <option v-for="type in sportTypes" :key="type.value" :value="type.value">
            {{ type.label }}
          </option>
        </select>
      </div>
      <div class="filter-item">
        <label>比赛状态：</label>
        <select v-model="selectedStatus" @change="filterEvents" class="filter-select">
          <option v-for="status in matchStatuses" :key="status.value" :value="status.value">
            {{ status.label }}
          </option>
        </select>
      </div>
    </div>

    <!-- 赛事列表 -->
    <div class="event-list">
      <div v-if="isLoading" class="loading-state">
        <p>加载中...</p>
      </div>
      <div v-else-if="filteredEvents.length === 0" class="empty-state">
        <p>暂无符合条件的赛事</p>
      </div>
      <div v-else>
        <div v-for="event in filteredEvents" :key="event.id" class="event-card">
          <div class="event-header">
            <h3 class="event-name">{{ event.name }}</h3>
            <span class="sport-tag" :class="event.sportType">
              {{ sportTypes.find(t => t.value === event.sportType)?.label || '未知' }}
            </span>
          </div>
          <div class="event-info">
            <div class="teams">
              <span class="team">{{ event.teamA }}</span>
              <span class="vs">VS</span>
              <span class="team">{{ event.teamB }}</span>
            </div>
            <div class="event-details">
              <p class="detail-item">
                <span class="label">时间：</span>
                <span class="value">{{ event.time }}</span>
              </p>
              <p class="detail-item">
                <span class="label">场地：</span>
                <span class="value">{{ event.venue }}</span>
              </p>
              <p class="detail-item">
                <span class="label">状态：</span>
                <span class="value" :class="event.status">
                  {{ event.status === 'not_started' ? '未开始' : '进行中' }}
                </span>
              </p>
            </div>
          </div>
          <button 
            class="collect-button"
            @click="goToDataCollection(event)"
          >
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
  background: linear-gradient(135deg, #f8fafc 0%, #e0e7ff 25%, #f3e8ff 50%, #fce7f3 75%, #f8fafc 100%);
  background-size: 400% 400%;
  animation: gradientShift 20s ease infinite;
  display: flex;
  flex-direction: column;
}

@keyframes gradientShift {
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

.page-header {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 50%, #a855f7 100%);
  color: white;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(99, 102, 241, 0.3);
  backdrop-filter: blur(10px);
}

.header-left h1 {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
  letter-spacing: 0.5px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-right {
  display: flex;
  gap: 12px;
}

.refresh-button,
.logout-button {
  background-color: rgba(255, 255, 255, 0.25);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s var(--transition-bezier);
  backdrop-filter: blur(10px);
}

.refresh-button:hover,
.logout-button:hover {
  background-color: rgba(255, 255, 255, 0.35);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.refresh-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.filter-bar {
  background-color: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  padding: 20px 24px;
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  border-bottom: 1px solid rgba(229, 231, 235, 0.8);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-item label {
  font-size: 13px;
  color: #6b7280;
  font-weight: 500;
  letter-spacing: 0.3px;
}

.filter-select {
  padding: 10px 14px;
  border: 2px solid #e5e7eb;
  border-radius: 10px;
  font-size: 14px;
  background-color: white;
  color: #374151;
  font-weight: 500;
  transition: all 0.3s var(--transition-bezier);
  cursor: pointer;
}

.filter-select:hover {
  border-color: #d1d5db;
}

.filter-select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.event-list {
  flex: 1;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #9ca3af;
  font-size: 15px;
}

.event-card {
  background-color: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.04), 0 1px 3px rgba(0, 0, 0, 0.06);
  position: relative;
  border: 1px solid rgba(229, 231, 235, 0.8);
  transition: all 0.3s var(--transition-bezier);
  overflow: hidden;
}

.event-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, #6366f1, #8b5cf6, #a855f7);
  opacity: 0;
  transition: opacity 0.3s;
}

.event-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1), 0 4px 12px rgba(0, 0, 0, 0.08);
  border-color: rgba(99, 102, 241, 0.3);
}

.event-card:hover::before {
  opacity: 1;
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.event-name {
  font-size: 17px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  flex: 1;
  padding-right: 12px;
  letter-spacing: 0.2px;
}

.sport-tag {
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  color: white;
  background-color: #6366f1;
  letter-spacing: 0.3px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.sport-tag.football {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

.sport-tag.basketball {
  background: linear-gradient(135deg, #10b981, #059669);
}

.sport-tag.badminton {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.sport-tag.volleyball {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.sport-tag.water {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
}

.event-info {
  margin-bottom: 20px;
}

.teams {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 16px;
  padding: 16px 0;
  background: linear-gradient(135deg, #f9fafb, #f3f4f6);
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.team {
  font-size: 17px;
  font-weight: 600;
  color: #1f2937;
  letter-spacing: 0.3px;
}

.vs {
  font-size: 14px;
  color: #9ca3af;
  font-weight: 600;
  letter-spacing: 2px;
}

.event-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.detail-item {
  font-size: 14px;
  margin: 0;
  display: flex;
  gap: 6px;
}

.detail-item .label {
  color: #6b7280;
  font-weight: 500;
}

.detail-item .value {
  color: #1f2937;
  font-weight: 600;
}

.detail-item .value.not_started {
  color: #f59e0b;
}

.detail-item .value.in_progress {
  color: #10b981;
}

.collect-button {
  position: absolute;
  bottom: 20px;
  right: 20px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s var(--transition-bezier);
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.35);
  letter-spacing: 0.3px;
}

.collect-button:hover {
  background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  box-shadow: 0 6px 20px rgba(99, 102, 241, 0.45);
  transform: translateY(-2px);
}

.collect-button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    padding: 16px 20px;
  }

  .header-left h1 {
    font-size: 18px;
  }

  .filter-bar {
    flex-direction: column;
    padding: 16px 20px;
  }
  
  .filter-item {
    width: 100%;
  }
  
  .filter-select {
    width: 100%;
  }
  
  .teams {
    flex-direction: column;
    gap: 12px;
    padding: 14px 0;
  }
  
  .event-header {
    flex-direction: column;
    gap: 10px;
  }
  
  .sport-tag {
    align-self: flex-start;
  }

  .event-list {
    padding: 16px 20px;
  }

  .event-card {
    padding: 16px;
  }

  .collect-button {
    position: static;
    width: 100%;
    margin-top: 16px;
  }
}
</style>