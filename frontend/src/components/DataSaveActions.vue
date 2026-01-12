<script setup>
// 组件属性
const props = defineProps({
  isSubmitting: {
    type: Boolean,
    default: false
  }
})

// 组件事件
const emit = defineEmits(['submit-server', 'finish-match'])

// 提交服务器
const handleSubmitServer = () => {
  if (!props.isSubmitting) {
    emit('submit-server')
  }
}

// 自定义确认弹窗
import { ref } from 'vue'
const showFinishConfirm = ref(false)
const openFinishConfirm = () => {
  if (!props.isSubmitting) showFinishConfirm.value = true
}
const cancelFinish = () => {
  showFinishConfirm.value = false
}
const confirmFinish = () => {
  showFinishConfirm.value = false
  emit('finish-match')
}
</script>

<template>
  <div class="data-save-actions">
    <button 
      class="action-button finish-match-button"
      @click="openFinishConfirm"
      :disabled="isSubmitting"
    >
      结束比赛
    </button>
    <button 
      class="action-button submit-server-button"
      @click="handleSubmitServer"
      :disabled="isSubmitting"
    >
      <span v-if="isSubmitting">提交中...</span>
      <span v-else>提交服务器</span>
    </button>
    <!-- 自定义确认弹窗 -->
    <Teleport to="body">
      <div v-if="showFinishConfirm" class="modal-overlay" @click.self="cancelFinish">
        <div class="modal-card">
          <div class="modal-header">结束比赛</div>
          <div class="modal-body">
            确定要结束这场比赛吗？结束后的比赛将进入“已结束”列表。此操作不可撤销。
          </div>
          <div class="modal-actions">
            <button class="modal-btn cancel" @click="cancelFinish">取消</button>
            <button class="modal-btn confirm" @click="confirmFinish">确认结束</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
.data-save-actions {
  background-color: var(--background-primary);
  padding: var(--spacing-xl) var(--spacing-2xl);
  border-top: 2px solid var(--border-color);
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 100;
  box-shadow: var(--shadow-md);
  backdrop-filter: blur(5px);
  transition: all var(--transition-normal);
}

.action-button {
  padding: var(--spacing-md) var(--spacing-2xl);
  border: none;
  border-radius: var(--border-radius-lg);
  font-size: var(--font-size-base);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-bezier);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-sm);
}

/* 按钮发光效果 */
.action-button::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: var(--border-radius-full);
  background-color: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%);
  transition: width var(--transition-slow), height var(--transition-slow);
}

.action-button:hover::before {
  width: 300px;
  height: 300px;
}

.submit-server-button {
  background: linear-gradient(135deg, var(--success-color) 0%, var(--success-hover) 100%);
  color: var(--text-white);
  min-width: 220px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  box-shadow: var(--shadow-md);
}

.submit-server-button:hover:not(:disabled) {
  background: linear-gradient(135deg, var(--success-hover) 0%, var(--success-color) 100%);
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.finish-match-button {
  background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
  color: white;
  min-width: 150px;
  margin-right: 16px;
  box-shadow: var(--shadow-md);
}

.finish-match-button:hover:not(:disabled) {
  background: linear-gradient(135deg, #ff7875 0%, #ff4d4f 100%);
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.submit-server-button:active:not(:disabled) {
  transform: translateY(0);
  background-color: var(--success-active);
  box-shadow: var(--shadow-sm);
}

/* 提交中状态动画 */
.submit-server-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  background-color: var(--success-color);
  box-shadow: var(--shadow-sm);
}

/* 添加加载动画 */
.submit-server-button:disabled span::after {
  content: '';
  display: inline-block;
  width: 16px;
  height: 16px;
  margin-left: 8px;
  border: 2px solid transparent;
  border-top-color: var(--text-white);
  border-radius: 50%;
  animation: spin var(--transition-normal) linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* 脉冲效果 - 用于吸引注意 */
@keyframes pulse {
  0% {
    box-shadow: var(--shadow-md);
  }
  50% {
    box-shadow: var(--shadow-lg);
  }
  100% {
    box-shadow: var(--shadow-md);
  }
}

/* 移动端适配 - 增强响应式体验 */
@media (max-width: 768px) {
  .data-save-actions {
    flex-direction: row;
    padding: 12px 16px;
    gap: 12px;
  }
  
  .submit-server-button {
    width: auto;
    flex: 1;
    min-width: unset;
    padding: 10px;
    font-size: 14px;
  }

  .finish-match-button {
    width: auto;
    flex: 1;
    margin-right: 0;
    min-width: unset;
    padding: 10px;
    font-size: 14px;
  }
  
  .action-button {
    border-radius: var(--border-radius-md);
  }
}

/* 自定义弹窗样式 */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.modal-card {
  width: min(480px, 90vw);
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 12px 32px rgba(0,0,0,0.2);
  overflow: hidden;
}
.modal-header {
  padding: 14px 16px;
  font-weight: 700;
  font-size: 16px;
  color: var(--text-primary);
  background: rgba(0,0,0,0.04);
}
.modal-body {
  padding: 16px;
  font-size: 14px;
  color: var(--text-secondary);
}
.modal-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 12px 16px 16px;
}
.modal-btn {
  padding: 8px 16px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
}
.modal-btn.cancel {
  background: #f5f5f5;
  color: #333;
}
.modal-btn.confirm {
  background: var(--primary-color);
  color: #fff;
}

/* 平板适配 */
@media (min-width: 769px) and (max-width: 1024px) {
  .submit-server-button {
    min-width: 200px;
    font-size: var(--font-size-sm);
  }
}
</style>
