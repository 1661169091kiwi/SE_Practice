 <template>
    <div class="volleyball_history">
      <div class="table_head">
        <div class="match_name">赛事</div>
        <div class="match_time">赛季</div>
        <div class="team_name">冠军</div>
        <div class="score">积分</div>
      </div>
      <!-- 球队列表 -->
      <div class="table_body">
        <div class="list">
          <div class="match_name">逸仙杯</div>
          <div class="match_time">2025-2026</div>
          <div class="team_name">
            <!-- <span class="player_num">1</span> -->
            <span class="team_img"><img src="@/assets/d195641d8c27c510.webp" alt="" /></span>
            <span class="team_name">软件工程队</span>
          </div>
          <div class="score">78</div>
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
  .match_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
  }
  .match_time{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
  }
  .team_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 3;
    text-align: center;
  }
  .score{
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
    font-weight: 300;
    /* justify-content: space-between; */
  }
  .list .team_name{
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
  .list .team_img img{
    width: 4.3478vw;
    height: 4.3478vw;
    border-radius: 50%;
    vertical-align: middle;
  }

  .list .match_time{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
  }
  .list .match_name{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    /* text-align: center; */
  }
  .list .score{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  </style>
