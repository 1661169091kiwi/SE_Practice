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
          <input type="text" class="form-input" v-model="eventName" placeholder="请输入">
        </div>

        <!-- 时间选择 -->
        <div class="form-item">
          <label class="form-label">时间:</label>
          <div class="time-group">
            <input type="text" class="time-input" v-model="eventTime.year" >年
            <input type="text" class="time-input" v-model="eventTime.month" >月
            <input type="text" class="time-input" v-model="eventTime.day" >日
          </div>
        </div>

        <!-- 地点 -->
        <div class="form-item">
          <label class="form-label">地点:</label>
          <input type="text" class="form-input" v-model="eventLocation" placeholder="请输入">
        </div>

        <!-- 发布按钮 -->
        <div class="submit-btn-group">
          <button class="submit-btn" @click = "submitInfo">发布信息</button>
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
      activeMenu: 'info',
      // 表单数据绑定
      eventName: '',
      eventTime: {
        year: '',
        month: '',
        day: ''
      },
      eventLocation: ''
    }
  },
  methods: {
    // 处理菜单点击：更新激活状态 + 路由跳转
    handleMenuClick(menuKey, path) {
      this.activeMenu = menuKey // 更新激活的菜单
      this.$router.push(path) // 路由跳转
    },
    async submitInfo() {
      // 表单验证
      if (!this.eventName) {
        alert('请输入赛事名称')
        return
      }
      if (!this.eventTime.year || !this.eventTime.month || !this.eventTime.day) {
        alert('请完善日期信息')
        return
      }
      if (!this.eventLocation) {
        alert('请输入比赛地点')
        return
      }

      try {
        // 构造请求数据，与后端dto.PublishEventReq对应
        const requestData = {
          type: this.activeMatchType,
          name: this.eventName,
          time: {
            year: this.eventTime.year,
            month: this.eventTime.month,
            day: this.eventTime.day
          },
          location: this.eventLocation
        }

        // 调用后端发布赛事接口
        const response = await this.$http.post('/api/publish-event', requestData)

        // 处理成功响应
        if (response.data.code === 200) {
          alert('信息发布成功！赛事ID: ' + response.data.data.eventId)
          // 清空表单
          this.eventName = ''
          this.eventTime = { year: '', month: '', day: '' }
          this.eventLocation = ''
        } else {
          alert('发布失败: ' + (response.data.message || '未知错误'))
        }
      } catch (error) {
        console.error('发布失败', error)
        alert('发布失败，请检查网络或服务器状态')
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
  padding: 1.042vw; /* 20px→20/19.2≈1.042vw */
  box-sizing: border-box;
  position: relative; /* 新增：为按钮定位提供参考 */
}

/* 表单样式 */
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
  font-size: 0.729vw; /* 20px→20/19.2≈1.042vw */
  width: 10.417vw; /* 200px→200/19.2≈10.417vw */
  height: 1.458vw; /* 28px→28/19.2≈1.458vw */
  padding: 0 0.313vw; /* 6px→0.313vw */
  border: 0.052vw solid #ccc; /* 1px→0.052vw */
  box-sizing: border-box;
}

/* 比赛类型按钮组 */
.match-type-group {
  display: flex;
  gap: 0.521vw; /* 10px→0.521vw */
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
  font-size: 0.729vw; /* 20px→20/19.2≈1.042vw */
  display: flex;
  gap: 0.417vw; /* 8px→0.417vw */
  align-items: center;
}
.time-input {
  width: 2.083vw; /* 40px→40/19.2≈2.083vw */
  height: 1.458vw; /* 28px→1.458vw */
  padding: 0 0.313vw; /* 6px→0.313vw */
  border: 0.052vw solid #ccc; /* 1px→0.052vw */
  box-sizing: border-box;
  text-align: center;
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
