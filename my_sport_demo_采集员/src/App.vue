<script setup></script>

<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<style>
.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 50%, #f1f5f9 100%);
  color: var(--text-primary);
}

/* 页面切换的淡入效果，让整体过渡更柔和 */
.app-container > * {
  animation: app-fade-in 0.35s var(--transition-bezier);
}

@keyframes app-fade-in {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
