<template>
  <div class="data">
    <div class="chance_index">
      <ul @click="handleClickclass">
        <li :class="{ active: currentindex == 1 }" data-index="1">热门</li>
        <li :class="{ active: currentindex == 2 }" data-index="2">足球</li>
        <li :class="{ active: currentindex == 3 }" data-index="3">篮球</li>
        <li :class="{ active: currentindex == 4 }" data-index="4">羽毛球</li>
        <li :class="{ active: currentindex == 5 }" data-index="5">排球</li>
        <li :class="{ active: currentindex == 6 }" data-index="6">水上运动</li>
      </ul>
    </div>

    <div class="data_index">
      <div class="date_chance">
        <van-dropdown-menu active-color="#3775F6">
          <van-dropdown-item
            v-model="value1"
            :options="option1"
            background-color="#f5f5f5"
            @change="updateData"
          />
        </van-dropdown-menu>
      </div>
      <div class="my_data_opction">
        <ul @click="handleClickoption">
          <li :class="{ active: currentoptionindex == 1 }" data-index="1">积分</li>
          <li :class="{ active: currentoptionindex == 2 }" data-index="2">球员榜</li>
          <li :class="{ active: currentoptionindex == 3 }" data-index="3">球队榜</li>
          <li :class="{ active: currentoptionindex == 4 }" data-index="4">赛程</li>
          <li :class="{ active: currentoptionindex == 5 }" data-index="5">历史</li>
        </ul>
      </div>
    </div>

    <div class="table_container">
      <div class="table_head">
        <div class="qiudui">球队</div>
        <ul>
          <li>赛</li>
          <li>胜</li>
          <li>平</li>
          <li>负</li>
          <li>进/失</li>
          <li>积分</li>
        </ul>
      </div>
      <div class="table_body">
        <div class="list">
          <div class="qiudui">
            <span>1</span>
            <span><img src="/src/assets/d195641d8c27c510.webp" alt="" /></span>
            <span>软件工程</span>
          </div>
          <ul>
            <li>10</li>
            <li>8</li>
            <li>2</li>
            <li>2</li>
            <li>14-6</li>
            <li>26</li>
          </ul>
        </div>
        <div class="list">
          <div class="qiudui">
            <span>1</span>
            <span><img src="/src/assets/d195641d8c27c510.webp" alt="" /></span>
            <span>软件工程</span>
          </div>
          <ul>
            <li>10</li>
            <li>8</li>
            <li>2</li>
            <li>2</li>
            <li>14-6</li>
            <li>26</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const currentindex = ref(1)
// 点击导航栏触发事件
const handleClickclass = (e) => {
  // currentindex.value = index.dataset.index;
  if (e.target.tagName == 'LI') {
    currentindex.value = e.target.dataset.index
    //更新数据
    updateData()
  }
}
const currentoptionindex = ref(1)
// 点击数据类型导航栏触发事件
const handleClickoption = (e) => {
  // currentindex.value = index.dataset.index;
  if (e.target.tagName == 'LI') {
    currentoptionindex.value = e.target.dataset.index
    //更新数据
    updateData()
  }
}
const value1 = ref(2025)
const option1 = [
  {
    text: '2025',
    value: 2025,
  },
  {
    text: '2024',
    value: 2024,
  },
  {
    text: '2023',
    value: 2023,
  },
  {
    text: '2022',
    value: 2022,
  },
  {
    text: '2021',
    value: 2021,
  },
  {
    text: '2020',
    value: 2020,
  },
  {
    text: '2019',
    value: 2019,
  },
]
// 更新数据
const updateData = () => {
  // 这里可以根据 currentindex.value 和 currentoptionindex.value 来更新数据
  console.log('更新数据', currentindex.value, currentoptionindex.value)
}
</script>
<style scoped>
/* 运动类型导航栏 */
.chance_index {
  width: 100%;
  height: 16.9082vw;
  /* background-color: #f5f5f5; */
  padding: 4.8309vw;
}
.chance_index ul {
  display: flex;
  justify-content: space-between;
}
.chance_index ul li {
  font-size: 4.8309vw;
  color: #333;
  flex: auto;
}
.chance_index ul .active {
  color: #3775f6;
}
/* 数据类型导航栏 */
.data_index {
  width: 100%;
  height: 12.0773vw;
  background-color: #f5f5f5;
  /* padding:4.8309vw; */
}
.date_chance {
  width: 16%;
  display: inline-block;
}
.my_data_opction {
  width: 84%;
  display: inline-block;
  padding: 0 0 0 4.8309vw;
}
.my_data_opction ul {
  display: flex;
  justify-content: space-between;
}
.my_data_opction ul li {
  font-size: 4.8309vw;
  color: #333;
  flex: auto;
}
.my_data_opction ul li.active {
  color: #3775f6;
}
/* 表格 */
.table_container {
  width: 100%;
  /* padding: 2.4155vw; */
}
.table_head {
  width: 100%;
  /* height: 9.6618vw; */
  background-color: #cecece;
  padding: 2.4155vw;
  display: flex;
  justify-content: space-between;
}
.qiudui {
  font-size: 3.8647vw;
  color: #000;
}
.table_head ul {
  display: flex;
  /* justify-content: space-between; */
  gap: 3.8647vw;
}
.table_head ul li {
  font-size: 3.8647vw;
  color: #000;
  flex: auto;
}
.table_body {
  width: 100%;
  background-color: #54bcbd;
  padding: 2.4155vw;
}
.list {
  width: 100%;
  height: 9.6618vw;
  /* background-color: #fff; */
  /* padding: 2.4155vw; */
  display: flex;
  justify-content: space-between;
}
.list ul {
  display: flex;
  /* justify-content: space-between; */
  gap: 4.8309vw;
  /* background-color: #fff; */
}
.list ul li {
  font-size: 4.8309vw;
  color: #000;
  flex: auto;
}
.qiudui {
  display: flex;
  /* align-items: center; */
  gap: 2.4155vw;
}
.qiudui span {
  font-size: 4.8309vw;
  color: #000;
  vertical-align: middle;
  /* text-align: center; */
}
.qiudui img {
  width: 4.3478vw;
  height: 4.3478vw;
  vertical-align: middle;
}
</style>
