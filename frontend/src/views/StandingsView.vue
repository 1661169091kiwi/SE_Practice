<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/utils/http'

const route = useRoute()
const router = useRouter()

const isLoading = ref(false)
const events = ref([])
const searchText = ref('')
const selectedEventId = ref(route.params.eventId ? String(route.params.eventId) : '')
const isDropdownOpen = ref(false)

const eventId = ref(0)
const formatType = ref('')
const leagueStandings = ref([])
const overview = ref(null)
const activeTab = ref('groups')
const knockoutViewMode = ref('bracket')

// 响应式布局参数
const bracketConfig = ref({
  matchHeight: 64,
  gap: 14,
  colWidth: 240,
  colGap: 40, // 增加列间距以容纳连接线
  topPadding: 34,
  isMobile: false
})

// 监听窗口大小变化
const updateLayoutConfig = () => {
  const width = window.innerWidth
  // Force mobile/pyramid layout on all screens as requested
  if (width < 99999) {
    // 移动端紧凑布局
    bracketConfig.value = {
      matchHeight: 70, // 增加高度以容纳更多信息
      gap: 10,
      colWidth: 160,
      colGap: 24,
      topPadding: 20,
      isMobile: true
    }
  } else {
    // 桌面端标准布局
    bracketConfig.value = {
      matchHeight: 64,
      gap: 14,
      colWidth: 240,
      colGap: 40,
      topPadding: 34,
      isMobile: false
    }
  }
}

const currentEventName = computed(() => {
  const e = events.value.find(x => String(x.event_id) === selectedEventId.value)
  return e ? e.event_name : '请选择赛事'
})

const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
}

const selectEvent = (id) => {
  selectedEventId.value = String(id)
  isDropdownOpen.value = false
}

// Close dropdown when clicking outside
const closeDropdown = (e) => {
  if (isDropdownOpen.value && !e.target.closest('.custom-select')) {
    isDropdownOpen.value = false
  }
}

const formatMatchTime = (timeStr) => {
  if (!timeStr || timeStr.startsWith('0001') || timeStr.startsWith('1000')) {
    return 'TBD'
  }
  return timeStr
}

onMounted(async () => {
  document.addEventListener('click', closeDropdown)
  window.addEventListener('resize', updateLayoutConfig)
  updateLayoutConfig() // 初始化
  await fetchEvents()
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
  window.removeEventListener('resize', updateLayoutConfig)
})

const knockoutStages = computed(() => (Array.isArray(overview.value?.knockout_stages) ? overview.value.knockout_stages : []))

const bracketLayout = computed(() => {
  const stages = knockoutStages.value
  const { matchHeight, gap, colWidth, colGap, topPadding, isMobile } = bracketConfig.value
  
  // 移动端垂直布局逻辑 (金字塔/聚合布局)
  if (isMobile) {
    if (!stages || stages.length === 0) return { matches: [], height: 0, isPyramid: true, upperRows: [], lowerRows: [], finalMatch: null }
    
    // 假设最后一个阶段是决赛
    const finalStage = stages[stages.length - 1]
    const finalMatch = finalStage.matches && finalStage.matches.length > 0 ? finalStage.matches[0] : null
    
    // 构建上半区和下半区
    // 假设 stages 顺序为：16强 -> 8强 -> 4强 -> 决赛
    // 我们需要将每一轮的比赛分为上半区和下半区
    // 上半区：每一轮的前半部分比赛
    // 下半区：每一轮的后半部分比赛
    
    const upperRows = []
    const lowerRows = []
    
    // 遍历除决赛外的所有阶段
    for (let i = 0; i < stages.length - 1; i++) {
      const stage = stages[i]
      const matches = stage.matches || []
      const mid = Math.ceil(matches.length / 2)
      
      const upperMatches = matches.slice(0, mid)
      const lowerMatches = matches.slice(mid)
      
      upperRows.push({
        stage_name: stage.stage_name,
        matches: upperMatches
      })
      
      lowerRows.push({
        stage_name: stage.stage_name,
        matches: lowerMatches
      })
    }
    
    // 下半区需要反转顺序吗？
    // 上半区渲染顺序：Row 0 (R16) -> Row 1 (QF) -> Row 2 (SF) -> Final
    // 视觉上：Top -> Bottom.
    // 下半区渲染顺序：Final -> Row 2 (SF) -> Row 1 (QF) -> Row 0 (R16)
    // 视觉上：Top -> Bottom.
    // 所以 lowerRows 需要反转，使得 SF 在最前（紧邻决赛），R16 在最后
    lowerRows.reverse()
    
    return {
      isPyramid: true,
      upperRows,
      lowerRows,
      finalMatch,
      // 保留旧的 matches 属性以防万一，但主要使用上面三个
      matches: []
    }
  }

  // 桌面端水平布局逻辑 (保持不变)
  const cols = stages.map((st, colIndex) => {
    const matches = Array.isArray(st.matches) ? st.matches : []
    const baseStep = matchHeight + gap
    // 计算当前列每场比赛的垂直间距步长
    const step = baseStep * (2 ** colIndex)
    // 计算起始偏移量，确保树形垂直居中
    const offset = (baseStep * ((2 ** colIndex) - 1)) / 2
    
    const placed = matches.map((m, i) => {
      // 计算连接线高度 (上一列两场比赛中心点的距离)
      // 对于第i列(i>0)，连接线高度等于上一列的step
      // 上一列是 colIndex - 1, step = baseStep * 2^(colIndex-1)
      const prevStep = colIndex > 0 ? baseStep * (2 ** (colIndex - 1)) : 0
      
      return {
        ...m,
        _top: topPadding + i * step + offset,
        _connectorHeight: prevStep
      }
    })
    return { ...st, _matches: placed }
  })

  const firstCount = cols[0]?._matches?.length || 0
  const baseStep = matchHeight + gap
  const height = firstCount
    ? topPadding + (firstCount - 1) * baseStep + matchHeight + 24
    : 0
  const width = cols.length
    ? cols.length * colWidth + (cols.length - 1) * colGap
    : 0

  return { columns: cols, height, width }
})

const filteredEvents = computed(() => {
  const q = searchText.value.trim().toLowerCase()
  const list = events.value
  if (!q) return list
  return list.filter(e => {
    const name = String(e.event_name || '').toLowerCase()
    return name.includes(q)
  })
})

const fetchEvents = async () => {
  try {
    const res = await get('/events', { params: { _t: Date.now() } })
    if (res.code === 200 && Array.isArray(res.data)) {
      events.value = res.data
    } else {
      events.value = []
    }
  } catch {
    events.value = []
  }
}

const loadByEventId = async (eid) => {
  if (!eid) {
    eventId.value = 0
    formatType.value = ''
    leagueStandings.value = []
    overview.value = null
    return
  }

  isLoading.value = true
  eventId.value = Number(eid)

  try {
    const ovRes = await get(`/events/${eid}/standings/overview`, { params: { _t: Date.now() } })
    const ov = ovRes?.code === 200 ? ovRes.data : null
    const ft = ov?.format_type || ''
    formatType.value = ft

    if (ft === 'points') {
      activeTab.value = 'league'
      leagueStandings.value = Array.isArray(ov?.league) ? ov.league : []
      overview.value = null
      return
    }
    if (ft === 'group_knockout') {
      activeTab.value = 'groups'
      overview.value = ov
      leagueStandings.value = []
      knockoutViewMode.value = 'bracket'
      return
    }
    if (ft === 'knockout') {
      activeTab.value = 'knockout'
      overview.value = ov
      leagueStandings.value = []
      knockoutViewMode.value = 'bracket'
      return
    }

    leagueStandings.value = []
    overview.value = null
  } catch {
    leagueStandings.value = []
    overview.value = null
  } finally {
    isLoading.value = false
  }
}

const goBack = () => {
  router.back()
}

watch(
  () => route.params.eventId,
  (newId) => {
    if (newId) {
      selectedEventId.value = String(newId)
    } else {
      selectedEventId.value = ''
    }
  },
  { immediate: true }
)

watch(
  () => selectedEventId.value,
  (v) => {
    // Only replace route if it differs from current param to avoid redundant navigation
    if (v && v !== route.params.eventId) {
      router.replace(`/student/standings/${v}`)
    }
    loadByEventId(v)
  },
  { immediate: true }
)


</script>

<template>
  <div class="standings-container">
    <header class="page-header">
      <button class="back-button" @click="goBack">
        ←
      </button>
      <div class="header-main">
        <h1>积分榜</h1>
        <div class="header-controls">
          <input
            v-model="searchText"
            class="search-input"
            type="text"
            placeholder="搜索赛事"
          />
          <div class="custom-select" :class="{ open: isDropdownOpen }">
            <div class="select-trigger" @click.stop="toggleDropdown">
              <span class="trigger-text">{{ currentEventName }}</span>
              <span class="arrow" :class="{ rotated: isDropdownOpen }">▼</span>
            </div>
            <div v-if="isDropdownOpen" class="select-options">
              <div 
                v-for="e in filteredEvents" 
                :key="e.event_id" 
                class="select-option" 
                @click="selectEvent(e.event_id)"
                :class="{ selected: String(e.event_id) === selectedEventId }"
              >
                {{ e.event_name }}
              </div>
              <div v-if="filteredEvents.length === 0" class="select-option disabled">
                无匹配赛事
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="placeholder"></div>
    </header>

    <div class="content-area">
      <div v-if="isLoading" class="loading-state">
        <p>加载中...</p>
      </div>

      <div v-else-if="!selectedEventId" class="empty-state">
        <p>请先选择一个赛事</p>
      </div>

      <div v-else-if="formatType === 'points'">
        <div v-if="leagueStandings.length === 0" class="empty-state">
          <p>暂无积分数据</p>
        </div>
        <div v-else class="table-container">
          <table class="standings-table">
            <thead>
              <tr>
                <th class="rank-col">排名</th>
                <th class="team-col">队名</th>
                <th>场次</th>
                <th>胜/平/负</th>
                <th>进/失</th>
                <th>积分</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in leagueStandings" :key="index">
                <td class="rank-col">
                  <span class="rank-badge" :class="`rank-${team.rank || index + 1}`">{{ team.rank || index + 1 }}</span>
                </td>
                <td class="team-col">{{ team.team_name }}</td>
                <td>{{ team.played }}</td>
                <td>{{ team.won }}/{{ team.drawn }}/{{ team.lost }}</td>
                <td>{{ team.goals_for }}/{{ team.goals_against }}</td>
                <td class="points-col">{{ team.points }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div v-else-if="formatType === 'group_knockout'">
        <div class="cup-tabs">
          <button class="tab-btn" :class="{ active: activeTab === 'groups' }" @click="activeTab = 'groups'">小组赛</button>
          <button class="tab-btn" :class="{ active: activeTab === 'knockout' }" @click="activeTab = 'knockout'">淘汰赛</button>
        </div>

        <div v-if="!overview" class="empty-state">
          <p>暂无杯赛数据</p>
        </div>

        <div v-else-if="activeTab === 'groups'">
          <div v-if="!overview.groups || overview.groups.length === 0" class="empty-state">
            <p>暂无小组赛积分数据</p>
          </div>
          <div v-else class="groups-container">
            <div v-for="g in overview.groups" :key="g.group" class="group-block">
              <div class="group-title">{{ g.group }}</div>
              <div v-if="!g.standings || g.standings.length === 0" class="empty-state">
                <p>暂无积分数据</p>
              </div>
              <div v-else class="table-container">
                <table class="standings-table">
                  <thead>
                    <tr>
                      <th class="rank-col">排名</th>
                      <th class="team-col">队名</th>
                      <th>场次</th>
                      <th>胜/平/负</th>
                      <th>进/失</th>
                      <th>积分</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(team, index) in g.standings" :key="index">
                      <td class="rank-col">
                        <span class="rank-badge" :class="`rank-${team.rank || index + 1}`">{{ team.rank || index + 1 }}</span>
                      </td>
                      <td class="team-col">{{ team.team_name }}</td>
                      <td>{{ team.played }}</td>
                      <td>{{ team.won }}/{{ team.drawn }}/{{ team.lost }}</td>
                      <td>{{ team.goals_for }}/{{ team.goals_against }}</td>
                      <td class="points-col">{{ team.points }}</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        <div v-else>
          <div v-if="!overview.knockout_stages || overview.knockout_stages.length === 0" class="empty-state">
            <p>暂无淘汰赛对阵数据</p>
          </div>
          <div v-else>
            <div class="knockout-view-tabs">
              <button class="mode-btn" :class="{ active: knockoutViewMode === 'bracket' }" @click="knockoutViewMode = 'bracket'">对阵图</button>
              <button class="mode-btn" :class="{ active: knockoutViewMode === 'list' }" @click="knockoutViewMode = 'list'">列表</button>
            </div>

            <div v-if="knockoutViewMode === 'bracket'">
              <!-- Mobile Pyramid View -->
              <div v-if="bracketLayout.isPyramid" class="bracket-pyramid">
                <!-- Upper Section -->
                <div class="pyramid-section upper-section">
                  <div v-for="(row, idx) in bracketLayout.upperRows" :key="'up-'+idx" class="pyramid-row">
                    <div class="row-matches">
                      <div v-for="m in row.matches" :key="m.knockout_match_id" class="match-card">
                        <div class="mc-content">
                          <div class="mc-team">
                            <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                            <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_a }}</span>
                          </div>
                          <div class="mc-team">
                            <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                            <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_b }}</span>
                          </div>
                        </div>
                        <div v-if="idx === bracketLayout.upperRows.length - 1" class="connector-to-final-top"></div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Final Match -->
                <div v-if="bracketLayout.finalMatch" class="pyramid-final">
                  <div class="final-card">
                    <div class="fc-header">FINAL</div>
                    <div class="fc-content">
                      <div class="fc-team">
                        <span class="fc-name">{{ bracketLayout.finalMatch.team_a_name || 'TBD' }}</span>
                      </div>
                      <div class="fc-score">
                        <span>{{ bracketLayout.finalMatch.status === 'not_started' ? '-' : bracketLayout.finalMatch.score_team_a }}</span>
                        <span class="divider">:</span>
                        <span>{{ bracketLayout.finalMatch.status === 'not_started' ? '-' : bracketLayout.finalMatch.score_team_b }}</span>
                      </div>
                      <div class="fc-team">
                        <span class="fc-name">{{ bracketLayout.finalMatch.team_b_name || 'TBD' }}</span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Lower Section -->
                <div class="pyramid-section lower-section">
                  <div v-for="(row, idx) in bracketLayout.lowerRows" :key="'low-'+idx" class="pyramid-row">
                    <div class="row-matches">
                      <div v-for="m in row.matches" :key="m.knockout_match_id" class="match-card">
                        <div class="mc-content">
                          <div class="mc-team">
                            <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                            <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_a }}</span>
                          </div>
                          <div class="mc-team">
                            <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                            <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_b }}</span>
                          </div>
                        </div>
                        <div v-if="idx === 0" class="connector-to-final-bot"></div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Desktop Bracket View -->
              <div v-else class="bracket-scroll">
                <div class="bracket" :style="{ height: `${bracketLayout.height}px`, width: `${bracketLayout.width}px` }">
                  <div
                    v-for="(stage, idx) in bracketLayout.columns"
                    :key="stage.stage_id"
                    class="bracket-col"
                    :style="{ left: `${idx * (bracketConfig.colWidth + bracketConfig.colGap)}px`, width: `${bracketConfig.colWidth}px` }"
                  >
                    <div class="bracket-col-title">{{ stage.stage_name }}</div>
                    <div
                      v-for="m in stage._matches"
                      :key="m.knockout_match_id"
                      class="bracket-match"
                      :style="{ top: `${m._top}px`, height: `${bracketConfig.matchHeight}px` }"
                    >
                      <!-- 连接线 -->
                      <div v-if="idx > 0" class="connector" 
                           :style="{ 
                             height: `${m._connectorHeight}px`, 
                             left: `-${bracketConfig.colGap}px`, 
                             width: `${bracketConfig.colGap}px`
                           }">
                        <span class="line-v"></span>
                        <span class="line-h-parent"></span>
                        <span class="line-h-child-top"></span>
                        <span class="line-h-child-bot"></span>
                      </div>

                      <div class="bm-row">
                        <span class="bm-team" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                        <span class="bm-score">{{ m.status === 'not_started' ? 'VS' : `${m.score_team_a}-${m.score_team_b}` }}</span>
                        <span class="bm-team" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                      </div>
                      <div class="bm-meta">
                        <span>{{ formatMatchTime(m.match_time) }}</span>
                        <span class="bm-status">{{ m.status }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="knockout-container">
              <div v-for="stage in overview.knockout_stages" :key="stage.stage_id" class="stage-block">
                <div class="stage-title">{{ stage.stage_name }}</div>
                <div v-if="!stage.matches || stage.matches.length === 0" class="empty-state">
                  <p>暂无对阵</p>
                </div>
                <div v-else class="stage-matches">
                  <div v-for="m in stage.matches" :key="m.knockout_match_id" class="knockout-match">
                    <div class="km-row">
                      <span class="km-team">{{ m.team_a_name || 'TBD' }}</span>
                      <span class="km-score">{{ m.status === 'not_started' ? 'VS' : `${m.score_team_a}-${m.score_team_b}` }}</span>
                      <span class="km-team">{{ m.team_b_name || 'TBD' }}</span>
                    </div>
                    <div class="km-meta">
                      <span>{{ formatMatchTime(m.match_time) }}</span>
                      <span class="km-status">{{ m.status }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="formatType === 'knockout'">
        <div v-if="!overview" class="empty-state">
          <p>暂无淘汰赛对阵数据</p>
        </div>
        <div v-else-if="!overview.knockout_stages || overview.knockout_stages.length === 0" class="empty-state">
          <p>暂无淘汰赛对阵数据</p>
        </div>
        <div v-else>
          <div class="knockout-view-tabs">
            <button class="mode-btn" :class="{ active: knockoutViewMode === 'bracket' }" @click="knockoutViewMode = 'bracket'">对阵图</button>
            <button class="mode-btn" :class="{ active: knockoutViewMode === 'list' }" @click="knockoutViewMode = 'list'">列表</button>
          </div>

          <div v-if="knockoutViewMode === 'bracket'">
            <!-- Mobile Pyramid View -->
            <div v-if="bracketLayout.isPyramid" class="bracket-pyramid">
              <!-- Upper Section -->
              <div class="pyramid-section upper-section">
                <div v-for="(row, idx) in bracketLayout.upperRows" :key="'up-'+idx" class="pyramid-row">
                  <div class="row-matches">
                    <div v-for="m in row.matches" :key="m.knockout_match_id" class="match-card">
                      <div class="mc-content">
                        <div class="mc-team">
                          <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                          <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_a }}</span>
                        </div>
                        <div class="mc-team">
                          <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                          <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_b }}</span>
                        </div>
                      </div>
                      <div v-if="idx === bracketLayout.upperRows.length - 1" class="connector-to-final-top"></div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Final Match -->
              <div v-if="bracketLayout.finalMatch" class="pyramid-final">
                <div class="final-card">
                  <div class="fc-header">FINAL</div>
                  <div class="fc-content">
                    <div class="fc-team">
                      <span class="fc-name">{{ bracketLayout.finalMatch.team_a_name || 'TBD' }}</span>
                    </div>
                    <div class="fc-score">
                      <span>{{ bracketLayout.finalMatch.status === 'not_started' ? '-' : bracketLayout.finalMatch.score_team_a }}</span>
                      <span class="divider">:</span>
                      <span>{{ bracketLayout.finalMatch.status === 'not_started' ? '-' : bracketLayout.finalMatch.score_team_b }}</span>
                    </div>
                    <div class="fc-team">
                      <span class="fc-name">{{ bracketLayout.finalMatch.team_b_name || 'TBD' }}</span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Lower Section -->
              <div class="pyramid-section lower-section">
                <div v-for="(row, idx) in bracketLayout.lowerRows" :key="'low-'+idx" class="pyramid-row">
                  <div class="row-matches">
                    <div v-for="m in row.matches" :key="m.knockout_match_id" class="match-card">
                      <div class="mc-content">
                        <div class="mc-team">
                          <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                          <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_a }}</span>
                        </div>
                        <div class="mc-team">
                          <span class="mc-team-name" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                          <span class="mc-score" v-if="m.status !== 'not_started'">{{ m.score_team_b }}</span>
                        </div>
                      </div>
                      <div v-if="idx === 0" class="connector-to-final-bot"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Desktop Bracket View -->
            <div v-else class="bracket-scroll">
              <div class="bracket" :style="{ height: `${bracketLayout.height}px`, width: `${bracketLayout.width}px` }">
                <div
                  v-for="(stage, idx) in bracketLayout.columns"
                  :key="stage.stage_id"
                  class="bracket-col"
                  :style="{ left: `${idx * (bracketConfig.colWidth + bracketConfig.colGap)}px`, width: `${bracketConfig.colWidth}px` }"
                >
                  <div class="bracket-col-title">{{ stage.stage_name }}</div>
                  <div
                    v-for="m in stage._matches"
                    :key="m.knockout_match_id"
                    class="bracket-match"
                    :style="{ top: `${m._top}px`, height: `${bracketConfig.matchHeight}px` }"
                  >
                    <!-- 连接线 -->
                    <div v-if="idx > 0" class="connector" 
                         :style="{ 
                           height: `${m._connectorHeight}px`, 
                           left: `-${bracketConfig.colGap}px`, 
                           width: `${bracketConfig.colGap}px`
                         }">
                      <span class="line-v"></span>
                      <span class="line-h-parent"></span>
                      <span class="line-h-child-top"></span>
                      <span class="line-h-child-bot"></span>
                    </div>

                    <div class="bm-row">
                      <span class="bm-team" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_a_id }">{{ m.team_a_name || 'TBD' }}</span>
                      <span class="bm-score">{{ m.status === 'not_started' ? 'VS' : `${m.score_team_a}-${m.score_team_b}` }}</span>
                      <span class="bm-team" :class="{ winner: m.winner_team_id && m.winner_team_id === m.team_b_id }">{{ m.team_b_name || 'TBD' }}</span>
                    </div>
                    <div class="bm-meta">
                      <span>{{ formatMatchTime(m.match_time) }}</span>
                      <span class="bm-status">{{ m.status }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="knockout-container">
            <div v-for="stage in overview.knockout_stages" :key="stage.stage_id" class="stage-block">
              <div class="stage-title">{{ stage.stage_name }}</div>
              <div v-if="!stage.matches || stage.matches.length === 0" class="empty-state">
                <p>暂无对阵</p>
              </div>
              <div v-else class="stage-matches">
                <div v-for="m in stage.matches" :key="m.knockout_match_id" class="knockout-match">
                  <div class="km-row">
                    <span class="km-team">{{ m.team_a_name || 'TBD' }}</span>
                    <span class="km-score">{{ m.status === 'not_started' ? 'VS' : `${m.score_team_a}-${m.score_team_b}` }}</span>
                    <span class="km-team">{{ m.team_b_name || 'TBD' }}</span>
                  </div>
                  <div class="km-meta">
                    <span>{{ formatMatchTime(m.match_time) }}</span>
                    <span class="km-status">{{ m.status }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <p>未知赛制</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.standings-container {
  min-height: 100vh;
  background-color: #f3f4f6; /* Lighter cool gray background */
  display: flex;
  flex-direction: column;
  font-family: 'Inter', system-ui, -apple-system, sans-serif;
}

.page-header {
  background: linear-gradient(135deg, #1e3a8a 0%, #2563eb 100%); /* Deep Royal Blue to Bright Blue */
  color: white;
  padding: 12px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.header-main {
  flex: 1;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
  gap: 32px;
}

.header-controls {
  display: flex;
  gap: 12px;
  width: auto;
  min-width: 420px;
}

.search-input {
  flex: 1;
  padding: 8px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px; /* More rounded */
  background: rgba(255, 255, 255, 0.15);
  color: white;
  outline: none;
  transition: all 0.2s;
  font-size: 14px;
}

.search-input:focus {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.5);
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.7);
}

.custom-select {
  position: relative;
  width: 200px;
  font-size: 14px;
  user-select: none;
}

.select-trigger {
  padding: 8px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.15);
  color: white;
  cursor: pointer;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s;
}

.select-trigger:hover {
  background: rgba(255, 255, 255, 0.25);
  border-color: rgba(255, 255, 255, 0.5);
}

.trigger-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow {
  font-size: 10px;
  transition: transform 0.2s;
  margin-left: 8px;
}

.arrow.rotated {
  transform: rotate(180deg);
}

.select-options {
  position: absolute;
  top: 110%;
  left: 0;
  width: 100%;
  background: white;
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 200;
  max-height: 300px;
  overflow-y: auto;
  padding: 8px 0;
  color: #1f2937;
  animation: slideDown 0.2s ease-out;
}

.select-option {
  padding: 10px 16px;
  cursor: pointer;
  transition: background 0.1s;
}

.select-option:hover {
  background-color: #f3f4f6;
}

.select-option.selected {
  background-color: #eff6ff;
  color: #2563eb;
  font-weight: 600;
}

.select-option.disabled {
  color: #9ca3af;
  cursor: default;
}

@keyframes slideDown {
  from { opacity: 0; transform: translateY(-10px); }
  to { opacity: 1; transform: translateY(0); }
}

.back-button {
  background: rgba(255, 255, 255, 0.1);
  border: none;
  color: white;
  font-size: 18px;
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s;
}

.back-button:hover {
  background: rgba(255, 255, 255, 0.25);
}

.page-header h1 {
  font-size: 24px;
  font-weight: 800;
  margin: 0;
  white-space: nowrap;
  letter-spacing: -0.5px;
  text-shadow: 0 2px 4px rgba(0,0,0,0.15);
  color: #fbbf24; /* Amber-400 */
}

.placeholder {
  width: 36px;
}

.cup-tabs, .knockout-view-tabs {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  justify-content: center;
}

.tab-btn, .mode-btn {
  padding: 8px 24px;
  border: none;
  border-radius: 20px;
  background: white;
  color: #6b7280;
  cursor: pointer;
  font-weight: 600;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
  transition: all 0.2s;
}

.tab-btn.active, .mode-btn.active {
  background: #2563eb;
  color: white;
  box-shadow: 0 4px 6px -1px rgba(37, 99, 235, 0.3);
}

.group-title, .stage-title {
  font-size: 18px;
  font-weight: 700;
  margin: 16px 0 12px;
  color: #111827;
  border-left: 4px solid #2563eb;
  padding-left: 12px;
}

.bracket-match, .knockout-match {
  background: white;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.02);
  transition: transform 0.2s, box-shadow 0.2s;
}

.bracket-match:hover, .knockout-match:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0,0,0,0.06);
}

.bm-team.winner, .km-team.winner {
  color: #2563eb;
  font-weight: 700;
}

.bm-score, .km-score {
  background: #eff6ff;
  color: #2563eb;
  padding: 2px 8px;
  border-radius: 6px;
}

.bracket-scroll {
  overflow-x: auto;
  padding-bottom: 12px;
}

.bracket {
  position: relative;
}

.bracket-col {
  position: absolute;
  top: 0;
}

.bracket-col-title {
  height: 28px;
  line-height: 28px;
  font-weight: 700;
  text-align: center;
  color: #4b5563;
  margin-bottom: 12px;
}

.bm-row, .km-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  font-weight: 500;
}

.bm-team, .km-team {
  flex: 1;
  min-width: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: #1f2937;
}

.bm-meta, .km-meta {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  color: #9ca3af;
  font-size: 12px;
  font-weight: 500;
}

.group-block, .stage-block {
  margin-bottom: 24px;
}

.content-area {
  padding: 24px;
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 60px 0;
  color: #9ca3af;
  font-size: 16px;
  font-weight: 500;
}

.table-container {
  overflow-x: auto;
  border-radius: 16px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  border: none;
  background: white;
}

.standings-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  min-width: 600px;
}

.standings-table th {
  background-color: #f8fafc;
  color: #64748b;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 12px;
  letter-spacing: 0.5px;
  padding: 16px;
  border-bottom: 2px solid #e2e8f0;
}

.standings-table td {
  padding: 16px;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
}

.standings-table tr:hover td {
  background-color: #f8fafc;
}

.rank-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #f1f5f9;
  color: #64748b;
  font-weight: 700;
  font-size: 13px;
}

.rank-1 {
  background: linear-gradient(135deg, #fbbf24 0%, #d97706 100%);
  color: white;
  box-shadow: 0 2px 4px rgba(217, 119, 6, 0.3);
}

.rank-2 {
  background: linear-gradient(135deg, #94a3b8 0%, #475569 100%);
  color: white;
  box-shadow: 0 2px 4px rgba(71, 85, 105, 0.3);
}

.rank-3 {
  background: linear-gradient(135deg, #fca5a5 0%, #b91c1c 100%); /* Bronze-ish Red */
  background: linear-gradient(135deg, #d6d3d1 0%, #a8a29e 100%); /* Actually bronze is brownish, let's try a copper */
  background: linear-gradient(135deg, #fdba74 0%, #c2410c 100%);
  color: white;
  box-shadow: 0 2px 4px rgba(194, 65, 12, 0.3);
}

.points-col {
  color: #2563eb;
  font-weight: 800;
  font-size: 16px;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .header-main {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
  
  .header-controls {
    min-width: 0;
    width: 100%;
    flex-direction: column;
  }

  .search-input, .custom-select {
    width: 100%;
  }
  
  .page-header {
    height: auto;
    padding: 12px;
  }

  .page-header h1 {
    font-size: 20px;
    text-align: center;
    margin-bottom: 8px;
  }

  .content-area {
    padding: 8px;
  }

  .standings-table {
    min-width: auto;
    font-size: 12px;
  }

  .standings-table th,
  .standings-table td {
    padding: 8px 2px;
    text-align: center;
  }

  .standings-table .team-col {
    text-align: left;
    max-width: 100px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    padding-left: 4px;
  }

  .rank-col {
    width: 30px;
  }

  .rank-badge {
    width: 22px;
    height: 22px;
    font-size: 11px;
    line-height: 22px;
  }

  .points-col {
    font-size: 14px;
  }
}

/* 连接线样式 */
.connector {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
}

.connector .line-v {
  position: absolute;
  left: 50%;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: #cbd5e1;
}

.connector .line-h-parent {
  position: absolute;
  left: 50%;
  right: 0;
  top: 50%;
  height: 2px;
  background-color: #cbd5e1;
  transform: translateY(-50%);
}

.connector .line-h-child-top {
  position: absolute;
  left: 0;
  width: 50%;
  top: 0;
  height: 2px;
  background-color: #cbd5e1;
}

.connector .line-h-child-bot {
  position: absolute;
  left: 0;
  width: 50%;
  bottom: 0;
  height: 2px;
  background-color: #cbd5e1;
}

/* 移动端金字塔/聚合布局 */
.bracket-pyramid {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  padding: 20px 0;
  background-color: #f8fafc;
}

.pyramid-section {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.upper-section {
  /* 上半区：从上到下 (R16 -> QF -> SF) */
}

.lower-section {
  /* 下半区：从上到下 (SF -> QF -> R16) */
}

.pyramid-row {
  display: flex;
  justify-content: center;
  width: 100%;
  margin-bottom: 24px; /* 行间距 */
  position: relative;
}

.row-matches {
  display: flex;
  justify-content: center; /* 居中，保持卡片间距固定以维持连接线 */
  width: 100%;
  padding: 0 10px;
}

/* 比赛卡片 */
.match-card {
  flex: 1;
  max-width: 160px; /* 卡片最大宽度 */
  background: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  padding: 8px;
  margin: 0 4px; /* 卡片间距 */
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
  min-height: 56px;
}

.mc-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.mc-team {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 11px;
}

.mc-team-name {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 80%;
  color: #64748b;
}

.mc-team-name.winner {
  font-weight: 700;
  color: #0f172a;
}

.mc-score {
  font-weight: 600;
  color: #334155;
}

/* 决赛卡片 */
.pyramid-final {
  margin: 20px 0;
  width: 100%;
  display: flex;
  justify-content: center;
  z-index: 10;
}

.final-card {
  background: white;
  border: 2px solid #fbbf24; /* 金色边框 */
  border-radius: 12px;
  padding: 16px;
  width: 80%;
  max-width: 300px;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.fc-header {
  font-size: 10px;
  font-weight: 700;
  color: #fbbf24;
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.fc-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.fc-team {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.fc-name {
  font-weight: 700;
  font-size: 14px;
  text-align: center;
}

.fc-score {
  display: flex;
  gap: 8px;
  font-size: 24px;
  font-weight: 800;
  color: #1e293b;
  margin: 0 16px;
}

.fc-score .divider {
  color: #cbd5e1;
}

/* 连接线逻辑 */

/* 上半区连接线：向下 */
/* 我们使用伪元素在 row-matches 上方或下方绘制线条 */
/* 简化方案：每个 match-card 如果不是最后一排，就在底部画一个"Fork"的一半 */

/* 
   更简单的方案：
   Match Card 底部中心伸出一根短线。
   下一排的 Match Card 顶部中心伸出一根长线，分叉去连接上面的两个。
   但是 Flexbox 布局下，下一排的一个卡片对应上一排的两个卡片。
   
   我们可以利用 .match-card::after 来画线。
   
   对于 Upper Section:
   Row i (matches)
   Row i+1 (matches)
   
   Row i 中的每两个 match 需要汇聚到 Row i+1 中的一个 match。
   
   我们在 Row i 的 match-card 底部画线：
   偶数索引(左): └─
   奇数索引(右): ─┘
   
   这需要 nth-child 选择器。
*/

.upper-section .match-card {
  margin-bottom: 16px; /* 给连接线留空间 */
}

/* 偶数个子元素（假设是 0,1 索引，即 1st, 2nd child） */
/* match-card:nth-child(odd) 是第 1, 3, 5... 个元素 (Left of pair) */
/* match-card:nth-child(even) 是第 2, 4, 6... 个元素 (Right of pair) */

/* 左侧卡片：向下，向右 */
.upper-section .match-card:nth-child(odd)::after {
  content: '';
  position: absolute;
  bottom: -16px; /* 延伸到下个空间 */
  left: 50%;
  width: 50%; /* 向右延伸 50% + 间距的一半 */
  /* 这里很难精确计算到下一个卡片的中心，除非间距固定 */
  /* 近似：向右延伸 55% */
  width: calc(50% + 6px);
  height: 16px;
  border-bottom: 2px solid #cbd5e1;
  border-left: 2px solid #cbd5e1;
  border-bottom-left-radius: 6px;
}

/* 右侧卡片：向下，向左 */
.upper-section .match-card:nth-child(even)::after {
  content: '';
  position: absolute;
  bottom: -16px;
  right: 50%;
  width: calc(50% + 6px);
  height: 16px;
  border-bottom: 2px solid #cbd5e1;
  border-right: 2px solid #cbd5e1;
  border-bottom-right-radius: 6px;
}

/* 修正：最后一排（Semi Final）不需要这些线，而是连接到决赛 */
.upper-section .pyramid-row:last-child .match-card::after {
  display: none;
}

/* 上半区最后一排连接到决赛 */
.upper-section .pyramid-row:last-child .match-card .connector-to-final-top {
  position: absolute;
  bottom: -24px;
  left: 50%;
  width: 2px;
  height: 24px;
  background-color: #cbd5e1;
  transform: translateX(-50%);
}


/* 下半区连接线：向上 */
.lower-section .match-card {
  margin-top: 16px;
}

/* 左侧卡片：向上，向右 */
.lower-section .match-card:nth-child(odd)::before {
  content: '';
  position: absolute;
  top: -16px;
  left: 50%;
  width: calc(50% + 6px);
  height: 16px;
  border-top: 2px solid #cbd5e1;
  border-left: 2px solid #cbd5e1;
  border-top-left-radius: 6px;
}

/* 右侧卡片：向上，向左 */
.lower-section .match-card:nth-child(even)::before {
  content: '';
  position: absolute;
  top: -16px;
  right: 50%;
  width: calc(50% + 6px);
  height: 16px;
  border-top: 2px solid #cbd5e1;
  border-right: 2px solid #cbd5e1;
  border-top-right-radius: 6px;
}

/* 修正：下半区第一排（Semi Final）不需要这些线，而是连接到决赛 */
.lower-section .pyramid-row:first-child .match-card::before {
  display: none;
}

/* 下半区第一排连接到决赛 */
.lower-section .pyramid-row:first-child .match-card .connector-to-final-bot {
  position: absolute;
  top: -24px;
  left: 50%;
  width: 2px;
  height: 24px;
  background-color: #cbd5e1;
  transform: translateX(-50%);
}

/* 下半区每排之间的垂直线？ */
/* 上面的伪元素已经画了横线和转角，缺中间的垂直线连接到下一层 */
/* 
   逻辑：
   Upper: Row i Matches (Pairs) -> Converge Line -> Vertical Line -> Row i+1 Match
   我们用伪元素画了 Converge Line (└─ ─┘)。
   还需要一根垂直线从 Converge Point 指向 Row i+1 的 Match。
   这可以在 Row i+1 的 Match 上画。
*/

.upper-section .pyramid-row:not(:first-child) .match-card::before {
  content: '';
  position: absolute;
  top: -16px; /* 延伸上去 */
  left: 50%;
  width: 2px;
  height: 16px;
  background-color: #cbd5e1;
  transform: translateX(-50%);
}

.lower-section .pyramid-row:not(:last-child) .match-card::after {
    content: '';
    position: absolute;
    bottom: -16px; /* 延伸下去 */
    left: 50%;
    width: 2px;
    height: 16px;
    background-color: #cbd5e1;
    transform: translateX(-50%);
  }
  
  /* 列表视图优化 */
  .stage-matches {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
  }
  
  .knockout-match {
    padding: 16px;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  
  @media (max-width: 768px) {
    .stage-matches {
      grid-template-columns: 1fr;
      gap: 12px;
    }
  
    .knockout-match {
      padding: 12px;
    }
  }
  
  </style>
