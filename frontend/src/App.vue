<script setup></script>

<template>
  <div class="app-container">
    <router-view />
  </div>
</template>

<style>
.app-container {
  min-height: 100vh;
  background-color: var(--background-secondary);
  color: var(--text-primary);
}

/* 页面切换的淡入效果，让整体过渡更柔和 */
.app-container > * {
  animation: app-fade-in 0.25s var(--transition-bezier);
}

@keyframes app-fade-in {
  from {
    opacity: 0;
    transform: translateY(4px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
