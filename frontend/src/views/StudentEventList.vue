<script setup>
import { ref, onMounted, onUnmounted, computed, nextTick, watch } from 'vue'
import { useRouter, onBeforeRouteLeave } from 'vue-router'
import { get, post } from '@/utils/http'
import { useAuthStore } from '@/stores/auth'
import { marked } from 'marked'

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

// Kiwi Assistant State
const chatInput = ref('')
const chatMessages = ref([
  { 
    role: 'assistant', 
    content: '你好！我是中大体育智能助手Kiwi。我可以为你解答关于中山大学体育课、体测标准的问题，也可以帮你查询近期的体育赛事。有什么我可以帮你的吗？', 
    type: 'text' 
  }
])
const isChatLoading = ref(false)
const chatMessagesRef = ref(null)
// Initialize from sessionStorage directly
const selectedModel = ref(sessionStorage.getItem('kiwi_selected_model') || 'pro') 
const isModelDropdownOpen = ref(false)
const modelDropdownRef = ref(null)

// Persist selection
watch(selectedModel, (newVal) => {
  console.log('[Kiwi] Model changed to:', newVal)
  sessionStorage.setItem('kiwi_selected_model', newVal)
})

const toggleModelDropdown = () => {
  isModelDropdownOpen.value = !isModelDropdownOpen.value
}

const selectModel = (model) => {
  selectedModel.value = model
  isModelDropdownOpen.value = false
}

const handleClickOutside = (event) => {
  if (modelDropdownRef.value && !modelDropdownRef.value.contains(event.target)) {
    isModelDropdownOpen.value = false
  }
}

const scrollToBottom = () => {
  nextTick(() => {
    if (chatMessagesRef.value) {
      chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
    }
  })
}

// System Prompt with SYSU PE Info
const SYSTEM_PROMPT = `你叫Kiwi，是中山大学的体育智能助手。
你负责解答中山大学体育相关事项，以及提供运动指导。
你需要知道中大体测相关的事项：

一、体测成绩的组成与影响
1. 体育课权重：
   - 大一秋季：体测占体育课总分 30%
   - 大二、大三秋季：体测占体育课总分 80%
2. 保研门槛：大一至大三学年体测平均分需达及格（60分）及以上。
3. 毕业要求：体测成绩达不到 50 分者按结业处理。
4. 特色加分：体测总分达到良好（80分）或优秀（90分），可获体育课“课外积分”奖励（5-10分）。

二、核心测试项目与权重
- BMI (15%)
- 肺活量 (15%)
- 50米跑 (20%)
- 立定跳远 (10%)
- 坐位体前屈 (10%)
- 引体向上(男)/仰卧起坐(女) (10%)
- 1000米(男)/800米(女) (20%) [有额外加分，最高20分]

三、评分标准（简要）
- 及格(60分)：BMI正常范围，肺活量(男3100/女2000)，50米(男9.1/女10.3)，跳远(男208/女151)，体前屈(男3.7/女6.0)，引体10/仰卧26，1000米4'32"/800米4'34"。
- 优秀(90分)：指标更高，如1000米3'27"，引体17个等。

四、中大政策
- 加分规则：1000m/800m和引体/仰卧满分后可加分，最高20分。
- 奖学金与保研：必须及格。

工具能力：
工具能力：
1. search_events(query): 搜索相关的体育比赛 (Matches)。仅当用户询问具体比赛、赛程或结果时使用。如果用户询问“积分榜”、“排名”、“对阵图”，绝对不要使用此工具。
2. get_my_subscriptions(): 查询我关注的比赛。
3. set_sport_filter(sport_type): 筛选比赛类型 (football, basketball, badminton, volleyball, all)。
4. get_standings(query): 查询赛事积分榜或对阵信息（query为赛事名称）。
重要提示：
- 如果用户询问“积分榜”、“排名”、“第几名”、“对阵表”，必须使用 get_standings(query)。
- 如果用户询问“什么时候有比赛”、“查询篮球赛”，使用 search_events(query)。
界面会根据你的 tool_calls 自动显示相应卡片或执行操作。`

const renderMarkdown = (text) => {
  if (!text) return ''
  try {
    return marked(text)
  } catch (e) {
    return text
  }
}

const processChatResponse = async (messagesPayload) => {
  isChatLoading.value = true
  scrollToBottom()

  try {
    const response = await fetch('https://api.chatanywhere.org/v1/chat/completions', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer sk-yDfjbHsHoDsCyPox4EaSwlr5rw8HPZLYXx4gvhVb9zgC5LRK'
      },
      body: JSON.stringify({
        model: selectedModel.value === 'pro' ? 'gemini-3-pro-preview' : 'gemini-3-flash-preview',
        messages: messagesPayload,
        tools: [
            {
            type: 'function',
            function: {
              name: 'search_events',
              description: 'Search for specific MATCHES or GAMES (schedule, score). DO NOT use for standings/rankings/tables.',
              parameters: {
                type: 'object',
                properties: {
                  query: {
                    type: 'string',
                    description: 'Name of the team or event (e.g. "Lakers", "Finals")'
                  }
                },
                required: ['query']
              }
            }
          },
          {
            type: 'function',
            function: {
              name: 'get_my_subscriptions',
              description: '获取用户当前关注的比赛列表',
              parameters: { type: 'object', properties: {} }
            }
          },
          {
            type: 'function',
            function: {
              name: 'set_sport_filter',
              description: '设置比赛筛选类型',
              parameters: {
                type: 'object',
                properties: {
                  sport_type: {
                    type: 'string',
                    enum: ['football', 'basketball', 'badminton', 'volleyball', 'all'],
                    description: '运动类型'
                  }
                },
                required: ['sport_type']
              }
            }
          },

          {
             type: 'function',
             function: {
               name: 'get_standings',
               description: 'Get STANDINGS, RANKINGS, SCORE TABLE, or BRACKET for a tournament.',
               parameters: {
                 type: 'object',
                 properties: {
                   query: {
                     type: 'string',
                     description: 'Tournament name (e.g. "Super League", "Freshman Cup")'
                   }
                 },
                 required: ['query']
               }
             }
          }
        ]
      })
    })

    const data = await response.json()
    
    if (data.error) {
       console.error('API Error:', data.error)
       if (data.error.code === 'model_not_found') {
          chatMessages.value.push({ role: 'assistant', content: '抱歉，当前AI模型不可用，请联系管理员。', type: 'text' })
       } else {
          chatMessages.value.push({ role: 'assistant', content: '抱歉，我遇到了一些问题，请稍后再试。', type: 'text' })
       }
       return
    }

    const choice = data.choices[0]
    const message = choice.message

    if (message.content) {
      chatMessages.value.push({ role: 'assistant', content: message.content, type: 'text' })
    }

    if (message.tool_calls) {
      const toolCalls = message.tool_calls
      const toolResponses = []

      for (const toolCall of toolCalls) {
        const args = JSON.parse(toolCall.function.arguments)
        let result = ''

        if (toolCall.function.name === 'search_events') {
           const searchResults = searchEventsTool(args.query)
           if (searchResults.length > 0) {
              searchResults.forEach(event => {
                 chatMessages.value.push({ role: 'assistant', type: 'match-card', data: event })
              })
           }
           result = JSON.stringify(searchResults.map(e => ({ id: e.id, name: e.eventName + ' ' + e.name, time: e.time })))
        } else if (toolCall.function.name === 'get_my_subscriptions') {
           const subs = await getMySubscriptionsTool()
           result = JSON.stringify(subs)
        } else if (toolCall.function.name === 'set_sport_filter') {
           const success = setSportFilterTool(args.sport_type)
           result = success ? `已成功筛选为 ${args.sport_type}` : '筛选失败，类型无效'

        } else if (toolCall.function.name === 'get_standings') {
           const standingData = await getStandingsTool(args.query)
           if (standingData) {
             chatMessages.value.push({ role: 'assistant', type: 'standings-card', data: standingData })
             result = `已展示 ${standingData.eventName} 的${standingData.format_type === 'points' ? '积分榜' : '赛制信息'}。`
           } else {
             result = '未找到相关赛事的积分榜信息。'
           }
        }

        toolResponses.push({
           role: 'tool',
           tool_call_id: toolCall.id,
           name: toolCall.function.name,
           content: result
        })
      }

      const newMessages = [...messagesPayload, message, ...toolResponses]

      try {
        const secondResponse = await fetch('https://api.chatanywhere.org/v1/chat/completions', {
           method: 'POST',
           headers: {
             'Content-Type': 'application/json',
             'Authorization': 'Bearer sk-yDfjbHsHoDsCyPox4EaSwlr5rw8HPZLYXx4gvhVb9zgC5LRK'
           },
           body: JSON.stringify({
             model: selectedModel.value === 'pro' ? 'gemini-3-pro-preview' : 'gemini-3-flash-preview',
             messages: newMessages
           })
        })
        const secondData = await secondResponse.json()
        if (secondData.error) {
           console.error('[Kiwi] Second API Error:', secondData.error)
           chatMessages.value.push({ role: 'assistant', content: '处理您的请求时遇到了一些问题。', type: 'text' })
        } else if (secondData.choices && secondData.choices[0].message.content) {
           chatMessages.value.push({ role: 'assistant', content: secondData.choices[0].message.content, type: 'text' })
        }
      } catch (err2) {
         console.error('[Kiwi] Second API Network Error:', err2)
         chatMessages.value.push({ role: 'assistant', content: '网络连接不稳定，无法获取完整回复。', type: 'text' })
      }
    }

  } catch (error) {
    console.error('Chat Error:', error)
    chatMessages.value.push({ role: 'assistant', content: '网络连接异常，请检查网络。', type: 'text' })
  } finally {
    isChatLoading.value = false
    scrollToBottom()
  }
}

const sendMessage = async () => {
  if (!chatInput.value.trim() || isChatLoading.value) return
  
  const userMsg = chatInput.value.trim()
  chatMessages.value.push({ role: 'user', content: userMsg, type: 'text' })
  chatInput.value = ''
  
  const messages = [
    { role: 'system', content: SYSTEM_PROMPT },
    ...chatMessages.value.filter(m => m.type === 'text').map(m => ({ role: m.role, content: m.content }))
  ]
  
  await processChatResponse(messages)
}

const resetChat = () => {
  chatMessages.value = [
    { 
      role: 'assistant', 
      content: '你好！我是中大体育智能助手Kiwi。我可以为你解答关于中山大学体育课、体测标准的问题，也可以帮你查询近期的体育赛事。有什么我可以帮你的吗？', 
      type: 'text' 
    }
  ]
  sessionStorage.removeItem('kiwi_chat_messages')
}

const regenerateResponse = async () => {
  if (isChatLoading.value) return
  
  // Find the last user message index
  let lastUserIndex = -1
  for (let i = chatMessages.value.length - 1; i >= 0; i--) {
    if (chatMessages.value[i].role === 'user') {
      lastUserIndex = i
      break
    }
  }
  
  if (lastUserIndex === -1) return

  // Remove everything after the last user message (including the user message itself if we want to "re-send", but actually we just want to remove the *assistant's* response to it)
  // Actually, usually "regenerate" means "generate again for the last prompt".
  // So we keep the last user message, but remove subsequent assistant messages.
  
  chatMessages.value = chatMessages.value.slice(0, lastUserIndex + 1)
  
  const messages = [
    { role: 'system', content: SYSTEM_PROMPT },
    ...chatMessages.value.filter(m => m.type === 'text').map(m => ({ role: m.role, content: m.content }))
  ]
  
  await processChatResponse(messages)
}

const searchEventsTool = (query) => {
  if (!query) return []
  const lowerQuery = query.toLowerCase()
  
  // Anti-pattern check: if query contains "standings" or "rank" related keywords, return empty to avoid noise
  if (['积分', '排名', '榜', 'standings', 'rank'].some(k => lowerQuery.includes(k))) {
    return []
  }

  const results = events.value.filter(e => 
    e.name.toLowerCase().includes(lowerQuery) || 
    e.eventName.toLowerCase().includes(lowerQuery) ||
    e.teamA.toLowerCase().includes(lowerQuery) ||
    e.teamB.toLowerCase().includes(lowerQuery) ||
    (e.sportType && sportTypes.find(t => t.value === e.sportType)?.label.includes(lowerQuery))
  )
  
  // Sort: In progress > Not Started (nearest first) > Finished (recent first)
  results.sort((a, b) => {
      const statusMap = { 'in_progress': 3, 'ongoing': 3, 'not_started': 2, 'finished': 1, 'completed': 1 }
      const sA = statusMap[normalizeStatus(a.status)] || 0
      const sB = statusMap[normalizeStatus(b.status)] || 0
      if (sA !== sB) return sB - sA
      
      // If both finished, show recent first (desc)
      if (sA === 1) return (b.timestamp || 0) - (a.timestamp || 0)
      // If both future/ongoing, show nearest first (asc)
      return (a.timestamp || 0) - (b.timestamp || 0)
  })
  
  return results.slice(0, 3) 
}

const getMySubscriptionsTool = async () => {
  // Ensure latest subscriptions
  await fetchSubscribedEvents()
  if (subscribedEvents.value.length === 0) return '您目前没有关注任何比赛。'
  return subscribedEvents.value.map(e => ({
    id: e.id,
    name: `${e.eventName}: ${e.name}`,
    time: e.time,
    status: getStatusText(e.status)
  }))
}

const setSportFilterTool = (type) => {
  const validTypes = ['football', 'basketball', 'badminton', 'volleyball', 'all']
  if (!validTypes.includes(type)) return false
  
  selectedSportType.value = type
  if (activeTab.value !== 'all') {
    activeTab.value = 'all' // Switch to list view to show results
  }
  filterEvents()
  return true
}



const getStandingsTool = async (query) => {
  console.log('[Kiwi] getStandingsTool called with query:', query)
  if (!query) return null
  
  if (Object.keys(eventMap.value).length === 0) {
    console.log('[Kiwi] eventMap is empty, fetching events...')
    await fetchAllEvents()
  }

  const lowerQ = query.toLowerCase()
  let targetEventId = null
  let targetEventName = ''
  
  // Search in eventMap
  for (const [id, name] of Object.entries(eventMap.value)) {
    if (name.toLowerCase().includes(lowerQ)) {
      targetEventId = id
      targetEventName = name
      console.log('[Kiwi] Found matching event:', name, 'ID:', id)
      break
    }
  }
  
  if (!targetEventId) {
    console.log('[Kiwi] No matching event found for:', query)
    return null
  }
  
  try {
    const url = `/events/${targetEventId}/standings/overview`
    console.log('[Kiwi] Fetching standings from:', url)
    const res = await get(url)
    console.log('[Kiwi] Standings API Response:', res)

    if (res.code === 200 && res.data) {
       return {
         eventId: targetEventId,
         eventName: targetEventName,
         format_type: res.data.format_type,
         data: res.data 
       }
    } else {
       console.warn('[Kiwi] Standings API returned invalid data or non-200 code')
    }
  } catch (e) {
    console.error('[Kiwi] Standings tool error', e)
  }
  return null
}

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

const userAvatarUrl = computed(() => {
  if (!authStore.avatar) return null
  if (authStore.avatar.startsWith('http')) return authStore.avatar
  return `${API_BASE_URL}${authStore.avatar}`
})

// Watch Chat Messages
watch(chatMessages, (newVal) => {
  sessionStorage.setItem('kiwi_chat_messages', JSON.stringify(newVal))
}, { deep: true })

// 生命周期钩子
onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  window.addEventListener('click', closeEventDropdown)
  window.addEventListener('click', handleClickOutside)

  const savedTab = sessionStorage.getItem('student_event_active_tab')
  if (savedTab && ['home', 'all', 'subscribed', 'kiwi'].includes(savedTab)) {
    activeTab.value = savedTab
  }

  // Restore Chat
  const savedChat = sessionStorage.getItem('kiwi_chat_messages')
  if (savedChat) {
    try {
      chatMessages.value = JSON.parse(savedChat)
      scrollToBottom()
    } catch (e) {
      console.error('Failed to parse chat history', e)
    }
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
  window.removeEventListener('click', handleClickOutside)
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
        <button 
          class="tab-btn" 
          :class="{ active: activeTab === 'kiwi' }"
          @click="switchTab('kiwi')"
        >
          Kiwi助手
        </button>
      </div>
    </header>

    <!-- Kiwi 助手内容 -->
    <div v-if="activeTab === 'kiwi'" class="kiwi-view fade-in">
      <div class="chat-container glass-card">
        <div class="chat-header">
          <div class="header-left">
            <div class="avatar-kiwi">
              <svg xmlns="http://www.w3.org/2000/svg" width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="10" rx="2"></rect>
                <circle cx="12" cy="5" r="2"></circle>
                <path d="M12 7v4"></path>
                <line x1="8" y1="16" x2="8" y2="16"></line>
                <line x1="16" y1="16" x2="16" y2="16"></line>
              </svg>
            </div>
            <div class="chat-info">
              <h3>Kiwi助手</h3>
              <div class="status-container">
                <span class="status-dot"></span>
                <span class="status-text">在线</span>
              </div>
            </div>
          </div>

            <div class="header-actions" style="display: flex; gap: 8px; align-items: center;">
             <button class="action-icon-btn" style="width: 32px; height: 32px;" @click="resetChat" title="重置对话">
               <span style="font-size: 14px;">↺</span>
             </button>
             <div class="model-selector" ref="modelDropdownRef">
              <div class="custom-select" :class="{ open: isModelDropdownOpen }">
                <div class="select-trigger" @click="toggleModelDropdown">
                  <span class="selected-text">{{ selectedModel === 'pro' ? 'Gemini 3.0 Pro' : 'Gemini 3.0 Flash' }}</span>
                  <svg class="arrow-icon" xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="6 9 12 15 18 9"></polyline>
                  </svg>
                </div>
                <transition name="fade-slide">
                  <div class="options-list" v-if="isModelDropdownOpen">
                    <div class="option-item" :class="{ selected: selectedModel === 'pro' }" @click="selectModel('pro')">
                      <div class="option-content">
                        <span class="option-title">Gemini 3.0 Pro</span>
                        <span class="option-desc">推理能力更强</span>
                      </div>
                      <span v-if="selectedModel === 'pro'" class="check-icon">✓</span>
                    </div>
                    <div class="option-item" :class="{ selected: selectedModel === 'flash' }" @click="selectModel('flash')">
                      <div class="option-content">
                        <span class="option-title">Gemini 3.0 Flash</span>
                        <span class="option-desc">响应速度更快</span>
                      </div>
                      <span v-if="selectedModel === 'flash'" class="check-icon">✓</span>
                    </div>
                  </div>
                </transition>
              </div>
             </div>
            </div>
        </div>
        
        <div class="chat-messages" ref="chatMessagesRef">
          <div v-for="(msg, index) in chatMessages" :key="index" class="message-wrapper" :class="msg.role">
            <div v-if="msg.role === 'assistant'" class="avatar-kiwi-small">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="10" rx="2"></rect>
                <circle cx="12" cy="5" r="2"></circle>
                <path d="M12 7v4"></path>
                <line x1="8" y1="16" x2="8" y2="16"></line>
                <line x1="16" y1="16" x2="16" y2="16"></line>
              </svg>
            </div>
            <div class="message-content">
              <div v-if="msg.type === 'text'" class="text-bubble markdown-body" v-html="renderMarkdown(msg.content)"></div>
              <div v-else-if="msg.type === 'match-card'" class="match-recommend-card" @click="goToMatchDetail(msg.data)">
                <div class="mini-card-header">
                  <span class="sport-badge-mini" :class="msg.data.sportType">
                    {{ sportTypes.find(t => t.value === msg.data.sportType)?.label }}
                  </span>
                  <span class="status-badge-mini" :class="normalizeStatus(msg.data.status)">
                    {{ getStatusText(msg.data.status) }}
                  </span>
                </div>
                <div class="mini-card-content">
                  <div class="team-mini">
                    <div class="avatar-mini">
                      <img v-if="msg.data.teamAAvatar" :src="msg.data.teamAAvatar" :alt="msg.data.teamA" />
                      <span v-else>{{ msg.data.teamA.charAt(0) }}</span>
                    </div>
                    <span class="team-name-mini">{{ msg.data.teamA }}</span>
                  </div>
                  <div class="score-mini">VS</div>
                  <div class="team-mini">
                    <div class="avatar-mini">
                      <img v-if="msg.data.teamBAvatar" :src="msg.data.teamBAvatar" :alt="msg.data.teamB" />
                      <span v-else>{{ msg.data.teamB.charAt(0) }}</span>
                    </div>
                    <span class="team-name-mini">{{ msg.data.teamB }}</span>
                  </div>
                </div>
                <div class="mini-card-footer">
                  {{ msg.data.eventName }} | {{ msg.data.time }}
                </div>
              </div>
              
              <!-- Standings Card -->
              <div v-else-if="msg.type === 'standings-card'" class="standings-card">
                 <div class="sc-header">
                   <span class="sc-title">{{ msg.data.eventName }}</span>
                   <span class="sc-tag">{{ msg.data.format_type === 'points' ? '积分赛' : '淘汰赛/混合' }}</span>
                 </div>
                 
                 <div class="sc-content">
                    <!-- Points Table -->
                    <div v-if="msg.data.format_type === 'points' && msg.data.data.league" class="sc-table-wrapper">
                       <table class="sc-table">
                         <thead>
                           <tr>
                             <th>排名</th>
                             <th>队伍</th>
                             <th>积分</th>
                           </tr>
                         </thead>
                         <tbody>
                           <tr v-for="(team, i) in msg.data.data.league.slice(0, 5)" :key="i">
                             <td><span class="sc-rank" :class="'rank-'+(team.rank||i+1)">{{ team.rank||i+1 }}</span></td>
                             <td>{{ team.team_name }}</td>
                             <td class="sc-points">{{ team.points }}</td>
                           </tr>
                         </tbody>
                       </table>
                       <div v-if="msg.data.data.league.length > 5" class="sc-more">...</div>
                    </div>
                    
                    <!-- Knockout/Group Info -->
                    <div v-else class="sc-info-block">
                       <p>该赛事包含 {{ msg.data.data.groups ? '小组赛' : '' }} {{ msg.data.data.groups && msg.data.data.knockout_stages ? '+' : '' }} {{ msg.data.data.knockout_stages ? '淘汰赛' : '' }} 阶段。</p>
                       <p class="sc-hint">赛制较复杂，建议查看详情图表。</p>
                    </div>
                 </div>
                 
                 <div class="sc-footer">
                    <button class="sc-btn" @click="viewStandings({eventId: msg.data.eventId})">查看完整积分榜</button>
                 </div>
              </div>
              
              <!-- Regenerate Button -->
              <div v-if="msg.role === 'assistant' && index === chatMessages.length - 1 && !isChatLoading" 
                   class="regenerate-actions">
                 <button class="regenerate-btn" @click="regenerateResponse" title="重新生成" :disabled="index === 0">
                   <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                     <path d="M21.5 2v6h-6M2.5 22v-6h6M2 11.5a10 10 0 0 1 18.8-4.3M22 12.5a10 10 0 0 1-18.8 4.2"/>
                   </svg>
                   <span>重新生成</span>
                 </button>
              </div>
            </div>
            <div v-if="msg.role === 'user'" class="avatar-user-small">
              <img v-if="userAvatarUrl" :src="userAvatarUrl" alt="Me" class="user-avatar-img" />
              <span v-else>👤</span>
            </div>
          </div>
          <div v-if="isChatLoading" class="message-wrapper assistant">
            <div class="avatar-kiwi-small">
              <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="#1890ff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="3" y="11" width="18" height="10" rx="2"></rect>
                <circle cx="12" cy="5" r="2"></circle>
                <path d="M12 7v4"></path>
                <line x1="8" y1="16" x2="8" y2="16"></line>
                <line x1="16" y1="16" x2="16" y2="16"></line>
              </svg>
            </div>
            <div class="message-content">
              <div class="typing-indicator">
                <span></span><span></span><span></span>
              </div>
            </div>
          </div>
        </div>

        <div class="chat-input-area">
          <input 
            v-model="chatInput" 
            @keyup.enter="sendMessage"
            type="text" 
            placeholder="问问Kiwi关于中大体育的事情..." 
            :disabled="isChatLoading"
          />
          <button class="send-btn" @click="sendMessage" :disabled="!chatInput.trim() || isChatLoading">
            发送
          </button>
        </div>
      </div>
    </div>

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
    <div v-if="activeTab !== 'home' && activeTab !== 'kiwi'" class="event-list">
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

/* Kiwi Assistant Chat Styles */
.kiwi-view {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
  height: calc(100vh - 140px);
}

.chat-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.chat-header {
  padding: 16px 20px;
  background: white;
  border-bottom: 1px solid rgba(0,0,0,0.05);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.model-selector {
  position: relative;
  z-index: 100;
}

.custom-select {
  position: relative;
  width: 160px;
  font-size: 13px;
}

.select-trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--text-primary);
}

.select-trigger:hover, .custom-select.open .select-trigger {
  border-color: var(--primary-color);
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.15);
}

.arrow-icon {
  color: var(--text-secondary);
  transition: transform 0.3s ease;
}

.custom-select.open .arrow-icon {
  transform: rotate(180deg);
}

.options-list {
  position: absolute;
  top: calc(100% + 6px);
  right: 0;
  width: 200px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
  padding: 4px;
}

.option-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.2s;
}

.option-item:hover {
  background: #f5f7fa;
}

.option-item.selected {
  background: #e6f7ff;
  color: var(--primary-color);
}

.option-content {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.option-title {
  font-weight: 500;
  color: inherit;
}

.option-desc {
  font-size: 11px;
  color: var(--text-secondary);
}

.option-item.selected .option-desc {
  color: rgba(24, 144, 255, 0.8);
}

.check-icon {
  color: var(--primary-color);
  font-weight: bold;
}

/* Transitions */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.2s ease;
}

.fade-slide-enter-from,
.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.user-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.status-container {
  display: flex;
  align-items: center;
  margin-top: 2px;
}

.avatar-kiwi {
  font-size: 28px;
  background: #e6f7ff;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-kiwi-small {
  font-size: 20px;
  background: #e6f7ff;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-user-small {
  font-size: 20px;
  background: #f0f2f5;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.chat-info h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
}

.status-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #52c41a;
  border-radius: 50%;
  margin-right: 6px;
}

.status-text {
  font-size: 12px;
  color: var(--text-secondary);
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  scroll-behavior: smooth;
}

.message-wrapper {
  display: flex;
  gap: 12px;
  max-width: 85%;
}

.message-wrapper.assistant {
  align-self: flex-start;
}

.message-wrapper.user {
  align-self: flex-end;
}

.text-bubble {
  padding: 12px 16px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.5;
  word-break: break-word;
}

.assistant .text-bubble {
  background: white;
  color: var(--text-primary);
  border-top-left-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.user .text-bubble {
  background: #e6f7ff;
  color: var(--text-primary);
  border-top-right-radius: 4px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
  border: 1px solid rgba(24, 144, 255, 0.2);
}

.chat-input-area {
  padding: 16px 20px;
  background: white;
  border-top: 1px solid rgba(0,0,0,0.05);
  display: flex;
  gap: 12px;
}

.chat-input-area input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 24px;
  outline: none;
  font-size: 14px;
  transition: all 0.2s;
}

.chat-input-area input:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(79, 172, 254, 0.1);
}

.send-btn {
  padding: 0 24px;
  border-radius: 24px;
  background: var(--primary-color);
  color: white;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.send-btn:hover {
  background: var(--primary-active);
  transform: translateY(-1px);
}

.send-btn:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.match-recommend-card {
  background: white;
  border-radius: 16px;
  padding: 12px;
  width: 260px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid rgba(0,0,0,0.05);
}

.match-recommend-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.12);
  border-color: var(--primary-color);
}

.typing-indicator span {
  display: inline-block;
  width: 6px;
  height: 6px;
  background: #ccc;
  border-radius: 50%;
  margin: 0 2px;
  animation: typing 1.4s infinite ease-in-out both;
}

.typing-indicator span:nth-child(1) { animation-delay: -0.32s; }
.typing-indicator span:nth-child(2) { animation-delay: -0.16s; }

@keyframes typing {
  0%, 80%, 100% { transform: scale(0); }
  40% { transform: scale(1); }
}

.status-badge-mini {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  background: #f0f0f0;
  color: #666;
}

.status-badge-mini.ongoing { background: #e6ffec; color: #52c41a; }
.status-badge-mini.finished { background: #f5f5f5; color: #999; }

/* Markdown Styles */
.markdown-body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
}
.markdown-body p {
  margin-bottom: 8px;
}
.markdown-body p:last-child {
  margin-bottom: 0;
}
.markdown-body ul, .markdown-body ol {
  padding-left: 20px;
  margin-bottom: 8px;
}
.markdown-body code {
  background: rgba(0,0,0,0.05);
  padding: 2px 4px;
  border-radius: 4px;
  font-family: monospace;
}
.markdown-body pre {
  background: #f6f8fa;
  padding: 10px;
  border-radius: 8px;
  overflow-x: auto;
  margin-bottom: 8px;
}
.user .markdown-body code {
  background: rgba(255,255,255,0.2);
}
.user .markdown-body pre {
  background: rgba(0,0,0,0.1);
}

/* Regenerate Button */
.regenerate-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 4px;
}

.regenerate-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  background: transparent;
  border: none;
  color: var(--text-tertiary);
  font-size: 12px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 12px;
  transition: all 0.2s;
  opacity: 0.6;
}

.message-wrapper:hover .regenerate-btn {
  opacity: 1;
}

.regenerate-btn:hover {
  background: rgba(0,0,0,0.05);
  color: var(--primary-color);
}

.regenerate-btn:disabled {
  opacity: 0.3 !important;
  cursor: not-allowed;
  color: #ccc !important;
  background: transparent !important;
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

/* Standings Card */
.standings-card {
  background: white;
  border-radius: 16px;
  padding: 16px;
  width: 280px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.08);
  border: 1px solid rgba(0,0,0,0.05);
  margin-top: 8px;
}

.sc-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
  padding-bottom: 8px;
}

.sc-title {
  font-weight: 700;
  font-size: 14px;
  color: var(--text-primary);
}

.sc-tag {
  font-size: 10px;
  padding: 2px 6px;
  background: #e6f7ff;
  color: #1890ff;
  border-radius: 4px;
}

.sc-table {
  width: 100%;
  font-size: 12px;
  border-collapse: collapse;
}

.sc-table th {
  text-align: left;
  color: var(--text-tertiary);
  font-weight: 400;
  padding-bottom: 4px;
}

.sc-table td {
  padding: 4px 0;
  color: var(--text-primary);
}

.sc-rank {
  display: inline-block;
  width: 16px;
  height: 16px;
  text-align: center;
  line-height: 16px;
  border-radius: 50%;
  background: #f0f0f0;
  font-size: 10px;
  color: #666;
}

.sc-rank.rank-1 { background: #fff1b8; color: #faad14; }
.sc-rank.rank-2 { background: #e6f7ff; color: #1890ff; }
.sc-rank.rank-3 { background: #fff0f6; color: #eb2f96; }

.sc-points {
  font-weight: 700;
  text-align: right;
}

.sc-table th:last-child { text-align: right; }

.sc-more {
  text-align: center;
  color: #ccc;
  font-size: 12px;
}

.sc-info-block {
  font-size: 13px;
  color: var(--text-secondary);
  padding: 8px 0;
  line-height: 1.5;
}

.sc-hint {
  font-size: 12px;
  color: var(--text-tertiary);
  margin-top: 4px;
}

.sc-footer {
  margin-top: 12px;
}

.sc-btn {
  width: 100%;
  padding: 8px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.sc-btn:hover {
  background: var(--primary-active);
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
