<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { get } from '@/utils/http'

const route = useRoute()
const router = useRouter()
const eventId = route.params.eventId
const standings = ref([])
const isLoading = ref(false)
const eventInfo = ref({}) // 存储赛事信息（如果有API获取）

// 获取积分榜数据
const fetchStandings = async () => {
  isLoading.value = true
  try {
    const res = await get(`/events/${eventId}/standings`)
    if (res.code === 200 && res.data) {
      standings.value = res.data
    } else {
      // 如果没有数据，模拟一些数据用于展示UI（在开发阶段）
      // 实际发布时应删除
      // console.warn('Failed to fetch standings, using mock data')
      // standings.value = mockStandings
      console.error('Failed to fetch standings:', res.msg)
    }
  } catch (error) {
    console.error('Error fetching standings:', error)
  } finally {
    isLoading.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.back()
}

onMounted(() => {
  fetchStandings()
})
</script>

<template>
  <div class="standings-container">
    <header class="page-header">
      <button class="back-button" @click="goBack">
        ←
      </button>
      <h1>积分榜</h1>
      <div class="placeholder"></div>
    </header>

    <div class="content-area">
      <div v-if="isLoading" class="loading-state">
        <p>加载中...</p>
      </div>
      <div v-else-if="standings.length === 0" class="empty-state">
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
            <tr v-for="(team, index) in standings" :key="index">
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
</template>

<style scoped>
.standings-container {
  min-height: 100vh;
  background-color: var(--background-primary);
  display: flex;
  flex-direction: column;
}

.page-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  padding: 16px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: var(--shadow-md);
}

.back-button {
  background: none;
  border: none;
  color: white;
  font-size: 20px;
  cursor: pointer;
  padding: 0 10px 0 0;
}

.page-header h1 {
  font-size: 18px;
  font-weight: 500;
  margin: 0;
  flex: 1;
  text-align: center;
}

.placeholder {
  width: 30px;
}

.content-area {
  padding: 16px;
  flex: 1;
}

.loading-state,
.empty-state {
  text-align: center;
  padding: 40px 0;
  color: var(--text-tertiary);
}

.table-container {
  overflow-x: auto;
  border-radius: var(--border-radius-lg);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-color);
  background: white;
}

.standings-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
  min-width: 500px; /* Ensure table doesn't get squashed on very small screens */
}

.standings-table th,
.standings-table td {
  padding: 12px 8px;
  text-align: center;
  border-bottom: 1px solid var(--border-color);
}

.standings-table th {
  background-color: #f8fafc;
  color: var(--text-secondary);
  font-weight: 500;
}

.standings-table tr:last-child td {
  border-bottom: none;
}

.team-col {
  text-align: left;
  font-weight: 500;
  color: var(--text-primary);
  min-width: 120px;
}

.rank-col {
  width: 50px;
}

.points-col {
  font-weight: bold;
  color: var(--primary-color);
}

.rank-badge {
  display: inline-block;
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  border-radius: 50%;
  background-color: #f0f0f0;
  color: var(--text-secondary);
  font-size: 12px;
}

.rank-1 {
  background-color: #ffe58f;
  color: #d46b08;
}

.rank-2 {
  background-color: #d3adf7;
  color: #531dab;
}

.rank-3 {
  background-color: #ffccc7;
  color: #cf1322;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    padding: 12px 16px;
  }
  
  .content-area {
    padding: 12px;
  }

  .table-container {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch; /* iOS 平滑滚动 */
  }

  .standings-table {
    min-width: 100%; /* 允许表格自适应宽度，必要时才滚动 */
    width: auto;
  }

  .standings-table th, 
  .standings-table td {
    padding: 10px 4px;
    font-size: 13px;
    white-space: nowrap; /* 防止换行 */
  }
  
  .rank-col {
    width: 36px;
    min-width: 36px;
    padding-left: 0;
    padding-right: 0;
  }
  
  .team-col {
    min-width: 90px;
    max-width: 120px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;
  }

  /* 缩小排名徽章 */
  .rank-badge {
    width: 20px;
    height: 20px;
    line-height: 20px;
    font-size: 11px;
  }
}
</style>
