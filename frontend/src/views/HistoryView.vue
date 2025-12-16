<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MatchHeader from '../components/MatchHeader.vue'

const router = useRouter()

// 赛事信息（仅用于头部显示）
const eventInfo = ref({
  id: '',
  name: '历史记录',
  teamA: '',
  teamB: '',
  status: '',
  collectorName: '采集员姓名'
})

// 历史记录数据
const historyData = ref([])

// 筛选和分页参数
const filterOptions = ref({
  sportType: 'all', // all, football, basketball, badminton, volleyball, water_sports
  dateRange: '',    // 日期范围
  status: 'all'     // all, completed, pending, rejected
})

const pagination = ref({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

// 加载状态
const isLoading = ref(true)
const isRefreshing = ref(false)

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

// 获取状态名称和样式
const getStatusInfo = (status) => {
  const statusMap = {
    completed: { name: '已完成', class: 'status-completed' },
    pending: { name: '待审核', class: 'status-pending' },
    rejected: { name: '已驳回', class: 'status-rejected' },
    approved: { name: '已通过', class: 'status-approved' }
  }
  return statusMap[status] || { name: status, class: '' }
}

// 加载历史记录
const loadHistoryData = (refresh = false) => {
  if (refresh) {
    isRefreshing.value = true
  } else {
    isLoading.value = true
  }
  
  // 模拟API请求
  setTimeout(() => {
    const mockData = [
      {
        id: '1001',
        eventName: '计算机学院 vs 机械学院',
        sportType: 'football',
        date: '2023-11-15',
        submitTime: '2023-11-15 16:30:45',
        status: 'completed',
        eventId: '101'
      },
      {
        id: '1002',
        eventName: '商学院 vs 外语学院',
        sportType: 'basketball',
        date: '2023-11-14',
        submitTime: '2023-11-14 20:15:30',
        status: 'pending',
        eventId: '201'
      },
      {
        id: '1003',
        eventName: '数学学院 vs 物理学院',
        sportType: 'badminton',
        date: '2023-11-13',
        submitTime: '2023-11-13 15:45:20',
        status: 'approved',
        eventId: '301'
      },
      {
        id: '1004',
        eventName: '化工学院 vs 材料学院',
        sportType: 'volleyball',
        date: '2023-11-12',
        submitTime: '2023-11-12 17:20:10',
        status: 'rejected',
        rejectReason: '数据不完整，请补充比赛事件信息',
        eventId: '401'
      },
      {
        id: '1005',
        eventName: '龙舟比赛 - 小组赛',
        sportType: 'water_sports',
        date: '2023-11-10',
        submitTime: '2023-11-10 12:30:50',
        status: 'completed',
        eventId: '501'
      },
      {
        id: '1006',
        eventName: '自动化学院 vs 电子学院',
        sportType: 'football',
        date: '2023-11-09',
        submitTime: '2023-11-09 19:45:30',
        status: 'approved',
        eventId: '102'
      },
      {
        id: '1007',
        eventName: '医学院 vs 文学院',
        sportType: 'basketball',
        date: '2023-11-08',
        submitTime: '2023-11-08 18:15:20',
        status: 'completed',
        eventId: '202'
      },
      {
        id: '1008',
        eventName: '法学院 vs 管理学院',
        sportType: 'badminton',
        date: '2023-11-07',
        submitTime: '2023-11-07 16:30:15',
        status: 'pending',
        eventId: '302'
      },
      {
        id: '1009',
        eventName: '环境学院 vs 生命学院',
        sportType: 'volleyball',
        date: '2023-11-06',
        submitTime: '2023-11-06 17:45:05',
        status: 'completed',
        eventId: '402'
      },
      {
        id: '1010',
        eventName: '皮艇比赛 - 半决赛',
        sportType: 'water_sports',
        date: '2023-11-05',
        submitTime: '2023-11-05 11:30:40',
        status: 'approved',
        eventId: '502'
      },
      {
        id: '1011',
        eventName: '建筑学院 vs 土木学院',
        sportType: 'football',
        date: '2023-11-04',
        submitTime: '2023-11-04 16:20:35',
        status: 'rejected',
        rejectReason: '比分数据有误，请核对',
        eventId: '103'
      },
      {
        id: '1012',
        eventName: '软件学院 vs 通信学院',
        sportType: 'basketball',
        date: '2023-11-03',
        submitTime: '2023-11-03 19:10:25',
        status: 'completed',
        eventId: '203'
      }
    ]
    
    historyData.value = mockData
    pagination.value.total = mockData.length
    
    isLoading.value = false
    isRefreshing.value = false
  }, 1500)
}

// 筛选后的数据
const filteredData = computed(() => {
  let data = [...historyData.value]
  
  // 按运动类型筛选
  if (filterOptions.value.sportType !== 'all') {
    data = data.filter(item => item.sportType === filterOptions.value.sportType)
  }
  
  // 按日期范围筛选
  if (filterOptions.value.dateRange) {
    const dates = filterOptions.value.dateRange.split(' 至 ')
    if (dates.length === 2) {
      data = data.filter(item => {
        const itemDate = new Date(item.date)
        const startDate = new Date(dates[0])
        const endDate = new Date(dates[1])
        return itemDate >= startDate && itemDate <= endDate
      })
    }
  }
  
  // 按状态筛选
  if (filterOptions.value.status !== 'all') {
    data = data.filter(item => item.status === filterOptions.value.status)
  }
  
  return data
})

// 分页后的数据
const paginatedData = computed(() => {
  const startIndex = (pagination.value.currentPage - 1) * pagination.value.pageSize
  const endIndex = startIndex + pagination.value.pageSize
  return filteredData.value.slice(startIndex, endIndex)
})

// 刷新数据
const handleRefresh = () => {
  loadHistoryData(true)
}

// 筛选变更
const handleFilterChange = () => {
  pagination.value.currentPage = 1
  // 在实际应用中可能需要重新加载数据
}

// 分页变更
const handlePageChange = (page) => {
  pagination.value.currentPage = page
}

// 查看详情
const handleViewDetails = (item) => {
  // 根据运动类型跳转到对应的数据录入页面
  const routes = {
    football: '/football-data/',
    basketball: '/basketball-data/',
    badminton: '/badminton-data/',
    volleyball: '/volleyball-data/',
    water_sports: '/water-sports-data/'
  }
  
  const routePath = routes[item.sportType] || '/data-collection/'
  router.push(`${routePath}${item.eventId}`)
}

// 预览数据
const handlePreview = (item) => {
  router.push(`/data-preview/${item.sportType}/${item.eventId}`)
}

// 编辑数据（仅对已驳回和待审核状态）
const handleEdit = (item) => {
  if (item.status === 'rejected' || item.status === 'pending') {
    const routes = {
      football: '/football-data/',
      basketball: '/basketball-data/',
      badminton: '/badminton-data/',
      volleyball: '/volleyball-data/',
      water_sports: '/water-sports-data/'
    }
    
    const routePath = routes[item.sportType] || '/data-collection/'
    router.push(`${routePath}${item.eventId}`)
  }
}

// 生命周期钩子
onMounted(() => {
  loadHistoryData()
})
</script>

<template>
  <div class="history-view">
    <!-- 赛事信息头部 -->
    <MatchHeader :event="eventInfo" />
    
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">历史记录</h1>
      <button @click="handleRefresh" class="refresh-btn" :disabled="isRefreshing">
        <span v-if="!isRefreshing">刷新</span>
        <span v-else>刷新中...</span>
      </button>
    </div>
    
    <!-- 筛选栏 -->
    <div class="filter-section">
      <div class="filter-row">
        <div class="filter-item">
          <label class="filter-label">运动类型</label>
          <select v-model="filterOptions.sportType" @change="handleFilterChange" class="filter-select">
            <option value="all">全部</option>
            <option value="football">足球</option>
            <option value="basketball">篮球</option>
            <option value="badminton">羽毛球</option>
            <option value="volleyball">排球</option>
            <option value="water_sports">水上运动</option>
          </select>
        </div>
        
        <div class="filter-item">
          <label class="filter-label">日期范围</label>
          <input 
            v-model="filterOptions.dateRange"
            type="text"
            placeholder="选择日期范围"
            class="filter-input"
            @focus="$event.target.type='date'"
            @blur="$event.target.type='text'"
            @change="handleFilterChange"
          />
        </div>
        
        <div class="filter-item">
          <label class="filter-label">状态</label>
          <select v-model="filterOptions.status" @change="handleFilterChange" class="filter-select">
            <option value="all">全部</option>
            <option value="completed">已完成</option>
            <option value="pending">待审核</option>
            <option value="approved">已通过</option>
            <option value="rejected">已驳回</option>
          </select>
        </div>
      </div>
    </div>
    
    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>
    
    <!-- 历史记录列表 -->
    <div v-else class="history-list">
      <div v-if="filteredData.length === 0" class="empty-state">
        <div class="empty-icon">📋</div>
        <p>暂无符合条件的历史记录</p>
      </div>
      
      <div v-else class="history-cards">
        <div 
          v-for="item in paginatedData" 
          :key="item.id"
          class="history-card"
          @click="handleViewDetails(item)"
        >
          <div class="card-header">
            <h3 class="event-name">{{ item.eventName }}</h3>
            <span 
              class="status-badge"
              :class="getStatusInfo(item.status).class"
            >
              {{ getStatusInfo(item.status).name }}
            </span>
          </div>
          
          <div class="card-body">
            <div class="info-row">
              <span class="info-label">运动类型:</span>
              <span class="info-value">{{ getSportTypeName(item.sportType) }}</span>
            </div>
            
            <div class="info-row">
              <span class="info-label">比赛日期:</span>
              <span class="info-value">{{ item.date }}</span>
            </div>
            
            <div class="info-row">
              <span class="info-label">提交时间:</span>
              <span class="info-value">{{ item.submitTime }}</span>
            </div>
            
            <!-- 驳回原因 -->
            <div v-if="item.status === 'rejected' && item.rejectReason" class="reject-reason">
              <span class="reject-label">驳回原因:</span>
              <span class="reject-text">{{ item.rejectReason }}</span>
            </div>
          </div>
          
          <div class="card-actions">
            <button @click.stop="handlePreview(item)" class="action-btn preview-btn">
              预览
            </button>
            <button 
              v-if="item.status === 'rejected' || item.status === 'pending'"
              @click.stop="handleEdit(item)"
              class="action-btn edit-btn"
            >
              编辑
            </button>
          </div>
        </div>
      </div>
      
      <!-- 分页 -->
      <div v-if="filteredData.length > 0" class="pagination">
        <div class="pagination-info">
          共 {{ filteredData.length }} 条记录，第 {{ pagination.currentPage }} / {{ Math.ceil(filteredData.length / pagination.pageSize) }} 页
        </div>
        
        <div class="pagination-controls">
          <button 
            @click="handlePageChange(pagination.currentPage - 1)"
            :disabled="pagination.currentPage <= 1"
            class="page-btn"
          >
            上一页
          </button>
          
          <div class="page-numbers">
            <button
              v-for="page in Math.min(5, Math.ceil(filteredData.length / pagination.pageSize))"
              :key="page"
              @click="handlePageChange(page)"
              :class="['page-btn', 'page-number', { active: page === pagination.currentPage }]"
            >
              {{ page }}
            </button>
            
            <span v-if="Math.ceil(filteredData.length / pagination.pageSize) > 5" class="page-ellipsis">...</span>
          </div>
          
          <button 
            @click="handlePageChange(pagination.currentPage + 1)"
            :disabled="pagination.currentPage >= Math.ceil(filteredData.length / pagination.pageSize)"
            class="page-btn"
          >
            下一页
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-view {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
  padding-bottom: 40px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  background-color: white;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.refresh-btn {
  padding: 10px 20px;
  background-color: #1890ff;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.refresh-btn:hover:not(:disabled) {
  background-color: #40a9ff;
}

.refresh-btn:disabled {
  background-color: #d9d9d9;
  cursor: not-allowed;
}

/* 筛选栏样式 */
.filter-section {
  background-color: white;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.filter-row {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.filter-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 200px;
}

.filter-label {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.filter-select,
.filter-input {
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  background-color: white;
}

.filter-select:focus,
.filter-input:focus {
  outline: none;
  border-color: #1890ff;
}

/* 加载状态 */
.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 16px;
  color: #666;
  padding: 40px;
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

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  background-color: white;
  border-radius: 12px;
  margin: 0 20px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state p {
  color: #999;
  font-size: 16px;
  margin: 0;
}

/* 历史记录卡片 */
.history-list {
  flex: 1;
  padding: 0 20px;
}

.history-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.history-card {
  background-color: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
}

.history-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16px;
}

.event-name {
  font-size: 18px;
  font-weight: 500;
  color: #333;
  margin: 0;
  flex: 1;
  margin-right: 12px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.status-completed {
  background-color: #f6ffed;
  color: #52c41a;
  border: 1px solid #b7eb8f;
}

.status-pending {
  background-color: #fff7e6;
  color: #fa8c16;
  border: 1px solid #ffd591;
}

.status-approved {
  background-color: #e6f7ff;
  color: #1890ff;
  border: 1px solid #91d5ff;
}

.status-rejected {
  background-color: #fff2f0;
  color: #f5222d;
  border: 1px solid #ffccc7;
}

.card-body {
  margin-bottom: 16px;
}

.info-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  font-size: 14px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  color: #999;
  margin-right: 8px;
  min-width: 70px;
}

.info-value {
  color: #333;
  font-weight: 500;
}

.reject-reason {
  margin-top: 12px;
  padding: 12px;
  background-color: #fff2f0;
  border: 1px solid #ffccc7;
  border-radius: 6px;
  font-size: 14px;
}

.reject-label {
  color: #f5222d;
  font-weight: 500;
  margin-right: 8px;
}

.reject-text {
  color: #333;
  word-break: break-word;
}

.card-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.action-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.preview-btn {
  background-color: #e6f7ff;
  color: #1890ff;
}

.preview-btn:hover {
  background-color: #bae7ff;
}

.edit-btn {
  background-color: #fff7e6;
  color: #fa8c16;
}

.edit-btn:hover {
  background-color: #ffe7ba;
}

/* 分页样式 */
.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 20px;
  padding: 20px;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.pagination-info {
  font-size: 14px;
  color: #666;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.page-btn {
  padding: 8px 16px;
  border: 1px solid #ddd;
  background-color: white;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.page-btn:hover:not(:disabled) {
  border-color: #1890ff;
  color: #1890ff;
}

.page-btn:disabled {
  background-color: #f5f5f5;
  color: #d9d9d9;
  cursor: not-allowed;
}

.page-number {
  padding: 8px 12px;
  min-width: 40px;
  text-align: center;
}

.page-number.active {
  background-color: #1890ff;
  color: white;
  border-color: #1890ff;
}

.page-ellipsis {
  padding: 0 8px;
  color: #999;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 16px;
  }
  
  .page-title {
    font-size: 20px;
  }
  
  .refresh-btn {
    width: 100%;
  }
  
  .filter-section {
    padding: 16px;
  }
  
  .filter-row {
    flex-direction: column;
    gap: 16px;
  }
  
  .filter-item {
    min-width: auto;
  }
  
  .history-list {
    padding: 0 16px;
  }
  
  .history-cards {
    grid-template-columns: 1fr;
    gap: 16px;
  }
  
  .history-card {
    padding: 16px;
  }
  
  .event-name {
    font-size: 16px;
  }
  
  .pagination {
    flex-direction: column;
    align-items: center;
    gap: 16px;
    padding: 16px;
  }
  
  .pagination-info {
    order: 2;
  }
  
  .pagination-controls {
    order: 1;
    flex-wrap: wrap;
    justify-content: center;
  }
  
  .page-btn {
    padding: 6px 12px;
    font-size: 12px;
  }
  
  .page-number {
    min-width: 36px;
    padding: 6px 10px;
  }
}
</style>