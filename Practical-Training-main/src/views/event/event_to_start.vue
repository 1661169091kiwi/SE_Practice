<template>
  <div class="match-detail">
    <!-- 顶部导航 -->
    <div class="top-nav">
      <button class="back-btn" @click="golastpage">←</button>
      <h1 class="match-title">比赛标题</h1>
      <span class="match-type">篮球</span>
    </div>

    <!-- 比赛状态 -->
    <div class="match-status">未开始</div>

    <!-- 比赛信息 -->
    <div class="match-info">
      <div class="info-item">
        <i class="icon">📅</i>
        <span>10月28日 14:00</span>
      </div>
      <div class="info-item">
        <i class="icon">📍</i>
        <span>主体育馆2号场地</span>
      </div>
    </div>

    <!-- 对阵信息 -->
    <div class="vs-container">
      <div class="team team-left">
        <div class="team-logo">🏀</div>
        <div class="team-name">软工队</div>
      </div>
      <div class="vs">VS</div>
      <div class="team team-right">
        <div class="team-logo">🏀</div>
        <div class="team-name">计院队</div>
      </div>
    </div>

    <!-- 近期战绩 -->
    <div class="recent-records">近3场：1胜2负</div>

    <!-- 报名区域 -->
    <div class="registration">
      <div class="registration-info">
        <i class="icon">✏️</i>
        <span>报名通道（剩余23名额）</span>
      </div>
      <button class="register-btn">立即报名</button>
      <div class="deadline">
        距报名时间截止还有:2天14小时
      </div>
    </div>

    <!-- 比赛亮点 -->
    <div class="highlights">
      <h3 class="section-title">比赛亮点：</h3>
      <ul class="highlight-list">
        <li>1.近三场比分 1:3</li>
        <li>2.左队近期连胜不断</li>
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
  name: 'MatchDetail',
  // Vue 2 中通过 data 定义响应式数据
  data() {
    return {
      isSubscribed: false // 订阅状态（响应式）
    }
  },
  // Vue 2 中通过 methods 定义方法
  methods: {
    // 返回上一页（Vue 2 中使用 this.$router）
    golastpage() {
      this.$router.go(-1);
    },
    // 切换订阅状态（Vue 2 中直接操作 this 上的属性）
    toggleSubscribe() {
      this.isSubscribed = !this.isSubscribed;
    }
  }
}
</script>

<style scoped>
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
