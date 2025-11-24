<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'

const props = defineProps(['eventId', 'eventInfo'])
const route = useRoute()

// 实时比分
const scores = ref({
  teamA: 0,
  teamB: 0
})

// 首发阵容
const lineups = ref({
  teamA: [
    { id: 1, name: '', isStarting: true },
    { id: 2, name: '', isStarting: true },
    { id: 3, name: '', isStarting: true },
    { id: 4, name: '', isStarting: true },
    { id: 5, name: '', isStarting: true },
    { id: 6, name: '', isStarting: true },
    { id: 7, name: '', isStarting: true },
    { id: 8, name: '', isStarting: true },
    { id: 9, name: '', isStarting: true },
    { id: 10, name: '', isStarting: true },
    { id: 11, name: '', isStarting: true }
  ],
  teamB: [
    { id: 1, name: '', isStarting: true },
    { id: 2, name: '', isStarting: true },
    { id: 3, name: '', isStarting: true },
    { id: 4, name: '', isStarting: true },
    { id: 5, name: '', isStarting: true },
    { id: 6, name: '', isStarting: true },
    { id: 7, name: '', isStarting: true },
    { id: 8, name: '', isStarting: true },
    { id: 9, name: '', isStarting: true },
    { id: 10, name: '', isStarting: true },
    { id: 11, name: '', isStarting: true }
  ]
})





// 加载状态
const isLoading = ref(true)

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



// 更新比分
const updateScore = (team, increment) => {
  scores.value[team] = Math.max(0, scores.value[team] + increment)
}



// 定义获取当前数据的方法，供父组件调用
defineExpose({
  // 获取足球赛事数据的方法
  getFootballData() {
    return {
      scores: scores.value,
      lineups: lineups.value,
      // 这里可以添加更多足球特有的数据
      sportType: 'football'
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
  <div class="football-data-collection">
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 数据录入内容 -->
    <div v-else class="collection-content">
      <!-- 实时比分 -->
      <div class="form-section score-section">
        <h2 class="section-title">实时比分</h2>
        <div class="score-display">
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamA }}</span>
            <div class="score-controls">
              <button @click="updateScore('teamA', -1)" class="score-btn decrease">-</button>
              <span class="score">{{ scores.teamA }}</span>
              <button @click="updateScore('teamA', 1)" class="score-btn increase">+</button>
            </div>
          </div>
          <span class="vs">:</span>
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamB }}</span>
            <div class="score-controls">
              <button @click="updateScore('teamB', -1)" class="score-btn decrease">-</button>
              <span class="score">{{ scores.teamB }}</span>
              <button @click="updateScore('teamB', 1)" class="score-btn increase">+</button>
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
                :disabled="lineups.teamA.length <= 11"
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
                :disabled="lineups.teamB.length <= 11"
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
.football-data-collection {
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

.score-display {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  padding: 20px 0;
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
  margin-bottom: 8px;
}

.score-controls {
  display: flex;
  align-items: center;
  gap: 16px;
}

.score-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 50%;
  font-size: 20px;
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

.score {
  font-size: 32px;
  font-weight: bold;
  color: #333;
  min-width: 50px;
  text-align: center;
}

.vs {
  font-size: 28px;
  font-weight: bold;
  color: #666;
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

/* 事件样式 */
.form-row {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.form-group {
  flex: 1;
  min-width: 150px;
}

.form-label {
  display: block;
  margin-bottom: 6px;
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.form-control {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
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

.event-type.goal {
  background-color: #52c41a;
}

.event-type.red_card {
  background-color: #f5222d;
}

.event-type.yellow_card {
  background-color: #faad14;
}

.event-type.penalty {
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

/* 排行榜样式 */
.ranking-table-container {
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 16px;
}

.ranking-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.ranking-table th,
.ranking-table td {
  padding: 12px 8px;
  text-align: center;
  border-bottom: 1px solid #f0f0f0;
}

.ranking-table th {
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

.goal-diff {
  font-weight: 500;
}

.points {
  font-weight: bold;
  color: #1890ff;
}

.remove-ranking-btn {
  padding: 4px 8px;
  background-color: #f5f5f5;
  color: #f5222d;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.add-ranking-btn {
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
  
  .score-display {
    flex-direction: column;
    gap: 16px;
  }
  
  .team-score {
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
  
  .ranking-table {
    font-size: 12px;
  }
  
  .ranking-table th,
  .ranking-table td {
    padding: 8px 4px;
  }
  
  .table-input.small {
    max-width: 40px;
    padding: 4px;
  }
}
</style>