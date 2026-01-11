<template>
    <div class="dragonboat">
      <div class="table_head">
        <div class="match_name">赛事名称</div>
        <div class="team_name">队伍</div>
        <div class="boatnum">人数</div>
        <div class="match_class">赛事</div>
        <div class="match_time">时间</div>
      </div>
      <!-- 球员列表 -->
      <div class="table_body">
        <div class="list" v-for="item in dataList" :key="item.team_num">
          <div class="match_name">{{item.match_name}}</div>
          <div class="team_name">
            <span class="team_num">{{item.team_num}}</span>
            <span class="team_img"><img :src="item.team_img" alt="" /></span>
            <span class="team_names">{{item.team_name}}</span>
          </div>
          <div class="boatnum">{{item.boatnum}}</div>
          <div class="match_class">{{item.match_class}}</div>
          <div class="match_time">{{item.match_time}}</div>
        </div>
      </div>


      </div>
  </template>
  <script setup lang="ts">
  import { ref, onBeforeMount } from 'vue'
  onBeforeMount(() => {
    // 初始化数据
    // updateData()
  })
  const dataList=ref([
    {
      match_name:'逸仙杯',
      team_num:'1',
      team_img: '/src/assets/d195641d8c27c510.webp',
      team_name:'软件工程队',
      boatnum:'22人',
      match_class:'200米',
      match_time:'1分50秒'
    },
    {
      match_name:'逸仙杯',
      team_num:'2',
      team_img: '/src/assets/d195641d8c27c510.webp',
      team_name:'软件工程队',
      boatnum:'22人',
      match_class:'200米',
      match_time:'1分50秒'
    },
  ])
  </script>
  <style scoped>
  /* 球队头部 */
  .table_head {
    width: 100%;
    /* height: 39.9999px; */
    background-color: #cecece;
    padding: 2.4155vw;
    display: flex;
    align-items: center;
    text-align: center;
    /* justify-content: space-between; */
  }
  .match_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 3;
    text-align: left;
  }
  .team_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 4;
    text-align: left;
  }
  .boatnum{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
  }
  .match_class{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .match_time{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
    text-align: center;
  }
  /* 球员列表 */
  .table_body {
    width: 100%;
    padding: 10.0002px;
  }
  .list{
    width: 100%;
    height: 39.9999px;
    background-color: #fff;
    /* padding: 12.0002px; */
    display: flex;
    font-family: '宋体';
    align-items: center;
    /* justify-content: space-between; */
  }
  .list .match_name{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 3;
    text-align: left;
  }
  .list .team_name{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 4;
    /* position: relative; */
    display: flex;
    align-items: center;
    justify-content: start;
    gap: 1.9324vw;

  }
  .list .team_name img{
    width: 4.3478vw;
    height: 4.3478vw;
    border-radius: 50%;
    vertical-align: middle;
  }

  .list .boatnum{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .match_class{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .list .match_time{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    text-align: center;
  }

  </style>

