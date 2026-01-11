<template>
  <div class="admin-container">
    <!-- 顶部标题栏 -->
    <div class="header">
      <h1 class="header-title">体育赛事系统管理员界面</h1>
      <div class="user-info">你好, xxx</div>
    </div>

    <!-- 主体布局 -->
    <div class="main-layout">
      <!-- 侧边导航栏 -->
      <div class="sidebar">
        <div
          class="menu-item"
          :class="{ active: activeMenu === 'info' }"
          @click="handleMenuClick('info', '/admin')"
        >
          <span class="menu-icon">&#xe63c;</span>
          <span class="menu-text">赛事基本信息发布</span>
        </div>

        <!-- 赛事编排管理 → /arrange -->
        <div
          class="menu-item"
          :class="{ active: activeMenu === 'arrange' }"
          @click="handleMenuClick('arrange', '/arrange')"
        >
          <span class="menu-icon">&#xe62a;</span>
          <span class="menu-text">赛事编排管理</span>
        </div>

        <!-- 参赛队伍管理 → /team -->
        <div
          class="menu-item"
          :class="{ active: activeMenu === 'team' }"
          @click="handleMenuClick('team', '/team')"
        >
          <span class="menu-icon">&#xe61e;</span>
          <span class="menu-text">参赛队伍管理</span>
        </div>

        <!-- 比赛规则配置 → /rule -->
        <div
          class="menu-item"
          :class="{ active: activeMenu === 'rule' }"
          @click="handleMenuClick('rule', '/rule')"
        >
          <span class="menu-icon">&#xe60e;</span>
          <span class="menu-text">比赛规则配置</span>
        </div>

        <!-- 成绩录入与管理 → /grade -->
        <div
          class="menu-item"
          :class="{ active: activeMenu === 'score' }"
          @click="handleMenuClick('score', '/grade')"
        >
          <span class="menu-icon">&#xe628;</span>
          <span class="menu-text">成绩录入与管理</span>
        </div>
      </div>

      <!-- 内容区域（赛事基本信息发布表单） -->
      <div class="content-area">
        <!-- 比赛类型选择 -->
        <div class="form-item">
          <label class="form-label">比赛类型:</label>
          <div class="match-type-group">
            <button
              v-for="(type, idx) in matchTypes"
              :key="idx"
              class="type-btn"
              :class="{ active: activeMatchType === type }"
              @click="activeMatchType = type"
            >
              {{ type }}
            </button>
          </div>
        </div>

        <!-- 赛事名称 -->
        <div class="form-item">
          <label class="form-label">赛事名称:</label>
          <input type="text" class="form-input" placeholder="请输入">
        </div>

        <!-- 时间选择 -->
        <div class="form-item">
          <label class="form-label">时间:</label>
          <div class="time-group">
            <input type="text" class="time-input" placeholder="">月
            <input type="text" class="time-input" placeholder="">日
          </div>
        </div>

        <!-- 地点 -->
        <div class="form-item">
          <label class="form-label">地点:</label>
          <input type="text" class="form-input" placeholder="请输入">
        </div>

        <!-- 发布按钮 -->
        <div class="submit-btn-group">
          <button class="submit-btn">发布信息</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      // 侧边菜单数据
      menuList: [
        { label: '赛事基本信息发布', icon: '&#xe63c;' },
        { label: '赛事编排管理', icon: '&#xe62a;' },
        { label: '参赛队伍管理', icon: '&#xe61e;' },
        { label: '比赛规则配置', icon: '&#xe60e;' },
        { label: '成绩录入与管理', icon: '&#xe628;' }
      ],
      activeMenuIndex: 0,
      // 比赛类型选项
      matchTypes: ['足球', '篮球', '羽毛球', '排球', '水上运动'],
      // 当前选中的比赛类型
      activeMatchType: '足球',
      activeMenu: 'info'
    }
  },
  methods: {
    // 处理菜单点击：更新激活状态 + 路由跳转
    handleMenuClick(menuKey, path) {
      this.activeMenu = menuKey // 更新激活的菜单
      this.$router.push(path) // 路由跳转
    }
  },
  mounted() {
    // 初始化时根据当前路由设置激活菜单
    const currentPath = this.$route.path
    switch (currentPath) {
      case '/admin':
        this.activeMenu = 'info'
        break
      case '/arrange':
        this.activeMenu = 'arrange'
        break
      case '/team':
        this.activeMenu = 'team'
        break
      case '/rule':
        this.activeMenu = 'rule'
        break
      case '/grade':
        this.activeMenu = 'score'
        break
    }
  },
  watch: {
    // 监听路由变化，同步更新激活菜单
    $route(to) {
      switch (to.path) {
        case '/admin':
          this.activeMenu = 'info'
          break
        case '/arrange':
          this.activeMenu = 'arrange'
          break
        case '/team':
          this.activeMenu = 'team'
          break
        case '/rule':
          this.activeMenu = 'rule'
          break
        case '/grade':
          this.activeMenu = 'score'
          break
      }
    }
  }
}
</script>

<style scoped>
/* 全局容器 */
.admin-container {
  width: 100vw;
  height: 100vh;
  border: 2px solid #007bff;
  box-sizing: border-box;
  font-family: "微软雅黑", sans-serif;
}

/* 顶部标题栏 */
.header {
  background-color: #b3e5fc;
  padding: 8px 16px;
  position: relative;
  border-bottom: 1px solid #007bff;
}
.header-title {
  font-size: 18px;
  font-weight: bold;
  text-align: center;
  margin: 0;
}
.user-info {
  position: absolute;
  right: 16px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
}

/* 主体布局 */
.main-layout {
  display: flex;
  height: calc(100% - 42px);
}

/* 侧边导航栏 */
.sidebar {
  width: 150px;
  background-color: #f0f0f0;
  border-right: 1px solid #007bff;
}
.menu-item {
  padding: 10px 12px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
}
.menu-item.active {
  background-color: #e3f2fd;
  color: #2196f3;
}
.menu-icon {
  font-family: "iconfont";
}

/* 内容区域 */
.content-area {
  flex: 1;
  padding: 20px;
  box-sizing: border-box;
}

/* 表单样式 */
.form-item {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}
.form-label {
  width: 80px;
  text-align: left;
  font-size: 14px;
  margin-right: 10px;
}
.form-input {
  width: 200px;
  height: 28px;
  padding: 0 6px;
  border: 1px solid #ccc;
  box-sizing: border-box;
}

/* 比赛类型按钮组 */
.match-type-group {
  display: flex;
  gap: 10px;
}
.type-btn {
  padding: 4px 12px;
  border: 1px solid #007bff;
  background-color: #fff;
  cursor: pointer;
  font-size: 14px;
}
.type-btn.active {
  background-color: #007bff;
  color: #fff;
}

/* 时间选择组 */
.time-group {
  display: flex;
  gap: 8px;
  align-items: center;
}
.time-input {
  width: 40px;
  height: 28px;
  padding: 0 6px;
  border: 1px solid #ccc;
  box-sizing: border-box;
  text-align: center;
}

/* 发布按钮 */
.submit-btn-group {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}
.submit-btn {
  padding: 6px 20px;
  background-color: #ccc;
  border: none;
  cursor: pointer;
  font-size: 14px;
}
</style>
