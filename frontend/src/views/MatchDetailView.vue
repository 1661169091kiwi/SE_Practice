<template>
  <div class="match-detail-container">
    <!-- Header -->
    <header class="detail-header">
      <button class="back-btn" @click="goBack">
        <span class="icon">←</span>
      </button>
      <h1 class="header-title">赛事详情</h1>
      <button class="share-btn">
        <span class="icon">🔗</span>
      </button>
    </header>

    <!-- Match Score Board -->
    <div class="score-board glass-card">
      <div class="match-status-bar">
        <span class="status-tag" :class="matchInfo.status">
          {{ formatStatus(matchInfo.status) }}
        </span>
        <span class="match-time">{{ formatTime(matchInfo.match_time) }}</span>
      </div>

      <div class="teams-score">
        <div class="team-col home">
          <div class="team-logo">{{ matchInfo.team_a_name?.charAt(0) || 'A' }}</div>
          <span class="team-name">{{ matchInfo.team_a_name }}</span>
        </div>
        
        <div class="score-col">
          <div class="score">
            <span class="score-num">{{ matchInfo.score_a || 0 }}</span>
            <span class="colon">:</span>
            <span class="score-num">{{ matchInfo.score_b || 0 }}</span>
          </div>
          <div class="venue">{{ matchInfo.venue || '校体育馆' }}</div>
        </div>

        <div class="team-col away">
          <div class="team-logo">{{ matchInfo.team_b_name?.charAt(0) || 'B' }}</div>
          <span class="team-name">{{ matchInfo.team_b_name }}</span>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs-nav">
      <button 
        v-for="tab in tabs" 
        :key="tab.id"
        class="tab-item"
        :class="{ active: currentTab === tab.id }"
        @click="currentTab = tab.id"
      >
        {{ tab.label }}
      </button>
    </div>

    <!-- Content Area -->
    <div class="tab-content glass-card fade-in">
      
      <!-- 1. 数据统计 (Mock) -->
      <div v-if="currentTab === 'stats'" class="stats-panel">
        <h3 class="panel-title">技术统计</h3>
        <div v-for="(stat, index) in mockStats" :key="index" class="stat-row">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-bars">
            <span class="stat-val home">{{ stat.home }}</span>
            <div class="progress-bar">
              <div class="fill home" :style="{ width: getStatWidth(stat, 'home') }"></div>
              <div class="fill away" :style="{ width: getStatWidth(stat, 'away') }"></div>
            </div>
            <span class="stat-val away">{{ stat.away }}</span>
          </div>
        </div>
        <div class="empty-hint" v-if="mockStats.length === 0">暂无详细数据</div>
      </div>

      <!-- 2. 阵容 (Mock) -->
      <div v-if="currentTab === 'lineup'" class="lineup-panel">
        <h3 class="panel-title">首发阵容</h3>
        <div class="lineup-cols">
          <div class="lineup-list home">
            <h4 class="team-title">{{ matchInfo.team_a_name }}</h4>
            <table class="lineup-table">
              <thead>
                <tr>
                  <th>球员姓名</th>
                  <th>球衣号</th>
                  <th>位置</th>
                  <th>首发</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="p in mockLineup.home" :key="p.id">
                  <td class="name">{{ p.name || '-' }}</td>
                  <td class="number">{{ p.number || '-' }}</td>
                  <td class="pos">{{ p.pos || '-' }}</td>
                  <td class="starter">{{ p.isStarting ? '是' : '否' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="lineup-list away">
            <h4 class="team-title">{{ matchInfo.team_b_name }}</h4>
            <table class="lineup-table">
              <thead>
                <tr>
                  <th>球员姓名</th>
                  <th>球衣号</th>
                  <th>位置</th>
                  <th>首发</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="p in mockLineup.away" :key="p.id">
                  <td class="name">{{ p.name || '-' }}</td>
                  <td class="number">{{ p.number || '-' }}</td>
                  <td class="pos">{{ p.pos || '-' }}</td>
                  <td class="starter">{{ p.isStarting ? '是' : '否' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <!-- 3. 评论区 (Mock) -->
      <div v-if="currentTab === 'comments'" class="comments-panel">
        <div class="comment-input-area">
          <input 
            v-model="newComment" 
            type="text" 
            placeholder="发表你的看法..." 
            @keyup.enter="postComment"
          >
          <button class="send-btn" @click="postComment" :disabled="!newComment.trim()">
            发送
          </button>
        </div>
        
        <div class="comments-list">
          <div v-for="comment in comments" :key="comment.id" class="comment-item">
            <div class="avatar">{{ comment.user.charAt(0) }}</div>
            <div class="comment-content">
              <div class="comment-header">
                <span class="username">{{ comment.user }}</span>
                <span class="time">{{ comment.time }}</span>
              </div>
              <p class="text">{{ comment.text }}</p>
            </div>
          </div>
          <div v-if="comments.length === 0" class="empty-hint">
            还没有人评论，快来抢沙发！
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post } from '@/utils/http'

const route = useRoute()
const router = useRouter()
const matchId = route.params.id

const matchInfo = ref({})
const stats = ref(null)
const lineup = ref({ home: [], away: [] })
const comments = ref([])
const currentTab = ref('stats')
const newComment = ref('')

const tabs = [
  { id: 'stats', label: '数据' },
  { id: 'lineup', label: '阵容' },
  { id: 'comments', label: '讨论' }
]

// Stats 展示映射
const mockStats = ref([])

const getStatWidth = (stat, side) => {
  const total = stat.home + stat.away
  if (total === 0) return '50%'
  return `${(stat[side] / total) * 100}%`
}

const mockLineup = ref({ home: [], away: [] })

const normalizePlayer = (x) => {
  const id = String(x.lineup_id ?? x.id ?? '').trim()
  const studentId = String(x.student_id ?? '').trim()
  const name = String(x.player_name ?? x.name ?? '').trim()
  const number = String(x.jersey_number ?? x.number ?? '').trim()
  const pos = String(x.position ?? x.pos ?? '').trim()
  const isStarting = !!(x.is_starting ?? x.isStarting)
  return { id, name, number, pos, isStarting, teamId: x.team_id }
}

const dedupByKey = (arr) => {
  const out = []
  const seenId = new Set()
  const seenName = new Set()
  // 先按 student_id 去重并记录姓名
  for (const x of arr) {
    const sid = String(x.student_id ?? '').trim()
    const nameLower = (x.player_name?.trim()?.toLowerCase() || x.name?.trim()?.toLowerCase() || '').trim()
    if (sid) {
      const key = `id:${sid}`
      if (seenId.has(key)) continue
      seenId.add(key)
      if (nameLower) seenName.add(nameLower)
      out.push(x)
    }
  }
  // 再处理无 student_id 的外部条目：如果姓名已出现则跳过；否则按 name+number+position 去重
  const seenExt = new Set()
  for (const x of arr) {
    const sid = String(x.student_id ?? '').trim()
    if (sid) continue
    const nameLower = (x.player_name?.trim()?.toLowerCase() || x.name?.trim()?.toLowerCase() || '').trim()
    if (nameLower && seenName.has(nameLower)) continue
    const extKey = `ext:${nameLower}|${String(x.jersey_number ?? x.number ?? '').trim()}|${String(x.position ?? x.pos ?? '').trim().toLowerCase()}`
    if (seenExt.has(extKey)) continue
    seenExt.add(extKey)
    if (nameLower) seenName.add(nameLower)
    out.push(x)
  }
  return out
}

const fetchComments = async () => {
  const res = await get(`/matches/${matchId}/comments`)
  if (res.code === 200) {
    comments.value = (res.data || []).map(c => ({
      id: c.id,
      user: c.user_name || c.student_id,
      text: c.content,
      time: new Date(c.created_at).toLocaleString('zh-CN')
    }))
  }
}

const fetchMatchDetail = async () => {
  const res = await get(`/matches/${matchId}`)
  if (res.code === 200 && res.data) {
    // 兼容不同字段命名
    const d = res.data
    matchInfo.value = {
      team_a_name: d.team_a?.name || d.team_a?.team_name || d.team_a_name,
      team_b_name: d.team_b?.name || d.team_b?.team_name || d.team_b_name,
      team_a_id: d.team_a?.id || d.team_a_id,
      team_b_id: d.team_b?.id || d.team_b_id,
      score_a: d.score_team_a ?? d.score_a ?? 0,
      score_b: d.score_team_b ?? d.score_b ?? 0,
      status: d.status,
      match_time: d.match_time
    }
  }
}

const fetchStats = async () => {
  const res = await get(`/matches/${matchId}/stats`)
  if (res.code === 200 && res.data) {
    stats.value = res.data
    mockStats.value = [
      { label: '控球率', home: res.data.home_possession, away: res.data.away_possession, max: 100 },
      { label: '射门', home: res.data.home_shots, away: res.data.away_shots, max: 20 },
      { label: '射正', home: res.data.home_shots_on_target, away: res.data.away_shots_on_target, max: 10 },
      { label: '犯规', home: res.data.home_fouls, away: res.data.away_fouls, max: 20 },
      { label: '角球', home: res.data.home_corners, away: res.data.away_corners, max: 10 }
    ]
  } else {
    mockStats.value = []
  }
}

const fetchLineup = async () => {
  const res = await get(`/matches/${matchId}/lineups`)
  if (res.code === 200 && res.data) {
    // 支持两种返回格式：数组(list) 或 { home, away }
    if (Array.isArray(res.data)) {
      const raw = dedupByKey(res.data || [])
      const homeRaw = matchInfo.value.team_a_id
        ? raw.filter(x => x.team_id === matchInfo.value.team_a_id)
        : raw.filter(x => x.is_starting)
      const awayRaw = matchInfo.value.team_b_id
        ? raw.filter(x => x.team_id === matchInfo.value.team_b_id)
        : []
      mockLineup.value.home = homeRaw.map(normalizePlayer)
      mockLineup.value.away = awayRaw.map(normalizePlayer)
    } else {
      const homeRaw = dedupByKey(res.data.home || [])
      const awayRaw = dedupByKey(res.data.away || [])
      mockLineup.value.home = homeRaw.map(normalizePlayer)
      mockLineup.value.away = awayRaw.map(normalizePlayer)
    }
  }
}

const postComment = async () => {
  if (!newComment.value.trim()) return
  const res = await post(`/matches/${matchId}/comments`, { content: newComment.value })
  if (res.code === 200) {
    await fetchComments()
    newComment.value = ''
  }
}

const goBack = () => {
  router.back()
}

const formatStatus = (status) => {
  const map = {
    'not_started': '未开始',
    'in_progress': '进行中',
    'finished': '已结束'
  }
  return map[status] || status
}

const formatTime = (timeStr) => {
  if (!timeStr) return ''
  return new Date(timeStr).toLocaleString('zh-CN', {
    month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}

onMounted(async () => {
  await fetchMatchDetail()
  await fetchStats()
  await fetchLineup()
  await fetchComments()
})
</script>

<style scoped>
.match-detail-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
  padding-bottom: 40px;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.back-btn, .share-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: none;
  background: white;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.header-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.glass-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  border-radius: 20px;
  box-shadow: 0 8px 32px 0 rgba(31, 38, 135, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.4);
}

.score-board {
  padding: 24px;
  margin-bottom: 24px;
  text-align: center;
}

.match-status-bar {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  margin-bottom: 20px;
  font-size: 12px;
}

.status-tag {
  padding: 2px 8px;
  border-radius: 4px;
  background: #eee;
  color: #666;
}
.status-tag.in_progress { background: #e6ffed; color: #52c41a; }
.status-tag.ongoing { background: #e6ffed; color: #52c41a; }

.teams-score {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.team-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.team-logo {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: linear-gradient(135deg, #e0eafc 0%, #cfdef3 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: var(--primary-color);
}

.team-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.score-col {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.score {
  font-size: 36px;
  font-weight: 800;
  color: var(--text-primary);
  letter-spacing: 2px;
}

.venue {
  font-size: 12px;
  color: var(--text-secondary);
  margin-top: 4px;
}

/* Tabs */
.tabs-nav {
  display: flex;
  justify-content: center;
  gap: 30px;
  margin-bottom: 20px;
}

.tab-item {
  background: none;
  border: none;
  font-size: 16px;
  color: var(--text-secondary);
  padding-bottom: 8px;
  cursor: pointer;
  position: relative;
  transition: all 0.3s;
}

.tab-item.active {
  color: var(--primary-color);
  font-weight: 600;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: var(--primary-color);
  border-radius: 2px;
}

/* Tab Content */
.tab-content {
  padding: 20px;
  min-height: 300px;
}

.panel-title {
  font-size: 16px;
  margin-bottom: 16px;
  color: var(--text-primary);
}

/* Lineup */
.lineup-cols {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}
.team-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 10px;
}
.lineup-table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(255, 255, 255, 0.85);
  border-radius: 12px;
  overflow: hidden;
}
.lineup-table th,
.lineup-table td {
  padding: 10px 12px;
  border-bottom: 1px solid rgba(0,0,0,0.06);
  font-size: 13px;
  color: var(--text-secondary);
}
.lineup-table thead th {
  text-align: left;
  background: rgba(0,0,0,0.03);
  color: var(--text-primary);
}
.lineup-table tbody tr:last-child td {
  border-bottom: none;
}
.lineup-table td.name {
  font-weight: 600;
  color: var(--text-primary);
}
.lineup-table td.number {
  width: 80px;
  text-align: center;
}
.lineup-table td.pos {
  width: 120px;
}
.lineup-table td.starter {
  width: 80px;
  text-align: center;
}

@media (max-width: 768px) {
  .lineup-cols {
    grid-template-columns: 1fr;
  }
}
/* Stats */
.stat-row {
  margin-bottom: 16px;
}

.stat-label {
  text-align: center;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 4px;
}

.stat-bars {
  display: flex;
  align-items: center;
  gap: 10px;
}

.stat-val {
  width: 30px;
  font-size: 14px;
  font-weight: 600;
  text-align: center;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: #f0f0f0;
  border-radius: 4px;
  display: flex;
  overflow: hidden;
}

.fill {
  height: 100%;
}
.fill.home { background: var(--primary-color); }
.fill.away { background: var(--accent-color, #ff9a9e); }

/* Comments */
.comment-input-area {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.comment-input-area input {
  flex: 1;
  padding: 10px 16px;
  border-radius: 20px;
  border: 1px solid rgba(0,0,0,0.1);
  background: rgba(255,255,255,0.8);
}

.send-btn {
  padding: 0 20px;
  border-radius: 20px;
  background: var(--primary-color);
  color: white;
  border: none;
  cursor: pointer;
}

.send-btn:disabled {
  background: #ccc;
}

.comment-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #eee;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #666;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 4px;
  font-size: 12px;
}

.username { font-weight: 600; color: var(--text-primary); }
.time { color: var(--text-tertiary); }
.text { font-size: 14px; color: var(--text-secondary); line-height: 1.4; }

.empty-hint {
  text-align: center;
  color: var(--text-tertiary);
  padding: 20px;
  font-size: 14px;
}
</style>
