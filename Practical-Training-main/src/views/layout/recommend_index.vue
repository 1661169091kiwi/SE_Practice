<template>
  <div class="recommend">
    <div class="search">
      <van-search
        v-model="value"
        placeholder="羽毛球"
        shape="round"
        background="#fff"
        border
        @click="goToSearch"
      />
      <!-- <van-button round type="primary"  @click="handleSearch">搜索</van-button> -->
      <button type="button" @click="handleSearch">搜索</button>
    </div>
    <!-- 导航栏 -->
    <div class="nav">
      <ul @click="handleClickclass">
        <li :class="{ active: currentindex == 1 }" data-index="1">推荐</li>
        <li :class="{ active: currentindex == 2 }" data-index="2">足球</li>
        <li :class="{ active: currentindex == 3 }" data-index="3">篮球</li>
        <li :class="{ active: currentindex == 4 }" data-index="4">羽毛球</li>
        <li :class="{ active: currentindex == 5 }" data-index="5">排球</li>
        <li :class="{ active: currentindex == 6 }" data-index="6">水上运动</li>
      </ul>
    </div>
    <!-- 轮播图 -->
    <div class="swiper">
      <van-swipe class="my-swipe" :autoplay="3000" indicator-color="white">
        <van-swipe-item>1</van-swipe-item>
        <van-swipe-item>2</van-swipe-item>
        <van-swipe-item>3</van-swipe-item>
        <van-swipe-item>4</van-swipe-item>
      </van-swipe>
    </div>
    <!-- 推荐运动 -->
    <div class="myrecommend">
      <!-- <div class="recommend-item">
      </div> -->
      <div class="recommend-item">
        <van-image src="#" fit="cover" style="width: 100%; height: 80%" />
        <div class="recommend-item-text-container">
          <div class="recommend-item-text">
            我觉得有一个分析里面说的挺对的，就是他认为所谓a药的效果实际上就是所谓的杀死这个人，指的是把这个人整个人的细胞全部杀死，然后在细胞死亡以后胶囊里面的药物就会立刻复制他身上的所有细胞。最后让那个人像是重获新生一样活过来，相当于他已经死了，现在还是有药物复制他原本的细胞，从而恢复的人。这也能说是为什么说他是能够让死人复活的神药了，如果真是这个效果的话，他的确能让死人复活。只要有对方的部分细胞配合上这个药物，就能在短时间重塑对方的身体。而且他的逻辑的确是有的，就比如说像柯南跟灰原来他们两个基本上没受什么外伤，在吃药的时候也没有什么症状，比如感冒之类的，所以他们恢复而不得变小是非常正常的，没有什么咳嗽之类的。但赤金玛丽在去的时候，他本身就有点要感冒的症状，再加上又掉进湖里面，掉进那个河里面呛了好几口水，如果说按照这个逻辑的话，细胞把它那个咳嗽或者感冒的病毒也默认成了他的细胞，并且复制下来。，在重塑他的身体的时候，把这个细胞也带上了那么再结合这个药物会让人永生，而永生最主要的就是让他的细胞不死，也就是说这种情况下他把感冒的病毒也搞成了永生不死的类型，也就是不会消亡的。，所以他的感冒到现在都没有恢复，也很合理。
          </div>
        </div>
      </div>
      <div class="recommend-item" style="background-color: pink"></div>
      <div class="recommend-item" style="background-color: pink"></div>
      <div class="recommend-item" style="background-color: pink"></div>
      <div class="recommend-item" style="background-color: pink"></div>
    </div>
  </div>
  .
</template>

<script setup lang="ts">
import router from '@/router'
import { ref } from 'vue'
const value = ref('')
const currentindex = ref(1)
// 点击导航栏触发事件
const handleClickclass = (e) => {
  // currentindex.value = index.dataset.index;
  if (e.target.tagName == 'LI') {
    currentindex.value = e.target.dataset.index
  }
}
const handleSearch = () => {
  router.push({
    path: '/search_result',
    query: {
      keyword: value.value,
    },
  })
}
const goToSearch = () => {
  router.push({
    path: '/search',
  })
}
</script>

<style scoped>
.recommend {
  width: 100%;
  background-color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  /* position: relative; */
}
/* 输入框样式 */
.search {
  width: 100%;
  height: 13.0435vw;
  display: flex;
  justify-content: space-between;
  /* position: absolute; */
  /* top: 0;
  left: 0; */
  /* padding: 0 2.4155vw; */
  /* border-bottom:.2415vw solid #4095E5; */

  /* background-color: #f5f5f500; */
}
.search .van-search {
  width: 80%;
  height: 100%;
  /* background-color: #fff; */
}
/* 搜索按钮样式 */
.search button {
  width: 18%;
  height: 50%;
  margin: auto;
  background-color: #4095e5;
  border: 0.2415vw solid #4095e5;
  border-radius: 3.6232vw;
}
/* 导航栏样式 */
.nav {
  padding-top: 2.4155vw;
  width: 100%;
  height: 8.6957vw;
  display: flex;
  justify-content: space-between;
  position: sticky;
  top: 0;
  /* padding: 0 2.4155vw; */
  background-color: #fff;
  /* opacity: 0.9; */
  z-index: 100;
  /* background-color: #f5f5f500; */
}
.nav ul {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 2.4155vw;
}
.nav ul li {
  font-size: 3.6232vw;
  /* width: 10%; */
  height: 100%;
  flex: 1;
  /* justify-content: center; */
  /* align-items: center; */
}
.nav ul .active {
  color: #4095e5;
}
/* 轮播图样式 */
.swiper {
  width: 100%;
  height: 36.715vw;
  /* margin-top: 4vh; */
  padding: 0 2.4155vw;
}

.swiper .my-swipe .van-swipe-item {
  color: #fff;
  font-size: 5.314vw;
  line-height: 36.2319vw;
  text-align: center;
  background-color: #39a9ed;
}
/* 推荐运动样式 */
.myrecommend {
  width: 100%;
  padding: 2.4155vw;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  grid-gap: 2.4155vw;
}
/* 推荐运动项目样式 */
.recommend-item {
  width: 100%;
  height: 43.4783vw;
  background-color: pink;
}
.recommend-item .van-image {
  width: 100%;
  height: 80%;
}
.recommend-item-text-container {
  width: 100%;
  height: 20%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 1.2077vw; /* 添加一点内边距防止文本紧贴边缘 */
}
.recommend-item-text {
  width: 100%;
  /* height: 20%; */
  /* line-height: 100%; */
  font-size: 3.1401vw;
  white-space: normal;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  word-break: break-all;
  line-height: 1.4; /* 调整行高以改善文本间距 */
}
</style>
