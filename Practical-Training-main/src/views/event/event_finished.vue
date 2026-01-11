<template>
  <div class="match-detail">
    <!-- 顶部导航 -->
    <header class="top-nav">
      <button class="back-btn" @click="goBack">←</button>
      <h1 class="match-title">比赛标题</h1>
      <span class="match-type">足球</span>
    </header>

    <!-- 比赛状态 -->
    <div class="match-status">已结束</div>

    <!-- 比分区域 -->
    <div class="score-section">
      <div class="team team-left">
        <div class="team-logo">⚽</div>
        <div class="team-name">软工队</div>
      </div>

      <div class="score-info">
        <div class="main-score">2 VS 1</div>
        <div class="half-scores">
          <span>上半场1-1</span> | <span>下半场1-0</span>
        </div>
        <div class="match-time">10月28日 14:00</div>
        <div class="match-venue">主体育馆</div>
      </div>

      <div class="team team-right">
        <div class="team-logo">⚽</div>
        <div class="team-name">计院队</div>
      </div>
    </div>

    <!-- 详细战报 -->
    <div class="report-section">
      <h3 class="section-title">详细战报：</h3>
      <p class="report-content">
        计算机学院凭借李明补时阶段的绝杀，3-2战胜文学院晋级四强。第35分钟，文学院刘强扳平比分；第88分钟，李明接张伟助攻完成绝杀
      </p>
    </div>

    <!-- 技术统计 -->
    <div class="stats-section">
      <h3 class="section-title">技术统计</h3>
      <div class="stats-grid">
        <div class="stat-item">
          <span class="stat-name">射门</span>
          <span class="stat-value">1-1</span>
        </div>
        <div class="stat-item">
          <span class="stat-name">射正</span>
          <span class="stat-value">5-8</span>
        </div>
        <div class="stat-item">
          <span class="stat-name">角球</span>
          <span class="stat-value">2-3</span>
        </div>
        <div class="stat-item">
          <span class="stat-name">犯规</span>
          <span class="stat-value">2-1</span>
        </div>
      </div>
    </div>

    <!-- 球员评分 -->
    <div class="players-rating">
      <!-- 软工队评分 -->
      <div class="team-ratings">
        <h3 class="section-title">软工队球员评分</h3>
        <div class="player-item">
          <div class="player-info">
            <i class="player-icon">🥅</i>
            <span class="player-name">门将 张三</span>
            <span class="player-score">7分</span>
          </div>
          <div class="rating-stars">
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star empty">☆</i>
          </div>
        </div>
      </div>

      <!-- 计院队评分 -->
      <div class="team-ratings">
        <h3 class="section-title">计院队球员评分</h3>
        <div class="player-item">
          <div class="player-info">
            <i class="player-icon">🧑⚽</i>
            <span class="player-name">前锋 王五</span>
            <span class="player-score">8分</span>
          </div>
          <div class="rating-stars">
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star filled">★</i>
            <i class="star empty">☆</i>
          </div>
        </div>
      </div>
    </div>

    <!-- 精彩瞬间 -->
    <div class="highlights-section">
      <h3 class="section-title">精彩瞬间（比赛集锦）</h3>
      <div class="highlights-tabs">
        <button class="tab-btn active">图片</button>
        <button class="tab-btn">视频</button>
      </div>

      <!-- 图片轮播 -->
      <div class="highlights-carousel">
        <img
          src="https://p26-flow-imagex-download-sign.byteimg.com/tos-cn-i-a9rns2rl98/7d636ccdf0ec47e99d7dea77aa5dbc82.png~tplv-a9rns2rl98-24:720:720.png"
          alt="比赛精彩瞬间"
          class="highlight-img"
        >
        <!-- 轮播指示器 -->
        <div class="carousel-dots">
          <span class="dot active"></span>
          <span class="dot"></span>
          <span class="dot"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FinishedMatchDetail',
  methods: {
    // 返回上一页（Vue 2 中使用 this.$router）
    goBack() {
      this.$router.go(-1);
    }
  }
};
</script>

<style scoped>
.match-detail {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #fff;
  min-height: 100vh;
  padding-bottom: 20px;
}

/* 顶部导航 */
.top-nav {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #eee;
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
  margin: 0;
  flex: 1;
  text-align: center;
}

.match-type {
  color: #4caf50; /* 足球绿色 */
  font-size: 14px;
  position: absolute;
  right: 16px;
}

/* 比赛状态 */
.match-status {
  color: #666; /* 已结束灰色 */
  font-size: 16px;
  text-align: center;
  padding: 12px 0;
  font-weight: 500;
  border-bottom: 1px solid #eee;
}

/* 比分区域 */
.score-section {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 20px 0;
  border-bottom: 1px solid #eee;
}

.team {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 30%;
}

.team-logo {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-bottom: 8px;
}

.team-name {
  font-size: 16px;
  font-weight: 500;
}

.score-info {
  text-align: center;
}

.main-score {
  font-size: 36px;
  font-weight: bold;
  margin: 0 0 8px 0;
}

.half-scores {
  font-size: 14px;
  color: #666;
  margin: 0 0 4px 0;
}

.match-time, .match-venue {
  font-size: 12px;
  color: #888;
  margin: 2px 0;
}

/* 通用区域标题 */
.section-title {
  font-size: 16px;
  font-weight: 600;
  padding: 16px 16px 8px;
  margin: 0;
  color: #333;
}

/* 详细战报 */
.report-content {
  padding: 0 16px 16px;
  margin: 0;
  font-size: 14px;
  line-height: 1.6;
  color: #555;
  border-bottom: 1px solid #eee;
}

/* 技术统计 */
.stats-section {
  border-bottom: 1px solid #eee;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1px;
  background-color: #f5f5f5;
}

.stat-item {
  background-color: #fff;
  padding: 12px 16px;
  display: flex;
  justify-content: space-between;
  font-size: 14px;
}

.stat-name {
  color: #666;
}

.stat-value {
  font-weight: 500;
}

/* 球员评分 */
.players-rating {
  border-bottom: 1px solid #eee;
}

.team-ratings {
  margin-bottom: 16px;
}

.player-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
}

.player-info {
  display: flex;
  align-items: center;
}

.player-icon {
  margin-right: 8px;
  font-size: 16px;
}

.player-name {
  font-size: 14px;
  margin-right: 16px;
}

.player-score {
  font-size: 14px;
  font-weight: 500;
  color: #1677ff;
}

.rating-stars {
  display: flex;
}

.star {
  color: #ddd;
  margin-left: 2px;
}

.star.filled {
  color: #ffcc00;
}

/* 精彩瞬间 */
.highlights-tabs {
  display: flex;
  padding: 0 16px;
  margin-bottom: 12px;
}

.tab-btn {
  flex: 1;
  padding: 8px 0;
  background-color: #f5f5f5;
  border: none;
  font-size: 14px;
  cursor: pointer;
  margin-right: 8px;
  border-radius: 4px;
}

.tab-btn:last-child {
  margin-right: 0;
}

.tab-btn.active {
  background-color: #1677ff;
  color: white;
}

.highlights-carousel {
  position: relative;
  padding: 0 16px;
}

.highlight-img {
  width: 100%;
  border-radius: 8px;
  display: block;
}

.carousel-dots {
  display: flex;
  justify-content: center;
  margin-top: 12px;
}

.dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #ddd;
  margin: 0 4px;
}

.dot.active {
  background-color: #1677ff;
}
</style>
