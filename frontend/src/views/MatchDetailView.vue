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
          <div class="team-logo">
            <img v-if="matchInfo.team_a_avatar" :src="matchInfo.team_a_avatar" class="team-logo-img" />
            <span v-else>{{ matchInfo.team_a_name?.charAt(0) || 'A' }}</span>
          </div>
          <span class="team-name">{{ matchInfo.team_a_name }}</span>
        </div>
        
        <div class="score-col">
          <div class="score">
            <span class="score-num">{{ matchInfo.score_a || 0 }}</span>
            <span class="colon">:</span>
            <span class="score-num">{{ matchInfo.score_b || 0 }}</span>
          </div>
        </div>

        <div class="team-col away">
          <div class="team-logo">
            <img v-if="matchInfo.team_b_avatar" :src="matchInfo.team_b_avatar" class="team-logo-img" />
            <span v-else>{{ matchInfo.team_b_name?.charAt(0) || 'B' }}</span>
          </div>
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
      
      <!-- 1. 数据统计 (Real) -->
      <div v-if="currentTab === 'stats'" class="stats-panel">
        <h3 class="panel-title">队伍概况 & 预测</h3>
        
        <!-- Standings Comparison -->
        <div class="standings-comparison" v-if="(teamAStats && teamAStats.played > 0) || (teamBStats && teamBStats.played > 0)">
            <div class="versus-container">
                <!-- Points -->
                <div class="versus-row">
                    <div class="versus-label">积分</div>
                    <div class="versus-bar-container">
                        <span class="val home">{{ teamAStats?.points || 0 }}</span>
                        <div class="bar-track">
                            <div class="bar-fill home" :style="{ width: getVersusWidth(teamAStats?.points, teamBStats?.points, 'home') }"></div>
                            <div class="bar-fill away" :style="{ width: getVersusWidth(teamAStats?.points, teamBStats?.points, 'away') }"></div>
                        </div>
                        <span class="val away">{{ teamBStats?.points || 0 }}</span>
                    </div>
                </div>

                <!-- Goals For -->
                <div class="versus-row">
                    <div class="versus-label">进球</div>
                    <div class="versus-bar-container">
                        <span class="val home">{{ teamAStats?.goals_for || 0 }}</span>
                        <div class="bar-track">
                            <div class="bar-fill home" :style="{ width: getVersusWidth(teamAStats?.goals_for, teamBStats?.goals_for, 'home') }"></div>
                            <div class="bar-fill away" :style="{ width: getVersusWidth(teamAStats?.goals_for, teamBStats?.goals_for, 'away') }"></div>
                        </div>
                        <span class="val away">{{ teamBStats?.goals_for || 0 }}</span>
                    </div>
                </div>

                <!-- Goals Against -->
                <div class="versus-row">
                    <div class="versus-label">失球</div>
                    <div class="versus-bar-container">
                        <span class="val home">{{ teamAStats?.goals_against || 0 }}</span>
                        <div class="bar-track">
                            <div class="bar-fill home" :style="{ width: getVersusWidth(teamAStats?.goals_against, teamBStats?.goals_against, 'home') }"></div>
                            <div class="bar-fill away" :style="{ width: getVersusWidth(teamAStats?.goals_against, teamBStats?.goals_against, 'away') }"></div>
                        </div>
                        <span class="val away">{{ teamBStats?.goals_against || 0 }}</span>
                    </div>
                </div>

                 <!-- Wins -->
                <div class="versus-row">
                    <div class="versus-label">胜场</div>
                    <div class="versus-bar-container">
                        <span class="val home">{{ teamAStats?.won || 0 }}</span>
                        <div class="bar-track">
                            <div class="bar-fill home" :style="{ width: getVersusWidth(teamAStats?.won, teamBStats?.won, 'home') }"></div>
                            <div class="bar-fill away" :style="{ width: getVersusWidth(teamAStats?.won, teamBStats?.won, 'away') }"></div>
                        </div>
                        <span class="val away">{{ teamBStats?.won || 0 }}</span>
                    </div>
                </div>

                <!-- Prediction Win Rate -->
                <div class="versus-row" v-if="predictionResult">
                    <div class="versus-label">预测指数 (胜/平/负概率)</div>
                    <div class="versus-bar-container prediction-bar-container">
                        <!-- Home Segment -->
                        <div class="pred-segment home" 
                             :style="{ width: predictionResult.probs.home + '%' }">
                             <span v-if="parseFloat(predictionResult.probs.home) > 10">{{ predictionResult.probs.home }}%</span>
                        </div>
                        <!-- Draw Segment -->
                        <div class="pred-segment draw" 
                             :style="{ width: predictionResult.probs.draw + '%' }">
                             <span v-if="parseFloat(predictionResult.probs.draw) > 10">{{ predictionResult.probs.draw }}%</span>
                        </div>
                        <!-- Away Segment -->
                        <div class="pred-segment away" 
                             :style="{ width: predictionResult.probs.away + '%' }">
                             <span v-if="parseFloat(predictionResult.probs.away) > 10">{{ predictionResult.probs.away }}%</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div v-else class="empty-hint">暂无积分榜数据</div>

        <!-- Prediction Details -->
        <div class="prediction-box" v-if="predictionResult">
            <h4 class="sub-title">AI 预测详情 (xGD模型)</h4>
            <div class="prediction-content">
                <div class="prediction-result-text" :class="getScoreClass(predictionResult.score)">
                   {{ predictionResult.outcome }}
                </div>

                <!-- Probability Breakdown -->
                <div class="prediction-probs-row">
                    <div class="prob-item home">
                        <span class="prob-label">主胜率</span>
                        <span class="prob-val">{{ predictionResult.probs.home }}%</span>
                    </div>
                    <div class="prob-item draw">
                        <span class="prob-label">平局率</span>
                        <span class="prob-val">{{ predictionResult.probs.draw }}%</span>
                    </div>
                    <div class="prob-item away">
                        <span class="prob-label">客胜率</span>
                        <span class="prob-val">{{ predictionResult.probs.away }}%</span>
                    </div>
                </div>

                <div class="prediction-details-grid">
                    <div class="detail-item">
                        <span class="detail-label">A对B得分势能</span>
                        <span class="detail-val">{{ predictionResult.details.attackAB }}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">B对A得分势能</span>
                        <span class="detail-val">{{ predictionResult.details.attackBA }}</span>
                    </div>
                    <div class="detail-item">
                        <span class="detail-label">综合实力差</span>
                        <span class="detail-val">{{ predictionResult.details.deltaStrength }}</span>
                    </div>
                </div>
            </div>
            <div class="algorithm-note">
                * 基于期望进球差(xGD)与Elo修正算法计算
            </div>
        </div>
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
            <img v-if="comment.avatar" :src="comment.avatar" class="avatar" alt="用户头像">
            <div v-else class="avatar">{{ comment.user.charAt(0) }}</div>
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
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get, post } from '@/utils/http'

const route = useRoute()
const router = useRouter()
const matchId = route.params.id

const matchInfo = ref({})
const stats = ref(null)
const comments = ref([])
const currentTab = ref('stats')
const newComment = ref('')

// Standings & Prediction
const standings = ref([])
const teamAStats = ref(null)
const teamBStats = ref(null)
const predictionResult = ref(null)

const tabs = computed(() => {
  const allTabs = [
    { id: 'stats', label: '分析' },
    { id: 'lineup', label: '阵容' },
    { id: 'comments', label: '讨论' }
  ]
  
  // 如果比赛已结束，移除“分析”标签
  if (matchInfo.value.status === 'finished' || matchInfo.value.status === 'completed') {
    return allTabs.filter(t => t.id !== 'stats')
  }
  return allTabs
})

// 监听 tabs 变化，如果当前 tab 被移除，则自动切换到第一个可用 tab
watch(tabs, (newTabs) => {
  if (newTabs.length > 0 && !newTabs.find(t => t.id === currentTab.value)) {
    currentTab.value = newTabs[0].id
  }
})

// Stats 展示映射
const mockStats = ref([])

const getStatWidth = (stat, side) => {
  const total = stat.home + stat.away
  if (total === 0) return '50%'
  return `${(stat[side] / total) * 100}%`
}

const getScoreClass = (score) => {
  const s = parseFloat(score)
  if (s > 0.3) return 'text-home'
  if (s < -0.3) return 'text-away'
  return 'text-draw'
}

const getVersusWidth = (valA, valB, side) => {
    const a = parseFloat(valA) || 0
    const b = parseFloat(valB) || 0
    const total = a + b
    if (total === 0) return '50%'
    
    // Minimal width to show color
    const pct = (side === 'home' ? a : b) / total * 100
    return `${pct}%`
}

const getPredictionClass = (score) => {
    const s = parseFloat(score)
    if (s > 0.3) return 'bg-home'
    if (s < -0.3) return 'bg-away'
    return 'bg-draw'
}

const getPredictionBarStyle = (score) => {
    // Score range assumed roughly -2 to +2 for visualization clamping
    // We map 0 to center (50%).
    // Max width is 50% (from center to edge).
    const s = parseFloat(score)
    const maxScore = 2.0 // Clamp at 2.0
    const clamped = Math.min(Math.max(s, -maxScore), maxScore)
    const widthPct = (Math.abs(clamped) / maxScore) * 50
    
    if (s > 0) {
        // Home advantage: Start at 50%, grow right
        return { left: '50%', width: `${widthPct}%` }
    } else {
        // Away advantage: End at 50%, grow left
        return { right: '50%', width: `${widthPct}%` }
    }
}

const calculatePrediction = () => {
    const sA = teamAStats.value
    const sB = teamBStats.value
    
    if (!sA || !sB) {
        predictionResult.value = null
        return
    }

    const getMetrics = (s) => {
        const N = s.played || 0
        if (N === 0) return { avgGS: 0, avgGC: 0, winRate: 0, pointEff: 0 }
        return {
            avgGS: s.goals_for / N,
            avgGC: s.goals_against / N,
            winRate: (s.won / N) + 0.5 * (s.drawn / N),
            pointEff: s.points / (3 * N)
        }
    }

    const mA = getMetrics(sA)
    const mB = getMetrics(sB)
    
    if (sA.played === 0 || sB.played === 0) {
        // Data insufficient for prediction
        predictionResult.value = null
        return
    }

    const attackAB = mA.avgGS * mB.avgGC
    const attackBA = mB.avgGS * mA.avgGC
    const deltaStrength = (mA.pointEff - mB.pointEff) + (mA.winRate - mB.winRate)
    
    const w1 = 0.6
    const w2 = 0.4
    const scoreFinal = w1 * (attackAB - attackBA) + w2 * deltaStrength
    
    let outcome = '平局'
    if (scoreFinal > 0.3) outcome = `${matchInfo.value.team_a_name} 胜`
    else if (scoreFinal < -0.3) outcome = `${matchInfo.value.team_b_name} 胜`
    
    // Calculate Probabilities using Softmax-like approach with Draw bias
    // Score > 0 favors Home, Score < 0 favors Away
    // We map score to logits: Home=Score, Away=-Score, Draw=Constant
    // Draw Constant tuned so that at Score=0.3, Draw% is still significant but dropping
    
    // Heuristic: 
    // Home Logit = score * 1.5
    // Away Logit = -score * 1.5
    // Draw Logit = 0.5 (Base bias)
    
    const factor = 1.2
    const homeLogit = Math.exp(scoreFinal * factor)
    const awayLogit = Math.exp(-scoreFinal * factor)
    const drawLogit = Math.exp(0.5) // Bias for draw
    
    const sum = homeLogit + awayLogit + drawLogit
    const pHome = (homeLogit / sum * 100).toFixed(1)
    const pAway = (awayLogit / sum * 100).toFixed(1)
    const pDraw = (drawLogit / sum * 100).toFixed(1)

    predictionResult.value = {
        score: scoreFinal.toFixed(3),
        outcome,
        probs: { home: pHome, draw: pDraw, away: pAway },
        details: {
            attackAB: attackAB.toFixed(3),
            attackBA: attackBA.toFixed(3),
            deltaStrength: deltaStrength.toFixed(3)
        }
    }
}

const fetchStandings = async () => {
  if (!matchInfo.value.event_id) return
  const res = await get(`/events/${matchInfo.value.event_id}/standings`)
  if (res.code === 200 && res.data) {
    standings.value = res.data
    // Ensure ID comparison is safe (string vs number)
    const idA = Number(matchInfo.value.team_a_id)
    const idB = Number(matchInfo.value.team_b_id)
    teamAStats.value = standings.value.find(s => Number(s.team_id) === idA)
    teamBStats.value = standings.value.find(s => Number(s.team_id) === idB)
    calculatePrediction()
  }
}

const mockLineup = ref({ home: [], away: [] })

const normalizePlayer = (x) => {
  const id = String(x.lineup_id ?? x.id ?? '').trim()
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
      avatar: c.avatar_url ? `http://localhost:8080${c.avatar_url}` : '',
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
      team_a_avatar: d.team_a?.avatar ? `http://localhost:8080${d.team_a.avatar}` : '',
      team_b_avatar: d.team_b?.avatar ? `http://localhost:8080${d.team_b.avatar}` : '',
      team_a_id: d.team_a?.team_id || d.team_a?.id || d.team_a_id,
      team_b_id: d.team_b?.team_id || d.team_b?.id || d.team_b_id,
      event_id: d.event_id,
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
  await fetchStandings()
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
  width: 80px;
  height: 80px;
  background: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: bold;
  color: #333;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
  margin-bottom: 12px;
  overflow: hidden; /* Ensure image stays within circle */
}

.team-logo-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
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
/* Stats & Prediction */
.standings-comparison {
  margin-bottom: 24px;
  overflow-x: auto;
}

.comparison-table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 12px;
  overflow: hidden;
}

.comparison-table th,
.comparison-table td {
  padding: 12px;
  text-align: center;
  border-bottom: 1px solid rgba(0,0,0,0.05);
  font-size: 14px;
}

.comparison-table th {
  background: rgba(0,0,0,0.02);
  font-weight: 600;
  color: var(--text-secondary);
}

.comparison-table td.highlight {
  font-weight: 700;
  color: var(--primary-color);
  background: rgba(56, 189, 248, 0.05); /* fallback if variable not set */
}

.prediction-box {
  background: white;
  border-radius: 16px;
  padding: 20px;
  box-shadow: 0 4px 12px rgba(0,0,0,0.03);
  border: 1px solid rgba(0,0,0,0.04);
}

.sub-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.sub-title::before {
  content: '';
  display: block;
  width: 4px;
  height: 16px;
  background: var(--accent-color, #ff9a9e);
  border-radius: 2px;
}

.prediction-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.prediction-row {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding-bottom: 20px;
  border-bottom: 1px dashed rgba(0,0,0,0.1);
}

.prediction-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
}

.prediction-item .label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.prediction-item .value.large {
  font-size: 32px;
  font-weight: 800;
  font-family: 'DIN Alternate', sans-serif;
}

.result-tag {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  background: #f5f5f5;
}

.text-home { color: #8e44ad; }
.text-away { color: #2980b9; }
.text-draw { color: #95a5a6; }

.result-tag.text-home { background: rgba(142, 68, 173, 0.1); }
.result-tag.text-away { background: rgba(41, 128, 185, 0.1); }

/* Versus Mode Styles */
.versus-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 10px 0;
}

.versus-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.versus-label {
  text-align: center;
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.versus-bar-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.val {
  width: 40px;
  font-weight: 700;
  font-size: 16px;
  text-align: center;
}

.val.home { color: #8e44ad; }
.val.away { color: #2980b9; }

.bar-track {
  flex: 1;
  height: 10px;
  background: #f0f0f0;
  border-radius: 5px;
  display: flex;
  overflow: hidden;
  position: relative;
}

.bar-fill {
  height: 100%;
  transition: width 0.5s ease;
}

.bar-fill.home { background: #8e44ad; }
.bar-fill.away { background: #2980b9; }

/* Prediction Bar Special Styles */
.prediction-bar-container {
  height: 20px;
  background: #f0f0f0;
  border-radius: 10px;
  overflow: hidden;
  display: flex;
  width: 100%;
}

.pred-segment {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 10px;
  font-weight: 700;
  transition: width 0.5s ease;
  overflow: hidden;
  white-space: nowrap;
}

.pred-segment.home { background: #8e44ad; }
.pred-segment.draw { background: #95a5a6; color: #fff; }
.pred-segment.away { background: #2980b9; }

.prediction-result-text {
  text-align: center;
  font-size: 18px;
  font-weight: 800;
  margin-bottom: 10px;
}

.prediction-probs-row {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(0,0,0,0.05);
}

.prob-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.prob-label { font-size: 12px; color: var(--text-secondary); }
.prob-val { font-weight: 700; font-size: 16px; }
.prob-item.home .prob-val { color: #8e44ad; }
.prob-item.away .prob-val { color: #2980b9; }
.prob-item.draw .prob-val { color: #95a5a6; }

.prediction-details-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

@media (max-width: 480px) {
  .match-detail-container { padding: 12px; }
  .prediction-result-text { font-size: 16px; }
  .prob-val { font-size: 14px; }
  .detail-val { font-size: 13px; }
  .versus-label { font-size: 12px; }
  .versus-container { gap: 12px; }
}

.detail-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  background: #f8f9fa;
  padding: 10px;
  border-radius: 8px;
}

.detail-label {
  font-size: 11px;
  color: var(--text-tertiary);
}

.detail-val {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.algorithm-note {
  margin-top: 16px;
  font-size: 11px;
  color: var(--text-tertiary);
  text-align: right;
  font-style: italic;
}

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
  object-fit: cover;
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
