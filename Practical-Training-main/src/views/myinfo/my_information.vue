<template>
  <div class="MyInformation">
    <!-- 导航栏 -->
    <div class="tobber">
      <van-nav-bar title="我的资料" left-arrow @click-left="onClickLeft" />
    </div>

    <!-- 我的资料内容 -->
    <div class="myinfo_content">
      <!-- 头像 -->
      <div class="mypic" @click="changeAvatar">
        <span>头像</span>
        <img src="@/assets/default-avatar.png" alt="" />
        <van-icon name="arrow" />
      </div>

      <!-- 昵称 -->
      <div class="myname mystyle" @click="changeName">
        <span>姓名</span>
        <!-- <img src="@/assets/default-avatar.png" alt=""> -->
        <p>Yog_Sothoth</p>
        <van-icon name="arrow" />
      </div>

      <!-- 性别 -->
      <div class="mysex mystyle" @click="showsexpicker = true">
        <span>性别</span>
        <!-- <img src="@/assets/default-avatar.png" alt=""> -->
        <p v-if="checked === '1'">男</p>
        <p v-else>女</p>
        <van-icon name="arrow" />
      </div>

      <!-- 手机号 -->
      <div class="myphone mystyle">
        <span>学号</span>
        <!-- <img src="@/assets/default-avatar.png" alt=""> -->
        <p>23331031</p>
        <van-icon name="arrow" />
      </div>

      <!-- 学院 -->
      <div class="myid mystyle" @click="showPicker = true">
        <span>学院</span>

        <p>{{ selectedCollege }}</p>
        <van-icon name="arrow" />
      </div>

      <!-- 我的身份 -->
      <div class="myidentity mystyle">
        <span>身份</span>
        <!-- <img src="@/assets/default-avatar.png" alt=""> -->
        <p>观众</p>
        <van-icon name="arrow" />
      </div>
    </div>
    <!-- 修改密码按钮 -->
    <button class="mybtn changepassword" @click="changePassword">修改密码</button>
    <!-- 退出登录按钮 -->
    <button class="mybtn exitlogin" @click="exitLogin">退出登录</button>

    <van-popup v-model:show="showPicker" round position="bottom">
      <van-picker :columns="columns" @cancel="showPicker = false" @confirm="onConfirm" />
    </van-popup>
    <!-- 性别选择器 -->
    <div class="changesex" v-if="showsexpicker">
      <van-radio-group v-model="checked" direction="horizontal" @change="onChange">
        <van-radio name="1" label-disabled>男</van-radio>
        <van-radio name="2" label-disabled>女</van-radio>
      </van-radio-group>
    </div>
  </div>
</template>
<script setup>
import { useRouter } from 'vue-router'
import { ref } from 'vue'

const router = useRouter()
const onClickLeft = () => {
  router.back()
}
// 退出登录
const exitLogin = () => {
  // 清除用户登录状态
  localStorage.removeItem('token')
  // 跳转到登录页面
  router.push('/login_password')
}
// 跳转到修改密码页面
const changePassword = () => {
  router.push('/password_change')
}
//直接修改头像
const changeAvatar = () => {
  // 直接让用户选择图片文件
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  // 监听文件选择事件
  input.addEventListener('change', (e) => {
    const file = e.target.files[0]
    if (file) {
      // 预览图片
      const reader = new FileReader()
      reader.readAsDataURL(file)
      reader.onload = () => {
        // 更新头像预览
        document.querySelector('.mypic img').src = reader.result
      }
    }
  })
  input.click()
}
// 跳转到修改姓名页面
const changeName = () => {
  router.push('/change_myname')
}

const selectedCollege = ref('软件工程学院')
//更改学院
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
]
const showPicker = ref(false)

const onConfirm = ({ selectedOptions }) => {
  showPicker.value = false
  selectedCollege.value = selectedOptions[0].text
}
// 性别选择器
const showsexpicker = ref(false)
const checked = ref('1')
const onChange = ({ name }) => {
  showsexpicker.value = false
}
</script>
<style scoped>
.tobber {
  margin: 4.8309vw 0 0 0;
}

/* 我的资料 */

/* 我的资料头像 */
.mypic {
  width: 100%;
  height: 14.4928vw;
  padding: 2.4155vw 1.2077vw;
  position: relative;
  border-bottom: 0.2415vw solid #ccc;
  border-radius: 1.2077vw;
  margin-top: 4.8309vw;
  /* display: flex;
    justify-content: flex-start;
    align-items: center; */
}
.mypic span {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 2.4155vw;
  font-size: 3.8647vw;
  /* font-weight: bold; */
}
.mypic img {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 7.2464vw;
  width: 9.6618vw;
  height: 9.6618vw;
  border-radius: 2.8986vw;
}
.mypic .van-icon {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 2.4155vw;
}
/* 我的资料 */
.mystyle {
  width: 100%;
  height: 14.4928vw;
  padding: 2.4155vw 1.2077vw;
  position: relative;
  border-bottom: 0.2415vw solid #ccc;
  border-radius: 1.2077vw;
  margin-top: 4.8309vw;
  /* display: flex;
    justify-content: flex-start;
    align-items: center; */
}
.mystyle span {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  left: 2.4155vw;
  font-size: 3.8647vw;
  /* font-weight: bold; */
}
.mystyle .van-icon {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 2.4155vw;
}
.mystyle p {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 7.2464vw;
  font-size: 3.8647vw;
  /* font-weight: bold; */
}
/* 我的资料修改密码按钮 */
.mybtn {
  width: 50.7246vw;
  height: 10.8696vw;
  border-radius: 2.8986vw;
  font-size: 4.3478vw;
  margin-left: 24.1546vw;
}
/* 我的资料修改密码按钮 */
.changepassword {
  margin-top: 36.2319vw;
  /* background-color: #007AFF;
    color: #fff; */
}

/* 我的资料退出登录按钮 */
.exitlogin {
  /* margin-left: 0vw; */
  margin-top: 12.0773vw;
  color: #fff;
  background-color: #de868f;
}

.changesex {
  /* width: 36.2319vw; */
  height: 24.1546vw;
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
  right: 50%;
  font-size: 3.8647vw;
  /* font-weight: bold; */
}
:deep(.van-radio__icon) {
  width: 9.6618vw; /* 根据需要调整大小 */
  height: 9.6618vw;
  font-size: 4.8309vw; /* 调整内部圆点大小 */
}

:deep(.van-radio__label) {
  font-size: 7.2464vw; /* 调整文字大小以匹配 */
}
</style>
