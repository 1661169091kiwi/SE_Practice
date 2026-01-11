<template>
  <div class="hotsport_index">
    <GameDisplay v-for="match in matchList" v-bind="match" key="match.gameid" @update:issaved="handleUpdateIssaved"/>
  </div>
</template>

<script setup lang="ts">
import GameDisplay from '@/component/GameDisplay.vue';
import { ref, onBeforeMount } from 'vue';

// const gameinfo = {
//   datatime: '2023-08-01 10:00:00',
//   gameid: '0',
//   team1: '软件工程',
//   team2: '微电子科学与技术',
//   team1score: '0',
//   team2score: '0',
//   gamestate: 1,
//   issaved: true,
//   gameclass:'篮球',
//   location:'珠海校区',
// }
const matchList = ref([
  {
    datatime: '2023-08-01 10:00:00',
    gameid: '0',
    teamA: '软件工程',
    teamB: '微电子科学与技术',
    scoreA: 0,
    scoreB: 0,
    gamestate: 2,
    issaved: true,
    gameclass:'篮球',
    location:'珠海校区',
  },
  {
    datatime: '2023-08-01 12:00:00',
    gameid: '1',
    teamA: '软件工程',
    teamB: '微电子科学与技术',
    scoreA: 0,
    scoreB: 0,
    gamestate: 1,
    issaved: true,
    gameclass:'篮球',
    location:'珠海校区',
  },
])

function handleUpdateIssaved(gameid, issaved) {
  const match = matchList.value.find((match) => match.gameid === gameid);
  if (match) {
    match.issaved = issaved;
  }
  //将更新的参数传给服务器
}
onBeforeMount(() => {
  // 从服务器获取数据
  // 更新 matchList.value
})
</script>

<style scoped>
</style>
