<template>
  <div class="match-detail">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <button class="back-btn" @click="golastpage">←</button>
      <h1 class="match-title">{{ matchInfo.title }}</h1>
      <span class="match-type">{{ matchInfo.sportType }}</span>
    </div>

    <!-- 比赛状态 -->
    <div class="match-status">未开始</div>

    <!-- 比赛信息 -->
    <div class="match-info">
      <div class="info-item">
        <i class="icon">📅</i>
        <span>{{ matchInfo.dateTime }}</span>
      </div>
      <div class="info-item">
        <i class="icon">📍</i>
        <span>{{ matchInfo.venue }}</span>
      </div>
    </div>

    <!-- 对阵信息 -->
    <div class="vs-container">
      <div class="team team-left">
        <div class="team-logo">{{ matchInfo.teamA.logo }}</div>
        <div class="team-name">{{ matchInfo.teamA.name }}</div>
      </div>
      <div class="vs">VS</div>
      <div class="team team-right">
        <div class="team-logo">{{ matchInfo.teamB.logo }}</div>
        <div class="team-name">{{ matchInfo.teamB.name }}</div>
      </div>
    </div>

    <!-- 近期战绩 -->
    <div class="recent-records">近3场：{{ matchInfo.recentRecords }}</div>

    <!-- 报名区域 -->
    <div class="registration">
      <div class="registration-info">
        <i class="icon">✏️</i>
        <span>报名通道（剩余{{ matchInfo.remainingQuota }}名额）</span>
      </div>
      <button class="register-btn">立即报名</button>
      <div class="deadline">
        距报名时间截止还有:{{ matchInfo.deadline }}
      </div>
    </div>

    <!-- 比赛亮点 -->
    <div class="highlights">
      <h3 class="section-title">比赛亮点：</h3>
      <ul class="highlight-list">
        <li v-for="(highlight, index) in matchInfo.highlights" :key="index">
          {{ index + 1 }}.{{ highlight }}
        </li>
      </ul>
    </div>

    <!-- 订阅提醒 -->
    <div class="subscribe-section">
      <button class="subscribe-btn" @click="toggleSubscribe">
        <i class="heart-icon" :class="{ active: isSubscribed }">❤️</i>
        <span>订阅比赛提醒</span>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'EventToStart',
  props: {
    sportType: {
      type: String,
      required: true,
      validator: value => ['football', 'basketball', 'tabletennis', 'badminton', 'water'].includes(value)
    },
    matchType: {
      type: String,
      default: 'doubles',
      validator: value => ['singles', 'doubles', 'team'].includes(value)
    }
  },
  data() {
    return {
      isSubscribed: false
    }
  },
  created() {
    console.log('当前matchType:', this.matchType); // 若输出'doubles'，则说明参数未传递成功
  },
  computed: {
    matchInfo() {
      const sportConfigs = {
        football: {
          title: '足球联赛半决赛',
          sportType: '足球',
          dateTime: '11月5日 15:00',
          venue: '主体育场',
          teamA: { name: '软工队', logo: '⚽' },
          teamB: { name: '计院队', logo: '⚽' },
          recentRecords: '2胜1平',
          remainingQuota: 18,
          deadline: '3天8小时',
          highlights: [
            '上赛季交手2-2战平',
            '软工队主场不败纪录',
            '计院队新外援首秀'
          ]
        },
        basketball: {
          title: '篮球友谊赛',
          sportType: '篮球',
          dateTime: '10月28日 19:30',
          venue: '主体育馆2号场地',
          teamA: { name: '经管队', logo: '🏀' },
          teamB: { name: '外语队', logo: '🏀' },
          recentRecords: '1胜2负',
          remainingQuota: 23,
          deadline: '2天14小时',
          highlights: [
            '两队历史交锋10胜10负',
            '经管队三分球命中率联盟第一',
            '外语队内线优势明显'
          ]
        },
        tabletennis: {
          title: '乒乓球团体赛',
          sportType: '乒乓球',
          dateTime: '11月1日 10:00',
          venue: '体育中心乒乓球馆',
          teamA: { name: '研究生院队', logo: '🏓' },
          teamB: { name: '本科生队', logo: '🏓' },
          recentRecords: '3胜0负',
          remainingQuota: 5,
          deadline: '1天2小时',
          highlights: [
            '研究生队拥有校队主力',
            '本科生队平均年龄小5岁',
            '上届决赛重演'
          ]
        },
        badminton: {
          ...(this.matchType === 'singles'
            ? {
                title: '羽毛球单打决赛',
                sportType: '羽毛球',
                dateTime: '10月30日 14:00',
                venue: '羽毛球馆1号场',
                teamA: { name: '林丹', logo: '🏸' },
                teamB: { name: '李宗伟', logo: '🏸' },
                recentRecords: '2胜1负',
                remainingQuota: 0,
                deadline: '12小时30分',
                highlights: [
                  '世界排名前两位选手对决',
                  '职业生涯第40次交锋',
                  '胜者将代表学校参加省赛'
                ]
              }
            : {
                title: '羽毛球混双决赛',
                sportType: '羽毛球',
                dateTime: '10月30日 15:30',
                venue: '羽毛球馆1号场',
                teamA: { name: '赵阳/孙梅', logo: '🏸' },
                teamB: { name: '钱峰/周琳', logo: '🏸' },
                recentRecords: '1胜1负',
                remainingQuota: 3,
                deadline: '14小时',
                highlights: [
                  '两对组合均为校队主力',
                  '上月交手3局苦战分出胜负',
                  '关键分处理能力决定胜负'
                ]
              })
        },
        water: {
          title: '游泳4x100米接力',
          sportType: '水上运动',
          dateTime: '11月8日 9:00',
          venue: '游泳馆主池',
          teamA: { name: '红队', logo: '🏊' },
          teamB: { name: '蓝队', logo: '🏊' },
          recentRecords: '1胜2负',
          remainingQuota: 10,
          deadline: '5天6小时',
          highlights: [
            '上届纪录保持者参赛',
            '蓝队平均年龄更小',
            '红队第四棒为国家二级运动员'
          ]
        }
      }
      return sportConfigs[this.sportType]
    }
  },
  methods: {
    golastpage() {
      this.$router.go(-1);
    },
    toggleSubscribe() {
      this.isSubscribed = !this.isSubscribed;
    }
  }
}
</script>

<style scoped>
/* 保持原有样式不变 */
.match-detail {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  padding: 16px;
  background-color: #fff;
  min-height: 100vh;
  box-sizing: border-box;
}

/* 顶部导航 */
.top-nav {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  position: relative;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
}

.match-title {
  font-size: 18px;
  font-weight: 600;
  flex: 1;
  text-align: center;
  margin: 0;
}

.match-type {
  color: #4285f4;
  font-size: 14px;
  position: absolute;
  right: 0;
  top: 50%;
  transform: translateY(-50%);
}

/* 比赛状态 */
.match-status {
  color: #4285f4;
  font-size: 16px;
  text-align: center;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

/* 比赛信息 */
.match-info {
  margin-bottom: 24px;
}

.info-item {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
  color: #333;
}

.icon {
  margin-right: 8px;
  width: 20px;
  text-align: center;
}

/* 对阵信息 */
.vs-container {
  display: flex;
  align-items: center;
  justify-content: space-around;
  margin: 32px 0;
}

.team {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 30%;
}

.team-logo {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-bottom: 12px;
}

.team-name {
  font-size: 16px;
  font-weight: 500;
}

.vs {
  font-size: 24px;
  font-weight: bold;
  color: #666;
}

/* 近期战绩 */
.recent-records {
  text-align: center;
  color: #666;
  font-size: 14px;
  margin-bottom: 32px;
  padding: 8px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

/* 报名区域 */
.registration {
  text-align: center;
  margin-bottom: 32px;
  padding: 16px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.registration-info {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  font-size: 14px;
  color: #333;
}

.register-btn {
  background-color: #ffcc00;
  color: #333;
  border: none;
  border-radius: 20px;
  padding: 10px 32px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 12px;
}

.deadline {
  font-size: 12px;
  color: #666;
}

/* 比赛亮点 */
.highlights {
  margin-bottom: 32px;
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #333;
}

.highlight-list {
  padding-left: 20px;
  color: #666;
  font-size: 14px;
}

.highlight-list li {
  margin-bottom: 8px;
}

/* 订阅提醒 */
.subscribe-section {
  position: fixed;
  bottom: 20px;
  left: 0;
  right: 0;
  padding: 0 16px;
}

.subscribe-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  background-color: #f5f5f5;
  border: none;
  border-radius: 20px;
  padding: 10px 0;
  font-size: 14px;
  cursor: pointer;
}

.heart-icon {
  margin-right: 8px;
}

.heart-icon.active {
  color: #ff3b30;
}
</style>
