<script setup>
import { ref, onMounted, onUnmounted, computed, defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'
import MatchHeader from '../components/MatchHeader.vue'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'
import { get, post } from '../utils/http'

import FootballDataCollection from './FootballDataCollection.vue'
// 动态导入其余运动项目的数据采集组件
const BasketballDataCollection = defineAsyncComponent(
  () => import('./BasketballDataCollection.vue'),
)
const BadmintonDataCollection = defineAsyncComponent(() => import('./BadmintonDataCollection.vue'))
const VolleyballDataCollection = defineAsyncComponent(
  () => import('./VolleyballDataCollection.vue'),
)
const WaterSportsDataCollection = defineAsyncComponent(
  () => import('./WaterSportsDataCollection.vue'),
)

const route = useRoute()

// 赛事信息
const eventInfo = ref({
  id: '',
  name: '赛事名称',
  teamA: '主队',
  teamB: '客队',
  status: 'not_started',
  collectorName: '张三',
})

// 加载状态
const isLoading = ref(true)
const isSaving = ref(false)
const isSubmitting = ref(false)

const loadEventData = async () => {
  isLoading.value = true
  const sportType = route.params.sportType
  const eventId = route.params.eventId
  try {
    const res = await get(`/collector/${sportType}/events/${eventId}`)
    if (res.code === 200) {
      const data = res.data
      eventInfo.value = {
        id: data.id,
        name: data.name,
        teamA: data.teamA,
        teamB: data.teamB,
        teamAId: data.teamAId,
        teamBId: data.teamBId,
        status: data.status,
        collectorName: '当前用户',
        currentData: data.currentData,
        autoFinishAt: data.autoFinishAt,
        candidatesTeamA: data.candidatesTeamA || [],
        candidatesTeamB: data.candidatesTeamB || []
      }
    }
  } catch (error) {
    eventInfo.value = {
      id: eventId || '1',
      name: `${sportType}赛事数据采集 (Offline/Mock)`,
      teamA: '主队名称',
      teamB: '客队名称',
      status: 'in_progress',
      collectorName: '采集员姓名',
    }
  } finally {
    isLoading.value = false
  }
}

 

// 提交服务器
const performSubmit = async (statusOverride = null) => {
  isSubmitting.value = true
  try {
    const sportType = route.params.sportType
    const sportTypeNormalized = normalizeSportType(sportType)
    const eventId = route.params.eventId
    const url = `/collector/${sportType}/events/${eventId}/data`
    
    // 根据运动类型从子组件获取对应的数据
    let sportData = {}
    if (currentCollectionRef.value) {
      switch (sportTypeNormalized) {
        case 'football':
          sportData = currentCollectionRef.value.getFootballData ? currentCollectionRef.value.getFootballData() : {};
          break;
        case 'basketball':
          sportData = currentCollectionRef.value.getBasketballData ? currentCollectionRef.value.getBasketballData() : {};
          break;
        case 'badminton':
          sportData = currentCollectionRef.value.getBadmintonData ? currentCollectionRef.value.getBadmintonData() : {};
          break;
        case 'volleyball':
          sportData = currentCollectionRef.value.getVolleyballData ? currentCollectionRef.value.getVolleyballData() : {};
          break;
        case 'water':
          sportData = currentCollectionRef.value.getWaterSportsData ? currentCollectionRef.value.getWaterSportsData() : {};
          break;
        default:
          sportData = {};
      }
    }
    
    const finalStatus = statusOverride || eventInfo.value.status

    // 确保 eventId 是数字
    const numericEventId = typeof eventInfo.value.id === 'string' ? parseInt(eventInfo.value.id, 10) : eventInfo.value.id

    // 构建提交数据
    const normalizeIn = (p) => ({
      id: String(p.id ?? ''),
      name: (p.name ?? '').trim(),
      isStarting: !!p.isStarting,
      number: String(p.number ?? ''),
      position: (p.position ?? '').trim()
    })
    const dedupArr = (arr) => {
      const out = []
      const seen = new Set()
      for (const p of Array.isArray(arr) ? arr : []) {
        const np = normalizeIn(p)
        const key = np.id ? `id:${np.id}` : `ext:${np.name.toLowerCase()}|${np.number.trim()}|${np.position.toLowerCase()}`
        if (seen.has(key)) continue
        seen.add(key)
        out.push(np)
      }
      return out
    }
    if (sportData && sportData.lineups) {
      const rawTA = Array.isArray(sportData.lineups.teamA) ? sportData.lineups.teamA : []
      const rawTB = Array.isArray(sportData.lineups.teamB) ? sportData.lineups.teamB : []
      const hasEmptyName = rawTA.some(p => !String(p?.name ?? '').trim()) || rawTB.some(p => !String(p?.name ?? '').trim())
      if (hasEmptyName) {
        alert('球员姓名不能为空')
        isSubmitting.value = false
        return
      }
      const ta = dedupArr(rawTA)
      const tb = dedupArr(rawTB)
      sportData.lineups = { teamA: ta, teamB: tb }
    }
    const dataToSubmit = {
      eventId: numericEventId,
      eventInfo: eventInfo.value,
      sportType: sportType,
      timestamp: new Date().toISOString(),
      collectorName: eventInfo.value.collectorName,
      status: finalStatus,
      // 包含从子组件获取的具体运动数据
      data: sportData
    }
    
    const res = await post(url, dataToSubmit)
    if (res && res.code && res.code !== 200) {
      throw new Error(res.message || '提交失败')
    }
    
    if (statusOverride === 'finished') {
      alert('比赛已结束，感谢您的工作！')
      eventInfo.value.status = 'finished'
      await loadEventData()
      // 可以在这里跳转回列表页
      // router.push('/collector/events')
    } else {
      alert('数据已成功提交')
      await loadEventData()
    }
    
    isSubmitting.value = false
  } catch (error) {
    console.error('提交数据失败:', error)
    alert('提交失败，请重试')
    isSubmitting.value = false
  }
}

const handleSubmitServer = () => performSubmit()

const handleFinishMatch = () => performSubmit('finished')

// 动态组件映射
const componentsMap = {
  football: FootballDataCollection,
  basketball: BasketballDataCollection,
  badminton: BadmintonDataCollection,
  volleyball: VolleyballDataCollection,
  water: WaterSportsDataCollection,
  'water_sports': WaterSportsDataCollection,
  'water-sports': WaterSportsDataCollection,
}

// 计算当前应显示的组件
const normalizeSportType = (t) => {
  const s = String(t || '').toLowerCase().trim()
  if (s === 'water_sports' || s === 'water-sports' || s === 'watersports') return 'water'
  return s
}
const CurrentCollectionComponent = computed(() => {
  const sportType = normalizeSportType(route.params.sportType)
  return componentsMap[sportType] || null
})

// 组件实例引用
const matchHeaderRef = ref(null)
const fieldDescRefs = ref([])
const currentCollectionRef = ref(null)

// 生命周期钩子
onMounted(() => {
  loadEventData()
  // 初始化字段说明组件
  setTimeout(() => {
    fieldDescRefs.value.forEach((ref) => ref?.init?.())
  }, 100)
})

onUnmounted(() => {
  // 清理事件监听
  matchHeaderRef.value?.cleanup?.()
  fieldDescRefs.value.forEach((ref) => ref?.cleanup?.())
})
</script>

<template>
  <div class="data-collection-container">
    <MatchHeader ref="matchHeaderRef" :event="eventInfo" in-progress-label="采集中" />

    <!-- 数据录入内容区域 -->
    <div class="collection-content">
      <!-- 动态加载对应的具体运动项目数据采集组件 -->
      <component 
        v-if="CurrentCollectionComponent"
        ref="currentCollectionRef"
        :is="CurrentCollectionComponent" 
        :event-id="eventInfo.id" 
        :event-info="eventInfo"
      />
      <div v-else style="padding: 16px; color: #666;">
        未能加载采集组件，请刷新或检查运动类型：{{ (route.params && route.params.sportType) || '' }}
      </div>
      <!-- 加载覆盖层，不卸载子组件，避免本地录入数据被重置 -->
      <div v-if="isLoading" class="loading-overlay">
        <div class="loading-container">
          <div class="loading-spinner"></div>
          <p>加载中...</p>
        </div>
      </div>
    </div>

    <!-- 保存操作按钮 -->
    <DataSaveActions
      :is-submitting="isSubmitting"
      @submit-server="handleSubmitServer"
      @finish-match="handleFinishMatch"
    />
  </div>
</template>

<style scoped>
.data-collection-container {
  min-height: 100vh;
  background-color: var(--background-secondary);
  display: flex;
  flex-direction: column;
  padding-bottom: 80px; /* 为底部按钮留出空间 */
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

.loading-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(255, 255, 255, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
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
  padding: 32px 40px 120px;
  max-width: 1440px;
  width: 100%;
  margin: 0 auto;
}

.placeholder-content {
  background-color: white;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  color: #666;
}

.placeholder-content h2 {
  margin-bottom: 16px;
  color: #333;
}

.demo-field {
  margin-top: 32px;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.field-label {
  display: flex;
  align-items: center;
  font-weight: 500;
  color: #333;
}

.demo-input {
  width: 100%;
  max-width: 300px;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
}

.demo-input:focus {
  outline: none;
  border-color: #1890ff;
}

/* 通用表单样式 */
.form-section {
  background-color: var(--background-primary);
  border-radius: var(--border-radius-lg);
  padding: 24px 28px;
  margin-bottom: 24px;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
}

.section-title {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 2px solid var(--border-color);
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: var(--text-primary);
}

.form-control {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-size: 16px;
  box-sizing: border-box;
  background-color: var(--background-primary);
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px rgba(42, 122, 226, 0.18);
}

/* 移动端适配 */
@media (max-width: 768px) {
  .collection-content {
    padding: 20px 16px 100px;
  }

  .form-section {
    padding: 16px;
  }

  .demo-input {
    max-width: none;
  }
}
</style>
