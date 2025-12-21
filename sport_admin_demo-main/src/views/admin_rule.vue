<template>
  <div class="admin-container">
    <!-- 顶部标题栏 -->
    <div class="header">
      <h1 class="header-title">体育赛事系统管理员界面</h1>
      <div class="user-info">你好, xxx</div>
    </div>

    <!-- 主体布局 -->
    <div class="main-layout">
      <!-- 侧边导航栏（带路由跳转） -->
      <div class="sidebar">
        <!-- 赛事基本信息发布 → /admin -->
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

        <!-- 参赛队伍管理 → /team（当前选中） -->
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

      <!-- 参赛队伍管理内容区域 -->
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

        <!-- 队伍人数 -->
        <div class="form-item">
          <label class="form-label">队伍人数:</label>
          <input type="text" class="form-input" placeholder="请输入人数" >
        </div>

        <!-- 时间选择 -->
        <div class="form-item">
          <label class="form-label">时间:</label>
          <div class="time-group">
            <input type="text" class="time-input" placeholder="" v-model = "eventTime.year">年
            <input type="text" class="time-input" placeholder="" v-model = "eventTime.month">月
            <input type="text" class="time-input" placeholder="" v-model = "eventTime.day">日
          </div>
        </div>

        <!-- 获胜条件 -->
        <div class="form-item">
          <label class="form-label">获胜条件:</label>
          <input type="text" class="form-input form-input-long" placeholder="格式待定" v-model = "winCondition">
        </div>

        <!-- 保存按钮 -->
        <div class="submit-btn-group">
          <button class="submit-btn" @click = "saveInfo">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
axios.defaults.baseURL = 'http://localhost:5173'

export default {
  data() {
    return {
      // 当前激活的菜单（默认选中参赛队伍管理）
      activeMenu: 'team',
      // 比赛类型选项
      matchTypes: ['足球', '篮球', '羽毛球', '排球', '水上运动'],
      // 当前选中的比赛类型
      activeMatchType: '足球',
      eventTime:{
        year:'',
        month:'',
        day:''
      },
      winCondition:''
    }
  },
  methods: {
    // 处理菜单点击：更新激活状态 + 路由跳转
    handleMenuClick(menuKey, path) {
      this.activeMenu = menuKey
      this.$router.push(path)
    },
    async saveInfo(){
      if(this.eventTime.year == ''||this.eventTime.month == ''||this.eventTime.day==''){
        alert("请完整填写时间信息")
        return
      }
      if(this.winCondition == ''){
        alert("请填写获胜条件")
        return
      }
      try{
        const requestData = {
          Type: this.activeMatchType,
          MatchTime:{
            year:this.eventTime.year,
            month:this.eventTime.month,
            day:this.eventTime.day
          },
          WinCondition:this.winCondition
        }
        const response = await this.$axios.post('/api/save-rule',requestData)

        if(response.data.code === 200){
          alert("保存成功")
          this.winCondition = ''
          this.eventTime.year = ''
          this.eventTime.month = ''
          this.eventTime.day = ''
        } else{
          alert("保存失败：" + response.data.message)
        }
      }
      catch(error){
        console.error("保存规则信息时出错:", error)
        alert("保存失败，请检查网络或服务器状态")
      }
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
  border: 0.104vw solid #007bff; /* 2px → 2/19.2≈0.104vw */
  box-sizing: border-box;
  font-family: "微软雅黑", sans-serif;
}

/* 顶部标题栏 */
.header {
  background-color: #b3e5fc;
  padding: 0.417vw 0.833vw; /* 8px→0.417vw，16px→0.833vw */
  position: relative;
  border-bottom: 0.052vw solid #007bff; /* 1px→0.052vw */
}
.header-title {
  font-size: 0.938vw; /* 18px→18/19.2=0.9375≈0.938vw */
  font-weight: bold;
  text-align: center;
  margin: 0;
}
.user-info {
  position: absolute;
  right: 0.833vw; /* 16px→0.833vw */
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.729vw; /* 14px→14/19.2≈0.729vw */
}

/* 主体布局 */
.main-layout {
  display: flex;
  height: calc(100% - 2.188vw); /* 42px→42/19.2=2.1875≈2.188vw */
}

/* 侧边导航栏 */
.sidebar {
  width: 7.813vw; /* 150px→150/19.2=7.8125≈7.813vw */
  background-color: #f0f0f0;
  border-right: 0.052vw solid #007bff; /* 1px→0.052vw */
}
.menu-item {
  padding: 0.521vw 0.625vw; /* 10px→0.521vw，12px→0.625vw */
  font-size: 0.729vw; /* 14px→0.729vw */
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.313vw; /* 6px→6/19.2=0.3125≈0.313vw */
  transition: all 0.2s;
}
.menu-item.active {
  background-color: #e3f2fd;
  color: #2196f3;
  border-left: 0.156vw solid #007bff; /* 3px→3/19.2=0.15625≈0.156vw */
}
.menu-item:hover {
  background-color: #e8e8e8;
}
.menu-icon {
  font-family: "iconfont";
  font-size: 0.833vw; /* 16px→16/19.2≈0.833vw */
}

/* 内容区域 */
.content-area {
  flex: 1;
  padding: 1.042vw; /* 20px→20/19.2≈1.042vw */
  box-sizing: border-box;
  position: relative;
}

/* 表单项样式 */
.form-item {
  margin-bottom: 1.042vw; /* 20px→1.042vw */
  display: flex;
  align-items: center;
}
.form-label {
  width: 4.167vw; /* 80px→80/19.2≈4.167vw */
  text-align: left;
  font-size: 0.729vw; /* 14px→0.729vw */
  margin-right: 0.521vw; /* 10px→0.521vw */
}
.form-input {
  width: 6.25vw; /* 120px→120/19.2=6.25vw */
  height: 1.458vw; /* 28px→28/19.2≈1.458vw */
  padding: 0 0.313vw; /* 6px→0.313vw */
  border: 0.052vw solid #ccc; /* 1px→0.052vw */
  font-size: 0.729vw; /* 14px→0.729vw */
  box-sizing: border-box;
}
.form-input-long {
  width: 15.625vw; /* 300px→300/19.2=15.625vw */
}

/* 比赛类型按钮组 */
.match-type-group {
  display: flex;
  gap: 0.26vw; /* 5px→5/19.2≈0.26vw */
}
.type-btn {
  padding: 0.208vw 0.625vw; /* 4px→0.208vw，12px→0.625vw */
  border: 0.052vw solid #007bff; /* 1px→0.052vw */
  background-color: #fff;
  cursor: pointer;
  font-size: 0.729vw; /* 14px→0.729vw */
}
.type-btn.active {
  background-color: #007bff;
  color: #fff;
}

/* 时间选择组 */
.time-group {
  display: flex;
  gap: 0.417vw; /* 8px→0.417vw */
  align-items: center;
  font-size: 0.729vw; /* 14px→0.729vw */
}
.time-input {
  font-size: 0.729vw; /* 14px→0.729vw */
  width: 2.083vw; /* 40px→40/19.2≈2.083vw */
  height: 1.458vw; /* 28px→1.458vw */
  padding: 0 0.313vw; /* 6px→0.313vw */
  border: 0.052vw solid #ccc; /* 1px→0.052vw */
  box-sizing: border-box;
  text-align: center;
}

/* 保存按钮 */
.btn-group {
  display: flex;
  justify-content: center;
  margin-top: 2.083vw; /* 40px→40/19.2≈2.083vw */
}
.save-btn {
  padding: 0.313vw 1.042vw; /* 6px→0.313vw，20px→1.042vw */
  background-color: #ccc;
  border: none;
  cursor: pointer;
  font-size: 0.729vw; /* 14px→0.729vw */
}
/* 发布按钮 */
.submit-btn-group {
  position: absolute; /* 新增：绝对定位 */
  right: 1.042vw; /* 新增：右距离 */
  bottom: 1.042vw; /* 新增：下距离 */
}
.submit-btn {
  padding: 0.313vw 1.042vw; /* 6px→0.313vw，20px→1.042vw */
  background-color: #ccc;
  border: none;
  cursor: pointer;
  font-size: 0.729vw; /* 14px→0.729vw */
}
</style>
