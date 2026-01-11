<template>
  <div class="GameDisplay">
    <div class="topline">
      <div class="gameinfo">
        <span>2023-01-01 12:00</span>
        <span>珠海校区</span>
      </div>
      <div class="shouchang">
        <van-icon name="star" :color="shouchang ? '#bd3124' : '#bbbbbb'" @click="toggleShouchang" />
      </div>
      <div class="gameclass">篮球</div>
    </div>
    <div class="middleline">
      <div class="matchinfo">
        <img src="@/assets/d195641d8c27c510.webp" alt="" />
        <span>软件工程</span>
      </div>
      <div class="gameres">
        <p>VS</p>
        <p>3:2</p>
      </div>
      <div class="matchinfo">
        <img src="@/assets/d195641d8c27c510.webp" alt="" />
        <span>微电子科学与技术</span>
      </div>
    </div>
    <!-- 最下面一行 -->
    <div class="bottomline">
      <button class="joinbtn finished">已结束</button>
      <button class="detailbtn">查看比赛详情</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const shouchang = ref(false)

const toggleShouchang = () => {
  if (shouchang.value) {
    shouchang.value = false
  } else {
    shouchang.value = true
  }
}
</script>
<style scoped>
.GameDisplay {
  width: 330px;
  height: 180px;
  border-radius: 15px;
  border: 1px solid #bbbbbb;
  margin: 10px auto;
  padding: 10px 30px 20px 20px;
}
/* 最上面一行 */
.topline {
  height: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.gameinfo {
  font-size: 12px;
  font-family: '宋体';
  color: #101010;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
}
.gameclass {
  font-size: 14px;
  font-family: '黑体';
  color: #7f83f7;
}
/* 中间一行 */
.middleline {
  height: 80px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 10px auto;
  /* padding: 0 20px; */
  /* border-bottom: 1px solid #bbbbbb; */
}
.matchinfo {
  font-size: 12px;
  font-family: '宋体';
  font-weight: 300;
  color: #101010;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-direction: column;
}
.matchinfo img {
  width: 53px;
  height: 53px;
}
.gameres {
  font-size: 14px;
  font-family: '宋体';
  color: #101010;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 10px;
  flex-direction: column;
}
/* 最下面一行 */
.bottomline {
  /* height: 40px; */
  display: flex;
  justify-content: space-between;
  align-items: space-between;
  margin: 10px auto;
}
.joinbtn {
  width: 70px;
  height: 30px;
  border-radius: 15px;
  /* border: 1px solid #7f83f7; */
  font-size: 14px;
  font-family: '宋体';
  /* color: #7f83f7; */
}
.bottomline .finished {
  background-color: #a4adb3;
  color: #000;
  border: 1px solid #a4adb3;
}
.bottomline .underway {
  background-color: rgba(162, 239, 77, 0.43);
  color: #81b337;
  border: 1px solid rgba(162, 239, 77, 0.43);
}
.bottomline .tostart {
  background-color: rgba(84, 188, 189, 0.73);
  color: #0f40f5;
  border: 1px solid rgba(84, 188, 189, 0.73);
}
.detailbtn {
  width: 100px;
  height: 30px;
  border-radius: 20px;
  border: 1px solid #e99d42;
  background-color: #e99d42;
  font-size: 14px;
  font-family: '宋体';
  color: #fff;
}
</style>
