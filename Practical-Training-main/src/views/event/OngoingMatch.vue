<template>
  <base-match :match-info="matchInfo">
    <!-- 比分区域 -->
    <div class="score-section">
      <div class="team team-left">
        <div class="team-logo">{{ matchInfo.teamA.logo }}</div>
        <div class="team-name">{{ matchInfo.teamA.name }}</div>
      </div>

      <div class="score-info">
        <div class="main-score">{{ matchInfo.scoreA }} VS {{ matchInfo.scoreB }}</div>
        <div class="half-score">{{ periodScore }}</div>
        <div class="extra-time" v-if="matchInfo.extraTime">{{ matchInfo.extraTime }}</div>
      </div>

      <div class="team team-right">
        <div class="team-logo">{{ matchInfo.teamB.logo }}</div>
        <div class="team-name">{{ matchInfo.teamB.name }}</div>
      </div>
    </div>

    <!-- 比赛事件 -->
    <div class="match-events">
      <div class="section-title">
        比赛事件
        <i class="expand-icon">▼</i>
      </div>
      <ul class="event-list">
        <li class="event-item" :class="event.type" v-for="(event, i) in matchInfo.events" :key="i">
          <span class="event-time">{{ event.time }}</span>
          <span class="event-content">{{ event.content }}</span>
        </li>
      </ul>
    </div>

    <!-- 技术统计 -->
    <div class="tech-stats">
      <div class="section-title">技术统计</div>
      <div class="stats-grid">
        <div class="stat-item" v-for="(stat, i) in matchInfo.stats" :key="i">
          <span class="stat-name">{{ stat.name }}</span>
          <span class="stat-value">{{ stat.teamA }}-{{ stat.teamB }}</span>
        </div>
      </div>
    </div>

    <!-- 其他通用区域 -->
    <div class="lineup-section">
      <div class="section-title">首发阵容</div>
      <div class="lineup-placeholder">
        <div class="empty-hint">点击查看详细阵容</div>
      </div>
    </div>

    <div class="comments-section">
      <div class="section-title">观众评论 ({{ comments.length }})</div>
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
  </base-match>
</template>

<script>
import BaseMatch from './BaseMatch.vue'
import { sportIcons, sportColors } from '@/utils/sportConfig.js'

export default {
  components: { BaseMatch },
  props: {
    sportType: {
      type: String,
      required: true,
      validator: value => ['football', 'basketball', 'tabletennis', 'badminton', 'water'].includes(value)
    },
    // 添加比赛类型属性，默认为混双
    matchType: {
      type: String,
      default: 'doubles',
      validator: value => ['singles', 'doubles'].includes(value)
    }
  },
  created() {
  console.log('当前matchType:', this.matchType); // 若输出'doubles'，则说明参数未传递成功
},
  data() {
    return {
      newComment: '',
      sportIcons,
      sportColors,
      comments: [
        { author: '球迷小四', content: '太精彩了！', likes: 10, time: '10:42' },
        { author: '球迷小五', content: '期待更多精彩表现', likes: 5, time: '9:28' }
      ]
    }
  },
  computed: {
    matchInfo() {
      // 不同运动项目的配置
      const sportConfigs = {
        football: {
          title: '足球联赛半决赛',
          sportType: '足球',
          status: 'ongoing',
          teamA: { name: '软工队', logo: '⚽' },
          teamB: { name: '计院队', logo: '⚽' },
          scoreA: 1,
          scoreB: 1,
          extraTime: '伤停补时3分钟',
          events: [
            { time: '75\'', content: '李明进球 (1-1)', type: 'green' },
            { time: '60\'', content: '张三黄牌', type: 'yellow' },
            { time: '46\'', content: '换人 李四→王五', type: 'blue' },
            { time: '45\'', content: '上半场结束', type: 'gray' }
          ],
          stats: [
            { name: '射门', teamA: 8, teamB: 6 },
            { name: '射正', teamA: 3, teamB: 2 },
            { name: '角球', teamA: 4, teamB: 2 },
            { name: '犯规', teamA: 8, teamB: 10 }
          ]
        },
        basketball: {
          title: '篮球友谊赛',
          sportType: '篮球',
          status: 'ongoing',
          teamA: { name: '经管队', logo: '🏀' },
          teamB: { name: '外语队', logo: '🏀' },
          scoreA: 56,
          scoreB: 48,
          extraTime: '',
          events: [
            { time: 'Q3 02:45', content: '张三三分球', type: 'green' },
            { time: 'Q3 05:10', content: '李四犯规', type: 'yellow' },
            { time: 'Q2 10:30', content: '暂停', type: 'blue' }
          ],
          stats: [
            { name: '得分', teamA: 56, teamB: 48 },
            { name: '篮板', teamA: 22, teamB: 18 },
            { name: '助攻', teamA: 8, teamB: 12 },
            { name: '抢断', teamA: 5, teamB: 3 }
          ]
        },
        // 乒乓球、羽毛球、水上运动配置类似，省略...
        tabletennis: {
          title: '乒乓球团体赛',
          sportType: '乒乓球',
          status: 'ongoing',
          teamA: { name: '研究生院队', logo: '🏓' },
          teamB: { name: '本科生队', logo: '🏓' },
          scoreA: 2,
          scoreB: 1,
          extraTime: '',
          events: [
            { time: '第3场', content: '王浩 3-1 战胜 张伟', type: 'green' },
            { time: '第2场', content: '李娜 0-3 不敌 赵敏', type: 'red' },
            { time: '第1场', content: '双打组合获胜', type: 'green' }
          ],
          stats: [
            { name: '总胜场', teamA: 2, teamB: 1 },
            { name: '总得分', teamA: 58, teamB: 42 },
            { name: '失误数', teamA: 12, teamB: 18 },
            { name: '扣杀得分', teamA: 15, teamB: 9 }
          ]
        },
        badminton: {
          // 根据比赛类型返回不同配置
          ...(this.matchType === 'singles'
            ? {
                title: '羽毛球单打决赛',
                sportType: '羽毛球',
                status: 'ongoing',
                teamA: { name: '林丹', logo: '🏸' },  // 单打为个人
                teamB: { name: '李宗伟', logo: '🏸' },
                scoreA: 1,
                scoreB: 1,
                extraTime: '决胜局进行中',
                events: [
                  { time: '第3局 19-18', content: '林丹正手劈杀得分', type: 'green' },
                  { time: '第3局 15-15', content: '李宗伟网前放小球得分', type: 'red' },
                  { time: '第2局', content: '李宗伟 21-18 扳平', type: 'red' },
                  { time: '第1局', content: '林丹 21-16 获胜', type: 'green' },
                  { time: '第1局 10-8', content: '林丹连续扣杀得分', type: 'green' }
                ],
                stats: [
                  { name: '局分', teamA: 1, teamB: 1 },
                  { name: '发球得分', teamA: 6, teamB: 5 },
                  { name: '网前得分', teamA: 9, teamB: 12 },
                  { name: '扣杀得分', teamA: 18, teamB: 14 },  // 单打扣杀更多
                  { name: '失误', teamA: 7, teamB: 6 }
                ]
              }
            : {
                title: '羽毛球混双决赛',
                sportType: '羽毛球',
                status: 'ongoing',
                teamA: { name: '赵阳/孙梅', logo: '🏸' },
                teamB: { name: '钱峰/周琳', logo: '🏸' },
                scoreA: 1,
                scoreB: 1,
                extraTime: '决胜局进行中',
                events: [
                  { time: '第3局 18-17', content: '赵阳扣杀得分', type: 'green' },
                  { time: '第2局', content: '钱峰/周琳 21-19 扳平', type: 'red' },
                  { time: '第1局', content: '赵阳/孙梅 21-15 获胜', type: 'green' }
                ],
                stats: [
                  { name: '局分', teamA: 1, teamB: 1 },
                  { name: '发球得分', teamA: 8, teamB: 6 },
                  { name: '网前得分', teamA: 12, teamB: 15 },
                  { name: '失误', teamA: 5, teamB: 7 }
                ]
              })
        },
        water: {
          title: '游泳4x100米接力',
  sportType: '水上运动',
  status: 'ongoing',
  teamA: { name: '红队', logo: '🏊' },
  teamB: { name: '蓝队', logo: '🏊' },
  scoreA: 3,
  scoreB: 2,
  extraTime: '',
  events: [
    { time: '第3棒', content: '红队反超1.2秒', type: 'green' },
    { time: '第2棒', content: '蓝队领先0.8秒', type: 'red' },
    { time: '第1棒', content: '双方持平', type: 'blue' }
  ],
  stats: [
    { name: '当前棒次', teamA: '第4棒', teamB: '第4棒' },
    { name: '累计时间', teamA: '2:58.3', teamB: '2:59.1' },
    { name: '交接次数', teamA: 3, teamB: 3 },
    { name: '最佳单棒', teamA: '42.1秒', teamB: '41.8秒' }
  ]
        }
      }
      return sportConfigs[this.sportType]
    },
    periodScore() {
      // 根据不同运动显示不同的时段比分
      const periodMap = {
        football: '半场1-1',
        basketball: '第三节 56-48',
        tabletennis: '局分 2-1',
        badminton: this.matchType === 'singles'
          ? '决胜局 19-18'
          : '局分 1-1',
        water: '第3圈 1:24.5-1:26.3'
      }
      return periodMap[this.sportType]
    }
  },
  methods: {
    publishComment() {
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
/* 原有样式保持不变，根据需要调整运动项目特有样式 */
/* 比分区域 */
.score-section {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 20px 0;
  border-bottom: 1px solid #eee;
}
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
