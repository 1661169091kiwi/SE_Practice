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
      collectorName: '采集员',
    }),
  },
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
  },
})
</script>

<template>
  <div class="match-header">
    <!-- 离线提示 -->
    <div v-if="isOffline" class="offline-indicator">离线模式，数据暂存本地</div>

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
          {{
            event.status === 'not_started'
              ? '未开始'
              : event.status === 'in_progress'
                ? '进行中'
                : '已结束'
          }}
        </span>
        <span class="collector">采集员：{{ event.collectorName }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.match-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
}

/* 离线提示 - 改进动画效果 */
.offline-indicator {
  background-color: var(--danger-color);
  color: var(--text-white);
  text-align: center;
  padding: 6px;
  font-size: 13px;
  font-weight: 500;
  animation: pulse 2s infinite;
  transition: all 0.3s ease;
}

/* 脉冲动画 */
@keyframes pulse {
  0% {
    opacity: 0.8;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.8;
  }
}

.match-info {
  padding: 20px 24px;
  transition: all 0.3s ease;
}

/* 赛事名称 - 增强视觉效果 */
.event-name {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 16px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

/* 队伍信息 - 改进布局和样式 */
.teams {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 16px;
  padding: 16px;
  background-color: rgba(255, 255, 255, 0.12);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.team {
  font-size: 20px;
  font-weight: 700;
  min-width: 80px;
  text-align: center;
  padding: 8px 16px;
  border-radius: 8px;
  background-color: rgba(255, 255, 255, 0.15);
  transition: all 0.3s ease;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.team:hover {
  transform: translateY(-2px);
  background-color: rgba(255, 255, 255, 0.2);
}

.vs {
  font-size: 16px;
  font-weight: 600;
  opacity: 0.9;
  color: var(--text-white);
}

/* 状态信息 - 改进样式和视觉层次 */
.status-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.status {
  padding: 6px 12px;
  border-radius: 16px;
  font-weight: 600;
  text-transform: uppercase;
  font-size: 13px;
  letter-spacing: 0.5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  animation: fadeIn 0.5s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.status.not_started {
  background-color: var(--warning-color);
  color: var(--text-white);
}

.status.in_progress {
  background-color: var(--success-color);
  color: var(--text-white);
  animation: blink 2s infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0.8;
  }
}

.status.ended {
  background-color: #666;
  color: var(--text-white);
}

.collector {
  opacity: 0.9;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.collector::before {
  content: '👤';
  font-size: 16px;
}

/* 过渡动画 */
.match-header * {
  transition: all var(--transition-bezier);
}

/* 移动端适配 - 改进响应式设计 */
@media (max-width: 768px) {
  .match-info {
    padding: 16px 20px;
  }

  .teams {
    flex-direction: column;
    gap: 12px;
    padding: 12px;
  }

  .team {
    font-size: 18px;
    min-width: 100%;
  }

  .event-name {
    font-size: 16px;
    text-align: center;
  }

  .status-info {
    flex-direction: column;
    gap: 10px;
    align-items: center;
  }

  .vs {
    font-size: 14px;
  }
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1024px) {
  .teams {
    gap: 16px;
  }

  .team {
    font-size: 19px;
    min-width: 70px;
  }
}
</style>
