<template>
  <div class="match-detail-container">
    <!-- 顶部导航 -->
    <header class="top-header">
      <button class="back-btn" @click="goBack">←</button>
      <h1 class="match-title">比赛标题</h1>
      <span class="match-type">羽毛球</span>
    </header>

    <!-- 比赛状态 -->
    <div class="match-status">进行中</div>

    <!-- 比分区域 -->
    <div class="score-section">
      <div class="team team-left">
        <div class="team-logo">🏸</div>
        <div class="team-name">软工队</div>
      </div>

      <div class="score-info">
        <div class="main-score">1 VS 1</div>
        <div class="half-score">半场1-1</div>
        <div class="extra-time">伤停补时3分钟</div>
      </div>

      <div class="team team-right">
        <div class="team-logo">🏸</div>
        <div class="team-name">计院队</div>
      </div>
    </div>

    <!-- 比赛事件 -->
    <div class="match-events">
      <div class="section-title">
        比赛事件
        <i class="expand-icon">▼</i>
      </div>
      <ul class="event-list">
        <li class="event-item green">
          <span class="event-time">75'</span>
          <span class="event-content">李明进球 (1-1)</span>
        </li>
        <li class="event-item yellow">
          <span class="event-time">60'</span>
          <span class="event-content">张三黄牌</span>
        </li>
        <li class="event-item blue">
          <span class="event-time">46'</span>
          <span class="event-content">换人 李四→王五</span>
        </li>
        <li class="event-item gray">
          <span class="event-time">45'</span>
          <span class="event-content">上半场结束</span>
        </li>
      </ul>
    </div>

    <!-- 技术统计 -->
    <div class="tech-stats">
      <div class="section-title">技术统计</div>
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

    <!-- 首发阵容 -->
    <div class="lineup-section">
      <div class="section-title">首发阵容</div>
      <div class="lineup-placeholder">
        <!-- 实际项目中可动态渲染首发球员列表 -->
        <div class="empty-hint">点击查看详细阵容</div>
      </div>
    </div>

    <!-- 观众评论 -->
    <div class="comments-section">
      <div class="section-title">观众评论 (28)</div>
      <div class="comments-list">
        <div class="comment-item" v-for="(comment, index) in comments" :key="index">
          <div class="comment-author">{{ comment.author }}：</div>
          <div class="comment-content">{{ comment.content }}</div>
          <div class="comment-actions">
            <span class="like-btn">👍 {{ comment.likes }}</span>
            <span class="comment-time">{{ comment.time }}</span>
          </div>
        </div>
      </div>

      <!-- 评论输入框 -->
      <div class="comment-input-area">
        <input
          type="text"
          placeholder="说点什么..."
          v-model="newComment"
          class="comment-input"
        >
        <button class="publish-btn" @click="publishComment">发布</button>
      </div>
    </div>
  </div>
</template>

<script>


export default {
  name: 'OngoingMatchDetail',
  data(){
    return{
        newComment: '',
        comments:[
          {
            author: '球迷小四',
            content: '球员太努力了',
            likes: 10,
            time: '10:42'
          },
          {
            author: '球迷小五',
            content: '可以踢一辈子球吗？',
            likes: 10,
            time: '9:28'
          }
        ]
    }
  },
  methods:{
    // 返回上一页
    goBack(){
      this.$router.back()
    },
    publishComment(){
      if (this.newComment.trim()) {
        this.comments.unshift({
          author: '我',
          content: this.newComment.trim(),
          likes: 0,
          time: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
        })
        this.newComment = ''
      }
    }
  }

}
</script>

<style scoped>
.match-detail-container {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  background-color: #fff;
  min-height: 100vh;
  padding-bottom: 80px; /* 预留底部评论区空间 */
}

/* 顶部导航 */
.top-header {
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
  position: absolute;
  right: 16px;
  color: #888;
  font-size: 14px;
}

/* 比赛状态 */
.match-status {
  color: #f5222d; /* 红色表示进行中 */
  font-size: 16px;
  text-align: center;
  padding: 12px 0;
  font-weight: 500;
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
  font-size: 32px;
  font-weight: bold;
  margin: 0 0 4px 0;
}

.half-score {
  font-size: 14px;
  color: #666;
  margin: 0 0 4px 0;
}

.extra-time {
  font-size: 12px;
  color: #f5222d;
  margin: 0;
}

/* 通用区域标题 */
.section-title {
  font-size: 16px;
  font-weight: 600;
  padding: 12px 16px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.expand-icon {
  font-size: 14px;
  color: #888;
}

/* 比赛事件 */
.event-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.event-item {
  padding: 12px 16px;
  border-bottom: 1px solid #f5f5f5;
  display: flex;
  align-items: center;
}

.event-time {
  min-width: 40px;
  font-weight: 500;
  margin-right: 12px;
  color: #666;
}

.event-content {
  flex: 1;
}

/* 事件类型颜色 */
.green {
  background-color: rgba(76, 175, 80, 0.05);
}
.green .event-content {
  color: #2e7d32;
}

.yellow {
  background-color: rgba(255, 193, 7, 0.05);
}
.yellow .event-content {
  color: #ff8f00;
}

.blue {
  background-color: rgba(33, 150, 243, 0.05);
}
.blue .event-content {
  color: #1565c0;
}

.gray {
  background-color: rgba(158, 158, 158, 0.05);
}
.gray .event-content {
  color: #616161;
}

/* 技术统计 */
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
}

.stat-name {
  color: #666;
}

.stat-value {
  font-weight: 500;
}

/* 首发阵容 */
.lineup-placeholder {
  padding: 24px 16px;
  text-align: center;
  color: #888;
  font-size: 14px;
}

/* 观众评论 */
.comments-list {
  padding: 16px;
}

.comment-item {
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f5f5f5;
}

.comment-author {
  font-weight: 500;
  margin-bottom: 4px;
}

.comment-content {
  margin-bottom: 8px;
  color: #333;
}

.comment-actions {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #888;
}

.like-btn {
  cursor: pointer;
}

/* 评论输入区 */
.comment-input-area {
  display: flex;
  padding: 12px 16px;
  border-top: 1px solid #eee;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #fff;
}

.comment-input {
  flex: 1;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 20px;
  outline: none;
  font-size: 14px;
}

.publish-btn {
  margin-left: 8px;
  padding: 8px 16px;
  background-color: #1677ff;
  color: #fff;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
}
</style>
