<template>
  <!-- 保持原有模板布局不变，仅修改表格第三列（位置）和第五列（年级）的渲染逻辑 -->
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
        <div class="form-item">
          <label class="form-label">队伍名称:</label>
          <input type="text" class="form-input" v-model="teamName" placeholder="请输入队伍名称">
        </div>
        <div class="form-item">
          <label class="form-label">队长:</label>
          <input type="text" class="form-input" v-model="captain" placeholder="请输入队长姓名">
        </div>
        <div class="table-container">
          <table class="team-table">
            <thead>
              <tr>
                <th>学号</th>
                <th>姓名</th>
                <th>位置</th>
                <th>球衣号码</th>
                <th>年级</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(row, rowIdx) in tableData" :key="rowIdx">
                <td v-for="(cell, colIdx) in row" :key="colIdx"
                    @click="enterEditMode(rowIdx, colIdx)"
                    class="editable-cell"
                    :class="{ disabled: colIdx === 2 && (activeMatchType === '羽毛球' || activeMatchType === '水上运动') }">
                  <!-- 学号、姓名、球衣号码：输入框 -->
                  <template v-if="colIdx === 0 || colIdx === 1 || colIdx === 3">
                    <span v-if="!isEditing(rowIdx, colIdx)">{{ cell.value || '-' }}</span>
                    <input v-else ref="editInput" type="text" v-model="editValue"
                           @blur="saveEdit(rowIdx, colIdx)"
                           @keyup.enter="saveEdit(rowIdx, colIdx)"
                           @keyup.esc="cancelEdit">
                  </template>

                  <!-- 位置列（第2列）：下拉框，羽毛球/水上运动时禁用且为空 -->
                  <template v-else-if="colIdx === 2">
                    <span v-if="!isEditing(rowIdx, colIdx) || (activeMatchType === '羽毛球' || activeMatchType === '水上运动')">
                      {{ (activeMatchType === '羽毛球' || activeMatchType === '水上运动') ? '-' : (cell.value || '-') }}
                    </span>
                    <select v-else
                            ref="editSelect"
                            v-model="editValue"
                            @blur="saveEdit(rowIdx, colIdx)"
                            @change="saveEdit(rowIdx, colIdx)"
                            @keyup.esc="cancelEdit"
                            :disabled="activeMatchType === '羽毛球' || activeMatchType === '水上运动'">
                      <option value="">请选择位置</option>
                      <option v-for="position in getPositionOptions()" :key="position" :value="position">{{ position }}</option>
                    </select>
                  </template>

                  <!-- 年级列（第4列）：下拉框 -->
                  <template v-else>
                    <span v-if="!isEditing(rowIdx, colIdx)">{{ cell.value || '-' }}</span>
                    <select v-else ref="editGradeSelect" v-model="editValue"
                            @blur="saveEdit(rowIdx, colIdx)"
                            @change="saveEdit(rowIdx, colIdx)"
                            @keyup.esc="cancelEdit">
                      <option value="">请选择年级</option>
                      <option v-for="grade in gradeOptions" :key="grade" :value="grade">{{ grade }}</option>
                    </select>
                  </template>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="submit-btn-group">
          <button class="submit-btn" @click="saveTeam">保存队伍信息</button>
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
      activeMenu: 'team',
      matchTypes: ['足球', '篮球', '羽毛球', '排球', '水上运动'],
      activeMatchType: '足球',
      teamName: '',
      captain: '',
      // 年级下拉选项
      gradeOptions: ['大一', '大二', '大三', '大四', '研一', '研二', '研三'],
      // 各运动类型对应的位置选项（可扩展）
      positionMap: {
        足球: ['前锋', '中场', '后卫', '门将'],
        篮球: ['控球后卫', '得分后卫', '小前锋', '大前锋', '中锋'],
        排球: ['主攻', '副攻', '二传', '自由人'],
        羽毛球: [], // 空，禁用
        水上运动: [] // 空，禁用
      },
      tableData: Array(10).fill().map(() =>
        Array(5).fill().map(() => ({ value: '' }))
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
    // 根据当前比赛类型获取位置选项
    getPositionOptions() {
      return this.positionMap[this.activeMatchType] || []
    },
    enterEditMode(row, col) {
      // 羽毛球/水上运动时，禁止编辑位置列
      if (col === 2 && (this.activeMatchType === '羽毛球' || this.activeMatchType === '水上运动')) {
        return
      }
      // 初始化编辑值为当前单元格值
      this.editValue = this.tableData[row][col].value
      this.editing.row = row
      this.editing.col = col
      this.$nextTick(() => {
        // 区分不同列的焦点处理
        if (col === 2) {
          // 位置列：下拉框
          this.$refs.editSelect?.focus()
        } else if (col === 4) {
          // 年级列：下拉框
          this.$refs.editGradeSelect?.focus()
        } else {
          // 其他列：输入框
          this.$refs.editInput?.focus()
        }
      })
    },
    saveEdit(row, col) {
      if (this.isEditing(row, col)) {
        // 羽毛球/水上运动时，强制清空位置列值
        if (col === 2 && (this.activeMatchType === '羽毛球' || this.activeMatchType === '水上运动')) {
          this.tableData[row][col].value = ''
        } else {
          this.tableData[row][col].value = this.editValue
        }
        // 退出编辑模式
        this.editing.row = -1
        this.editing.col = -1
      }
    },
    cancelEdit() {
      // 取消编辑，清空编辑值并退出编辑模式
      this.editValue = ''
      this.editing.row = -1
      this.editing.col = -1
    },
    // 保存队伍信息（原有逻辑+位置列校验）
    async saveTeam() {
      if (!this.teamName.trim()) {
        alert('请填写队伍名称！')
        return
      }
      if (!this.captain.trim()) {
        alert('请填写队长姓名！')
        return
      }
      const validData = this.tableData.filter(row =>
        row.some(cell => cell.value.trim() !== '')
      )
      if (validData.length === 0) {
        alert('请填写至少一名队员信息！')
        return
      }

      // 非羽毛球/水上运动时，校验位置列必填
      if (this.activeMatchType !== '羽毛球' && this.activeMatchType !== '水上运动') {
        const invalidPosition = validData.some(row => row[2].value === '')
        if (invalidPosition) {
          alert('请为所有队员选择位置！')
          return
        }
      }

      try {
        const requestData = {
          type: this.activeMatchType,
          teamName: this.teamName.trim(),
          captain: this.captain.trim(),
          members: this.tableData.map(row =>
            row.map(cell => ({ value: cell.value.trim() }))
          )
        }
        const response = await axios.post('/api/save-team', requestData)

        if (response.data.code === 200) {
          alert('队伍信息保存成功！队伍ID：' + (response.data.data?.teamId || ''))
          this.teamName = ''
          this.captain = ''
          this.tableData = Array(10).fill().map(() =>
            Array(5).fill().map(() => ({ value: '' }))
          )
        } else {
          alert('保存失败：' + (response.data.message || '未知错误'))
        }
      } catch (error) {
        console.error('队伍保存失败：', error)
        alert('保存失败，请检查网络或后端服务！')
      }
    }
  },
  // 监听比赛类型变化，自动清空位置列（羽毛球/水上运动）
  watch: {
    activeMatchType(newType) {
      if (newType === '羽毛球' || newType === '水上运动') {
        // 清空所有位置列的值
        this.tableData.forEach(row => {
          row[2].value = ''
        })
      }
    },
    $route(to) {
      switch (to.path) {
        case '/admin': this.activeMenu = 'info'; break
        case '/arrange': this.activeMenu = 'arrange'; break
        case '/team': this.activeMenu = 'team'; break
        case '/rule': this.activeMenu = 'rule'; break
        case '/grade': this.activeMenu = 'score'; break
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
  }
}
</script>

<style scoped>
/* 保持原有样式不变，新增位置列禁用样式 */
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
.form-input {
  font-size: 0.729vw;
  width: 10.417vw;
  height: 1.458vw;
  padding: 0 0.313vw;
  border: 0.052vw solid #ccc;
  box-sizing: border-box;
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
.team-table {
  width: 100%;
  border-collapse: collapse;
}
.team-table th, .team-table td {
  border: 0.052vw solid #ccc;
  padding: 0.521vw;
  text-align: center;
  font-size: 0.729vw;
}
.team-table th {
  background-color: #f5f5f5;
}
.editable-cell {
  cursor: pointer;
}
/* 禁用状态样式 */
.editable-cell.disabled {
  cursor: not-allowed;
  background-color: #f8f8f8;
}
/* 统一输入框和下拉框样式 */
.editable-cell input, .editable-cell select {
  width: 100%;
  padding: 0.208vw;
  border: 0.052vw solid #007bff;
  box-sizing: border-box;
  font-size: 0.729vw;
  height: 1.8vw;
}
/* 禁用的下拉框样式 */
.editable-cell select:disabled {
  background-color: #eee;
  border-color: #ddd;
  color: #999;
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
</style>
