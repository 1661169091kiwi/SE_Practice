 <template>
    <div class="table_container">
      <div class="table_head">
        <div class="qiudui">球队</div>
        <ul>
          <li>赛</li>
          <li>胜/平/负</li>
          <li>胜率</li>
          <li>积分</li>
        </ul>
      </div>
      <div class="table_body">
        <div class="list" v-for="team in matchdata" :key="team.gameId" @click="handleClickTeam(team)">
          <div class="qiudui">
            <span>{{team.teamrank}}</span>
            <span><img :src="team.teamlogo" alt="" /></span>
            <span>{{team.teamname}}</span>
          </div>
          <ul>
            <li class="gamenum">{{team.gamenum}}</li>
            <li class="winlose">{{team.winlose}}</li>
            <li class="score">{{team.winrate}}</li>
            <li class="points">{{team.points}}</li>
          </ul>
        </div>
      </div>
    </div>
  </template>
  <script setup lang="ts">
  import router from '@/router';
  import { ref ,onBeforeMount} from 'vue';
  const matchdata = ref([
    {
      teamrank:1,
      teamname:'软件工程',
      teamlogo:'/src/assets/d195641d8c27c510.webp',
      gamenum:10,
      winlose:'8/2/2',
      winrate:'98.6%',
      points:26,
      gameId:'1',
    },
    {
      teamrank:2,
      teamname:'微电子科学与技术',
      teamlogo:'/src/assets/d195641d8c27c510.webp',
      gamenum:1,
      winlose:'18/2/2',
      winrate:'1/6',
      points:2,
      gameId:'2',
    },
  ])
  onBeforeMount(()=>{
    // 从服务器获取数据
    // 更新 matchdata.value
    console.log(matchdata.value)
  })
  const handleClickTeam = (team) => {
    // router.push({
    //   name: 'badmintaton_team',
    //   params: { teamname: team.teamname }
    // });
  }
  </script>
  <style scoped>
    /* 表格 */
.table_container {
  width: 100%;
  /* padding: 2.4155vw; */
}
.table_head {
  width: 100%;
  /* height: 9.6618vw; */
  background-color: #cecece;
  padding: 2.8986vw;
  display: flex;
  justify-content: space-between;
}
.table_head .qiudui {
  font-size: 3.8647vw;
  color: #000;
}
.table_head ul {
  display: flex;
  /* justify-content: space-between; */
  gap: 4.8309vw;
}
.table_head ul li {
  font-size: 3.6232vw;
  color: #000;
  flex: auto;
  text-align: center;
}
.table_body {
  width: 100%;
  /* background-color: #54bcbd; */
  padding:0 12px;
  font-size: 16px;
  font-family: '宋体';
  font-weight: 300;
}
.list {
  width: 100%;
  height: 9.6618vw;
  /* background-color: #fff; */
  /* padding: 2.4155vw; */
  display: flex;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
}
.list .qiudui{
  text-align: center;
  /* vertical-align: middle; */
  display: flex;
  align-items: center;
}
.list ul {
  text-align: center;
  position: relative;
}
.list ul li {
  position: absolute;
  top: 50%;
  transform: translateY(-50%);
}
.gamenum {
  width: 15px;
 right: 175px;
}
.winlose {
  width: 53.34px;

right: 105px;
}
.score {
  width: 34.18px;
  right:50px;
}
.points {
width: 30px;
 right: 0;
}


.list .qiudui {
  display: flex;
  /* align-items: center; */
  gap: 2.4155vw;
}
.list .qiudui span {
  /* font-size: 4.8309vw; */
  color: #000;
  vertical-align: middle;
  /* text-align: center; */
}
.list .qiudui img {
  width: 4.3478vw;
  height: 4.3478vw;
  vertical-align: middle;
}
  </style>
