<template>
  <div class="PasswordChange">
    <div class="tobber">
      <van-nav-bar title="修改密码" left-arrow @click-left="onClickLeft" />
    </div>

    <!-- 输入手机号 -->
    <input type="text" placeholder="请输入学号" v-model="myId" />
    <!-- 输入验证码 -->
    <!-- <div class="getcode">
        <input type="text" placeholder="请输入验证码">

       </div> -->
    <!-- 输入旧密码 -->
    <input type="password" placeholder="请输入旧密码" v-model="oldPassword" />
    <!-- 输入新密码 -->
    <input type="password" placeholder="请输入新密码" v-model="newPassword" />
    <!-- 确认新密码 -->
    <input type="password" placeholder="请确认新密码" v-model="confirmPassword" />
    <!-- 修改密码按钮 -->
    <!-- <van-button type="primary" block font-size="4.3478vw">修改密码</van-button> -->
    <button class="mybtn" @click="changePassword">修改密码</button>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
const router = useRouter()
const onClickLeft = () => {
  router.back()
}
const changePassword = () => {
  if (myId.value == '') {
    showToast({
      message: '请输入学号',
      duration: 1000,
      position: 'top',
      className: 'mytoast',
    })
    return
  }
  if (oldPassword.value == '') {
    showToast({
      message: '请输入旧密码',
      duration: 1000,
      position: 'top',
      className: 'mytoast',
    })
    return
  }
  if (newPassword.value == '') {
    showToast({
      message: '请输入新密码',
      duration: 1000,
      position: 'top',
      className: 'mytoast',
    })
    return
  }
  if (confirmPassword.value != newPassword.value) {
    showToast({
      message: '两次输入密码不一致',
      duration: 1000,
      position: 'top',
      className: 'mytoast',
    })
    return
  }
}
// 学号
const myId = ref('')
const oldPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
</script>
<style scoped>
.PasswordChange {
  padding: 4.8309vw;
  /* background-color: #f5f5f5; */
}
.tobber {
  margin-bottom: 4.8309vw;
}
input {
  width: 100%;
  height: 13.285vw;
  /* border-radius: 2.8986vw; */
  font-size: 3.8647vw;
  /* margin-top: 4.8309vw; */
  padding: 2.4155vw 2.4155vw;
  position: relative;
  border: none;
  border-bottom: 0.2415vw solid #ccc;
  /* border-radius: 1.2077vw; */
  /* margin-top: 4.8309vw; */
}
/* .getcode{
  position: relative;
}
.getcode:after{
  content: '获取验证码';
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 2.4155vw;
  /* width: 2.4155vw;
  height: 2.4155vw; */

/* background-color: #007AFF; */
/* border-radius: 50%;} */

/* 按钮样式 */
.mybtn {
  width: 50.7246vw;
  height: 10.8696vw;
  border-radius: 2.8986vw;
  font-size: 4.3478vw;
  margin-top: 12.0773vw;
  margin-left: 16.9082vw;
  background-color: #007aff;
  color: #fff;
}
</style>
