<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'

const props = defineProps(['eventId', 'eventInfo'])
const route = useRoute()

// 单节比分
const quarterScores = ref([
  { quarter: 1, teamA: 0, teamB: 0 },
  { quarter: 2, teamA: 0, teamB: 0 },
  { quarter: 3, teamA: 0, teamB: 0 },
  { quarter: 4, teamA: 0, teamB: 0 }
])



// 首发阵容
const lineups = ref({
  teamA: [
    { id: 1, name: '', isStarting: true },
    { id: 2, name: '', isStarting: true },
    { id: 3, name: '', isStarting: true },
    { id: 4, name: '', isStarting: true },
    { id: 5, name: '', isStarting: true }
  ],
  teamB: [
    { id: 1, name: '', isStarting: true },
    { id: 2, name: '', isStarting: true },
    { id: 3, name: '', isStarting: true },
    { id: 4, name: '', isStarting: true },
    { id: 5, name: '', isStarting: true }
  ]
})





// 加载状态
const isLoading = ref(true)

// 计算总比分
const totalScore = computed(() => {
  return {
    teamA: quarterScores.value.reduce((sum, quarter) => sum + quarter.teamA, 0),
    teamB: quarterScores.value.reduce((sum, quarter) => sum + quarter.teamB, 0)
  }
})

// 比赛结果
const matchResult = ref(null)

// 计算并设置比赛结果
const setMatchResult = () => {
  const scoreA = totalScore.value.teamA
  const scoreB = totalScore.value.teamB
  
  if (scoreA > scoreB) {
    matchResult.value = 'teamA'
  } else if (scoreB > scoreA) {
    matchResult.value = 'teamB'
  } else {
    matchResult.value = null
  }
}

// 加载赛事数据
const loadEventData = () => {
  isLoading.value = true
  setTimeout(() => {
    // 赛事信息已通过props传入，这里可以进行额外的数据加载
    isLoading.value = false
  }, 1000)
}

// 添加球员
const addPlayer = (team) => {
  const newPlayer = {
    id: Date.now(),
    name: '',
    isStarting: false
  }
  lineups.value[team].push(newPlayer)
}

// 移除球员
const removePlayer = (team, playerId) => {
  lineups.value[team] = lineups.value[team].filter(p => p.id !== playerId)
}



// 更新单节比分
const updateQuarterScore = (quarterIndex, team, increment) => {
  quarterScores.value[quarterIndex][team] = Math.max(0, quarterScores.value[quarterIndex][team] + increment)
  setMatchResult()
}





// 定义获取当前数据的方法，供父组件调用
defineExpose({
  // 获取篮球赛事数据的方法
  getBasketballData() {
    // 这里返回篮球特有的数据结构
    // 具体字段需要根据实际的篮球数据模型来定义
    return {
      // scores, players, statistics等篮球特有的数据
      sportType: 'basketball'
    }
  }
})

// 组件实例引用
const fieldDescRefs = ref([])

// 生命周期钩子
onMounted(() => {
  loadEventData()
  setTimeout(() => {
    fieldDescRefs.value.forEach(ref => ref?.init?.())
  }, 100)
})

onUnmounted(() => {
  fieldDescRefs.value.forEach(ref => ref?.cleanup?.())
})
</script>

<template>
  <div class="basketball-data-collection">
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 数据录入内容 -->
    <div v-else class="collection-content">
      <!-- 单节比分 -->
      <div class="form-section score-section">
        <h2 class="section-title">比分统计</h2>
        
        <!-- 总比分 -->
        <div class="total-score">
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamA }}</span>
            <span class="score">{{ totalScore.teamA }}</span>
          </div>
          <span class="vs">:</span>
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamB }}</span>
            <span class="score">{{ totalScore.teamB }}</span>
          </div>
        </div>
        
        <!-- 比赛结果 -->
        <div class="match-result" v-if="matchResult">
          <span class="result-label">比赛结果：</span>
          <span class="result-winner">{{ matchResult === 'teamA' ? eventInfo.teamA : eventInfo.teamB }} 胜</span>
        </div>
        <div class="match-result" v-else-if="totalScore.teamA > 0 || totalScore.teamB > 0">
          <span class="result-label">比赛结果：</span>
          <span class="result-draw">待定</span>
        </div>
        
        <!-- 单节比分 -->
        <div class="quarter-scores">
          <div 
            v-for="(quarter, index) in quarterScores" 
            :key="quarter.quarter" 
            class="quarter-score"
          >
            <div class="quarter-label">第{{ quarter.quarter }}节</div>
            <div class="quarter-controls">
              <div class="team-quarter-score">
                <button 
                  @click="updateQuarterScore(index, 'teamA', -1)" 
                  class="score-btn decrease"
                >-</button>
                <span class="quarter-score-value">{{ quarter.teamA }}</span>
                <button 
                  @click="updateQuarterScore(index, 'teamA', 1)" 
                  class="score-btn increase"
                >+</button>
              </div>
              <div class="team-quarter-score">
                <button 
                  @click="updateQuarterScore(index, 'teamB', -1)" 
                  class="score-btn decrease"
                >-</button>
                <span class="quarter-score-value">{{ quarter.teamB }}</span>
                <button 
                  @click="updateQuarterScore(index, 'teamB', 1)" 
                  class="score-btn increase"
                >+</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 首发阵容 -->
      <div class="form-section lineup-section">
        <h2 class="section-title">首发阵容</h2>
        
        <!-- 主队阵容 -->
        <div class="team-lineup">
          <h3 class="team-title">{{ eventInfo.teamA }}</h3>
          <div class="player-list">
            <div v-for="player in lineups.teamA" :key="player.id" class="player-item">
              <input 
                v-model="player.name"
                type="text" 
                placeholder="球员姓名"
                class="player-input"
              />
              <label class="checkbox-label">
                <input 
                  v-model="player.isStarting" 
                  type="checkbox"
                  class="starting-checkbox"
                />
                首发
              </label>
              <button 
                @click="removePlayer('teamA', player.id)" 
                class="remove-btn"
                :disabled="lineups.teamA.length <= 5"
              >
                删除
              </button>
            </div>
            <button @click="addPlayer('teamA')" class="add-player-btn">
              + 添加替补
            </button>
          </div>
        </div>
        
        <!-- 客队阵容 -->
        <div class="team-lineup">
          <h3 class="team-title">{{ eventInfo.teamB }}</h3>
          <div class="player-list">
            <div v-for="player in lineups.teamB" :key="player.id" class="player-item">
              <input 
                v-model="player.name"
                type="text" 
                placeholder="球员姓名"
                class="player-input"
              />
              <label class="checkbox-label">
                <input 
                  v-model="player.isStarting" 
                  type="checkbox"
                  class="starting-checkbox"
                />
                首发
              </label>
              <button 
                @click="removePlayer('teamB', player.id)" 
                class="remove-btn"
                :disabled="lineups.teamB.length <= 5"
              >
                删除
              </button>
            </div>
            <button @click="addPlayer('teamB')" class="add-player-btn">
              + 添加替补
            </button>
          </div>
        </div>
      </div>




    </div>
    
    <!-- 保存操作由父组件BaseDataCollection统一处理 -->
  </div>
</template>

<style scoped>
.basketball-data-collection {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  padding-bottom: 80px;
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

.collection-content {
  flex: 1;
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-section {
  background-color: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.section-title {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
}

/* 比分样式 */
.score-section {
  text-align: center;
}

.match-result {
  margin-top: 12px;
  padding: 12px;
  background-color: #f0f9ff;
  border-radius: 6px;
  display: inline-block;
  font-size: 16px;
}

.result-label {
  color: #666;
  font-weight: 500;
}

.result-winner {
  color: #1890ff;
  font-weight: bold;
  margin-left: 6px;
}

.result-draw {
  color: #faad14;
  font-weight: 500;
  margin-left: 6px;
}

.total-score {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 20px 0;
  margin-bottom: 20px;
  background-color: #fafafa;
  border-radius: 8px;
}

.team-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 120px;
}

.team-name {
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.score {
  font-size: 36px;
  font-weight: bold;
  color: #1890ff;
  min-width: 60px;
}

.vs {
  font-size: 32px;
  font-weight: bold;
  color: #666;
}

.quarter-scores {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.quarter-score {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
}

.quarter-label {
  font-weight: 500;
  color: #333;
}

.quarter-controls {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.team-quarter-score {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.score-btn {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.score-btn.decrease {
  background-color: #f5f5f5;
  color: #666;
}

.score-btn.increase {
  background-color: #1890ff;
  color: white;
}

.quarter-score-value {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  min-width: 30px;
  text-align: center;
}

/* 阵容样式 */
.team-lineup {
  margin-bottom: 24px;
}

.team-lineup:last-child {
  margin-bottom: 0;
}

.team-title {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  margin-bottom: 16px;
}

.player-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.player-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.player-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #666;
  cursor: pointer;
}

.starting-checkbox {
  width: 16px;
  height: 16px;
}

.remove-btn,
.add-player-btn {
  padding: 8px 12px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.remove-btn {
  background-color: #f5f5f5;
  color: #f5222d;
}

.remove-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.add-player-btn {
  background-color: #52c41a;
  color: white;
  align-self: flex-start;
  margin-top: 8px;
}



.add-event-btn {
  padding: 8px 16px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  align-self: flex-end;
  height: fit-content;
}

.event-list {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.empty-list {
  text-align: center;
  color: #999;
  padding: 20px;
}

.event-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px;
  background-color: #fafafa;
  border-radius: 8px;
}

.event-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.event-time {
  font-weight: bold;
  color: #1890ff;
}

.event-team {
  font-weight: 500;
  color: #333;
}

.event-player {
  color: #333;
}

.event-type {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  color: white;
}

.event-type.score {
  background-color: #52c41a;
}

.event-type.foul {
  background-color: #f5222d;
}

.event-type.timeout {
  background-color: #1890ff;
}

.event-type.substitution {
  background-color: #722ed1;
}

.event-note {
  color: #666;
  font-size: 13px;
}

.remove-event-btn {
  padding: 4px 8px;
  background-color: #f5f5f5;
  color: #f5222d;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

/* 球队数据样式 */
.stats-table-container {
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.stats-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.stats-table th,
.stats-table td {
  padding: 12px 8px;
  text-align: center;
  border-bottom: 1px solid #f0f0f0;
}

.stats-table th {
  background-color: #fafafa;
  font-weight: 500;
  color: #333;
  position: sticky;
  top: 0;
}

.table-input {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
}

.table-input.small {
  max-width: 60px;
}

.remove-stat-btn {
  padding: 4px 8px;
  background-color: #f5f5f5;
  color: #f5222d;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.add-stat-btn {
  padding: 8px 16px;
  background-color: #52c41a;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .collection-content {
    padding: 16px;
    gap: 16px;
  }
  
  .form-section {
    padding: 16px;
  }
  
  .total-score {
    flex-direction: column;
    gap: 16px;
  }
  
  .team-score {
    flex-direction: row;
    gap: 16px;
  }
  
  .quarter-scores {
    grid-template-columns: 1fr;
    gap: 8px;
  }
  
  .quarter-score {
    flex-direction: row;
    justify-content: space-between;
    padding: 12px;
  }
  
  .quarter-controls {
    flex-direction: row;
    gap: 16px;
  }
  
  .form-row {
    flex-direction: column;
    gap: 16px;
  }
  
  .form-group {
    min-width: auto;
  }
  
  .add-event-btn {
    align-self: stretch;
  }
  
  .event-info {
    flex-wrap: wrap;
  }
  
  .stats-table {
    font-size: 12px;
  }
  
  .stats-table th,
  .stats-table td {
    padding: 8px 4px;
  }
  
  .table-input.small {
    max-width: 40px;
    padding: 4px;
  }
  
  .match-result {
    margin-top: 12px;
    padding: 10px;
    font-size: 14px;
  }
}
</style>