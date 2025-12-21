<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MatchHeader from '../components/MatchHeader.vue'

const route = useRoute()
const router = useRouter()

// 赛事信息
const eventInfo = ref({
  id: '',
  name: '赛事数据预览',
  teamA: '',
  teamB: '',
  status: 'in_progress',
  collectorName: '采集员姓名'
})

// 运动类型
const sportType = ref('football')

// 加载状态
const isLoading = ref(true)

// 足球预览数据
const footballData = ref({
  matchInfo: {
    homeTeamScore: 2,
    awayTeamScore: 1,
    status: '已结束'
  },
  lineups: {
    home: [
      { name: '张三', number: 1, position: '门将', isStarter: true },
      { name: '李四', number: 4, position: '后卫', isStarter: true },
      { name: '王五', number: 5, position: '后卫', isStarter: true },
      { name: '赵六', number: 6, position: '中场', isStarter: true },
      { name: '钱七', number: 7, position: '中场', isStarter: true },
      { name: '孙八', number: 8, position: '中场', isStarter: true },
      { name: '周九', number: 9, position: '前锋', isStarter: true },
      { name: '吴十', number: 10, position: '前锋', isStarter: true },
      { name: '郑十一', number: 11, position: '前锋', isStarter: true },
      { name: '王十二', number: 14, position: '后卫', isStarter: true },
      { name: '李十三', number: 16, position: '后卫', isStarter: true }
    ],
    away: [
      { name: '陈一', number: 1, position: '门将', isStarter: true },
      { name: '陈二', number: 2, position: '后卫', isStarter: true },
      { name: '陈三', number: 3, position: '后卫', isStarter: true },
      { name: '陈四', number: 4, position: '中场', isStarter: true },
      { name: '陈五', number: 5, position: '中场', isStarter: true },
      { name: '陈六', number: 6, position: '中场', isStarter: true },
      { name: '陈七', number: 7, position: '前锋', isStarter: true },
      { name: '陈八', number: 8, position: '前锋', isStarter: true },
      { name: '陈九', number: 9, position: '前锋', isStarter: true },
      { name: '陈十', number: 14, position: '后卫', isStarter: true },
      { name: '陈十一', number: 16, position: '后卫', isStarter: true }
    ]
  },
  events: [
    { time: '15', type: '进球', player: '周九', team: '主队', description: '左脚射门' },
    { time: '35', type: '黄牌', player: '陈四', team: '客队', description: '铲球犯规' },
    { time: '58', type: '进球', player: '陈七', team: '客队', description: '头球破门' },
    { time: '75', type: '进球', player: '吴十', team: '主队', description: '点球命中' },
    { time: '85', type: '红牌', player: '李四', team: '主队', description: '严重犯规' }
  ],
  rankings: [
    { team: '主队', matches: 10, wins: 7, draws: 2, losses: 1, goalsFor: 21, goalsAgainst: 8, points: 23 },
    { team: '客队', matches: 10, wins: 6, draws: 1, losses: 3, goalsFor: 18, goalsAgainst: 12, points: 19 },
    { team: '其他队A', matches: 10, wins: 4, draws: 3, losses: 3, goalsFor: 15, goalsAgainst: 15, points: 15 },
    { team: '其他队B', matches: 10, wins: 2, draws: 4, losses: 4, goalsFor: 10, goalsAgainst: 17, points: 10 },
    { team: '其他队C', matches: 10, wins: 1, draws: 2, losses: 7, goalsFor: 8, goalsAgainst: 22, points: 5 }
  ]
})

// 篮球预览数据
const basketballData = ref({
  matchInfo: {
    homeTeamScore: 88,
    awayTeamScore: 82,
    homeQuarterScores: [22, 20, 26, 20],
    awayQuarterScores: [18, 24, 22, 18]
  },
  lineups: {
    home: [
      { name: '张三', number: 1, position: '控球后卫', isStarter: true },
      { name: '李四', number: 2, position: '得分后卫', isStarter: true },
      { name: '王五', number: 3, position: '小前锋', isStarter: true },
      { name: '赵六', number: 4, position: '大前锋', isStarter: true },
      { name: '钱七', number: 5, position: '中锋', isStarter: true }
    ],
    away: [
      { name: '陈一', number: 1, position: '控球后卫', isStarter: true },
      { name: '陈二', number: 2, position: '得分后卫', isStarter: true },
      { name: '陈三', number: 3, position: '小前锋', isStarter: true },
      { name: '陈四', number: 4, position: '大前锋', isStarter: true },
      { name: '陈五', number: 5, position: '中锋', isStarter: true }
    ]
  },
  events: [
    { quarter: 1, time: '05:30', type: '得分', player: '张三', team: '主队', scoreType: '2分球' },
    { quarter: 1, time: '04:15', type: '得分', player: '陈一', team: '客队', scoreType: '3分球' },
    { quarter: 2, time: '07:20', type: '犯规', player: '李四', team: '主队', description: '打手犯规' },
    { quarter: 2, time: '02:45', type: '暂停', team: '客队', description: '常规暂停' },
    { quarter: 3, time: '06:10', type: '换人', playerIn: '赵六', playerOut: '钱七', team: '主队' },
    { quarter: 4, time: '01:30', type: '得分', player: '王五', team: '主队', scoreType: '2分球' }
  ],
  teamStats: [
    { team: '主队', matches: 12, wins: 9, losses: 3, winRate: 75, behind: 0, recent: 'W-W-W-W-L' },
    { team: '客队', matches: 12, wins: 7, losses: 5, winRate: 58.3, behind: 2, recent: 'W-W-L-L-W' },
    { team: '其他队A', matches: 12, wins: 6, losses: 6, winRate: 50, behind: 3, recent: 'L-W-W-L-L' },
    { team: '其他队B', matches: 12, wins: 5, losses: 7, winRate: 41.7, behind: 4, recent: 'L-L-W-W-L' },
    { team: '其他队C', matches: 12, wins: 4, losses: 8, winRate: 33.3, behind: 5, recent: 'L-L-L-W-W' }
  ]
})

// 羽毛球预览数据
const badmintonData = ref({
  matchInfo: {
    homeTeamScore: 2,
    awayTeamScore: 0,
    gameScores: [
      { home: 21, away: 18 },
      { home: 21, away: 16 }
    ],
    format: '单打'
  },
  players: {
    home: [
      { name: '张三', position: '单打选手', isStarter: true },
      { name: '李四', position: '替补选手', isStarter: false }
    ],
    away: [
      { name: '陈一', position: '单打选手', isStarter: true },
      { name: '陈二', position: '替补选手', isStarter: false }
    ]
  },
  events: [
    { game: 1, point: 5, type: '得分', player: '张三', team: '主队' },
    { game: 1, point: 10, type: '失误', player: '陈一', team: '客队' },
    { game: 1, point: 15, type: '暂停', team: '主队' },
    { game: 2, point: 8, type: '得分', player: '张三', team: '主队' },
    { game: 2, point: 12, type: '失误', player: '张三', team: '主队' }
  ],
  teamStats: [
    { team: '主队', matches: 8, wins: 6, draws: 1, losses: 1, winRate: 75, points: 19 },
    { team: '客队', matches: 8, wins: 5, draws: 0, losses: 3, winRate: 62.5, points: 15 },
    { team: '其他队A', matches: 8, wins: 4, draws: 1, losses: 3, winRate: 50, points: 13 },
    { team: '其他队B', matches: 8, wins: 3, draws: 1, losses: 4, winRate: 37.5, points: 10 },
    { team: '其他队C', matches: 8, wins: 1, draws: 1, losses: 6, winRate: 12.5, points: 4 }
  ]
})

// 排球预览数据
const volleyballData = ref({
  matchInfo: {
    homeTeamScore: 3,
    awayTeamScore: 1,
    gameScores: [
      { home: 25, away: 20 },
      { home: 23, away: 25 },
      { home: 25, away: 18 },
      { home: 25, away: 22 }
    ]
  },
  lineups: {
    home: [
      { name: '张三', number: 1, position: '主攻', isStarter: true },
      { name: '李四', number: 2, position: '副攻', isStarter: true },
      { name: '王五', number: 3, position: '二传', isStarter: true },
      { name: '赵六', number: 4, position: '接应', isStarter: true },
      { name: '钱七', number: 5, position: '副攻', isStarter: true },
      { name: '孙八', number: 6, position: '自由人', isStarter: true }
    ],
    away: [
      { name: '陈一', number: 1, position: '主攻', isStarter: true },
      { name: '陈二', number: 2, position: '副攻', isStarter: true },
      { name: '陈三', number: 3, position: '二传', isStarter: true },
      { name: '陈四', number: 4, position: '接应', isStarter: true },
      { name: '陈五', number: 5, position: '副攻', isStarter: true },
      { name: '陈六', number: 6, position: '自由人', isStarter: true }
    ]
  },
  events: [
    { game: 1, point: 5, type: '得分', player: '张三', team: '主队', skill: '扣球' },
    { game: 1, point: 10, type: '犯规', player: '陈二', team: '客队', description: '触网' },
    { game: 1, point: 15, type: '换人', playerIn: '王五', playerOut: '李四', team: '主队' },
    { game: 2, point: 8, type: '暂停', team: '客队' },
    { game: 3, point: 12, type: '得分', player: '赵六', team: '主队', skill: '发球直接得分' }
  ],
  teamStats: [
    { team: '主队', matches: 6, wins: 5, losses: 1, points: 15 },
    { team: '客队', matches: 6, wins: 4, losses: 2, points: 12 },
    { team: '其他队A', matches: 6, wins: 3, losses: 3, points: 9 },
    { team: '其他队B', matches: 6, wins: 2, losses: 4, points: 6 },
    { team: '其他队C', matches: 6, wins: 0, losses: 6, points: 0 }
  ]
})

// 水上运动预览数据
const waterSportsData = ref({
  category: '龙舟',
  eventData: [
    { rank: 1, team: '计算机学院', peopleCount: 20, eventName: '500米竞速', eventTime: '2023-11-10 10:00', finalResult: 118.5 },
    { rank: 2, team: '机械学院', peopleCount: 20, eventName: '500米竞速', eventTime: '2023-11-10 10:00', finalResult: 122.3 },
    { rank: 3, team: '商学院', peopleCount: 20, eventName: '500米竞速', eventTime: '2023-11-10 10:00', finalResult: 125.1 }
  ],
  segmentResults: [
    {
      team: '计算机学院',
      segments: [
        { name: '起航段', time: 30.2 },
        { name: '中途段', time: 60.8 },
        { name: '冲刺段', time: 27.5 },
        { total: 118.5 }
      ]
    },
    {
      team: '机械学院',
      segments: [
        { name: '起航段', time: 32.4 },
        { name: '中途段', time: 64.1 },
        { name: '冲刺段', time: 25.8 },
        { total: 122.3 }
      ]
    },
    {
      team: '商学院',
      segments: [
        { name: '起航段', time: 33.5 },
        { name: '中途段', time: 65.7 },
        { name: '冲刺段', time: 25.9 },
        { total: 125.1 }
      ]
    }
  ]
})

// 加载预览数据
const loadPreviewData = () => {
  isLoading.value = true
  setTimeout(() => {
    // 从路由参数获取运动类型和赛事ID
    const type = route.params.sportType || 'football'
    const eventId = route.params.eventId || '1'
    
    sportType.value = type
    
    // 设置赛事信息
    eventInfo.value = {
      id: eventId,
      name: `赛事数据预览 - ${getSportTypeName(type)}`,
      teamA: '主队',
      teamB: '客队',
      status: '已结束',
      collectorName: '采集员姓名'
    }
    
    isLoading.value = false
  }, 1000)
}

// 获取运动类型名称
const getSportTypeName = (type) => {
  const typeMap = {
    football: '足球',
    basketball: '篮球',
    badminton: '羽毛球',
    volleyball: '排球',
    water_sports: '水上运动'
  }
  return typeMap[type] || type
}

// 打印预览
const handlePrint = () => {
  window.print()
}

// 导出数据
const handleExport = () => {
  alert('数据导出功能将在后续实现')
}

// 返回上一页
const handleBack = () => {
  router.back()
}

// 生命周期钩子
onMounted(() => {
  loadPreviewData()
})
</script>

<template>
  <div class="data-preview">
    <!-- 赛事信息头部 -->
    <MatchHeader :event="eventInfo" />
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 预览内容 -->
    <div v-else class="preview-content">
      <div class="preview-header">
        <h1 class="preview-title">赛事数据预览</h1>
        <p class="preview-subtitle">{{ eventInfo.name }}</p>
        <p class="preview-date">生成时间: {{ new Date().toLocaleString() }}</p>
      </div>
      
      <!-- 操作按钮 -->
      <div class="action-buttons">
        <button @click="handlePrint" class="action-btn print-btn">
          打印预览
        </button>
        <button @click="handleExport" class="action-btn export-btn">
          导出数据
        </button>
        <button @click="handleBack" class="action-btn back-btn">
          返回上一页
        </button>
      </div>
      
      <!-- 足球数据预览 -->
      <div v-if="sportType === 'football'" class="sport-preview">
        <div class="preview-section">
          <h2 class="section-title">比赛结果</h2>
          <div class="score-display">
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamA }}</span>
              <span class="score">{{ footballData.matchInfo.homeTeamScore }}</span>
            </div>
            <span class="score-separator">-</span>
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamB }}</span>
              <span class="score">{{ footballData.matchInfo.awayTeamScore }}</span>
            </div>
            <span class="match-status">{{ footballData.matchInfo.status }}</span>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">首发阵容</h2>
          
          <div class="lineups-container">
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamA }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in footballData.lineups.home" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamB }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in footballData.lineups.away" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">比赛事件</h2>
          <table class="events-table">
            <thead>
              <tr>
                <th>时间</th>
                <th>事件类型</th>
                <th>球员</th>
                <th>队伍</th>
                <th>描述</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(event, index) in footballData.events" :key="index">
                <td>{{ event.time }}'</td>
                <td>{{ event.type }}</td>
                <td>{{ event.player }}</td>
                <td>{{ event.team }}</td>
                <td>{{ event.description || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">积分排行榜</h2>
          <table class="rankings-table">
            <thead>
              <tr>
                <th>排名</th>
                <th>球队</th>
                <th>场次</th>
                <th>胜</th>
                <th>平</th>
                <th>负</th>
                <th>进球</th>
                <th>失球</th>
                <th>积分</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in footballData.rankings" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ team.team }}</td>
                <td>{{ team.matches }}</td>
                <td>{{ team.wins }}</td>
                <td>{{ team.draws }}</td>
                <td>{{ team.losses }}</td>
                <td>{{ team.goalsFor }}</td>
                <td>{{ team.goalsAgainst }}</td>
                <td>{{ team.points }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- 篮球数据预览 -->
      <div v-else-if="sportType === 'basketball'" class="sport-preview">
        <div class="preview-section">
          <h2 class="section-title">比赛结果</h2>
          <div class="score-display">
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamA }}</span>
              <span class="score">{{ basketballData.matchInfo.homeTeamScore }}</span>
            </div>
            <span class="score-separator">-</span>
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamB }}</span>
              <span class="score">{{ basketballData.matchInfo.awayTeamScore }}</span>
            </div>
          </div>
          
          <div class="quarter-scores">
            <div class="quarter-score">
              <span class="quarter-label">第1节</span>
              <span class="score">{{ basketballData.matchInfo.homeQuarterScores[0] }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ basketballData.matchInfo.awayQuarterScores[0] }}</span>
            </div>
            <div class="quarter-score">
              <span class="quarter-label">第2节</span>
              <span class="score">{{ basketballData.matchInfo.homeQuarterScores[1] }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ basketballData.matchInfo.awayQuarterScores[1] }}</span>
            </div>
            <div class="quarter-score">
              <span class="quarter-label">第3节</span>
              <span class="score">{{ basketballData.matchInfo.homeQuarterScores[2] }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ basketballData.matchInfo.awayQuarterScores[2] }}</span>
            </div>
            <div class="quarter-score">
              <span class="quarter-label">第4节</span>
              <span class="score">{{ basketballData.matchInfo.homeQuarterScores[3] }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ basketballData.matchInfo.awayQuarterScores[3] }}</span>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">首发阵容</h2>
          
          <div class="lineups-container">
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamA }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in basketballData.lineups.home" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamB }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in basketballData.lineups.away" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">比赛事件</h2>
          <table class="events-table">
            <thead>
              <tr>
                <th>节次</th>
                <th>时间</th>
                <th>事件类型</th>
                <th>球员</th>
                <th>队伍</th>
                <th>描述</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(event, index) in basketballData.events" :key="index">
                <td>第{{ event.quarter }}节</td>
                <td>{{ event.time }}</td>
                <td>{{ event.type }}</td>
                <td>{{ event.player || '-' }}</td>
                <td>{{ event.team || '-' }}</td>
                <td>{{ event.scoreType || event.description || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">球队数据统计</h2>
          <table class="rankings-table">
            <thead>
              <tr>
                <th>排名</th>
                <th>球队</th>
                <th>场次</th>
                <th>胜</th>
                <th>负</th>
                <th>胜率(%)</th>
                <th>胜差</th>
                <th>近5场</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in basketballData.teamStats" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ team.team }}</td>
                <td>{{ team.matches }}</td>
                <td>{{ team.wins }}</td>
                <td>{{ team.losses }}</td>
                <td>{{ team.winRate }}</td>
                <td>{{ team.behind }}</td>
                <td>{{ team.recent }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- 羽毛球数据预览 -->
      <div v-else-if="sportType === 'badminton'" class="sport-preview">
        <div class="preview-section">
          <h2 class="section-title">比赛结果</h2>
          <div class="score-display">
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamA }}</span>
              <span class="score">{{ badmintonData.matchInfo.homeTeamScore }}</span>
            </div>
            <span class="score-separator">-</span>
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamB }}</span>
              <span class="score">{{ badmintonData.matchInfo.awayTeamScore }}</span>
            </div>
            <span class="match-format">{{ badmintonData.matchInfo.format }}</span>
          </div>
          
          <div class="game-scores">
            <div v-for="(game, index) in badmintonData.matchInfo.gameScores" :key="index" class="game-score">
              <span class="game-label">第{{ index + 1 }}局</span>
              <span class="score">{{ game.home }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ game.away }}</span>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">参赛人员</h2>
          
          <div class="lineups-container">
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamA }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(player, index) in badmintonData.players.home" :key="index">
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamB }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(player, index) in badmintonData.players.away" :key="index">
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">比赛事件</h2>
          <table class="events-table">
            <thead>
              <tr>
                <th>局次</th>
                <th>比分</th>
                <th>事件类型</th>
                <th>球员</th>
                <th>队伍</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(event, index) in badmintonData.events" :key="index">
                <td>第{{ event.game }}局</td>
                <td>{{ event.point }}:{{ event.point }}</td>
                <td>{{ event.type }}</td>
                <td>{{ event.player || '-' }}</td>
                <td>{{ event.team || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">球队数据统计</h2>
          <table class="rankings-table">
            <thead>
              <tr>
                <th>排名</th>
                <th>球队</th>
                <th>场次</th>
                <th>胜</th>
                <th>平</th>
                <th>负</th>
                <th>胜率(%)</th>
                <th>积分</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in badmintonData.teamStats" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ team.team }}</td>
                <td>{{ team.matches }}</td>
                <td>{{ team.wins }}</td>
                <td>{{ team.draws }}</td>
                <td>{{ team.losses }}</td>
                <td>{{ team.winRate }}</td>
                <td>{{ team.points }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- 排球数据预览 -->
      <div v-else-if="sportType === 'volleyball'" class="sport-preview">
        <div class="preview-section">
          <h2 class="section-title">比赛结果</h2>
          <div class="score-display">
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamA }}</span>
              <span class="score">{{ volleyballData.matchInfo.homeTeamScore }}</span>
            </div>
            <span class="score-separator">-</span>
            <div class="team-score">
              <span class="team-name">{{ eventInfo.teamB }}</span>
              <span class="score">{{ volleyballData.matchInfo.awayTeamScore }}</span>
            </div>
          </div>
          
          <div class="game-scores">
            <div v-for="(game, index) in volleyballData.matchInfo.gameScores" :key="index" class="game-score">
              <span class="game-label">第{{ index + 1 }}局</span>
              <span class="score">{{ game.home }}</span>
              <span class="score-separator">-</span>
              <span class="score">{{ game.away }}</span>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">首发阵容</h2>
          
          <div class="lineups-container">
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamA }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in volleyballData.lineups.home" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            
            <div class="lineup-team">
              <h3 class="team-title">{{ eventInfo.teamB }}</h3>
              <table class="lineups-table">
                <thead>
                  <tr>
                    <th>球衣号</th>
                    <th>球员姓名</th>
                    <th>位置</th>
                    <th>是否首发</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="player in volleyballData.lineups.away" :key="player.number">
                    <td>{{ player.number }}</td>
                    <td>{{ player.name }}</td>
                    <td>{{ player.position }}</td>
                    <td>{{ player.isStarter ? '是' : '否' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">比赛事件</h2>
          <table class="events-table">
            <thead>
              <tr>
                <th>局次</th>
                <th>比分</th>
                <th>事件类型</th>
                <th>球员</th>
                <th>队伍</th>
                <th>描述</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(event, index) in volleyballData.events" :key="index">
                <td>第{{ event.game }}局</td>
                <td>{{ event.point }}:{{ event.point }}</td>
                <td>{{ event.type }}</td>
                <td>{{ event.player || '-' }}</td>
                <td>{{ event.team || '-' }}</td>
                <td>{{ event.skill || event.description || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">球队数据统计</h2>
          <table class="rankings-table">
            <thead>
              <tr>
                <th>排名</th>
                <th>球队</th>
                <th>场次</th>
                <th>胜</th>
                <th>负</th>
                <th>积分</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(team, index) in volleyballData.teamStats" :key="index">
                <td>{{ index + 1 }}</td>
                <td>{{ team.team }}</td>
                <td>{{ team.matches }}</td>
                <td>{{ team.wins }}</td>
                <td>{{ team.losses }}</td>
                <td>{{ team.points }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      
      <!-- 水上运动数据预览 -->
      <div v-else-if="sportType === 'water_sports'" class="sport-preview">
        <div class="preview-section">
          <h2 class="section-title">{{ waterSportsData.category }}赛事数据</h2>
          <table class="rankings-table">
            <thead>
              <tr>
                <th>排名</th>
                <th>队伍</th>
                <th>参赛人数</th>
                <th>赛事名称</th>
                <th>比赛时间</th>
                <th>最终成绩(秒)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="data in waterSportsData.eventData" :key="data.rank">
                <td>{{ data.rank }}</td>
                <td>{{ data.team }}</td>
                <td>{{ data.peopleCount }}</td>
                <td>{{ data.eventName }}</td>
                <td>{{ data.eventTime }}</td>
                <td>{{ data.finalResult }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <div class="preview-section">
          <h2 class="section-title">分段成绩详情</h2>
          <div v-for="teamResult in waterSportsData.segmentResults" :key="teamResult.team" class="segment-team-result">
            <h3 class="team-title">{{ teamResult.team }}</h3>
            <table class="segments-table">
              <thead>
                <tr>
                  <th>分段名称</th>
                  <th>时间(秒)</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(segment, index) in teamResult.segments" :key="index">
                  <td v-if="segment.name">{{ segment.name }}</td>
                  <td v-else><strong>总计</strong></td>
                  <td v-if="segment.time">{{ segment.time }}</td>
                  <td v-else><strong>{{ segment.total }}</strong></td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.data-preview {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  padding-bottom: 40px;
}

.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 16px;
  color: #666;
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

.preview-content {
  flex: 1;
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.preview-header {
  text-align: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 2px solid #1890ff;
}

.preview-title {
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin-bottom: 10px;
}

.preview-subtitle {
  font-size: 18px;
  color: #666;
  margin-bottom: 10px;
}

.preview-date {
  font-size: 14px;
  color: #999;
}

.action-buttons {
  display: flex;
  gap: 16px;
  margin-bottom: 30px;
  justify-content: center;
  flex-wrap: wrap;
}

.action-btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.print-btn {
  background-color: #1890ff;
  color: white;
}

.export-btn {
  background-color: #52c41a;
  color: white;
}

.back-btn {
  background-color: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
}

.sport-preview {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.preview-section {
  background-color: white;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 20px;
  font-weight: 500;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
}

/* 比分显示样式 */
.score-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 20px;
}

.team-score {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.team-name {
  font-size: 20px;
  font-weight: 500;
  color: #333;
}

.score {
  font-size: 36px;
  font-weight: 700;
  color: #1890ff;
}

.score-separator {
  font-size: 36px;
  font-weight: 700;
  color: #666;
}

.match-status,
.match-format {
  padding: 8px 16px;
  background-color: #f5f5f5;
  border-radius: 20px;
  font-size: 14px;
  color: #666;
}

/* 节次/局次比分样式 */
.quarter-scores,
.game-scores {
  display: flex;
  gap: 20px;
  justify-content: center;
  flex-wrap: wrap;
}

.quarter-score,
.game-score {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background-color: #fafafa;
  border-radius: 8px;
  font-size: 14px;
}

.quarter-label,
.game-label {
  font-weight: 500;
  color: #666;
}

/* 阵容容器样式 */
.lineups-container {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.lineup-team {
  flex: 1;
  min-width: 300px;
}

.team-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f0;
}

/* 表格样式 */
.lineups-table,
.events-table,
.rankings-table,
.segments-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.lineups-table th,
.events-table th,
.rankings-table th,
.segments-table th {
  background-color: #fafafa;
  padding: 12px 8px;
  text-align: left;
  font-weight: 500;
  color: #333;
  border-bottom: 2px solid #f0f0f0;
}

.lineups-table td,
.events-table td,
.rankings-table td,
.segments-table td {
  padding: 10px 8px;
  border-bottom: 1px solid #f0f0f0;
  color: #666;
}

.lineups-table tr:hover,
.events-table tr:hover,
.rankings-table tr:hover,
.segments-table tr:hover {
  background-color: #fafafa;
}

/* 分段成绩样式 */
.segment-team-result {
  margin-bottom: 24px;
}

.segment-team-result:last-child {
  margin-bottom: 0;
}

/* 打印样式 */
@media print {
  .action-buttons {
    display: none;
  }
  
  .preview-section {
    box-shadow: none;
    border: 1px solid #eee;
  }
  
  .data-preview {
    background-color: white;
  }
  
  @page {
    margin: 2cm;
  }
}

/* 移动端适配 */
@media (max-width: 768px) {
  .preview-content {
    padding: 16px;
  }
  
  .preview-title {
    font-size: 24px;
  }
  
  .preview-subtitle {
    font-size: 16px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
  
  .action-btn {
    width: 100%;
  }
  
  .preview-section {
    padding: 16px;
  }
  
  .section-title {
    font-size: 18px;
  }
  
  .score-display {
    gap: 12px;
  }
  
  .team-name {
    font-size: 16px;
  }
  
  .score {
    font-size: 28px;
  }
  
  .score-separator {
    font-size: 28px;
  }
  
  .quarter-scores,
  .game-scores {
    flex-direction: column;
    gap: 12px;
  }
  
  .lineups-container {
    flex-direction: column;
  }
  
  .lineup-team {
    min-width: auto;
  }
  
  .lineups-table,
  .events-table,
  .rankings-table,
  .segments-table {
    font-size: 12px;
  }
  
  .lineups-table th,
  .events-table th,
  .rankings-table th,
  .segments-table th,
  .lineups-table td,
  .events-table td,
  .rankings-table td,
  .segments-table td {
    padding: 8px 4px;
  }
  
  /* 表格横向滚动 */
  .lineups-table,
  .events-table,
  .rankings-table,
  .segments-table {
    display: block;
    overflow-x: auto;
    white-space: nowrap;
  }
}
</style>