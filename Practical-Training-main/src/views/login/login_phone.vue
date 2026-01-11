<template>
  <div class="login_password">
    <div class="navbar">
      <van-nav-bar
  title="手机短信登录"
  left-arrow
  @click-left="onClickLeft"
/>
    </div>
    <div class="login_form">
      <van-form @submit="onSubmit">
  <van-cell-group inset>
    <van-field
      v-model="username"
      name="手机号"
      label="手机号"
      placeholder="手机号"
      :rules="[{ required: true, message: '请输入手机号' }]"
    />
    <van-field
    v-model="sms"
    center
    clearable
    label="短信验证码"
    placeholder="请输入短信验证码"
    :rules="[{ required: true, message: '请输入短信验证码' }]"
  >
    <template #button>
      <van-button size="small" type="primary">发送验证码</van-button>
    </template>
  </van-field>
  </van-cell-group>
  <div style="margin: 16px;">
    <van-button round block type="primary" native-type="submit">
      登录
    </van-button>
  </div>
</van-form>

</div>
<div class="bottom_text">
    <div class="gotopassword" @click="onClickPassword">密码登录</div>|<div class="goto_register" @click="onClickRegister">去注册</div>
</div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import router from '@/router';

    // 顶部导航栏区域
    const onClickLeft = () => history.back();
    // const onClickRight = () => {};

    // 表单区域
    const username = ref('');
    const sms = ref('');
    const onSubmit = (values) => {
      console.log('submit', values);
    };

    // 底部文字区域
    const onClickRegister = () => {
      // 跳转到注册页面
      router.push('/register');
    };
    const onClickPassword = () => {
      // 跳转到密码登录页面
      router.push('/login_password');
    };
</script>

<style scoped>
 .navbar {
  width: 100%;
 }

.login_form {
  margin-top: 20px;
  width: 100%;
  padding: 10px;
}
.van-field {
  margin-bottom: 10px;
  font-size: 16px;
}
.van-button {
  font-size: 16px;
}

/* 底部文字区域 */
.bottom_text {
  /* margin-bottom: 50px; */
  /* display: absolute; */
  position: absolute;
  bottom: 80px;
  margin-left: 50%;
  transform: translateX(-50%);
  font-size: 14px;
}
.goto_register, .gotopassword {
  display: inline-block;
  color: #4095E5;
  margin: 0 5px;
}
</style>