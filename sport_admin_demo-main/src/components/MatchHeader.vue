<script setup>
import { ref, onMounted } from 'vue'

// 组件属性
const props = defineProps({
  event: {
    type: Object,
    default: () => ({
      name: '',
      teamA: '',
      teamB: '',
      status: 'not_started',
      collectorName: '采集员'
    })
  }
})

// 网络状态
const isOffline = ref(false)

// 监听网络状态
const updateNetworkStatus = () => {
  isOffline.value = !navigator.onLine
}

onMounted(() => {
  window.addEventListener('online', updateNetworkStatus)
  window.addEventListener('offline', updateNetworkStatus)
  updateNetworkStatus() // 初始检查
})

// 组件卸载时清理事件监听
defineExpose({
  cleanup: () => {
    window.removeEventListener('online', updateNetworkStatus)
    window.removeEventListener('offline', updateNetworkStatus)
  }
})
</script>

<template>
  <div class="match-header">
    <!-- 离线提示 -->
    <div v-if="isOffline" class="offline-indicator">
      离线模式，数据暂存本地
    </div>
    
    <!-- 赛事信息 -->
    <div class="match-info">
      <h1 class="event-name">{{ event.name }}</h1>
      <div class="teams">
        <span class="team">{{ event.teamA }}</span>
        <span class="vs">VS</span>
        <span class="team">{{ event.teamB }}</span>
      </div>
      <div class="status-info">
        <span class="status" :class="event.status">
          {{ event.status === 'not_started' ? '未开始' : 
             event.status === 'in_progress' ? '进行中' : '已结束' }}
        </span>
        <span class="collector">采集员：{{ event.collectorName }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.match-header {
  background-color: #1890ff;
  color: white;
  position: sticky;
  top: 0;
  z-index: 100;
}

.offline-indicator {
  background-color: #f5222d;
  color: white;
  text-align: center;
  padding: 4px;
  font-size: 12px;
  font-weight: 500;
}

.match-info {
  padding: 16px 20px;
}

.event-name {
  font-size: 16px;
  font-weight: 500;
  margin: 0 0 12px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.teams {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 12px;
  padding: 12px 0;
  background-color: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.team {
  font-size: 18px;
  font-weight: 600;
  min-width: 60px;
  text-align: center;
}

.vs {
  font-size: 14px;
  opacity: 0.8;
}

.status-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.status {
  padding: 4px 8px;
  border-radius: 4px;
  background-color: rgba(255, 255, 255, 0.2);
}

.status.not_started {
  background-color: #faad14;
}

.status.in_progress {
  background-color: #52c41a;
}

.status.ended {
  background-color: #666;
}

.collector {
  opacity: 0.9;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .teams {
    flex-direction: column;
    gap: 8px;
  }
  
  .event-name {
    font-size: 15px;
  }
  
  .status-info {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }
}
</style>