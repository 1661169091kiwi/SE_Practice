 <template>
    <div class="basketball_player">
      <div class="table_head">
        <div class="player_name">球员</div>
        <div class="team">队伍</div>
        <div class="shot">赛</div>
        <div class="lanban">胜/负</div>
        <div class="assist">胜率</div>
        <div class="times">连胜</div>
      </div>
      <!-- 球员列表 -->
      <div class="table_body">
        <div class="list">
          <div class="player_name">
            <span class="player_num">1</span>
            <span class="player_img"><img src="@/assets/d195641d8c27c510.webp" alt="" /></span>
            <span class="player_names">林丹</span>
          </div>
          <div class="team">软件工程队</div>
          <div class="shot">25</div>
          <div class="lanban">8/2</div>
          <div class="assist">98.6%</div>
          <div class="times">5</div>
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
  .player_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 3;
  }
  .team{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
  }
  .shot{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .lanban{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .assist{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .times{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  /* 球员列表 */
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
  .list .player_name{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 3;
    /* position: relative; */
    display: flex;
    align-items: center;
    justify-content: start;
    gap: 1.9324vw;

  }
  .list .player_name img{
    width: 4.3478vw;
    height: 4.3478vw;
    border-radius: 50%;
    vertical-align: middle;
  }

  .list .team{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
  }
  .list .shot{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .lanban{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .assist{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .times{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  </style>

