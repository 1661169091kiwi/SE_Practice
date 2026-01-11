<template>
  <div class="player-page">
    <header class="top-nav">
      <button class="back-btn" @click="goBack">←</button>
      <h1 class="page-title">软工队球员列表</h1>
    </header>

    <div class="team-summary">
      <div class="team-logo">⚽</div>
      <div class="team-info">
        <div class="team-name">软件工程学院足球队</div>
        <div class="team-stats">
          <span>成立时间：2022年</span>
          <span>成员数量：22人</span>
          <span>胜率：58%</span>
        </div>
      </div>
    </div>

    <div class="section-header">
      <h2>球员信息</h2>
      <span class="count">{{ players.length }}名球员</span>
    </div>

    <div class="players-list">
      <div class="player-card" v-for="(player, index) in players" :key="index">
        <div class="player-avatar">
          <img :src="player.avatar" alt="球员头像" class="avatar-img">
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
            <div class="stat-label">进球数</div>
            <div class="stat-value">{{ player.goals }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">助攻</div>
            <div class="stat-value">{{ player.assists }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">出场</div>
            <div class="stat-value">{{ player.appearances }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'FootballTeam',
  data() {
    return {
      players: [
        {
          name: '张明',
          number: 1,
          position: '门将',
          height: '192cm',
          weight: '85kg',
          age: 22,
          goals: 0,
          assists: 0,
          appearances: 15,
          avatar: 'https://via.placeholder.com/60?text=1'
        },
        {
          name: '李军',
          number: 5,
          position: '中后卫',
          height: '188cm',
          weight: '82kg',
          age: 21,
          goals: 2,
          assists: 1,
          appearances: 14,
          avatar: 'https://via.placeholder.com/60?text=5'
        },
        {
          name: '王浩',
          number: 10,
          position: '前腰',
          height: '178cm',
          weight: '70kg',
          age: 20,
          goals: 8,
          assists: 10,
          appearances: 15,
          avatar: 'https://via.placeholder.com/60?text=10'
        },
        {
          name: '赵峰',
          number: 9,
          position: '前锋',
          height: '185cm',
          weight: '78kg',
          age: 22,
          goals: 12,
          assists: 3,
          appearances: 13,
          avatar: 'https://via.placeholder.com/60?text=9'
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
