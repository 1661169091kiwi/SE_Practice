<script setup>
import { ref } from 'vue'

// 组件属性
defineProps({
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

/* 提示图标 - 改进样式和交互 */
.description-icon {
  width: 24px;
  height: 24px;
  border-radius: var(--border-radius-full);
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-hover) 100%);
  color: var(--text-white);
  border: none;
  font-size: var(--font-size-sm);
  font-weight: 700;
  cursor: help;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: var(--spacing-sm);
  flex-shrink: 0;
  transition: all var(--transition-normal);
  box-shadow: var(--shadow-sm);
  user-select: none;
}

/* 图标悬停效果 */
.description-icon:hover {
  transform: translateY(-2px) scale(1.1);
  box-shadow: var(--shadow-md);
  background: linear-gradient(135deg, var(--primary-hover) 0%, var(--primary-active) 100%);
}

/* 图标点击效果 */
.description-icon:active {
  transform: translateY(0) scale(0.95);
  box-shadow: var(--shadow-sm);
}

/* Tooltip - 改进样式和动画 */
.tooltip {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  background-color: rgba(0, 0, 0, 0.92);
  color: var(--text-white);
  padding: var(--spacing-md) var(--spacing-lg);
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  line-height: 1.5;
  white-space: normal;
  max-width: 320px;
  margin-top: var(--spacing-sm);
  z-index: var(--z-tooltip);
  box-shadow: var(--shadow-lg);
  backdrop-filter: blur(10px);
  /* 淡入淡出动画 */
  animation: fadeInScale 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  pointer-events: none; /* 防止点击tooltip */
}

/* Tooltip 淡入缩放动画 */
@keyframes fadeInScale {
  from {
    opacity: 0;
    transform: translateX(-50%) translateY(10px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateX(-50%) translateY(0) scale(1);
  }
}

/* Tooltip 箭头 */
.tooltip::before {
  content: '';
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  border-width: 6px;
  border-style: solid;
  border-color: transparent transparent rgba(0, 0, 0, 0.92) transparent;
  animation: fadeIn 0.3s var(--transition-bezier);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .tooltip {
    max-width: 280px;
    font-size: var(--font-size-xs);
    padding: var(--spacing-sm) var(--spacing-md);
  }
  
  .description-icon {
    width: 22px;
    height: 22px;
    font-size: var(--font-size-xs);
  }
}

@media (max-width: 480px) {
  .tooltip {
    max-width: calc(100vw - 40px);
    left: 0;
    transform: translateX(0);
  }
  
  .tooltip::before {
    left: 20px;
    transform: none;
  }
}

/* 高对比度模式支持 */
@media (prefers-contrast: high) {
  .tooltip {
    background-color: #000000;
    color: #ffffff;
    border: 1px solid #ffffff;
  }
  
  .tooltip::before {
    border-color: transparent transparent #000000 transparent;
  }
}

/* 减少动画模式支持 */
@media (prefers-reduced-motion: reduce) {
  .tooltip,
  .tooltip::before {
    animation: none !important;
  }
  
  .description-icon {
    transition: none !important;
  }
}
</style>
