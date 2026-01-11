<template>
  <div class="match-card">
    <div class="card-header">
      <span class="match-time">{{ match.time }}</span>
      <span class="match-title">{{ match.league }}</span>
      <span class="match-sports">{{ match.sportsType }}</span>
    </div>
    <div class="team-info">
      <div class="team-left">
        <div class="team-name" @click="gototeampage()">
          {{ match.home }} <span class="team-icon">{{ match.teamIcon }}</span>
        </div>
        <div class="player-info">{{ playerInfoText }}{{ match.num }}人</div>
      </div>
      <div class="score-info">
        <div class="actual-score">{{ match.score }}</div>
        <div class="match-status">{{ match.status }}</div>
        <button class="detail-btn" @click="gotonextpage()">{{ detailBtnText }}</button>
      </div>
      <div class="team-right">
        <div class="team-name" @click="gototeampage()">
          {{ match.away }} <span class="team-icon">{{ match.teamIcon }}</span>
        </div>
        <div class="player-info">{{ playerInfoText }}{{ match.num }}人</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MatchCard',
  props: {
    match: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      playerInfoText: "首发：",
      detailBtnText: "详情"
    };
  },
  watch: {
    match: {
      immediate: true,
      handler(newVal) {
        this.event_status = newVal.status;
      }
    }
  },
  methods: {
    gotonextpage() {
      let sportsType = this.match.sports;
      let singleType = this.match.single;
      if (this.event_status === "未开始") {
        if (sportsType === 'badminton'&&singleType) { // 假设match对象有sport字段标识项目
            this.$router.push({
            path: '/event-to-start/badminton',
            query: { matchType: 'singles' } // 携带单打类型参数
            });
        }
        else if(sportsType === 'badminton'&&!singleType){
          this.$router.push({
            path: '/event-to-start/badminton',
            query: { matchType: 'double' } // 携带双打类型参数
            });
        }
        else{
          this.$router.push({
            path: '/event-to-start/'+sportsType,
            });
        }
      } else if (this.event_status === "进行中") {
        if (sportsType === 'badminton'&&singleType) { // 假设match对象有sport字段标识项目
            this.$router.push({
            path: '/ongoing/badminton',
            query: { matchType: 'singles' } // 携带单打类型参数
            });
        }
        else if(sportsType === 'badminton'&&!singleType){
          this.$router.push({
            path: '/ongoing/badminton',
            query: { matchType: 'double' } // 携带双打类型参数
            });
        }
        else{
          this.$router.push({
            path: '/ongoing/'+sportsType,
            });
        }
      } else if (this.event_status === "已结束") {
          if (sportsType === 'badminton'&&singleType) { // 假设match对象有sport字段标识项目
            this.$router.push({
            path: '/finished/badminton',
            query: { matchType: 'singles' } // 携带单打类型参数
            });
        }
        else if(sportsType === 'badminton'&&!singleType){
          this.$router.push({
            path: '/finished/badminton',
            query: { matchType: 'double' } // 携带双打类型参数
            });
        }
        else{
          this.$router.push({
            path: '/finished/'+sportsType,
            });
        }
      }


    },
    gototeampage() {
      let sportsType = this.match.sports;
      this.$router.push('/event/teampage/'+sportsType)
    }
  }
};
</script>

<style scoped>
/* 样式部分保持不变 */
.match-card {
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 16px;
  background-color: #fff;
}

.card-header {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  color: #666;
  position: relative;
}
.match-time {
  font-size: 14px;
}
.match-title {
  font-size: 14px;
  position: absolute;
  font-weight: bold;
  left: 50%;
  transform: translateX(-50%);
}
.match-sports {
  font-size: 14px;
  margin-left: auto;
}
.team-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.team-left, .team-right {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 35%;
}
.team-name {
  font-size: 16px;
  margin: 0 0 8px 0;
  font-weight: 500;
}
.team-icon {
  margin-left: 4px;
  font-size: 14px;
  color: #ccc;
}
.player-info {
  font-size: 12px;
  color: #888;
  margin: 0;
}
.score-info {
  display:flex;
  flex-direction: column;
  align-items: center;
}
.actual-score {
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 8px;
}
.match-status {
  font-size: 14px;
  color: #f5222d;
  margin-bottom: 12px;
}
.detail-btn {
  border: 1px solid #4285f4;
  color: #4285f4;
  background-color: transparent;
  border-radius: 4px;
  padding: 6px 12px;
  font-size: 14px;
  cursor: pointer;
}
.detail-btn:hover {
  background-color: #1677ff;
  color: white;
}
</style>
