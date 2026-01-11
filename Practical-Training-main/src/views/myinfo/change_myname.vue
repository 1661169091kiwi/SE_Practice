<script setup>
import { useRouter } from 'vue-router'
import { ref, useTemplateRef } from 'vue'
const newName = ref('')
const router = useRouter()
const mybutton = useTemplateRef('mybutton')
// 点击返回按钮
const onClickLeft = () => {
  router.back()
}
// 处理输入框变化
const handleChange = () => {
  if (newName.value.trim() !== '') {
    mybutton.value.classList.add('active')
  } else {
    mybutton.value.classList.remove('active')
  }
}
// 点击修改姓名按钮
const changeName = () => {
  if (newName.value.trim() !== '') {
    // 发送修改姓名请求
    // ...
    // 跳转到我的信息页面
    router.push('/my_information')
  }
}
</script>
<template>
  <div class="change-myname">
    <van-nav-bar title="修改姓名" left-arrow @click-left="onClickLeft" />
    <input type="text" placeholder="请输入姓名" v-model="newName" @input="handleChange" />
    <button @click="changeName" ref="mybutton" class="mybutton">修改姓名</button>
  </div>
</template>

<style scoped>
input {
  display: block;
  width: 80%;
  height: 9.6618vw;
  margin: 7.2464vw auto;
  padding: 0 2.4155vw;
  border: 0.2415vw solid #000;
  border-radius: 4.8309vw;
}
.mybutton {
  display: block;
  width: 43.4783vw;
  height: 9.6618vw;
  margin: 60.3865vw auto;
  margin-bottom: 24.1546vw;
  border-radius: 4.8309vw;
  background-color: #ccc;
  color: #000;
}
.active {
  background-color: #1989fa;
  color: #fff;
}
</style>
