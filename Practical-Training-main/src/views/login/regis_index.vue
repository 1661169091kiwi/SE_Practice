<template>
 <div class="registre-continer">
   <!-- 顶部导航栏 -->
   <div class="regis">
     <van-nav-bar title="注册" left-arrow @click="onClickBack" />
   </div>

  <!-- 表单区域 -->
   <div class="myform">
    <van-form @submit="onSubmit">
  <van-cell-group inset>
    <van-field
      v-model="username"
      name="myId"
      label="学号"
      placeholder="学号"
      :rules="[{ required: true, message: '请填写学号' },
      {pattern: /^[0-9]{8}$/, message: '请输入正确的8位学号'}
      ]"
    />
<!-- <van-field
    v-model="sms"
    center
    clearable
    label="短信验证码"
    placeholder="请输入短信验证码"
  >
    <template #button>
      <van-button size="small" type="primary">发送验证码</van-button>
    </template>
  </van-field> -->

    <van-field
      v-model="password"
      type="password"
      name="myPassword"
      label="密码"
      placeholder="密码"
      :rules="[{ required: true, message: '请填写密码' },
      {pattern: /^[0-9a-zA-Z]{6,12}$/, message: '请输入6-12位密码'}
      ]"
    />
    <van-field
      v-model="newpassword"
      type="password"
      name="newpassword"
      label="确认密码"
      placeholder="请确认密码"
      :rules="[{ required: true, message: '请填写确认密码' },
      {pattern: /^[0-9a-zA-Z]{6,12}$/, message: '请输入6-12位密码'},
      {validator: (val) => val === password.valueOf(), message: '两次输入密码不一致'}
      ]"
    />
    <van-field
      v-model="myName"
      type="text"
      name="myName"
      label="姓名"
      placeholder="请输入姓名"
      :rules="[{ required: true, message: '请填写姓名' }]"
    />
    <van-field
  v-model="fieldValue"
  is-link
  readonly
  label="学院"
  placeholder="选择学院"
  @click="showPicker = true"
/>
<van-popup v-model:show="showPicker" round position="bottom">
  <van-picker
    :columns="columns"
    @cancel="showPicker = false"
    @confirm="onConfirm"
  />
</van-popup>

  </van-cell-group>
  <div style="margin: 16px;">
    <van-button round block type="primary" native-type="submit" @click="onSubmit">
      提交
    </van-button>
  </div>
</van-form>

   </div>
 </div>

</template>

<script setup lang="ts">
import router from '@/router';
import { showToast } from 'vant';

// import { ref } from 'vue';

// 顶部导航栏
const onClickBack = () => {
  // 跳转到登录页面
  router.back();
};

//表单部分
import { ref } from 'vue';

const username = ref('');
const password = ref('');
const newpassword = ref('');
const myName = ref('');

const onSubmit = async () => {
  if (fieldValue.value === '') {
    showToast({
      message: '请选择学院',
      position:'top',
      icon: 'none',
      duration: 2000,
      className:'mytoast',
      
    })
    return;
  }
}


// 学院选择器
const columns = [
      { text: '软件工程学院', value: '软件工程学院' },
      { text: '人工智能学院', value: '人工智能学院' },
      { text: '微电子科学与技术学院', value: '微电子学院' },
      { text: '中国语言文学系(珠海)', value: '中国语言文学系(珠海)' },
      { text: '历史学系(珠海)', value: '历史学系(珠海)' },
      { text: '哲学系(珠海)', value: '哲学系(珠海)' },
      { text: '国际金融学院', value: '国际金融学院' },
      { text: '国际翻译学院', value: '国际翻译学院' },
      { text: '国际关系学院', value: '国际关系学院' },
      { text: '旅游学院', value: '旅游学院' },
      { text: '数学学院(珠海)', value: '数学学院(珠海)' },
      { text: '物理与天文学院', value: '物理与天文学院' },
      { text: '大气科学学院', value: '大气科学学院' },
      { text: '海洋科学学院', value: '海洋科学学院' },
      { text: '地球科学与工程学院', value: '地球科学与工程学院' },
      { text: '化学工程与技术学院', value: '化学工程与技术学院' },
      { text: '海洋工程与技术学院', value: '海洋工程与技术学院' },
      { text: '中法核工程与技术学院', value: '中法核工程与技术学院' },
      { text: '土木工程学院', value: '土木工程学院' },
      { text: '遥感科学与技术学院', value: '遥感科学与技术学院' },
    ];
    const fieldValue = ref('');
    const showPicker = ref(false);

    const onConfirm = ({ selectedOptions }) => {
      showPicker.value = false;
      fieldValue.value = selectedOptions[0].text;
    };
</script>

<style scoped>
.myform {
  margin-top: 20px;
  width: 100%;
  padding: 10px;
}
.myform .van-field {
  margin-bottom: 10px;
  font-size: 16px;
  height: 60px;
}
</style>
