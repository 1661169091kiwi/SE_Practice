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
  background: radial-gradient(circle at top, #e6f1ff 0, var(--background-secondary) 45%, #edf0f5 100%);
  display: flex;
  flex-direction: column;
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-md);
}

.header-left h1 {
  font-size: 18px;
  font-weight: 500;
  margin: 0;
}

.header-right {
  display: flex;
  gap: 12px;
}

.refresh-button,
.logout-button {
  background-color: rgba(255, 255, 255, 0.2);
  color: var(--text-white);
  border: none;
  padding: 6px 12px;
  border-radius: var(--border-radius-sm);
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.refresh-button:hover,
.logout-button:hover {
  background-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-1px);
}

.refresh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.filter-bar {
  background-color: var(--background-primary);
  padding: 16px 20px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  border-bottom: 1px solid var(--border-color);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.03);
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.filter-item label {
  font-size: 14px;
  color: var(--text-secondary);
}

.filter-select {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-size: 14px;
  background-color: var(--background-primary);
}

.filter-select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(42, 122, 226, 0.18);
}

.event-list {
  flex: 1;
  padding: 16px 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 40px 0;
  color: var(--text-tertiary);
}

.event-card {
  background-color: var(--background-primary);
  border-radius: var(--border-radius-lg);
  padding: 16px;
  box-shadow: var(--shadow-md);
  position: relative;
  border: 1px solid var(--border-color);
  transition: transform 0.2s var(--transition-bezier), box-shadow 0.2s var(--transition-bezier);
}

.event-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
}

.event-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.event-name {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0;
  flex: 1;
  padding-right: 8px;
}

.sport-tag {
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 12px;
  color: white;
  background-color: var(--primary-color);
}

.sport-tag.football {
  background-color: #1890ff;
}

.sport-tag.basketball {
  background-color: #52c41a;
}

.sport-tag.badminton {
  background-color: #faad14;
}

.sport-tag.volleyball {
  background-color: #f5222d;
}

.sport-tag.water {
  background-color: #722ed1;
}

.event-info {
  margin-bottom: 16px;
}

.teams {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 12px;
  padding: 12px 0;
  background-color: var(--background-tertiary);
  border-radius: var(--border-radius-md);
}

.team {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
}

.vs {
  font-size: 14px;
  color: var(--text-tertiary);
  font-weight: 400;
}

.event-details {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.detail-item {
  font-size: 14px;
  margin: 0;
  display: flex;
  gap: 4px;
}

.detail-item .label {
  color: var(--text-secondary);
}

.detail-item .value {
  color: var(--text-primary);
  font-weight: 500;
}

.detail-item .value.not_started {
  color: #faad14;
}

.detail-item .value.in_progress {
  color: #52c41a;
}

.collect-button {
  position: absolute;
  bottom: 16px;
  right: 16px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  border: none;
  padding: 8px 16px;
  border-radius: var(--border-radius-md);
  font-size: 14px;
  cursor: pointer;
  transition: all var(--transition-normal);
  box-shadow: var(--shadow-sm);
}

.collect-button:hover {
  background: linear-gradient(135deg, var(--primary-hover) 0%, var(--primary-active) 100%);
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

 .collect-button:active {
  transform: translateY(0);
  box-shadow: var(--shadow-sm);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
  }
  
  .filter-item {
    width: 100%;
  }
  
  .filter-select {
    width: 100%;
  }
  
  .teams {
    flex-direction: column;
    gap: 8px;
  }
  
  .event-header {
    flex-direction: column;
    gap: 8px;
  }
  
  .sport-tag {
    align-self: flex-start;
  }
}
</style>