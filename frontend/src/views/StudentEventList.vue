<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick } from 'vue'
import { useRouter, onBeforeRouteLeave } from 'vue-router'
import { get, post } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'

// 视图模式：home (首页) | all (全部赛事) | subscribed (我关注的)
const activeTab = ref('home')

// 轮播图数据 (模拟数据，实际应由管理员设置)
const carouselImages = ref([
  { id: 1, src: 'https://placehold.co/1200x600/1890ff/ffffff?text=Sports+Festival', title: '校园体育节盛大开幕' },
  { id: 2, src: 'https://placehold.co/1200x600/52c41a/ffffff?text=Basketball+Finals', title: '篮球决赛精彩瞬间' },
  { id: 3, src: 'https://placehold.co/1200x600/faad14/ffffff?text=Badminton+Cup', title: '羽毛球新生杯报名中' },
  { id: 4, src: 'https://placehold.co/1200x600/f5222d/ffffff?text=Football+League', title: '足球联赛激战正酣' }
])

const currentCarouselIndex = ref(0)
let carouselInterval = null

const startCarousel = () => {
  stopCarousel()
  carouselInterval = setInterval(() => {
    currentCarouselIndex.value = (currentCarouselIndex.value + 1) % carouselImages.value.length
  }, 5000)
}

const stopCarousel = () => {
  if (carouselInterval) clearInterval(carouselInterval)
}

const setCarouselIndex = (index) => {
  currentCarouselIndex.value = index
  startCarousel() // 重置计时器
}

// 触摸滑动逻辑
const touchStartX = ref(0)
const touchEndX = ref(0)
const minSwipeDistance = 50

const onTouchStart = (e) => {
  touchStartX.value = e.changedTouches[0].screenX
  stopCarousel() // 用户交互时暂停轮播
}

const onTouchEnd = (e) => {
  touchEndX.value = e.changedTouches[0].screenX
  handleSwipe()
  startCarousel() // 交互结束后恢复轮播
}

const handleSwipe = () => {
  const distance = touchEndX.value - touchStartX.value
  
  if (Math.abs(distance) < minSwipeDistance) return

  if (distance > 0) {
    // 向右滑动，显示上一张
    setCarouselIndex((currentCarouselIndex.value - 1 + carouselImages.value.length) % carouselImages.value.length)
  } else {
    // 向左滑动，显示下一张
    setCarouselIndex((currentCarouselIndex.value + 1) % carouselImages.value.length)
  }
}

const handleCarouselClick = (img) => {
  if (img.link) {
    window.location.href = img.link
  }
}

// 首页数据计算属性
const recentFinishedEvents = computed(() => {
  return events.value
    .filter(e => normalizeStatus(e.status) === 'finished')
    .sort((a, b) => (b.timestamp || 0) - (a.timestamp || 0))
    .slice(0, 10)
})

const upcomingEvents = computed(() => {
  return events.value
    .filter(e => normalizeStatus(e.status) === 'not_started')
    .sort((a, b) => (a.timestamp || 0) - (b.timestamp || 0))
    .slice(0, 10)
})

// 筛选条件
const selectedSportType = ref('all')
const selectedStatus = ref('all')
const selectedEventId = ref('all')
const isEventDropdownOpen = ref(false)
const isLoading = ref(false)
const chipContainerRef = ref(null)
const finishedQueue = ref([])
const showFinished = ref(false)
const selectedMonthStr = ref('')
const showMonthPicker = ref(false)
const pickerYear = ref(new Date().getFullYear())

// 运动类型选项
const sportTypes = [
  { value: 'all', label: '全部 🏅' },
  { value: 'football', label: '足球 ⚽' },
  { value: 'basketball', label: '篮球 🏀' },
  { value: 'badminton', label: '羽毛球 🏸' },
  { value: 'volleyball', label: '排球 🏐' }
]

// 赛事数据
const events = ref([])
const subscribedEvents = ref([])
const eventMap = ref({})

// 赛事筛选选项
const eventOptions = computed(() => {
  const map = new Map()
  // Use events.value to only show events that have matches in the current list
  // or use eventMap if we want to show all possible events even if empty
  // Using events.value is better for context
  events.value.forEach(e => {
    if (e.eventId && e.eventName) {
      map.set(e.eventId, e.eventName)
    }
  })
  return Array.from(map.entries()).map(([id, name]) => ({ id, name }))
})

// Custom Dropdown Methods
const toggleEventDropdown = () => {
  isEventDropdownOpen.value = !isEventDropdownOpen.value
}

const selectEvent = (id) => {
  selectedEventId.value = id
  isEventDropdownOpen.value = false
  filterEvents()
}

const closeEventDropdown = (e) => {
  // Check if click target is outside the dropdown
  const target = e.target
  if (!target.closest('.event-filter-wrapper')) {
    isEventDropdownOpen.value = false
  }
}

// 过滤后的赛事列表
const filteredEvents = ref([])

// 运动ID映射
const sportIdMap = {
  1: 'football',
  2: 'basketball',
  3: 'badminton',
  4: 'volleyball'
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

const normalizeStatus = (s) => {
  const v = String(s ?? '').trim().toLowerCase()
  if (v === 'ongoing') return 'in_progress'
  if (v === 'completed') return 'finished'
  return v
}
const getStatusText = (s) => {
  const v = normalizeStatus(s)
  if (v === 'finished') return '已结束'
  if (v === 'not_started') return '未开始'
  return '进行中'
}

// 通用数据映射函数
const mapEventData = (match) => ({
  id: match.match_id,
  eventId: match.event_id,
  eventName: eventMap.value[match.event_id] || '赛事活动',
  name: match.match_name,
  teamA: match.team_a_name || '主队',
  teamB: match.team_b_name || '客队',
  teamAAvatar: match.team_a_avatar ? `http://localhost:8080${match.team_a_avatar}` : '',
  teamBAvatar: match.team_b_avatar ? `http://localhost:8080${match.team_b_avatar}` : '',
  scoreA: match.score_team_a,
  scoreB: match.score_team_b,
  time: formatTime(match.match_time),
  timestamp: new Date(match.match_time).getTime(),
  sportType: sportIdMap[match.sport_id] || 'football',
  status: normalizeStatus(match.status),
  isSubscribed: match.is_subscribed
})

// 获取所有赛事信息（用于映射名称）
const fetchAllEvents = async () => {
  try {
    const res = await get('/events')
    if (res.code === 200 && res.data) {
      const map = {}
      res.data.forEach(e => {
        map[e.event_id] = e.event_name
      })
      eventMap.value = map
    }
  } catch (error) {
    console.error('Fetch all events error:', error)
  }
}

// 获取赛事列表
const fetchEvents = async () => {
  isLoading.value = true
  try {
    // 先获取赛事字典
    if (Object.keys(eventMap.value).length === 0) {
      await fetchAllEvents()
    }
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
      
      // Ensure events are fetched to populate the map
      if (Object.keys(eventMap.value).length === 0) {
        await fetchAllEvents()
      }

      subscribedEvents.value = res.data.subscribedMatches.map(item => ({
        id: item.matchId, // item.matchId is string
        eventId: item.eventId,
        eventName: eventMap.value[item.eventId] || '赛事活动',
        name: `${item.homeTeam.name} VS ${item.awayTeam.name}`, // Construct name if not provided
        teamA: item.homeTeam.name,
        teamB: item.awayTeam.name,
        teamAAvatar: item.homeTeam.avatar ? `http://localhost:8080${item.homeTeam.avatar}` : '',
        teamBAvatar: item.awayTeam.avatar ? `http://localhost:8080${item.awayTeam.avatar}` : '',
        scoreA: item.scoreA,
        scoreB: item.scoreB,
        time: item.matchTime, // Already formatted string? Backend says string.
        timestamp: new Date(item.matchTime).getTime(),
         venue: item.matchVenue,
         sportType: sportIdMap[item.sportId] || 'football',
         status: normalizeStatus(item.matchState === '未开始' ? 'not_started' : (item.matchState === '进行中' ? 'in_progress' : 'finished')),
         isSubscribed: true
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
  sessionStorage.setItem('student_event_active_tab', tab)
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
    const statusMatch =
      selectedStatus.value === 'all' || normalizeStatus(event.status) === normalizeStatus(selectedStatus.value)
    const eventMatch = selectedEventId.value === 'all' || String(event.eventId) === String(selectedEventId.value)
    return sportMatch && statusMatch && eventMatch
  })

  const now = Date.now()

  // 拆分三类状态
  const inProgressList = filtered.filter(e => normalizeStatus(e.status) === 'in_progress')
  const finishedList = filtered.filter(e => normalizeStatus(e.status) === 'finished')
  // 其余归为未开始
  const notStartedList = filtered.filter(
    e => normalizeStatus(e.status) !== 'in_progress' && normalizeStatus(e.status) !== 'finished',
  )

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
  sessionStorage.setItem('student_event_show_finished', String(showFinished.value))
  filterEvents()
}

const openMonthPicker = () => {
  if (selectedMonthStr.value) {
    const [y] = selectedMonthStr.value.split('-')
    pickerYear.value = parseInt(y)
  } else {
    pickerYear.value = new Date().getFullYear()
  }
  showMonthPicker.value = true
}

const closeMonthPicker = () => {
  showMonthPicker.value = false
}

const currentYear = new Date().getFullYear()
const currentMonth = new Date().getMonth() + 1

const selectMonth = (month) => {
  // Prevent future month selection
  if (pickerYear.value === currentYear && month > currentMonth) return
  if (pickerYear.value > currentYear) return

  // month is 1-12
  const m = month.toString().padStart(2, '0')
  selectedMonthStr.value = `${pickerYear.value}-${m}`
  showMonthPicker.value = false
  filterEvents()
}

const changePickerYear = (delta) => {
  const nextYear = pickerYear.value + delta
  if (delta > 0 && nextYear > currentYear) return
  pickerYear.value = nextYear
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
  if (event.eventId) {
    router.push(`/student/standings/${event.eventId}`)
  } else {
    // Fallback if eventId is missing (should not happen with new logic)
    router.push(`/student/standings`)
  }
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
      event.isSubscribed = true
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
      event.isSubscribed = false
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
const isSubscribed = (event) => {
  // 优先使用 event.isSubscribed 属性，如果是 undefined 则回退到数组检查（兼容旧逻辑）
  if (event.isSubscribed !== undefined) {
    return event.isSubscribed
  }
  return subscribedEvents.value.some(e => String(e.id) === String(event.id))
}

// 检查是否为待定比赛
const isTBD = (event) => {
  if (!event) return false
  const teamA = String(event.teamA || '').trim().toUpperCase()
  const teamB = String(event.teamB || '').trim().toUpperCase()
  return teamA === 'TBD' || teamB === 'TBD'
}

// 切换关注状态
const toggleSubscribe = async (event) => {
  if (isSubscribed(event)) {
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

// 跳转到个人中心
const goToProfile = () => {
  router.push('/profile')
}

// 获取轮播图数据
const fetchCarouselImages = async () => {
  try {
    const res = await get('/public/carousel')
    if (res.code === 200 && res.data && res.data.length > 0) {
      carouselImages.value = res.data.map(img => {
        let src = img.image_url
        if (src && src.startsWith('/')) {
          src = `${API_BASE_URL}${src}`
        }
        return {
          id: img.id,
          src: src,
          title: img.title,
          link: img.link_url // Ensure link is mapped
        }
      })
    }
  } catch (error) {
    console.error('Failed to fetch carousel images:', error)
  }
}

// 滚动处理与返回顶部
const showBackToTop = ref(false)

const handleScroll = () => {
  showBackToTop.value = window.scrollY > 300
}

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const restoreScrollPosition = () => {
  const savedPos = sessionStorage.getItem('student_event_scroll_pos')
  if (savedPos) {
    const targetPos = parseInt(savedPos)
    
    // 使用递归重试机制，对抗路由的默认滚动行为或异步渲染延迟
    const tryScroll = (attempts = 0) => {
      nextTick(() => {
        window.scrollTo({
          top: targetPos,
          behavior: 'instant'
        })
        
        // 如果当前位置与目标位置差距较大，且尝试次数未耗尽，则继续重试
        // 增加容差值(10px)，避免因浏览器计算精度问题导致死循环
        if (Math.abs(window.scrollY - targetPos) > 10 && attempts < 5) {
          setTimeout(() => tryScroll(attempts + 1), 50 + (attempts * 50)) // 递增延迟: 50, 100, 150...
        }
      })
    }
    
    tryScroll()
  }
}

// 路由离开前保存滚动位置
onBeforeRouteLeave((to, from, next) => {
  // 只有跳转到比赛详情页或积分榜才保存位置
  if (to.path.startsWith('/match/') || to.path.startsWith('/student/standings')) {
    sessionStorage.setItem('student_event_scroll_pos', window.scrollY)
    sessionStorage.setItem('student_event_selected_id', selectedEventId.value)
  } else {
    sessionStorage.removeItem('student_event_scroll_pos')
    sessionStorage.removeItem('student_event_selected_id')
  }
  next()
})

// 生命周期钩子
onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  window.addEventListener('click', closeEventDropdown)

  const savedTab = sessionStorage.getItem('student_event_active_tab')
  if (savedTab && ['home', 'all', 'subscribed'].includes(savedTab)) {
    activeTab.value = savedTab
  }

  const savedEventId = sessionStorage.getItem('student_event_selected_id')
  if (savedEventId) {
    selectedEventId.value = savedEventId
  }

  const savedShowFinished = sessionStorage.getItem('student_event_show_finished')
  if (savedShowFinished !== null) {
    showFinished.value = savedShowFinished === 'true'
  }

  fetchCarouselImages()
  
  // 处理数据加载和滚动恢复
  const loadData = async () => {
    // 并行加载数据
    const promises = [fetchEvents()]
    // 无论当前Tab是什么，都加载订阅信息以备用，但如果是订阅Tab，我们需要确保它完成后再恢复滚动
    const subPromise = fetchSubscribedEvents()
    promises.push(subPromise)
    
    await Promise.all(promises)
    
    // 如果当前是订阅Tab，需要手动触发一次filterEvents，因为fetchSubscribedEvents本身不触发
    if (activeTab.value === 'subscribed') {
      filterEvents()
    }
    
    // 数据加载完成后恢复滚动位置
    restoreScrollPosition()
  }

  loadData()
  startCarousel()

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

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
  window.removeEventListener('click', closeEventDropdown)
  stopCarousel()
})
</script>

<template>
  <div class="event-selection-container">
    <!-- 页面头部 -->
    <header class="page-header glass-effect">
      <div class="header-top">
        <div class="header-left">
          <h1>中大赛事平台</h1>
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
          <button class="action-text-btn refresh-btn" @click="refreshEvents" :disabled="isLoading">
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              width="20" 
              height="20" 
              viewBox="0 0 24 24" 
              fill="none" 
              stroke="currentColor" 
              stroke-width="2" 
              stroke-linecap="round" 
              stroke-linejoin="round"
              class="refresh-icon"
              :class="{ 'spin': isLoading }"
            >
              <path d="M21 12a9 9 0 0 0-9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"></path>
              <path d="M3 3v5h5"></path>
              <path d="M3 12a9 9 0 0 0 9 9 9.75 9.75 0 0 0 6.74-2.74L21 16"></path>
              <path d="M16 21h5v-5"></path>
            </svg>
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
          :class="{ active: activeTab === 'home' }"
          @click="switchTab('home')"
        >
          首页
        </button>
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

    <!-- 首页内容 -->
    <div v-if="activeTab === 'home'" class="home-view fade-in">
      <!-- 轮播图 -->
      <div 
        class="carousel-container glass-card"
        @touchstart="onTouchStart"
        @touchend="onTouchEnd"
      >
        <div class="carousel-slides" :style="{ transform: `translateX(-${currentCarouselIndex * 100}%)` }">
            <div 
              v-for="img in carouselImages" 
              :key="img.id" 
              class="carousel-slide"
              @click="handleCarouselClick(img)"
              :style="{ cursor: img.link ? 'pointer' : 'default' }"
            >
                <img :src="img.src" :alt="img.title" draggable="false" />
                <div class="carousel-caption glass-effect">{{ img.title }}</div>
            </div>
        </div>
        <div class="carousel-indicators">
            <span 
                v-for="(img, index) in carouselImages" 
                :key="index"
                class="indicator"
                :class="{ active: currentCarouselIndex === index }"
                @click="setCarouselIndex(index)"
            ></span>
        </div>
        <button class="carousel-nav prev" @click="setCarouselIndex((currentCarouselIndex - 1 + carouselImages.length) % carouselImages.length)">&#10094;</button>
        <button class="carousel-nav next" @click="setCarouselIndex((currentCarouselIndex + 1) % carouselImages.length)">&#10095;</button>
      </div>

      <!-- 历史最近赛事 -->
      <div class="section-container">
        <div class="section-title">
          <h3>📅 历史最近赛事</h3>
          <span class="section-subtitle">最近结束的精彩对决</span>
        </div>
        
        <div v-if="recentFinishedEvents.length === 0" class="empty-mini-state">
          <p>暂无已结束的赛事</p>
        </div>
        <div v-else class="mini-card-grid">
           <div 
             v-for="event in recentFinishedEvents" 
             :key="event.id" 
             class="mini-event-card glass-card" 
             @click="goToMatchDetail(event)"
           >
              <div class="mini-card-header">
                <span class="sport-badge-mini" :class="event.sportType">
                  {{ sportTypes.find(t => t.value === event.sportType)?.label }}
                </span>
                <span class="time-mini">{{ event.time.split(' ')[0] }}</span>
              </div>
              <div class="mini-card-content">
                <div class="team-mini">
                   <div class="avatar-mini">
                     <img v-if="event.teamAAvatar" :src="event.teamAAvatar" />
                     <span v-else>{{ event.teamA.charAt(0) }}</span>
                   </div>
                   <span class="team-name-mini">{{ event.teamA }}</span>
                </div>
                <div class="score-mini">
                  {{ event.scoreA }} : {{ event.scoreB }}
                </div>
                <div class="team-mini">
                   <div class="avatar-mini">
                     <img v-if="event.teamBAvatar" :src="event.teamBAvatar" />
                     <span v-else>{{ event.teamB.charAt(0) }}</span>
                   </div>
                   <span class="team-name-mini">{{ event.teamB }}</span>
                </div>
              </div>
              <div class="mini-card-footer">
                {{ event.eventName }}
              </div>
           </div>
        </div>
      </div>

      <!-- 即将开始赛事 -->
      <div class="section-container">
        <div class="section-title">
          <h3>🔥 即将开始</h3>
          <span class="section-subtitle">精彩赛事不容错过</span>
        </div>
        
        <div v-if="upcomingEvents.length === 0" class="empty-mini-state">
          <p>暂无即将开始的赛事</p>
        </div>
        <div v-else class="mini-card-grid">
           <div 
             v-for="event in upcomingEvents" 
             :key="event.id" 
             class="mini-event-card glass-card" 
             @click="goToMatchDetail(event)"
           >
              <div class="mini-card-header">
                <span class="sport-badge-mini" :class="event.sportType">
                  {{ sportTypes.find(t => t.value === event.sportType)?.label }}
                </span>
                <span class="time-mini">{{ event.time }}</span>
              </div>
              <div class="mini-card-content">
                <div class="team-mini">
                   <div class="avatar-mini">
                     <img v-if="event.teamAAvatar" :src="event.teamAAvatar" />
                     <span v-else>{{ event.teamA.charAt(0) }}</span>
                   </div>
                   <span class="team-name-mini">{{ event.teamA }}</span>
                </div>
                <div class="vs-mini">VS</div>
                <div class="team-mini">
                   <div class="avatar-mini">
                     <img v-if="event.teamBAvatar" :src="event.teamBAvatar" />
                     <span v-else>{{ event.teamB.charAt(0) }}</span>
                   </div>
                   <span class="team-name-mini">{{ event.teamB }}</span>
                </div>
              </div>
              <div class="mini-card-footer">
                {{ event.eventName }}
              </div>
           </div>
        </div>
      </div>
    </div>

    <!-- 筛选栏 (仅在全部赛事Tab显示) -->
    <div v-if="activeTab === 'all'" class="filter-bar glass-effect fade-in">
      <div class="filter-row">
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
        
        <div class="event-filter-wrapper">
          <div class="custom-select" :class="{ open: isEventDropdownOpen }">
            <div class="custom-select-trigger" @click.stop="toggleEventDropdown">
              <span>{{ selectedEventId === 'all' ? '所有赛事' : (eventOptions.find(e => e.id === selectedEventId)?.name || '选择赛事') }}</span>
              <div class="arrow"></div>
            </div>
            <div class="custom-options">
              <div 
                class="custom-option" 
                :class="{ selected: selectedEventId === 'all' }"
                @click="selectEvent('all')"
              >
                所有赛事
              </div>
              <div 
                v-for="evt in eventOptions" 
                :key="evt.id" 
                class="custom-option"
                :class="{ selected: selectedEventId === evt.id }"
                @click="selectEvent(evt.id)"
              >
                {{ evt.name }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 赛事列表 -->
    <div v-if="activeTab !== 'home'" class="event-list">
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
              <button 
                class="year-nav-btn" 
                :disabled="pickerYear >= currentYear"
                @click="changePickerYear(1)"
              >▶</button>
            </div>
            <div class="picker-grid">
              <button 
                v-for="m in 12" 
                :key="m" 
                class="month-cell"
                :class="{ 
                  active: selectedMonthStr === `${pickerYear}-${m.toString().padStart(2, '0')}`,
                  disabled: pickerYear === currentYear && m > currentMonth
                }"
                :disabled="pickerYear === currentYear && m > currentMonth"
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
            <div class="card-status-bar" :class="normalizeStatus(event.status)"></div>
            
            <div class="event-header">
              <div class="sport-badge" :class="event.sportType">
                {{ sportTypes.find(t => t.value === event.sportType)?.label || '赛事' }}
              </div>
              <div class="status-badge" :class="normalizeStatus(event.status)">
                {{ getStatusText(event.status) }}
              </div>
            </div>
            
            <h3 class="event-tournament-name">{{ event.eventName }}</h3>
            <div class="event-match-name">{{ event.name }}</div>
            
            <div class="match-up">
              <div class="team-block">
                <div class="team-avatar">
                  <img v-if="event.teamAAvatar" :src="event.teamAAvatar" :alt="event.teamA" class="team-logo-img" />
                  <span v-else>{{ event.teamA.charAt(0) }}</span>
                </div>
                <span class="team-name">{{ event.teamA }}</span>
              </div>
              <div v-if="normalizeStatus(event.status) === 'not_started'" class="vs-badge">VS</div>
              <div v-else class="score-badge">
                <span class="score">{{ event.scoreA }}</span>
                <span class="divider">:</span>
                <span class="score">{{ event.scoreB }}</span>
              </div>
              <div class="team-block">
                <div class="team-avatar">
                  <img v-if="event.teamBAvatar" :src="event.teamBAvatar" :alt="event.teamB" class="team-logo-img" />
                  <span v-else>{{ event.teamB.charAt(0) }}</span>
                </div>
                <span class="team-name">{{ event.teamB }}</span>
              </div>
            </div>
            
            <div class="event-meta">
              <div class="meta-item">
                <span class="meta-icon">📅</span>
                <span>{{ event.time }}</span>
              </div>
            </div>
            
            <div class="card-actions">
              <button 
                class="action-btn secondary"
                :class="{ 
                  'is-subscribed': !isTBD(event) && isSubscribed(event),
                  'disabled-btn': isTBD(event)
                }"
                :disabled="isTBD(event)"
                @click.stop="toggleSubscribe(event)"
              >
                {{ isTBD(event) ? '敬请期待' : (isSubscribed(event) ? '已关注' : '关注') }}
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

    <!-- 返回顶部按钮 -->
    <button 
      class="back-to-top-btn" 
      :class="{ visible: showBackToTop }" 
      @click="scrollToTop"
      title="回到顶部"
    >
      <span class="icon">↑</span>
    </button>
  </div>
</template>

<style scoped>
.back-to-top-btn {
  position: fixed;
  bottom: 80px;
  right: 24px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.9);
  color: var(--primary-color);
  font-size: 20px;
  font-weight: bold;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  opacity: 0;
  visibility: hidden;
  transform: translateY(20px) scale(0.9);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 999;
  border: 1px solid rgba(255, 255, 255, 0.5);
  backdrop-filter: blur(8px);
}

.back-to-top-btn.visible {
  opacity: 1;
  visibility: visible;
  transform: translateY(0) scale(1);
}

@media (hover: hover) {
  .back-to-top-btn:hover {
    background: var(--primary-color);
    color: white;
    transform: translateY(-4px) scale(1.05);
    box-shadow: 0 8px 20px rgba(79, 172, 254, 0.3);
    border-color: transparent;
  }
}

.back-to-top-btn:active {
  transform: translateY(-2px) scale(0.95);
}

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
  position: relative;
  z-index: 300;
}

.filter-row {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: space-between;
}

.event-filter-wrapper {
  flex-shrink: 0;
  position: relative;
  z-index: 100;
}

.custom-select {
  position: relative;
  user-select: none;
  min-width: 140px;
}

.custom-select-trigger {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary);
  background: white;
  border: 1px solid rgba(0,0,0,0.05);
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  gap: 8px;
}

.custom-select-trigger:hover {
  background: #f8f9fa;
  box-shadow: 0 4px 8px rgba(0,0,0,0.05);
}

.custom-select.open .custom-select-trigger {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(79, 172, 254, 0.2);
}

.arrow {
  width: 0; 
  height: 0; 
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-top: 5px solid var(--text-secondary);
  transition: transform 0.3s ease;
  opacity: 0.6;
}

.custom-select.open .arrow {
  transform: rotate(180deg);
}

.custom-options {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 100%;
  max-width: 240px;
  max-height: 240px;
  overflow-y: auto;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  opacity: 0;
  visibility: hidden;
  transform: translateY(-10px);
  transition: all 0.3s cubic-bezier(0.165, 0.84, 0.44, 1);
  padding: 6px;
  border: 1px solid rgba(0,0,0,0.05);
  z-index: 101;
}

.custom-select.open .custom-options {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.custom-option {
  padding: 10px 12px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.2s;
  border-radius: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.custom-option:hover {
  background: #f0f7ff;
  color: var(--primary-color);
}

.custom-option.selected {
  background: var(--primary-color);
  color: white;
  font-weight: 500;
}

.custom-options::-webkit-scrollbar {
  width: 4px;
}

.custom-options::-webkit-scrollbar-track {
  background: transparent;
}

.custom-options::-webkit-scrollbar-thumb {
  background: #e0e0e0;
  border-radius: 4px;
}

.custom-options::-webkit-scrollbar-thumb:hover {
  background: #ccc;
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
  flex: 1;
  width: auto;
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

.card-status-bar.ongoing {
  background: var(--success-color);
}

.card-status-bar.not_started {
  background: var(--warning-active);
}

.card-status-bar.finished {
  background: #e0e0e0;
}

.card-status-bar.completed {
  background: #e0e0e0;
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

.status-badge.ongoing::before { background: var(--success-color); }
.status-badge.ongoing { color: var(--success-active); }

.status-badge.not_started::before { background: var(--warning-active); }
.status-badge.not_started { color: var(--warning-active); }

.status-badge.finished::before { background: #bfbfbf; }
.status-badge.finished { color: #8c8c8c; }

.status-badge.completed::before { background: #bfbfbf; }
.status-badge.completed { color: #8c8c8c; }

.event-tournament-name {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 4px;
  line-height: 1.3;
}

.event-match-name {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-secondary); /* 灰色 */
  margin-bottom: 20px;
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

.team-logo-img {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
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

.action-btn.disabled-btn {
  background: rgba(0, 0, 0, 0.1) !important;
  color: #999 !important;
  border-color: transparent !important;
  cursor: not-allowed !important;
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

.month-cell.disabled {
  opacity: 0.3;
  cursor: not-allowed;
  background: #f5f5f5;
  border-color: transparent;
  color: #999;
  pointer-events: none;
}

.year-nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
  pointer-events: none;
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

/* Home View Styles */
.home-view {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

/* Carousel */
.carousel-container {
  position: relative;
  width: 100%;
  height: 300px;
  overflow: hidden;
  border-radius: 20px;
  margin-bottom: 30px;
  padding: 0 !important; /* Override glass-card padding */
  touch-action: pan-y; /* Allow vertical scroll but capture horizontal swipe */
  user-select: none;
}

.carousel-slides {
  display: flex;
  width: 100%;
  height: 100%;
  transition: transform 0.5s ease-in-out;
}

.carousel-slide {
  min-width: 100%;
  height: 100%;
  position: relative;
}

.carousel-slide img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  pointer-events: none; /* Prevent browser drag behavior */
}

.carousel-caption {
  position: absolute;
  bottom: 20px;
  left: 20px;
  padding: 8px 16px;
  border-radius: 12px;
  color: #333;
  font-weight: 600;
  font-size: 16px;
  background: rgba(255, 255, 255, 0.85);
}

.carousel-indicators {
  position: absolute;
  bottom: 15px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 8px;
}

.indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: all 0.3s;
}

.indicator.active {
  background: white;
  width: 24px;
  border-radius: 4px;
}

.carousel-nav {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  background: rgba(0, 0, 0, 0.3);
  color: white;
  border: none;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  transition: all 0.3s;
  opacity: 0;
}

.carousel-container:hover .carousel-nav {
  opacity: 1;
}

.carousel-nav:hover {
  background: rgba(0, 0, 0, 0.6);
}

.carousel-nav.prev { left: 10px; }
.carousel-nav.next { right: 10px; }

/* Section Styles */
.section-container {
  margin-bottom: 30px;
}

.section-title {
  margin-bottom: 16px;
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.section-title h3 {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.section-subtitle {
  font-size: 14px;
  color: var(--text-tertiary);
}

/* Mini Card Grid */
.mini-card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.mini-event-card {
  padding: 16px;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  gap: 12px;
  border: 1px solid rgba(255, 255, 255, 0.5);
}

.mini-event-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: var(--primary-light);
}

.mini-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sport-badge-mini {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 8px;
  font-weight: 600;
  text-transform: uppercase;
}

.sport-badge-mini.football { background: #e6f7ff; color: #1890ff; }
.sport-badge-mini.basketball { background: #f6ffed; color: #52c41a; }
.sport-badge-mini.badminton { background: #fff7e6; color: #faad14; }
.sport-badge-mini.volleyball { background: #fff1f0; color: #f5222d; }

.time-mini {
  font-size: 12px;
  color: var(--text-tertiary);
}

.mini-card-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.team-mini {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  width: 35%;
}

.avatar-mini {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  font-weight: 600;
  color: var(--text-secondary);
  border: 2px solid white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.avatar-mini img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.team-name-mini {
  font-size: 12px;
  color: var(--text-primary);
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 100%;
  font-weight: 500;
}

.score-mini {
  font-size: 18px;
  font-weight: 700;
  color: var(--text-primary);
  font-family: 'Monaco', monospace;
}

.vs-mini {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-tertiary);
  font-style: italic;
}

.mini-card-footer {
  font-size: 12px;
  color: var(--text-secondary);
  text-align: center;
  padding-top: 8px;
  border-top: 1px solid rgba(0,0,0,0.04);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.empty-mini-state {
  text-align: center;
  padding: 30px;
  color: var(--text-tertiary);
  background: rgba(255,255,255,0.4);
  border-radius: 12px;
}

/* Refresh Button */
.refresh-btn {
  color: #888 !important;
  padding: 0 10px !important;
}

.refresh-btn:active {
  color: var(--primary-color) !important;
  opacity: 0.7;
}

@media (hover: hover) {
  .refresh-btn:hover {
    color: var(--primary-color) !important;
  }
}

.refresh-icon.spin {
  animation: spin 1s linear infinite;
}
</style>
