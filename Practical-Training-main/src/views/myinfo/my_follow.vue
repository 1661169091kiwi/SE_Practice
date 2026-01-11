<template>
  <div class="MyFollow">
    <div class="tobber">
      <van-nav-bar title="我的关注" left-arrow @click-left="onClickLeft" />
    </div>
    <GameDisplay />
  </div>
</template>
<script setup>
import { useRouter } from 'vue-router'
// // import { ref } from 'vue'
const router = useRouter()
const onClickLeft = () => {
  router.back()
}
</script>
<style scoped></style>
