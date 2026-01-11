<template>
  <div id="app">
    <!-- 顶部导航栏 -->
    <header class="top-header">
      <div class="tab-bar">
        <button
          class="tab-btn"
          :class="{ active: currentTab === 'all' }"
          @click="currentTab = 'all'"
        >
          全部
        </button>

        <button
          class="tab-btn"
          :class="{ active: currentTab === 'all' }"
          @click="currentTab = 'all'"
        >
          订阅
        </button>
        <button
          class="tab-btn"
          :class="{ active: currentTab === 'important' }"
          @click="currentTab = 'important'"
        >
          热门
        </button>
        <button
          class="tab-btn"
          :class="{ active: currentTab === 'football' }"
          @click="currentTab = 'football'"
        >
          足球
        </button>
        <button
          class="tab-btn"
          :class="{ active: currentTab === 'basketball' }"
          @click="currentTab = 'basketball'"
        >
          篮球
        </button>
        <button
        class="tab-btn"
        :class="{ active: currentTab === 'pingpong' }"
          @click="currentTab = 'pingpong'"
        >
        乒乓球
      </button>
      <button
        class="tab-btn"
        :class="{ active: currentTab === 'badminton' }"
          @click="currentTab = 'badminton'"
        >羽毛球
      </button>
      <button
        class="tab-btn"
        :class="{ active: currentTab === 'watersports' }"
          @click="currentTab = 'watersports'"
        >水上运动
      </button>
      </div>
    </header>

    <!-- 日期区域 -->
    <div class="date-section">今天 10月27日 星期一</div>

    <!-- 比赛列表（根据当前Tab显示对应内容） -->
    <div class="match-list">
      <MatchCard v-for="(item, index) in filteredMatches" :key="index" :match="item" />
    </div>

  </div>
</template>

<script>
import MatchCard from './matchCard.vue';

export default {
  name: 'App',
  components: { MatchCard },
  data() {
    return {
      currentTab: 'all', // 默认激活"关注"Tab
      showMore: false, // 控制更多Tab的显示/隐藏
      // 模拟比赛数据（包含分类信息）
      matchData: [//羽毛球要考虑单双打_其他运动的single属性无效
        { id: 1, type: 'follow', league: '校足球联赛', home: '计算机学院', away: '管理学院', num:11, score: '2-1', status: '已结束',teamIcon: '⚽',time:'2025-10-11',single: false,sports:'football',sportsType:'足球' },
        { id: 2, type: 'football', league: '校足球联赛', home: '电子学院', away: '文学院', num:10, score: '0-0', status: '进行中',teamIcon: '⚽',time:'2025-10-11',single: false,sports:'football' ,sportsType:'足球'},
        { id: 3, type: 'important', league: '省大学生篮球赛', home: '法学院', away: '医学院', num:9, score: '56-48', status: '已结束',teamIcon: '🏀' ,time:'2025-10-11',single: false,sports:'basketball',sportsType:'篮球'},
        { id: 4, type: 'basketball', league: '校篮球联赛', home: '外语学院', away: '工学院', num:8, score: '32-28', status: '未开始',teamIcon: '🏀' ,time:'2025-10-11',single: false,sports:'basketball',sportsType:'篮球'},
        { id: 5, type: 'badminton', league: '木之本曼飛', home:'友枝小学',away:'友枝中学', num:7, score: '12-22', status: '未开始',teamIcon: '🏸',time:'2025-10-11',single: true,sports:'badminton',sportsType:'羽毛球'}
      ]
    };
  },
  computed: {
    // 根据当前Tab过滤比赛数据
    filteredMatches() {
      if (this.currentTab === 'all') {
        return this.matchData; // 全部Tab显示所有比赛
      }
      return this.matchData.filter(item => item.type === this.currentTab);
    }
  }
};
</script>

<style scoped>
#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  background-color: #f8f9fa;
  min-height: 100vh;
  color: #333;
}

/* 顶部导航栏 */
.top-header {
  background-color: #fff;
  padding: 16px 0;
  border-bottom: 1px solid #eaeaea;
}
.top-header h1 {
  font-size: 20px;
  font-weight: 600;
  text-align: center;
  margin: 0 0 12px 0;
}

/* Tab栏样式 */
.tab-bar {
  display: flex;
  justify-content: center;
  gap: 12px;
  overflow-x: auto;
  padding: 0 16px;
}
.tab-btn {
  padding: 6px 12px;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #666;
  cursor: pointer;
  white-space: nowrap;
  position: relative;
  transition: color 0.2s;
}
/* 激活状态样式 */
.tab-btn.active {
  color: #1677ff; /* 主题色 */
  font-weight: 500;
}
.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 10%;
  width: 80%;
  height: 2px;
  background-color: #1677ff;
  border-radius: 1px;
}
/* 悬停效果 */
.tab-btn:not(.active):hover {
  color: #333;
}

.tab-more {
  border: none;
  background: transparent;
  font-size: 16px;
  color: #666;
  cursor: pointer;
  padding: 0 8px;
}

/* 更多Tab选项 */
.more-tabs {
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 8px 16px;
  background-color: #fff;
  border-top: 1px solid #f0f0f0;
}

/* 日期区域 */
.date-section {
  padding: 12px 16px;
  color: #888;
  font-size: 13px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #eaeaea;
}

/* 比赛列表 */
.match-list {
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

/* 底部导航 */
.footer-nav {
  display: flex;
  justify-content: space-around;
  padding: 8px 0;
  border-top: 1px solid #eee;
  background-color: #fff;
  margin-top: auto;
}
.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  cursor: pointer;
}
.nav-item.active .nav-label {
  color: #4285f4;
}
.icon {
  font-size: 20px;
}
.nav-label {
  font-size: 12px;
  margin-top: 4px;
}
</style>
