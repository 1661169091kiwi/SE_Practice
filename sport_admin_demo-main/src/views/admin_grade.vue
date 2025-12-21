<template>
  <!-- 保持原有模板布局不变 -->
  <div class="admin-container">
    <div class="header">
      <h1 class="header-title">体育赛事系统管理员界面</h1>
      <div class="user-info">你好, xxx</div>
    </div>
    <div class="main-layout">
      <div class="sidebar">
        <div class="menu-item" :class="{ active: activeMenu === 'info' }" @click="handleMenuClick('info', '/admin')">
          <span class="menu-icon">&#xe63c;</span>
          <span class="menu-text">赛事基本信息发布</span>
        </div>
        <div class="menu-item" :class="{ active: activeMenu === 'arrange' }" @click="handleMenuClick('arrange', '/arrange')">
          <span class="menu-icon">&#xe62a;</span>
          <span class="menu-text">赛事编排管理</span>
        </div>
        <div class="menu-item" :class="{ active: activeMenu === 'team' }" @click="handleMenuClick('team', '/team')">
          <span class="menu-icon">&#xe61e;</span>
          <span class="menu-text">参赛队伍管理</span>
        </div>
        <div class="menu-item" :class="{ active: activeMenu === 'rule' }" @click="handleMenuClick('rule', '/rule')">
          <span class="menu-icon">&#xe60e;</span>
          <span class="menu-text">比赛规则配置</span>
        </div>
        <div class="menu-item" :class="{ active: activeMenu === 'score' }" @click="handleMenuClick('score', '/grade')">
          <span class="menu-icon">&#xe628;</span>
          <span class="menu-text">成绩录入与管理</span>
        </div>
      </div>
      <div class="content-area">
        <div class="form-item">
          <label class="form-label">比赛类型:</label>
          <div class="match-type-group">
            <button v-for="(type, idx) in matchTypes" :key="idx" class="type-btn" :class="{ active: activeMatchType === type }" @click="activeMatchType = type">
              {{ type }}
            </button>
          </div>
        </div>
        <div class="table-container">
          <table class="grade-table">
            <thead>
              <tr>
                <th>赛事名</th>
                <th>参赛队伍1</th>
                <th>参赛队伍2</th>
                <th>时间</th>
                <th>地点</th>
                <th>比分</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, rowIdx) in tableData" :key="rowIdx">
                <td v-for="(cell, colIdx) in row" :key="colIdx" @click="enterEditMode(rowIdx, colIdx)" class="editable-cell">
                  <span class = "column-input" v-if="!isEditing(rowIdx, colIdx)">{{ cell.value || '-' }}</span>
                  <input class = "column-input" v-else ref="editInput" type="text" v-model="editValue" @blur="saveEdit(rowIdx, colIdx)" @keyup.enter="saveEdit(rowIdx, colIdx)" @keyup.esc="cancelEdit">
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="submit-btn-group">
          <button class="submit-btn" @click="saveGrade">保存成绩</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
axios.defaults.baseURL = 'http://localhost:8080'

export default {
  data() {
    return {
      menuList: [
        { label: '赛事基本信息发布', icon: '&#xe63c;' },
        { label: '赛事编排管理', icon: '&#xe62a;' },
        { label: '参赛队伍管理', icon: '&#xe61e;' },
        { label: '比赛规则配置', icon: '&#xe60e;' },
        { label: '成绩录入与管理', icon: '&#xe628;' }
      ],
      activeMenu: 'score',
      matchTypes: ['足球', '篮球', '羽毛球', '排球', '水上运动'],
      activeMatchType: '足球',
      tableData: Array(7).fill().map(() =>
        Array(6).fill().map(() => ({ value: '' }))
      ),
      editing: { row: -1, col: -1 },
      editValue: ''
    }
  },
  methods: {
    handleMenuClick(menuKey, path) {
      this.activeMenu = menuKey
      this.$router.push(path)
    },
    isEditing(row, col) {
      return this.editing.row === row && this.editing.col === col
    },
    enterEditMode(row, col) {
      this.editValue = this.tableData[row][col].value
      this.editing.row = row
      this.editing.col = col
      this.$nextTick(() => {
        this.$refs.editInput?.focus()
      })
    },
    saveEdit(row, col) {
      if (this.isEditing(row, col)) {
        this.tableData[row][col].value = this.editValue
        this.editing.row = -1
        this.editing.col = -1
      }
    },
    cancelEdit() {
      this.editing.row = -1
      this.editing.col = -1
    },
    // 核心修改：调用后端成绩录入接口
    async saveGrade() {
      const validData = this.tableData.filter(row =>
        row.some(cell => cell.value.trim() !== '')
      )
      if (validData.length === 0) {
        alert('请填写至少一行成绩数据！')
        return
      }

      try {
        const requestData = {
          type: this.activeMatchType,
          data: this.tableData.map(row =>
            row.map(cell => ({ value: cell.value.trim() }))
          )
        }
        // 调用后端 /api/save-grade 接口
        const response = await axios.post('/api/save-grade', requestData)

        if (response.data.code === 200) {
          alert('成绩保存成功！')
          this.tableData = Array(7).fill().map(() =>
            Array(6).fill().map(() => ({ value: '' }))
          )
        } else {
          alert('保存失败：' + (response.data.message || '未知错误'))
        }
      } catch (error) {
        console.error('成绩保存失败：', error)
        alert('保存失败，请检查网络或后端服务！')
      }
    }
  },
  mounted() {
    const currentPath = this.$route.path
    switch (currentPath) {
      case '/admin': this.activeMenu = 'info'; break
      case '/arrange': this.activeMenu = 'arrange'; break
      case '/team': this.activeMenu = 'team'; break
      case '/rule': this.activeMenu = 'rule'; break
      case '/grade': this.activeMenu = 'score'; break
    }
  },
  watch: {
    $route(to) {
      switch (to.path) {
        case '/admin': this.activeMenu = 'info'; break
        case '/arrange': this.activeMenu = 'arrange'; break
        case '/team': this.activeMenu = 'team'; break
        case '/rule': this.activeMenu = 'rule'; break
        case '/grade': this.activeMenu = 'score'; break
      }
    }
  }
}
</script>

<style scoped>
/* 保持原有样式完全不变 */
.admin-container {
  width: 100vw;
  height: 100vh;
  border: 0.104vw solid #007bff;
  box-sizing: border-box;
  font-family: "微软雅黑", sans-serif;
}
.header {
  background-color: #b3e5fc;
  padding: 0.417vw 0.833vw;
  position: relative;
  border-bottom: 0.052vw solid #007bff;
}
.header-title {
  font-size: 0.938vw;
  font-weight: bold;
  text-align: center;
  margin: 0;
}
.user-info {
  position: absolute;
  right: 0.833vw;
  top: 50%;
  transform: translateY(-50%);
  font-size: 0.729vw;
}
.main-layout {
  display: flex;
  height: calc(100% - 2.188vw);
}
.sidebar {
  width: 7.813vw;
  background-color: #f0f0f0;
  border-right: 0.052vw solid #007bff;
}
.menu-item {
  padding: 0.521vw 0.625vw;
  font-size: 0.729vw;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.313vw;
}
.menu-item.active {
  background-color: #e3f2fd;
  color: #2196f3;
}
.menu-icon {
  font-family: "iconfont";
}
.content-area {
  flex: 1;
  padding: 1.042vw;
  box-sizing: border-box;
  position: relative;
}
.form-item {
  margin-bottom: 1.042vw;
  display: flex;
  align-items: center;
}
.form-label {
  width: 4.167vw;
  text-align: left;
  font-size: 0.729vw;
  margin-right: 0.521vw;
}
.match-type-group {
  display: flex;
  gap: 0.521vw;
}
.type-btn {
  padding: 0.208vw 0.625vw;
  border: 0.052vw solid #007bff;
  background-color: #fff;
  cursor: pointer;
  font-size: 0.729vw;
}
.type-btn.active {
  background-color: #007bff;
  color: #fff;
}
.table-container {
  margin: 1.042vw 0;
  overflow: auto;
}
.grade-table {
  width: 100%;
  border-collapse: collapse;
}
.grade-table th, .grade-table td {
  border: 0.052vw solid #ccc;
  padding: 0.521vw;
  text-align: center;
  font-size: 0.729vw;
}
.grade-table th {
  background-color: #f5f5f5;
}
.editable-cell {
  cursor: pointer;
}
.editable-cell input {
  width: 100%;
  padding: 0.208vw;
  border: 0.052vw solid #007bff;
  box-sizing: border-box;
}
.submit-btn-group {
  position: absolute;
  right: 1.042vw;
  bottom: 1.042vw;
}
.submit-btn {
  padding: 0.313vw 1.042vw;
  background-color: #ccc;
  border: none;
  cursor: pointer;
  font-size: 0.729vw;
}
.column-input {
  font-size: 0.729vw; /* 14px→0.729vw */

  text-align: center;
}
</style>
