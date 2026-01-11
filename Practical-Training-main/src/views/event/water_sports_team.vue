<template>
  <div class="player-page">
    <header class="top-nav">
      <button class="back-btn" @click="goBack">←</button>
      <h1 class="page-title">软工队队员列表</h1>
    </header>

    <div class="team-summary">
      <div class="team-logo">🏊</div>
      <div class="team-info">
        <div class="team-name">软件工程学院水上运动队</div>
        <div class="team-stats">
          <span>成立时间：2019年</span>
          <span>成员数量：6人</span>
          <span>最佳成绩：校冠军</span>
        </div>
      </div>
    </div>

    <div class="section-header">
      <h2>队员信息</h2>
      <span class="count">{{ players.length }}名队员</span>
    </div>

    <div class="players-list">
      <div class="player-card" v-for="(player, index) in players" :key="index">
        <div class="player-avatar">
          <img :src="player.avatar" alt="队员头像" class="avatar-img">
          <div class="player-number">{{ player.number }}</div>
        </div>

        <div class="player-basic">
          <div class="player-name">
            {{ player.name }}
            <span class="player-position">{{ player.position }}</span>
          </div>
          <div class="player-details">
            <span>身高：{{ player.height }}</span>
            <span>体重：{{ player.weight }}</span>
            <span>年龄：{{ player.age }}岁</span>
          </div>
        </div>

        <div class="player-stats">
          <div class="stat-item">
            <div class="stat-label">主项</div>
            <div class="stat-value">{{ player.mainEvent }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">最佳成绩</div>
            <div class="stat-value">{{ player.bestTime }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">金牌数</div>
            <div class="stat-value">{{ player.golds }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'WaterSportsTeam',
  data() {
    return {
      players: [
        {
          name: '孙洋',
          number: 1,
          position: '自由泳',
          height: '188cm',
          weight: '75kg',
          age: 22,
          mainEvent: '100米',
          bestTime: '52.3s',
          golds: 5,
          avatar: 'https://via.placeholder.com/60?text=1'
        },
        {
          name: '吴敏',
          number: 2,
          position: '蛙泳',
          height: '172cm',
          weight: '62kg',
          age: 20,
          mainEvent: '200米',
          bestTime: '2:25.8',
          golds: 3,
          avatar: 'https://via.placeholder.com/60?text=2'
        }
      ]
    }
  },
  methods: {
    goBack() {
      this.$router.push('/event/menu')
    }
  }
}
</script>

<style scoped>
.player-page {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #f5f5f5;
  min-height: 100vh;
  padding-bottom: 20px;
}

/* 顶部导航 */
.top-nav {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: #fff;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
  margin-bottom: 16px;
}

.back-btn {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  padding: 4px 8px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  margin: 0 0 0 12px;
}

/* 球队信息摘要 */
.team-summary {
  display: flex;
  align-items: center;
  background-color: #fff;
  padding: 16px;
  border-radius: 8px;
  margin: 0 16px 16px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.team-logo {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  margin-right: 16px;
}

.team-name {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}

.team-stats {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #666;
}

/* 列表标题 */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px 12px;
  margin: 0 16px 12px;
  border-bottom: 1px solid #eee;
}

.section-header h2 {
  font-size: 16px;
  margin: 0;
}

.count {
  font-size: 12px;
  color: #888;
}

/* 球员列表 */
.players-list {
  padding: 0 16px;
}

/* 球员卡片（复用多个） */
.player-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 12px;
  margin-bottom: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

/* 球员头像与号码 */
.player-avatar {
  position: relative;
  margin-right: 16px;
}

.avatar-img {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  background-color: #f0f0f0;
}

.player-number {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 22px;
  height: 22px;
  background-color: #1677ff;
  color: white;
  border-radius: 50%;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 球员基本信息 */
.player-basic {
  flex: 1;
  margin-right: 16px;
}

.player-name {
  font-size: 16px;
  font-weight: 500;
  margin-bottom: 6px;
}

.player-position {
  font-size: 12px;
  color: white;
  background-color: #666;
  padding: 2px 6px;
  border-radius: 4px;
  margin-left: 8px;
}

.player-details {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #666;
}

/* 球员数据统计 */
.player-stats {
  display: flex;
  gap: 16px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 12px;
  color: #888;
  margin-bottom: 4px;
  display: block;
}

.stat-value {
  font-size: 14px;
  font-weight: 600;
  color: #333;
}
</style>

