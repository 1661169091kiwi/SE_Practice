<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'

const route = useRoute()

const props = defineProps(['eventId', 'eventInfo'])

// 赛制类型选项
const tournamentTypes = [
  { value: 'group', label: '单循环小组赛' },
  { value: 'knockout', label: '淘汰赛' },
  { value: 'group_knockout', label: '小组赛+淘汰赛' }
]

// 赛事阶段选项
const tournamentStages = [
  { value: 'group', label: '小组赛' },
  { value: 'semifinal', label: '半决赛' },
  { value: 'final', label: '决赛' }
]

// 小组信息
const groups = ref([
  { id: 'A', name: 'A组', teams: [] },
  { id: 'B', name: 'B组', teams: [] },
  { id: 'C', name: 'C组', teams: [] },
  { id: 'D', name: 'D组', teams: [] }
])

// 水上运动分项
const sportCategories = [
  { value: 'dragon_boat', label: '龙舟', segments: ['起航段', '中途段', '冲刺段'] },
  { value: 'canoe', label: '皮艇', segments: ['第一段', '第二段', '第三段', '冲刺段'] },
  { value: 'paddle_board', label: '桨板', segments: ['第一段', '第二段', '终点段'] },
]



// 参赛队伍
const teams = ref([
  {
    id: 1,
    teamName: '计算机学院',
    peopleCount: 20,
    groupId: 'A' // 所属小组
  },
  {
    id: 2,
    teamName: '机械学院',
    peopleCount: 20,
    groupId: 'A'
  },
  {
    id: 3,
    teamName: '商学院',
    peopleCount: 20,
    groupId: 'B'
  },
  {
    id: 4,
    teamName: '电子学院',
    peopleCount: 20,
    groupId: 'B'
  },
  {
    id: 5,
    teamName: '建筑学院',
    peopleCount: 20,
    groupId: 'C'
  },
  {
    id: 6,
    teamName: '环境学院',
    peopleCount: 20,
    groupId: 'C'
  },
])

// 更新小组信息
const updateGroups = () => {
  // 清空小组队伍
  groups.value.forEach(group => {
    group.teams = []
  })
  
  // 按小组分配队伍
  teams.value.forEach(team => {
    const group = groups.value.find(g => g.id === team.groupId)
    if (group) {
      group.teams.push(team)
    }
  })
}

// 赛事数据
const eventData = ref([
  {
    id: 1,
    teamId: 1,
    teamName: '计算机学院',
    peopleCount: 20,
    eventName: '500米竞速',
    eventTime: '2023-11-10 10:00',
    finishTime: '2023-11-10 10:01:58', // 完赛时间
    finalResult: 118.5, // 秒
  },
  {
    id: 2,
    teamId: 2,
    teamName: '机械学院',
    peopleCount: 20,
    eventName: '500米竞速',
    eventTime: '2023-11-10 10:00',
    finishTime: '2023-11-10 10:02:02', // 完赛时间
    finalResult: 122.3, // 秒
  },
  {
    id: 3,
    teamId: 3,
    teamName: '商学院',
    peopleCount: 20,
    eventName: '500米竞速',
    eventTime: '2023-11-10 10:00',
    finishTime: '2023-11-10 10:02:05', // 完赛时间
    finalResult: 125.1, // 秒
  },
])



// 加载状态
const isLoading = ref(true)

// 新队伍
const newTeam = ref({
  teamName: '',
  peopleCount: 0,
  groupId: 'A' // 默认小组
})

// 新赛事数据
const newEventData = ref({
  teamName: '',
  peopleCount: 0,
  eventName: '',
  eventTime: '',
  finishTime: '',
  finalResult: '',
})

// 加载赛事数据
const loadEventData = () => {
  isLoading.value = true
  setTimeout(() => {
    // 赛事信息已通过props传入，这里可以进行额外的数据加载
    isLoading.value = false
  }, 1000)
}





// 添加队伍
const addTeam = () => {
  if (!newTeam.value.teamName) {
    alert('请输入队伍名称')
    return
  }
  if (!newTeam.value.peopleCount || newTeam.value.peopleCount <= 0) {
    alert('请输入有效的人数')
    return
  }

  const team = {
    id: Date.now(),
    ...newTeam.value,
  }

  teams.value.push(team)



  // 更新小组信息
  updateGroups()

  // 重置表单
  newTeam.value = {
    teamName: '',
    peopleCount: 0,
    groupId: 'A'
  }
}

// 移除队伍
const removeTeam = (teamId) => {
  // 从队伍列表中移除
  teams.value = teams.value.filter((team) => team.id !== teamId)

  // 更新小组信息
  updateGroups()
}



// 添加赛事数据
const addEventData = () => {
  if (!newEventData.value.teamName) {
    alert('请输入队伍名称')
    return
  }
  if (!newEventData.value.eventName) {
    alert('请输入赛事名称')
    return
  }
  if (!newEventData.value.eventTime) {
    alert('请输入比赛时间')
    return
  }
  if (!newEventData.value.finishTime) {
    alert('请输入完赛时间')
    return
  }
  if (!newEventData.value.finalResult) {
    alert('请输入最终成绩')
    return
  }

  const data = {
    id: Date.now(),
    ...newEventData.value,
    finalResult: parseFloat(newEventData.value.finalResult) || 0,
  }

  eventData.value.push(data)

  // 重置表单
  newEventData.value = {
    teamName: '',
    peopleCount: 0,
    eventName: '',
    eventTime: '',
    finishTime: '',
    finalResult: '',
  }
}

// 移除赛事数据
const removeEventData = (id) => {
  eventData.value = eventData.value.filter((data) => data.id !== id)
}

// 排序赛事数据（按最终成绩升序）
const sortedEventData = computed(() => {
  return [...eventData.value].sort((a, b) => a.finalResult - b.finalResult)
})

// 定义获取当前数据的方法，供父组件调用
defineExpose({
  // 获取水上运动赛事数据的方法
  getWaterSportsData() {
    // 这里返回水上运动特有的数据结构
    return {
      // teams, events, results等水上运动特有的数据
      sportType: 'water_sports'
    }
  }
})



// 组件实例引用
const fieldDescRefs = ref([])

// 生命周期钩子
onMounted(() => {
  loadEventData()
  // 初始化小组信息
  updateGroups()
  setTimeout(() => {
    fieldDescRefs.value.forEach(ref => ref?.init?.())
  }, 100)
})

onUnmounted(() => {
  fieldDescRefs.value.forEach(ref => ref?.cleanup?.())
})
</script>

<template>
  <div class="water-sports-data-collection">

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 数据录入内容 -->
    <div v-else class="collection-content">
  

      <!-- 赛制和阶段选择 -->
      <div class="form-section tournament-section">
        <h2 class="section-title">赛制设置</h2>
        
        <div class="tournament-settings">
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">赛制类型</label>
              <select v-model="eventInfo.tournamentType" class="form-control">
                <option v-for="type in tournamentTypes" :key="type.value" :value="type.value">
                  {{ type.label }}
                </option>
              </select>
            </div>
            <div class="form-group">
              <label class="form-label">当前阶段</label>
              <select v-model="eventInfo.currentStage" class="form-control">
                <option v-for="stage in tournamentStages" :key="stage.value" :value="stage.value">
                  {{ stage.label }}
                </option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- 参赛队伍录入 -->
      <div class="form-section teams-section">
        <h2 class="section-title">参赛队伍</h2>

        <!-- 添加新队伍 -->
        <div class="add-team-form">
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">队伍名称</label>
              <input
                v-model="newTeam.teamName"
                type="text"
                class="form-control"
                placeholder="请输入队伍名称"
              />
            </div>
            <div class="form-group">
              <label class="form-label">参赛人数</label>
              <input
                v-model.number="newTeam.peopleCount"
                type="number"
                min="1"
                class="form-control"
                placeholder="请输入参赛人数"
              />
            </div>
            <div class="form-group">
              <label class="form-label">所属小组</label>
              <select v-model="newTeam.groupId" class="form-control">
                <option v-for="group in groups" :key="group.id" :value="group.id">
                  {{ group.name }}
                </option>
              </select>
            </div>
            <button @click="addTeam" class="add-team-btn">添加</button>
          </div>
        </div>

        <!-- 队伍列表 -->
        <div class="teams-list">
          <div v-if="teams.length === 0" class="empty-list">暂无参赛队伍，请添加</div>
          <div v-else>
            <!-- 小组赛模式下按小组分组显示 -->
            <div v-if="eventInfo.tournamentType !== 'knockout'" class="groups-container">
              <div v-for="group in groups" :key="group.id" class="group-section">
                <h3 class="group-title">{{ group.name }}</h3>
                <div class="teams-grid">
                  <div v-for="team in teams.filter(t => t.groupId === group.id)" :key="team.id" class="team-card">
                    <div class="team-info">
                      <span class="team-name">{{ team.teamName }}</span>
                      <span class="people-count">{{ team.peopleCount }}人</span>
                    </div>
                    <button @click="removeTeam(team.id)" class="remove-team-btn">删除</button>
                  </div>
                  <div v-if="teams.filter(t => t.groupId === group.id).length === 0" class="empty-group">该组暂无队伍</div>
                </div>
              </div>
            </div>
            
            <!-- 淘汰赛模式下直接显示所有队伍 -->
            <div v-else class="teams-grid">
              <div v-for="team in teams" :key="team.id" class="team-card">
                <div class="team-info">
                  <span class="team-name">{{ team.teamName }}</span>
                  <span class="people-count">{{ team.peopleCount }}人</span>
                </div>
                <button @click="removeTeam(team.id)" class="remove-team-btn">删除</button>
              </div>
            </div>
          </div>
        </div>
      </div>



      <!-- 赛事数据录入 -->
      <div class="form-section event-data-section">
        <h2 class="section-title">赛事数据录入</h2>

        <!-- 添加新赛事数据 -->
        <div class="add-event-data-form">
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">队伍名称</label>
              <input
                v-model="newEventData.teamName"
                type="text"
                class="form-control"
                placeholder="请输入队伍名称"
              />
            </div>
            <div class="form-group">
              <label class="form-label">参赛人数</label>
              <input
                v-model.number="newEventData.peopleCount"
                type="number"
                min="1"
                class="form-control"
                placeholder="请输入参赛人数"
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-group">
              <label class="form-label">赛事名称</label>
              <input
                v-model="newEventData.eventName"
                type="text"
                class="form-control"
                placeholder="如：500米竞速"
              />
            </div>
            <div class="form-group">
              <label class="form-label">比赛时间</label>
              <input v-model="newEventData.eventTime" type="datetime-local" class="form-control" />
            </div>
            <div class="form-group">
              <label class="form-label">完赛时间</label>
              <input v-model="newEventData.finishTime" type="datetime-local" class="form-control" />
            </div>
            <div class="form-group">
              <label class="form-label">最终成绩(秒)</label>
              <input
                v-model="newEventData.finalResult"
                type="number"
                step="0.01"
                min="0"
                class="form-control"
                placeholder="请输入成绩"
              />
            </div>
            <button @click="addEventData" class="add-event-data-btn">添加</button>
          </div>
        </div>

        <!-- 赛事数据列表 -->
        <div class="event-data-list">
          <div v-if="eventData.length === 0" class="empty-list">暂无赛事数据，请添加</div>
          <div v-else class="event-data-table-container">
            <table class="event-data-table">
              <thead>
                <tr>
                  <th>排名</th>
                  <th>队伍名称</th>
                  <th>参赛人数</th>
                  <th>赛事名称</th>
                  <th>比赛时间</th>
                  <th>完赛时间</th>
                  <th>最终成绩(秒)</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(data, index) in sortedEventData" :key="data.id">
                  <td>{{ index + 1 }}</td>
                  <td>
                    <input v-model="data.teamName" type="text" class="table-input" />
                  </td>
                  <td>
                    <input
                      v-model.number="data.peopleCount"
                      type="number"
                      min="1"
                      class="table-input small"
                    />
                  </td>
                  <td>
                    <input v-model="data.eventName" type="text" class="table-input" />
                  </td>
                  <td>
                    <input
                      v-if="editingId === data.id"
                      v-model="data.eventTime"
                      type="datetime-local"
                      class="table-input"
                    />
                    <span v-else>{{ data.eventTime }}</span>
                  </td>
                  <td>
                    <input
                      v-if="editingId === data.id"
                      v-model="data.finishTime"
                      type="datetime-local"
                      class="table-input"
                    />
                    <span v-else>{{ data.finishTime }}</span>
                  </td>
                  <td>
                    <input
                      v-if="editingId === data.id"
                      v-model.number="data.finalResult"
                      type="number"
                      step="0.01"
                      min="0"
                      class="table-input small"
                    />
                    <span v-else>{{ data.finalResult }}</span>
                  </td>
                  <td>
                    <button @click="removeEventData(data.id)" class="remove-data-btn">删除</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- 保存操作由父组件BaseDataCollection统一处理 -->
  </div>
</template>

<style scoped>
.water-sports-data-collection {
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

/* 分项选择器样式 */
.category-selector {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
}

.category-btn {
  flex: 1;
  padding: 12px;
  border: 2px solid #ddd;
  background-color: white;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s;
}

.category-btn.active {
  background-color: #1890ff;
  color: white;
  border-color: #1890ff;
}

.current-category-info {
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
  font-size: 14px;
  color: #666;
}

.current-category-info strong {
  color: #333;
}

/* 表单样式 */
.form-row {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
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

.add-team-btn,
.add-event-data-btn {
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

/* 队伍列表样式 */
.teams-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.team-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
  border: 1px solid #eee;
}

.team-info {
  flex: 1;
}

.team-name {
  display: block;
  font-weight: 500;
  color: #333;
  margin-bottom: 4px;
}

.people-count {
  font-size: 14px;
  color: #666;
}

.remove-team-btn {
  padding: 6px 12px;
  background-color: #f5f5f5;
  color: #f5222d;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.empty-list {
  text-align: center;
  color: #999;
  padding: 40px;
  background-color: #fafafa;
  border-radius: 8px;
}

/* 分段成绩表格样式 */
.segments-table-container,
.event-data-table-container {
  max-height: 400px;
  overflow-y: auto;
  margin-top: 16px;
}

.segments-table,
.event-data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.segments-table th,
.segments-table td,
.event-data-table th,
.event-data-table td {
  padding: 12px 8px;
  text-align: center;
  border-bottom: 1px solid #f0f0f0;
}

.segments-table th,
.event-data-table th {
  background-color: #fafafa;
  font-weight: 500;
  color: #333;
  position: sticky;
  top: 0;
  white-space: nowrap;
}

.segment-time-input {
  width: 80px;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
}

.team-cell {
  font-weight: 500;
  color: #333;
}

.total-time {
  font-weight: 500;
  color: #1890ff;
}

/* 赛制设置样式 */
.tournament-settings {
  margin-top: 16px;
}

/* 小组展示样式 */
.groups-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.group-section {
  padding: 16px;
  background-color: #fafafa;
  border-radius: 8px;
  border: 1px solid #eee;
}

.group-title {
  font-size: 16px;
  font-weight: 500;
  color: #1890ff;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e6f7ff;
}

.empty-group {
  text-align: center;
  color: #999;
  padding: 20px;
  background-color: white;
  border-radius: 4px;
  border: 1px dashed #ddd;
  grid-column: 1 / -1;
}

/* 赛事数据表样式 */
.table-input {
  width: 100%;
  padding: 6px 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  text-align: center;
}

.table-input.small {
  max-width: 80px;
}

.remove-data-btn {
  padding: 4px 8px;
  background-color: #f5f5f5;
  color: #f5222d;
  border: none;
  border-radius: 4px;
  font-size: 12px;
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

  .category-selector {
    flex-direction: column;
    gap: 8px;
  }

  .form-row {
    flex-direction: column;
    gap: 16px;
  }

  .form-group {
    min-width: auto;
  }

  .add-team-btn,
  .add-event-data-btn {
    align-self: stretch;
  }

  .teams-grid {
    grid-template-columns: 1fr;
  }

  .segments-table,
  .event-data-table {
    font-size: 12px;
  }

  .segments-table th,
  .segments-table td,
  .event-data-table th,
  .event-data-table td {
    padding: 8px 4px;
  }

  .segment-time-input,
  .table-input.small {
    width: 60px;
    padding: 4px;
  }

  .segments-table-container,
  .event-data-table-container {
    overflow-x: auto;
    max-height: 300px;
  }

  .segments-table,
  .event-data-table {
    min-width: 700px; /* 增加最小宽度以容纳完赛时间列 */
  }
}
</style>
