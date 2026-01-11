<template>
  <div class="InfoIndex">
    <!-- 我的信息 -->
    <div class="myinfo_container" @click="toInfoDetail">
      <div class="myinfo">
        <div class="info_img">
          <img src="@/assets/default-avatar.png" alt="" />
        </div>
        <div class="info_text">
          <div class="myname">杜楠</div>
          <div class="myidentity">观众</div>
        </div>
      </div>
      <div class="myinfoicon">
        <van-icon name="arrow" size="32" />
      </div>
    </div>

    <!-- 我的订阅 -->
    <div class="myfollows" @click="toFollows">
      <div class="myicon">
        <!-- <van-icon name="star-o" size="34"  color="#FCCA00"/> -->
        <svg
          t="1762138014460"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="22455"
          width="34"
          height="34"
        >
          <path
            d="M832.8 63.9H191.2c-17.8 0-32.3 14.5-32.3 32.3V878c0 23.3 23.9 38.9 45.3 29.6L489.8 782l331.4 128.4c21.2 8.2 44-7.4 44-30.1V96.2c-0.1-17.9-14.5-32.3-32.4-32.3z"
            fill="#FCCA00"
            p-id="22456"
          ></path>
        </svg>
      </div>
      <div class="myfollows_text">我的订阅</div>
      <div class="myfollowsicon">
        <van-icon name="arrow" size="32" />
      </div>
    </div>

    <!-- 历史浏览 -->
    <div class="myfollows" @click="toHistory">
      <div class="myicon">
        <!-- <van-icon name="star-o" size="34"  color="#FCCA00"/> -->
        <svg
          t="1762136707887"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="8677"
          width="34"
          height="34"
        >
          <path
            d="M511.850073 0.299912C229.132901 0.299912 0.000029 229.432783 0.000029 512.149956s229.132871 511.850044 511.850044 511.850044 511.850044-229.132871 511.850044-511.850044S794.467275 0.299912 511.850073 0.299912z m204.740018 561.935371H522.247027c-3.399004 0.699795-6.798008 1.099678-10.396954 1.099677-28.291711 0-51.185004-22.893293-51.185004-51.185004V205.03993c0-28.291711 22.893293-51.185004 51.185004-51.185005 28.291711 0 51.185004 22.893293 51.185005 51.185005V459.865274h153.555013c28.291711 0 51.185004 22.893293 51.185004 51.185004 0 28.291711-22.993264 51.185004-51.185004 51.185005z"
            fill="#BD3124"
            p-id="8678"
          ></path>
        </svg>
      </div>
      <div class="myfollows_text">历史浏览</div>
      <div class="myfollowsicon">
        <van-icon name="arrow" size="32" />
      </div>
    </div>

    <!-- 我的评分 -->
    <div class="myfollows" @click="toScore">
      <div class="myicon">
        <!-- <van-icon name="star-o" size="34"  color="#FCCA00"/> -->
        <svg
          t="1762138099127"
          class="icon"
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
          p-id="27224"
          width="34"
          height="34"
        >
          <path
            d="M962 404.1c-3-8.9-10.7-15.4-20-16.8l-281.2-42.4L534.9 80.1c-4.1-8.7-12.9-14.3-22.6-14.3s-18.4 5.6-22.6 14.3L363.8 344.9 82.6 387.3c-9.3 1.4-17 7.9-20 16.8s-0.7 18.8 5.9 25.5l204 206.4-48.2 291.6c-1.6 9.4 2.4 18.9 10.2 24.5 7.8 5.5 18.1 6.1 26.4 1.6l251.3-137.1 251.3 137.1c3.7 2 7.9 3.1 12 3.1 5.1 0 10.2-1.6 14.5-4.6 7.8-5.5 11.7-15 10.2-24.5L752 635.9l204-206.4c6.7-6.6 9-16.5 6-25.4z"
            fill="#93D2F3"
            p-id="27225"
          ></path>
        </svg>
      </div>
      <div class="myfollows_text">我的评分</div>
      <div class="myfollowsicon">
        <van-icon name="arrow" size="32" />
      </div>
    </div>
  </div>
</template>

<script setup>
import router from '@/router'
// import My_rating from './my_rating.vue'
// 点击跳转我的信息详情页
const toInfoDetail = () => {
  router.push({
    path: '/my_information',
  })
}
// 点击跳转我的订阅页
const toFollows = () => {
  router.push({
    path: '/my_follow',
  })
}
// 点击跳转我的评分页
const toScore = () => {
  router.push({
    path: '/my_rating',
  })
}
// 点击跳转我的历史浏览页
const toHistory = () => {
  router.push({
    path: '/my_history',
  })
}
</script>
<style scoped>
/* // 我的信息页面 */
.myinfo_container {
  width: 100%;
  height: 35.2657vw;
  /* background-color: #f5f5f5; */
  display: flex;
  /* flex-direction: column; */
  align-items: center;
  justify-content: space-between;
}
.myinfo {
  display: flex;
  /* flex-direction: column; */
  align-items: center;
  justify-content: start;
  margin-left: 4.8309vw;
}
.info_img img {
  width: 19.3237vw;
  height: 19.3237vw;
  border-radius: 4.8309vw;
}
.info_text {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-left: 4.8309vw;
}
.myname {
  font-size: 6.7633vw;
  font-weight: bold;
}
.myidentity {
  font-size: 4.3478vw;
  font-weight: 300;
  font-family: '宋体';
}
.myinfoicon {
  margin-right: 4.8309vw;
}

/* 我的订阅 */
.myfollows {
  width: 100%;
  height: 24.1546vw;
  background-color: #f5f5f5;
  position: relative;
  margin-top: 4.8309vw;
}
.myicon {
  position: absolute;
  top: 50%;
  left: 4.8309vw;
  transform: translateY(-50%);
}
.myfollows_text {
  position: absolute;
  top: 50%;
  left: 24.1546vw;
  transform: translateY(-50%);
  font-size: 6.7633vw;
}
.myfollowsicon {
  position: absolute;
  top: 50%;
  right: 4.8309vw;
  transform: translateY(-50%);
}
</style>
