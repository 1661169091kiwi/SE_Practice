<script setup>
// 组件属性
const props = defineProps({
  isSubmitting: {
    type: Boolean,
    default: false
  }
})

// 组件事件
const emit = defineEmits(['submit-server'])

// 提交服务器
const handleSubmitServer = () => {
  if (!props.isSubmitting) {
    emit('submit-server')
  }
}
</script>

<template>
  <div class="data-save-actions">
    <button 
      class="action-button submit-server-button"
      @click="handleSubmitServer"
      :disabled="isSubmitting"
    >
      <span v-if="isSubmitting">提交中...</span>
      <span v-else>提交服务器</span>
    </button>
  </div>
</template>

<style scoped>
.data-save-actions {
  background-color: white;
  padding: 16px 20px;
  border-top: 1px solid #f0f0f0;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}

.action-button {
  padding: 12px 16px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s;
}

.submit-server-button {
  background-color: #52c41a;
  color: white;
  min-width: 200px;
}

.submit-server-button:hover:not(:disabled) {
  background-color: #73d13d;
}

.action-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .data-save-actions {
    flex-direction: column;
    padding: 12px 16px;
  }
  
  .submit-server-button {
    width: 100%;
  }
}
</style>