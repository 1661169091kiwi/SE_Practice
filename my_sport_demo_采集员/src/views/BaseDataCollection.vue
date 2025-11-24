<script setup>
import { ref, onMounted, onUnmounted, computed, defineAsyncComponent } from 'vue'
import { useRoute } from 'vue-router'
import MatchHeader from '../components/MatchHeader.vue'
import DataSaveActions from '../components/DataSaveActions.vue'
import FieldDescription from '../components/FieldDescription.vue'

// 动态导入各个运动项目的数据采集组件
const FootballDataCollection = defineAsyncComponent(() => import('./FootballDataCollection.vue'))
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

// 模拟加载赛事数据
const loadEventData = () => {
  isLoading.value = true
  // 模拟API请求延迟
  setTimeout(() => {
    // 根据路由参数加载不同赛事数据
    const sportType = route.params.sportType
    const eventId = route.params.eventId

    // 这里可以根据sportType和eventId加载具体数据
    eventInfo.value = {
      id: eventId || '1',
      name: `${
        sportType === 'football'
          ? '足球'
          : sportType === 'basketball'
            ? '篮球'
            : sportType === 'badminton'
              ? '羽毛球'
              : sportType === 'volleyball'
                ? '排球'
                : '水上运动'
      }赛事数据采集`,
      teamA: '主队名称',
      teamB: '客队名称',
      status: 'in_progress',
      collectorName: '采集员姓名',
    }

    isLoading.value = false
  }, 1000)
}

// 获取运动类型对应的API路径
const getApiPathBySportType = (sportType) => {
  const apiMap = {
    football: '/api/collector/football/data',
    basketball: '/api/collector/basketball/data',
    badminton: '/api/collector/badminton/data',
    volleyball: '/api/collector/volleyball/data',
    water: '/api/collector/watersports/data'
  }
  return apiMap[sportType] || '/api/collector/data'
}

// 提交服务器
const handleSubmitServer = async () => {
  isSubmitting.value = true
  try {
    const sportType = route.params.sportType
    const apiPath = getApiPathBySportType(sportType)
    
    // 根据运动类型从子组件获取对应的数据
    let sportData = {}
    if (currentCollectionRef.value) {
      switch (sportType) {
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
    
    // 构建提交数据
    const dataToSubmit = {
      eventId: eventInfo.value.id,
      eventInfo: eventInfo.value,
      sportType: sportType,
      timestamp: new Date().toISOString(),
      collectorName: eventInfo.value.collectorName,
      status: eventInfo.value.status,
      // 包含从子组件获取的具体运动数据
      data: sportData
    }
    
    console.log(`提交到API: ${apiPath}`, dataToSubmit)
    
    // 模拟API请求
    // 实际项目中应该使用axios或fetch进行真实的API调用
    // const response = await fetch(apiPath, {
    //   method: 'POST',
    //   headers: {
    //     'Content-Type': 'application/json'
    //   },
    //   body: JSON.stringify(dataToSubmit)
    // })
    // const result = await response.json()
    
    // 模拟提交延迟
    setTimeout(() => {
      alert(`数据已成功提交到 ${apiPath}`)
      isSubmitting.value = false
    }, 1500)
  } catch (error) {
    console.error('提交数据失败:', error)
    alert('提交失败，请重试')
    isSubmitting.value = false
  }
}

// 动态组件映射
const componentsMap = {
  football: FootballDataCollection,
  basketball: BasketballDataCollection,
  badminton: BadmintonDataCollection,
  volleyball: VolleyballDataCollection,
  water: WaterSportsDataCollection,
}

// 计算当前应显示的组件
const CurrentCollectionComponent = computed(() => {
  const sportType = route.params.sportType
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
    <!-- 赛事信息头部 -->
    <MatchHeader ref="matchHeaderRef" :event="eventInfo" />

    <!-- 加载状态 -->
    <div v-if="isLoading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载中...</p>
    </div>

    <!-- 数据录入内容区域 -->
    <div v-else class="collection-content">
      <!-- 动态加载对应的具体运动项目数据采集组件 -->
      <component 
        ref="currentCollectionRef"
        :is="CurrentCollectionComponent" 
        :event-id="eventInfo.id" 
        :event-info="eventInfo"
      />
    </div>

    <!-- 保存操作按钮 -->
    <DataSaveActions
      :is-submitting="isSubmitting"
      @submit-server="handleSubmitServer"
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
