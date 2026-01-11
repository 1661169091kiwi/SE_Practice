 <template>
    <div class="basketball_player">
      <div class="table_head">
        <div class="team_name">球队</div>
        <div class="match_num">赛</div>
        <div class="win_lose">胜/负</div>
        <div class="win_rate">胜率</div>
        <div class="win_diff">胜差</div>
        <div class="curr">近况</div>
      </div>
      <!-- 球员列表 -->
      <div class="table_body">
        <div class="list" v-for="team in matchdata" :key="team.teamrank" @click="handleClickTeam(team)">
          <div class="team_name">
            <span class="team_num">{{team.teamrank}}</span>
            <span class="team_img"><img :src="team.teamlogo" alt="" /></span>
            <span class="team_names">{{team.teamname}}</span>
          </div>
          <div class="match_num">{{team.gamenum}}</div>
          <div class="win_lose">{{team.winlose}}</div>
          <div class="win_rate">{{team.winrate}}</div>
          <div class="win_diff">{{team.win_diff}}</div>
          <div class="curr">{{team.curr}}</div>
        </div>
      </div>


      </div>
  </template>
  <script setup lang="ts">
  import router from '@/router';
  import { ref ,onBeforeMount} from 'vue';
  const matchdata=ref([
    {
      teamrank:1,
      teamname:'软件工程',
      teamlogo:'/src/assets/d195641d8c27c510.webp',
      gamenum:10,
      winlose:'12/1',
      winrate:'92.4%',
      win_diff:'2.0',
      curr:'4连胜',
    },
    {
      teamrank:2,
      teamname:'软件测试',
      teamlogo:'/src/assets/d195641d8c27c510.webp',
      gamenum:10,
      winlose:'12/1',
      winrate:'92.4%',
      win_diff:'2.0',
      curr:'4连胜',
    },
  ])
  onBeforeMount(()=>{
    // 从服务器获取数据
    // 更新 matchdata.value
    console.log(matchdata.value)
  })
  const handleClickTeam=(team)=>{
    // console.log(team)
    // 前往相应的队伍页面
    // router.push({
    //   name:'basketball_team',
    //   params:{
    //     teamid:team.teamid,
    //   }
    // })
  }
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
  .team_name{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 4;
    text-align: left;
  }
  .match_num{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
  }
  .win_lose{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .win_rate{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .win_diff{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .curr{
    font-size: 3.8647vw;
    font-weight: bold;
    color: #000;
    flex: 2;
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

  .list .match_num{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .win_lose{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .list .win_rate{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    text-align: center;
  }
  .list .win_diff{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 1;
    text-align: center;
  }
  .list .curr{
    font-size: 3.8647vw;
    font-weight: 300;
    color: #000;
    flex: 2;
    text-align: center;
  }
  </style>
