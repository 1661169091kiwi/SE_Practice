<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { get, post } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

// 视图模式：all (全部赛事) | subscribed (我关注的)
const activeTab = ref('all')

// 筛选条件
const selectedSportType = ref('all')
const selectedStatus = ref('all')
const isLoading = ref(false)
const chipContainerRef = ref(null)
const finishedQueue = ref([])
const showFinished = ref(false)
const selectedMonthStr = ref('')
const showMonthPicker = ref(false)
const pickerYear = ref(new Date().getFullYear())

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

// 赛事数据
const events = ref([])
const subscribedEvents = ref([])

// 过滤后的赛事列表
const filteredEvents = ref([])

// 运动ID映射
const sportIdMap = {
  1: 'football',
  2: 'basketball',
  3: 'badminton',
  4: 'volleyball',
  5: 'water'
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

// 通用数据映射函数
const mapEventData = (match) => ({
  id: match.match_id,
  name: match.match_name,
  teamA: match.team_a_name || '主队',
  teamB: match.team_b_name || '客队',
  scoreA: match.score_team_a,
  scoreB: match.score_team_b,
  time: formatTime(match.match_time),
  timestamp: new Date(match.match_time).getTime(),
  venue: '校体育馆', 
  sportType: sportIdMap[match.sport_id] || 'football',
  status: match.status
})

// 获取赛事列表
const fetchEvents = async () => {
  isLoading.value = true
  try {
    const res = await get('/matches')
    if (res.code === 200 && res.data) {
      events.value = res.data.map(mapEventData)
    }
  } catch (error) {
    console.error('Error fetching events:', error)
  } finally {
    isLoading.value = false
    filterEvents()
  }
}

// 获取已关注赛事
const fetchSubscribedEvents = async () => {
  try {
    const studentId = authStore.studentId
    if (!studentId) {
      subscribedEvents.value = []
      return
    }
    const res = await get('/user/subscribed-matches', { params: { userId: studentId } })
    if (res.code === 200 && res.data && res.data.subscribedMatches) {
      // Backend returns SubscribedMatchResponse which contains subscribedMatches array
      // mapEventData might need adjustment because the structure is different
      // SubscribedMatchItem: { matchId, matchTime, matchVenue, homeTeam, awayTeam, matchStatus, matchState }
      
      // Let's create a specific mapper for subscribed events or adapt mapEventData
      // The current mapEventData expects: { match_id, match_name, team_a_name, team_b_name, match_time, sport_id, status }
      
      // SubscribedMatchItem has different keys.
      // We should map it to the same structure as events.value for consistency.
      
      subscribedEvents.value = res.data.subscribedMatches.map(item => ({
        id: item.matchId, // item.matchId is string
        name: `${item.homeTeam.name} VS ${item.awayTeam.name}`, // Construct name if not provided
        teamA: item.homeTeam.name,
        teamB: item.awayTeam.name,
        scoreA: item.scoreA,
        scoreB: item.scoreB,
        time: item.matchTime, // Already formatted string? Backend says string.
        timestamp: new Date(item.matchTime).getTime(),
         venue: item.matchVenue,
         sportType: sportIdMap[item.sportId] || 'football',
         status: item.matchState === '未开始' ? 'not_started' : (item.matchState === '进行中' ? 'in_progress' : 'finished')
       }))
    } else {
      subscribedEvents.value = []
    }
  } catch (error) {
    console.error('Error fetching subscribed events:', error)
    subscribedEvents.value = []
  }
}

// 切换 Tab
const switchTab = (tab) => {
  activeTab.value = tab
  filterEvents()
  if (tab === 'subscribed') {
    fetchSubscribedEvents().then(filterEvents)
  }
}

// 刷新赛事列表
const refreshEvents = () => {
  if (activeTab.value === 'all') {
    fetchEvents()
  } else {
    fetchSubscribedEvents().then(filterEvents)
  }
}

// 过滤赛事
const filterEvents = () => {
  // 根据当前 Tab 决定数据源
  let sourceEvents = activeTab.value === 'all' ? events.value : subscribedEvents.value

  const filtered = sourceEvents.filter(event => {
    const sportMatch = selectedSportType.value === 'all' || event.sportType === selectedSportType.value
    const statusMatch = selectedStatus.value === 'all' || event.status === selectedStatus.value
    return sportMatch && statusMatch
  })

  // 归一化状态
  const normalize = (s) => (s === 'ongoing' ? 'in_progress' : s)
  const now = Date.now()

  // 拆分三类状态
  const inProgressList = filtered.filter(e => normalize(e.status) === 'in_progress')
  const finishedList = filtered.filter(e => normalize(e.status) === 'finished')
  // 其余归为未开始
  const notStartedList = filtered.filter(e => normalize(e.status) !== 'in_progress' && normalize(e.status) !== 'finished')

  // 排序规则：
  // 进行中：按时间倒序（最新的在最前，或者按重要性）- 这里暂按开始时间倒序
  inProgressList.sort((a, b) => (b.timestamp || 0) - (a.timestamp || 0))

  // 已结束：按时间倒序（最近结束的在前）
  finishedList.sort((a, b) => (b.timestamp || 0) - (a.timestamp || 0))

  // 未开始：按时间正序（最近将开始的在前）
  notStartedList.sort((a, b) => (a.timestamp || 0) - (b.timestamp || 0))

  finishedQueue.value = finishedList

  // 组合列表：进行中 -> (已结束) -> 未开始
  if (showFinished.value) {
    // 筛选已结束的比赛
    let filteredFinished = finishedList
    if (!selectedMonthStr.value) {
      // 默认最近一个月 (30天)
      const thirtyDaysAgo = now - 30 * 24 * 60 * 60 * 1000
      filteredFinished = finishedList.filter(e => (e.timestamp || 0) >= thirtyDaysAgo)
    } else {
      // 指定月份筛选
      const [year, month] = selectedMonthStr.value.split('-').map(Number)
      filteredFinished = finishedList.filter(e => {
        const d = new Date(e.timestamp || 0)
        return d.getFullYear() === year && (d.getMonth() + 1) === month
      })
    }

    if (filteredFinished.length === 0) {
      filteredEvents.value = [...inProgressList, ...notStartedList]
    } else {
      let allMatches = [...filteredFinished, ...inProgressList, ...notStartedList]
      allMatches.sort((a, b) => (a.timestamp || 0) - (b.timestamp || 0))
      filteredEvents.value = allMatches
    }
  } else {
    filteredEvents.value = [...inProgressList, ...notStartedList]
  }
}

const toggleFinished = () => {
  showFinished.value = !showFinished.value
  filterEvents()
}

const onMonthChange = () => {
  if (showFinished.value) {
    filterEvents()
  }
}

const openMonthPicker = () => {
  if (selectedMonthStr.value) {
    const [y, _] = selectedMonthStr.value.split('-')
    pickerYear.value = parseInt(y)
  } else {
    pickerYear.value = new Date().getFullYear()
  }
  showMonthPicker.value = true
}

const closeMonthPicker = () => {
  showMonthPicker.value = false
}

const selectMonth = (month) => {
  // month is 1-12
  const m = month.toString().padStart(2, '0')
  selectedMonthStr.value = `${pickerYear.value}-${m}`
  showMonthPicker.value = false
  filterEvents()
}

const changePickerYear = (delta) => {
  pickerYear.value += delta
}

const resetDateFilter = () => {
  selectedMonthStr.value = ''
  showMonthPicker.value = false
  if (showFinished.value) {
    filterEvents()
  }
}

// 查看详情
const goToMatchDetail = (event) => {
  router.push(`/match/${event.id}`)
}

// 查看积分榜
const viewStandings = (event) => {
  router.push(`/student/standings/${event.id}`)
}

// 订阅/关注
const subscribeEvent = async (event) => {
  try {
    const studentId = authStore.studentId
    if (!studentId) {
      alert('无法获取用户信息，请重新登录')
      router.push('/login')
      return
    }

    // 调用后端接口
    const res = await post('/user/subscribe', { 
      studentId: studentId,
      matchId: String(event.id),
      operateType: 0
    })
    
    if (res.code === 200) {
      alert('关注成功！')
      // 如果当前在关注列表，刷新一下
      fetchSubscribedEvents()
    } else {
      alert(res.msg || '关注失败')
    }
  } catch (error) {
     console.error('Subscribe error:', error)
     alert('关注失败，请稍后重试')
  }
}

// 取消关注
const unsubscribeEvent = async (event) => {
  try {
    const studentId = authStore.studentId
    if (!studentId) {
      alert('无法获取用户信息，请重新登录')
      router.push('/login')
      return
    }

    const res = await post('/user/subscribe', { 
      studentId: studentId,
      matchId: String(event.id),
      operateType: 1
    })
    
    if (res.code === 200) {
      alert('已取消关注')
      fetchSubscribedEvents()
    } else {
      alert(res.msg || '取消关注失败')
    }
  } catch (error) {
    console.error('Unsubscribe error:', error)
    alert('取消关注失败，请稍后重试')
  }
}

// 检查是否已关注
const isSubscribed = (eventId) => {
  return subscribedEvents.value.some(e => String(e.id) === String(eventId))
}

// 切换关注状态
const toggleSubscribe = async (event) => {
  if (isSubscribed(event.id)) {
    await unsubscribeEvent(event)
  } else {
    await subscribeEvent(event)
  }
}

// 选中运动类型
const selectSportType = (type) => {
  selectedSportType.value = type
  filterEvents()
}

// 切换到采集员视图
const switchToCollectorMode = () => {
  router.push('/events')
}

// 检查是否具有采集员或管理员权限
const canSwitchToCollector = computed(() => {
  return ['collector', 'admin'].includes(authStore.role)
})

// 登出
const logout = () => {
  authStore.clearAuth()
  router.push('/login')
}

// 跳转到个人中心
const goToProfile = () => {
  router.push('/profile')
}

// 生命周期钩子
onMounted(() => {
  fetchEvents()
  fetchSubscribedEvents() // 预加载订阅信息

  // 让鼠标滚轮可横向滚动筛选条
  if (chipContainerRef.value) {
    const el = chipContainerRef.value
    const wheelHandler = (e) => {
      if (Math.abs(e.deltaY) > Math.abs(e.deltaX)) {
        el.scrollLeft += e.deltaY
        e.preventDefault()
      }
    }
    el.addEventListener('wheel', wheelHandler, { passive: false })
  }
})
</script>

<template>
  <div class="event-selection-container">
    <!-- 页面头部 -->
    <header class="page-header glass-effect">
      <div class="header-top">
        <div class="header-left">
          <h1>赛事列表</h1>
        </div>
        <div class="header-right">
          <button 
            v-if="canSwitchToCollector" 
            class="action-text-btn collector-btn" 
            @click="switchToCollectorMode"
            title="切换至采集版"
          >
            采集员
          </button>
          <button class="action-text-btn" @click="refreshEvents" :disabled="isLoading">
            <span :class="{ 'spin': isLoading }">刷新</span>
          </button>
          <button class="action-text-btn" @click="goToProfile">
            个人中心
          </button>
        </div>
      </div>
      
      <!-- Tab 切换 -->
      <div class="tabs-container">
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'all' }"
          @click="switchTab('all')"
        >
          全部赛事
        </button>
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'subscribed' }"
          @click="switchTab('subscribed')"
        >
          我关注的
        </button>
      </div>
    </header>

    <!-- 筛选栏 (仅在全部赛事Tab显示) -->
    <div v-if="activeTab === 'all'" class="filter-bar glass-effect fade-in">
      <div class="filter-group">
        <div class="filter-chip-container" ref="chipContainerRef">
          <button 
            v-for="type in sportTypes" 
            :key="type.value"
            class="filter-chip"
            :class="{ active: selectedSportType === type.value }"
            @click="selectSportType(type.value)"
          >
            {{ type.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- 赛事列表 -->
    <div class="event-list">
      <div class="finished-cta">
        <button class="load-finished-btn" @click="toggleFinished">
          {{ showFinished ? '隐藏已结束比赛' : `显示已结束比赛 · ${finishedQueue.length} 场` }}
        </button>
        
        <div v-if="showFinished" class="month-filter-container">
          <div class="month-filter fade-in" @click="openMonthPicker">
             <div class="date-trigger">
               <span class="filter-label">{{ selectedMonthStr || '最近30天' }}</span>
               <span class="calendar-icon">📅</span>
             </div>
             
             <button 
               v-if="selectedMonthStr" 
               class="reset-btn" 
               @click.stop="resetDateFilter"
               title="重置为最近30天"
             >
               ↺
             </button>
          </div>

          <div v-if="showMonthPicker" class="month-picker-modal fade-in" @click.stop>
            <div class="picker-header">
              <button class="year-nav-btn" @click="changePickerYear(-1)">◀</button>
              <span class="picker-year">{{ pickerYear }}年</span>
              <button class="year-nav-btn" @click="changePickerYear(1)">▶</button>
            </div>
            <div class="picker-grid">
              <button 
                v-for="m in 12" 
                :key="m" 
                class="month-cell"
                :class="{ active: selectedMonthStr === `${pickerYear}-${m.toString().padStart(2, '0')}` }"
                @click="selectMonth(m)"
              >
                {{ m }}月
              </button>
            </div>
            <div class="picker-footer">
              <button class="picker-close-btn" @click="closeMonthPicker">关闭</button>
            </div>
          </div>
          <div v-if="showMonthPicker" class="picker-overlay" @click="closeMonthPicker"></div>
        </div>
      </div>
      <div v-if="isLoading" class="loading-state fade-in">
        <div class="loading-spinner"></div>
        <p>加载精彩赛事中...</p>
      </div>
      
      <div v-else-if="filteredEvents.length === 0" class="empty-state fade-in">
        <div class="empty-icon">📭</div>
        <p>{{ activeTab === 'subscribed' ? '还没有关注任何比赛哦' : '暂无符合条件的赛事' }}</p>
        <button v-if="activeTab === 'subscribed'" class="empty-action-btn" @click="switchTab('all')">
          去看看全部赛事
        </button>
      </div>
      
      <div v-else>
        <div class="events-grid fade-in">
          <div 
            v-for="event in filteredEvents" 
            :key="event.id" 
            class="event-card glass-card"
            @click="goToMatchDetail(event)"
          >
            <div class="card-status-bar" :class="event.status"></div>
            
            <div class="event-header">
              <div class="sport-badge" :class="event.sportType">
                {{ sportTypes.find(t => t.value === event.sportType)?.label || '赛事' }}
              </div>
              <div class="status-badge" :class="event.status">
                {{ event.status === 'finished' ? '已结束' : (event.status === 'not_started' ? '未开始' : '进行中') }}
              </div>
            </div>
            
            <h3 class="event-name">{{ event.name }}</h3>
            
            <div class="match-up">
              <div class="team-block">
                <div class="team-avatar">{{ event.teamA.charAt(0) }}</div>
                <span class="team-name">{{ event.teamA }}</span>
              </div>
              <div v-if="event.status === 'not_started'" class="vs-badge">VS</div>
              <div v-else class="score-badge">
                <span class="score">{{ event.scoreA }}</span>
                <span class="divider">:</span>
                <span class="score">{{ event.scoreB }}</span>
              </div>
              <div class="team-block">
                <div class="team-avatar">{{ event.teamB.charAt(0) }}</div>
                <span class="team-name">{{ event.teamB }}</span>
              </div>
            </div>
            
            <div class="event-meta">
              <div class="meta-item">
                <span class="meta-icon">📅</span>
                <span>{{ event.time }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-icon">📍</span>
                <span>{{ event.venue }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <button 
                class="action-btn secondary"
                :class="{ 'is-subscribed': isSubscribed(event.id) }"
                @click.stop="toggleSubscribe(event)"
              >
                {{ isSubscribed(event.id) ? '已关注' : '关注' }}
              </button>
              <button 
                class="action-btn primary"
                @click.stop="viewStandings(event)"
              >
                积分榜
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.event-selection-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  display: flex;
  flex-direction: column;
  padding-bottom: 40px;
}

.glass-effect {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.3);
}

.page-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  padding: 0;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
}

.header-left h1 {
  font-size: 22px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(to right, var(--primary-color), var(--primary-active));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.header-right {
  display: flex;
  gap: 12px;
}

.action-icon-btn {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.action-text-btn {
  height: 40px;
  padding: 0 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
}

.action-icon-btn:hover,
.action-text-btn:hover {
  background: white;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.collector-btn {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.collector-btn:hover {
  background: linear-gradient(135deg, #00f2fe 0%, #4facfe 100%);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(79, 172, 254, 0.4);
}

.action-icon-btn:active,
.action-text-btn:active {
  transform: scale(0.95);
}

.icon {
  font-size: 18px;
}

.spin {
  animation: spin 1s linear infinite;
}

/* Tabs */
.tabs-container {
  display: flex;
  padding: 0 20px;
  gap: 24px;
}

.tab-btn {
  padding: 12px 4px;
  font-size: 16px;
  font-weight: 500;
  color: var(--text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  position: relative;
  transition: all 0.3s;
}

.tab-btn.active {
  color: var(--primary-color);
  font-weight: 600;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: var(--primary-color);
  border-radius: 3px 3px 0 0;
}

/* Filter Bar */
.filter-bar {
  padding: 16px 20px;
  margin-bottom: 8px;
}

.filter-chip-container {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding-bottom: 4px;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
  flex-wrap: nowrap; /* 确保不换行 */
  touch-action: pan-x;
  overscroll-behavior-x: contain;
  width: 100%;
}

.filter-chip-container::-webkit-scrollbar {
  display: none;
}

.filter-chip {
  padding: 8px 16px;
  border-radius: 20px;
  background: white;
  border: 1px solid rgba(0,0,0,0.05);
  color: var(--text-secondary);
  font-size: 14px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  flex-shrink: 0;
  user-select: none;
}

.filter-chip:hover {
  background: #f8f9fa;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
}

.filter-chip.active {
  background: var(--primary-color);
  color: white;
  box-shadow: 0 4px 10px rgba(79, 172, 254, 0.3);
  border-color: transparent;
}

/* Event List */
.event-list {
  padding: 16px 20px;
  flex: 1;
}

.events-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
}

.event-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.glass-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 20px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.4);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
  cursor: pointer;
}

.glass-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 12px 40px 0 rgba(31, 38, 135, 0.1);
}

.card-status-bar {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: #e0e0e0;
}

.card-status-bar.in_progress {
  background: var(--success-color);
}

.card-status-bar.not_started {
  background: var(--warning-active);
}

.event-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.sport-badge {
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.sport-badge.football { background: #e6f7ff; color: #1890ff; }
.sport-badge.basketball { background: #f6ffed; color: #52c41a; }
.sport-badge.badminton { background: #fff7e6; color: #faad14; }
.sport-badge.volleyball { background: #fff1f0; color: #f5222d; }
.sport-badge.water { background: #f9f0ff; color: #722ed1; }

.status-badge {
  font-size: 12px;
  font-weight: 500;
  display: flex;
  align-items: center;
}

.status-badge::before {
  content: '';
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-right: 6px;
  background: #ccc;
}

.status-badge.in_progress::before { background: var(--success-color); }
.status-badge.in_progress { color: var(--success-active); }

.status-badge.not_started::before { background: var(--warning-active); }
.status-badge.not_started { color: var(--warning-active); }

.event-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 20px;
  line-height: 1.3;
}

.match-up {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  background: rgba(255,255,255,0.5);
  padding: 16px;
  border-radius: 16px;
}

.team-block {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40%;
}

.team-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e0eafc 0%, #cfdef3 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 8px;
  box-shadow: 0 4px 10px rgba(0,0,0,0.05);
}

.team-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  text-align: center;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.vs-badge {
  font-size: 14px;
  font-weight: 800;
  color: var(--text-tertiary);
  font-style: italic;
}

.score-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 800;
  color: var(--text-primary);
  gap: 8px;
}

.score {
  min-width: 24px;
  text-align: center;
}

.divider {
  color: var(--text-tertiary);
  margin: 0 4px;
}

.event-meta {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
  color: var(--text-secondary);
  font-size: 13px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.card-actions {
  display: flex;
  gap: 12px;
}

.action-btn {
  flex: 1;
  height: 44px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn.primary {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-active) 100%);
  color: white;
  border: none;
  box-shadow: 0 4px 12px rgba(79, 172, 254, 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(79, 172, 254, 0.4);
}

.action-btn.secondary {
  background: white;
  border: 1px solid rgba(0,0,0,0.08);
  color: var(--text-secondary);
}

.action-btn.secondary:hover {
  background: #f8f9fa;
  border-color: rgba(0,0,0,0.15);
}

.action-btn.is-subscribed {
  background: #fff0f6;
  color: #eb2f96;
  border-color: #ffadd2;
}

/* Empty State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  color: var(--text-tertiary);
  margin-bottom: 24px;
}

.empty-action-btn {
  padding: 10px 24px;
  border-radius: 20px;
  background: white;
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.empty-action-btn:hover {
  background: var(--primary-light);
}

/* Finished CTA */
.finished-cta {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
  flex-wrap: wrap;
  position: relative; /* 为弹窗定位提供上下文 */
  z-index: 200;
}
.load-finished-btn {
  padding: 10px 18px;
  border-radius: 999px;
  background: linear-gradient(135deg, var(--primary-gradient-start) 0%, var(--primary-gradient-end) 100%);
  color: #fff;
  font-weight: 600;
  border: none;
  box-shadow: 0 6px 16px rgba(79, 172, 254, 0.35);
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}
.load-finished-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 22px rgba(79, 172, 254, 0.45);
}

.month-filter-container {
  position: relative;
}

.month-filter {
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  padding: 5px 6px 5px 16px;
  border-radius: 30px;
  box-shadow: 0 8px 24px rgba(31, 38, 135, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.6);
  transition: all 0.3s cubic-bezier(0.25, 0.8, 0.25, 1);
  cursor: pointer;
}

.month-filter:hover {
  background: rgba(255, 255, 255, 0.95);
  box-shadow: 0 12px 32px rgba(31, 38, 135, 0.1);
  transform: translateY(-1px);
}

.date-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  white-space: nowrap;
  letter-spacing: 0.3px;
  min-width: 60px;
  text-align: center;
}

.calendar-icon {
  font-size: 16px;
  opacity: 0.7;
  transition: transform 0.3s ease;
}

.month-filter:hover .calendar-icon {
  transform: scale(1.1);
  opacity: 1;
}

.reset-btn {
  background: rgba(0,0,0,0.04);
  border: none;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-secondary);
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  margin-left: -4px;
}

.reset-btn:hover {
  background: #ffecf2;
  color: #eb2f96;
  transform: rotate(180deg) scale(1.1);
}

/* 自定义月份选择器弹窗 */
.month-picker-modal {
  position: absolute;
  top: calc(100% + 12px);
  right: 0; /* 对齐右侧 */
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 20px;
  padding: 16px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
  border: 1px solid rgba(255, 255, 255, 0.6);
  z-index: 1000;
  width: 280px;
  animation: slideDown 0.2s ease-out;
}

.picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

.picker-year {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.year-nav-btn {
  background: none;
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
  font-size: 14px;
  padding: 4px 8px;
  border-radius: 8px;
  transition: all 0.2s;
}

.year-nav-btn:hover {
  background: rgba(0,0,0,0.05);
  color: var(--primary-color);
}

.picker-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 12px;
}

.month-cell {
  background: none;
  border: 1px solid rgba(0,0,0,0.05);
  border-radius: 12px;
  padding: 8px 0;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.month-cell:hover {
  background: #f0f7ff;
  color: var(--primary-color);
  border-color: var(--primary-light);
}

.month-cell.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
  box-shadow: 0 4px 10px rgba(79, 172, 254, 0.3);
}

.picker-footer {
  text-align: right;
}

.picker-close-btn {
  background: none;
  border: none;
  color: var(--text-tertiary);
  font-size: 12px;
  cursor: pointer;
}
.picker-close-btn:hover {
  color: var(--text-secondary);
}

.picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 999;
  background: transparent;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Loading */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(79, 172, 254, 0.2);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Mobile */
@media (max-width: 768px) {
  .events-grid {
    grid-template-columns: 1fr;
  }
  
  .filter-chip-container {
    padding-bottom: 0;
  }
}
</style>
