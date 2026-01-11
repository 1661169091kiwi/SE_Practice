 <template>
    <div class="volleyball_team">
      <div class="table_head">
        <div class="team">球队</div>
        <div class="wins">胜场</div>
        <div class="lose">负场</div>
        <div class="win_rate">胜率</div>
        <div class="score">得分</div>
      </div>
      <!-- 球队列表 -->
      <div class="table_body">
        <div class="list">
          <div class="team-name">
            <span class="player_num">1</span>
            <span class="player_img"><img src="@/assets/d195641d8c27c510.webp" alt="" /></span>
            <span class="player_names">软件工程队</span>
          </div>

          <div class="wins">25</div>
          <div class="lose">10</div>
          <div class="win_rate">98.14%</div>
          <div class="score">49</div>
        </div>
      </div>


      </div>
  </template>
  <script setup lang="ts">

  </script>
  <style scoped>
  /* 球队头部 */
  .table_head {
    width: 100%;
    /* height: 9.6618vw; */
    background-color: #cecece;
    padding: 2.8986vw;
    display: flex;
    /* justify-content: space-between; */
  }
  .team{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 3;
  }
  .wins{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
  }
  .lose{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .win_rate{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .score{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  /* 球队列表 */
  .table_body {
    width: 100%;
    padding: 2.4155vw;
  }
  .list{
    width: 100%;
    height: 9.6618vw;
    background-color: #fff;
    /* padding: 2.8986vw; */
    display: flex;
    font-family: '宋体';
    align-items: center;
    /* justify-content: space-between; */
  }
  .list .team-name{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 3;
    /* position: relative; */
    display: flex;
    align-items: center;
    justify-content: start;
    gap: 8px;

  }
  .list .team-name img{
    width: 18px;
    height: 18px;
    border-radius: 50%;
    vertical-align: middle;
  }

  .list .wins{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
  }
  .list .lose{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .win_rate{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .score{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  </style>

