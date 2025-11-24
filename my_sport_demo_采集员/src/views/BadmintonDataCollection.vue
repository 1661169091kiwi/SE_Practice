<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'

const props = defineProps(['eventId', 'eventInfo'])
const route = useRoute()

// 比赛类型：单打或双打
const matchType = ref('single') // 'single' 或 'double'

// 参赛人员
const players = ref({
  teamA: [
    { id: 1, name: '', position: 'player1' },
    { id: 2, name: '', position: 'player2' },
  ],
  teamB: [
    { id: 1, name: '', position: 'player1' },
    { id: 2, name: '', position: 'player2' },
  ],
})

// 局数（最多3局）
const maxGames = 3
// 当前选中的局次
const currentGame = ref(1)
// 局分
const gameScores = ref([
  { game: 1, teamA: 0, teamB: 0, winner: null },
  { game: 2, teamA: 0, teamB: 0, winner: null },
  { game: 3, teamA: 0, teamB: 0, winner: null },
])





// 加载状态
const isLoading = ref(true)
const isSubmitting = ref(false)

// 计算大比分
const matchScore = computed(() => {
  const winsA = gameScores.value.filter((game) => game.winner === 'teamA').length
  const winsB = gameScores.value.filter((game) => game.winner === 'teamB').length
  return { teamA: winsA, teamB: winsB }
})

// 计算比赛结果
const matchResult = computed(() => {
  const winsA = gameScores.value.filter((game) => game.winner === 'teamA').length
  const winsB = gameScores.value.filter((game) => game.winner === 'teamB').length

  if (winsA >= 2) {
    return 'teamA'
  } else if (winsB >= 2) {
    return 'teamB'
  } else {
    return null
  }
})

// 加载赛事数据
const loadEventData = () => {
  isLoading.value = true
  setTimeout(() => {
    // 赛事信息已通过props传入，这里可以进行额外的数据加载
    isLoading.value = false
  }, 1000)
}

// 切换比赛类型（单打/双打）
const switchMatchType = (type) => {
  matchType.value = type
}

// 更新参赛人员
const updatePlayerName = (team, playerId, name) => {
  const player = players.value[team].find((p) => p.id === playerId)
  if (player) {
    player.name = name
  }
}



// 更新局分
const updateGameScore = (gameIndex, team, increment) => {
  const game = gameScores.value[gameIndex]

  // 如果该局已经有胜利者，不允许再修改比分
  if (game.winner) {
    alert('该局已经结束，无法修改比分')
    return
  }

  game[team] = Math.max(0, game[team] + increment)

  // 检查是否满足胜利条件（21分，且领先至少2分）
  checkWinCondition(game)
}

// 检查胜利条件
const checkWinCondition = (game) => {
  const scoreA = game.teamA
  const scoreB = game.teamB

  // 标准规则：21分且领先至少2分
  if (scoreA >= 21 && scoreA - scoreB >= 2) {
    game.winner = 'teamA'
    // 检查是否需要结束后续局数
    checkMatchCompletion()
  } else if (scoreB >= 21 && scoreB - scoreA >= 2) {
    game.winner = 'teamB'
    // 检查是否需要结束后续局数
    checkMatchCompletion()
  }
}

// 检查比赛是否完成（一方先赢2局）
const checkMatchCompletion = () => {
  const winsA = gameScores.value.filter((game) => game.winner === 'teamA').length
  const winsB = gameScores.value.filter((game) => game.winner === 'teamB').length

  // 如果一方已经赢了2局，自动终止后续局数
  if (winsA === 2 || winsB === 2) {
    // 查找未完成的局数
    for (let i = 0; i < gameScores.value.length; i++) {
      if (!gameScores.value[i].winner) {
        // 标记为无需进行
        gameScores.value[i].teamA = 0
        gameScores.value[i].teamB = 0
        gameScores.value[i].winner = 'not_played'
      }
    }
  }
}

// 重置局分
const resetGameScore = (gameIndex) => {
  const game = gameScores.value[gameIndex]
  game.teamA = 0
  game.teamB = 0
  game.winner = null
  // 重置后续局数的'not_played'状态
  for (let i = gameIndex + 1; i < gameScores.value.length; i++) {
    if (gameScores.value[i].winner === 'not_played') {
      gameScores.value[i].winner = null
    }
  }
}

// 切换当前局次
const switchGame = (game) => {
  // 如果之前的局没有结束，不允许切换到新局
  for (let i = 0; i < game - 1; i++) {
    if (!gameScores.value[i].winner && gameScores.value[i].winner !== 'not_played') {
      alert(`第${i + 1}局尚未结束，无法切换到第${game}局`)
      return
    }
  }
  currentGame.value = game
}



// 定义获取当前数据的方法，供父组件调用
defineExpose({
  // 获取羽毛球赛事数据的方法
  getBadmintonData() {
    // 这里返回羽毛球特有的数据结构
    return {
      // scores, sets, players等羽毛球特有的数据
      sportType: 'badminton'
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
  <div class="badminton-data-collection">

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 数据录入内容 -->
    <div v-else class="collection-content">
      <!-- 局分统计 -->
      <div class="form-section score-section">
        <h2 class="section-title">比分统计</h2>

        <!-- 大比分 -->
        <div class="match-score">
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamA }}</span>
            <span class="score">{{ matchScore.teamA }}</span>
          </div>
          <span class="vs">:</span>
          <div class="team-info">
            <span class="team-name">{{ eventInfo.teamB }}</span>
            <span class="score">{{ matchScore.teamB }}</span>
          </div>
        </div>

        <!-- 比赛结果 -->
        <div class="match-result" v-if="matchResult">
          <span class="result-label">比赛结果：</span>
          <span class="result-winner"
            >{{ matchResult === 'teamA' ? eventInfo.teamA : eventInfo.teamB }} 胜</span
          >
        </div>
        <div class="match-result" v-else-if="matchScore.teamA > 0 || matchScore.teamB > 0">
          <span class="result-label">比赛结果：</span>
          <span class="result-draw">进行中</span>
        </div>

        <!-- 局分 -->
        <div class="game-scores">
          <div
            v-for="(game, index) in gameScores"
            :key="game.game"
            class="game-score"
            :class="{
              completed: game.winner && game.winner !== 'not_played',
              'not-played': game.winner === 'not_played',
              active: game.game === currentGame,
            }"
          >
            <div class="game-header">
              <span class="game-label">第{{ game.game }}局</span>
              <div class="game-actions">
                <button
                  @click="switchGame(game.game)"
                  class="switch-game-btn"
                  :disabled="
                    game.winner === 'not_played' ||
                    (game.game > 1 && !gameScores[game.game - 2].winner)
                  "
                >
                  切换
                </button>
                <button @click="resetGameScore(index)" class="reset-game-btn" v-if="game.winner">
                  重置
                </button>
              </div>
            </div>
            <div class="game-controls" v-if="game.winner !== 'not_played'">
              <div class="team-game-score">
                <button
                  @click="updateGameScore(index, 'teamA', -1)"
                  class="score-btn decrease"
                  :disabled="game.winner"
                >
                  -
                </button>
                <span class="game-score-value">{{ game.teamA }}</span>
                <button
                  @click="updateGameScore(index, 'teamA', 1)"
                  class="score-btn increase"
                  :disabled="game.winner"
                >
                  +
                </button>
              </div>
              <div class="team-game-score">
                <button
                  @click="updateGameScore(index, 'teamB', -1)"
                  class="score-btn decrease"
                  :disabled="game.winner"
                >
                  -
                </button>
                <span class="game-score-value">{{ game.teamB }}</span>
                <button
                  @click="updateGameScore(index, 'teamB', 1)"
                  class="score-btn increase"
                  :disabled="game.winner"
                >
                  +
                </button>
              </div>
            </div>
            <div class="game-winner" v-if="game.winner && game.winner !== 'not_played'">
              {{ game.winner === 'teamA' ? eventInfo.teamA : eventInfo.teamB }} 胜
            </div>
            <div class="game-winner not-played" v-else-if="game.winner === 'not_played'">
              无需进行
            </div>
          </div>
        </div>
      </div>

      <!-- 参赛人员 -->
      <div class="form-section players-section">
        <h2 class="section-title">参赛人员</h2>

        <!-- 比赛类型选择 -->
        <div class="match-type-selector">
          <button
            @click="switchMatchType('single')"
            class="match-type-btn"
            :class="{ active: matchType === 'single' }"
          >
            单打
          </button>
          <button
            @click="switchMatchType('double')"
            class="match-type-btn"
            :class="{ active: matchType === 'double' }"
          >
            双打
          </button>
        </div>

        <!-- 参赛人员列表 -->
        <div class="players-container">
          <!-- 主队 -->
          <div class="team-players">
            <h3 class="team-title">{{ eventInfo.teamA }}</h3>
            <div class="player-list">
              <div
                v-for="player in players.teamA"
                :key="player.id"
                class="player-item"
                :class="{ hidden: matchType === 'single' && player.position === 'player2' }"
              >
                <span class="player-position">
                  {{
                    matchType === 'double'
                      ? player.position === 'player1'
                        ? '选手1'
                        : '选手2'
                      : '选手'
                  }}
                </span>
                <input
                  v-model="player.name"
                  type="text"
                  placeholder="球员姓名"
                  class="player-input"
                />
              </div>
            </div>
          </div>

          <!-- 客队 -->
          <div class="team-players">
            <h3 class="team-title">{{ eventInfo.teamB }}</h3>
            <div class="player-list">
              <div
                v-for="player in players.teamB"
                :key="player.id"
                class="player-item"
                :class="{ hidden: matchType === 'single' && player.position === 'player2' }"
              >
                <span class="player-position">
                  {{
                    matchType === 'double'
                      ? player.position === 'player1'
                        ? '选手1'
                        : '选手2'
                      : '选手'
                  }}
                </span>
                <input
                  v-model="player.name"
                  type="text"
                  placeholder="球员姓名"
                  class="player-input"
                />
              </div>
            </div>
          </div>
        </div>
      </div>




    </div>

    <!-- 保存操作由父组件BaseDataCollection统一处理 -->
  </div>
</template>

<style scoped>
.badminton-data-collection {
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
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
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
  margin-top: 16px;
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

.match-score {
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

.game-scores {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.game-score {
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
  border: 2px solid transparent;
  transition: all 0.3s;
}

.game-score.active {
  border-color: #1890ff;
  background-color: #e6f7ff;
}

.game-score.completed {
  opacity: 0.8;
}

.game-score.not-played {
  opacity: 0.5;
  background-color: #f0f0f0;
}

.game-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.game-label {
  font-weight: 500;
  color: #333;
}

.game-actions {
  display: flex;
  gap: 8px;
}

.switch-game-btn,
.reset-game-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.switch-game-btn {
  background-color: #1890ff;
  color: white;
}

.switch-game-btn:disabled {
  background-color: #d9d9d9;
  color: #999;
  cursor: not-allowed;
}

.reset-game-btn {
  background-color: #f5222d;
  color: white;
}

.game-controls {
  display: flex;
  justify-content: center;
  gap: 40px;
}

.team-game-score {
  display: flex;
  align-items: center;
  gap: 12px;
}

.score-btn {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 4px;
  font-size: 18px;
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

.score-btn:disabled {
  background-color: #d9d9d9;
  color: #999;
  cursor: not-allowed;
}

.game-score-value {
  font-size: 24px;
  font-weight: bold;
  color: #333;
  min-width: 40px;
}

.game-winner {
  margin-top: 12px;
  font-weight: 500;
  color: #52c41a;
}

.game-winner.not-played {
  color: #999;
}

/* 参赛人员样式 */
.match-type-selector {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.match-type-btn {
  flex: 1;
  padding: 12px;
  border: 2px solid #ddd;
  background-color: white;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.match-type-btn.active {
  background-color: #1890ff;
  color: white;
  border-color: #1890ff;
}

.players-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.team-players {
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
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

.player-item.hidden {
  display: none;
}

.player-position {
  font-size: 14px;
  color: #666;
  width: 60px;
}

.player-input {
  flex: 1;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
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

  .match-score {
    flex-direction: column;
    gap: 16px;
  }

  .team-score {
    flex-direction: row;
    gap: 16px;
  }

  .game-controls {
    gap: 20px;
  }

  .players-container {
    grid-template-columns: 1fr;
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
