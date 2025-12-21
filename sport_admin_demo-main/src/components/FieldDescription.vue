<script setup>
import { ref } from 'vue'

// 组件属性
const props = defineProps({
  description: {
    type: String,
    required: true
  }
})

// 控制提示框显示
const isVisible = ref(false)

// 切换提示框显示状态
const toggleVisibility = (event) => {
  event.stopPropagation() // 阻止冒泡，防止点击提示图标时关闭提示框
  isVisible.value = !isVisible.value
}

// 关闭提示框
const closeTooltip = () => {
  isVisible.value = false
}

// 点击提示框外部关闭提示框
const handleClickOutside = (event) => {
  // 确保点击的不是提示图标或提示框本身
  const target = event.target
  if (!target.closest('.field-description')) {
    closeTooltip()
  }
}

// 组件挂载时添加点击事件监听
defineExpose({
  init: () => {
    document.addEventListener('click', handleClickOutside)
  },
  cleanup: () => {
    document.removeEventListener('click', handleClickOutside)
  }
})
</script>

<template>
  <div class="field-description">
    <button 
      class="description-icon"
      @click="toggleVisibility"
      type="button"
      aria-label="字段说明"
    >
      ?
    </button>
    <div v-if="isVisible" class="tooltip">
      {{ description }}
    </div>
  </div>
</template>

<style scoped>
.field-description {
  position: relative;
  display: inline-block;
}

.description-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background-color: #1890ff;
  color: white;
  border: none;
  font-size: 12px;
  font-weight: bold;
  cursor: help;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 8px;
  flex-shrink: 0;
}

.tooltip {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(0, 0, 0, 0.9);
  color: white;
  padding: 12px;
  border-radius: 8px;
  font-size: 14px;
  line-height: 1.4;
  white-space: normal;
  max-width: 300px;
  margin-top: 8px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.tooltip::before {
  content: '';
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  border-width: 6px;
  border-style: solid;
  border-color: transparent transparent rgba(0, 0, 0, 0.9) transparent;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .tooltip {
    max-width: 250px;
    font-size: 13px;
  }
}
</style>