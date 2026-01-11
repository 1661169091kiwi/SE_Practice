<template>
  <div class="search-page">
    <!-- 顶部搜索栏 -->
    <div class="search-box">
      <button class="back-btn" @click="goBack">
        <i class="arrow">←</i>
      </button>
      <div class="search-select">
        <select v-model="selectedType">
          <option value="赛事名称">赛事名称</option>
          <option value="运动员">运动员</option>
          <option value="球队">球队</option>
        </select>
      </div>
      <div class="search-input-group">
        <input
          type="text"
          placeholder="搜索"
          v-model="searchKeyword"
          @input="handleInput"
          autofocus
        />
        <button class="search-btn" @click="handleSearch">搜索</button>
      </div>
    </div>

    <!-- 历史记录区 -->
    <div class="history-section" v-if="searchHistory.length > 0">
      <div class="history-header">
        <span>历史记录</span>
        <button class="clear-btn" @click="clearHistory">
          <i class="trash">🗑️</i>
        </button>
      </div>
      <div class="history-tags">
        <button
          v-for="(item, index) in searchHistory"
          :key="index"
          class="history-tag"
          @click="searchByHistory(item)"
        >
          {{ item }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
//import { useRouter } from 'vue-router'

export default {
  name: "SearchPage",
  setup() {
    // const router = useRouter()

    // const goBack = () => {
    //   router.back()
    // }

    // return {
    //   goBack
    // }
  },
  data() {
    return {
      selectedType: "赛事名称", // 默认选中的筛选类型
      searchKeyword: "", // 搜索输入框内容
      searchHistory: ["羽毛球", "杨航"], // 历史记录默认数据
    };
  },
  methods: {
    // 处理输入变化
    handleInput() {
      // 可在此添加实时搜索逻辑
    },
    // 执行搜索
    handleSearch() {
      if (this.searchKeyword && !this.searchHistory.includes(this.searchKeyword)) {
        this.searchHistory.unshift(this.searchKeyword); // 新增历史记录到首位
      }
      console.log(`搜索类型：${this.selectedType}，关键词：${this.searchKeyword}`);
    },
    // 清空历史记录
    clearHistory() {
      this.searchHistory = [];
    },
    // 点击历史标签搜索
    searchByHistory(keyword) {
      this.searchKeyword = keyword;
      this.handleSearch();
    },
    goBack(){
        this.$router.go(-1) // 返回上一页
    }
  },
};
</script>

<style scoped>
.search-page {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  padding: 16px;
  min-height: 100vh;
  box-sizing: border-box;
}

/* 顶部搜索栏样式 */
.search-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
}
.back-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 20px;
  margin-right: 12px;
}
.search-select select {
  border: 1px solid #ddd;
  border-radius: 4px 0 0 4px;
  padding: 8px 12px;
  font-size: 14px;
  background-color: #ffffff;
}
.search-input-group {
  display: flex;
  flex: 1;
}
.search-input-group input {
  flex: 1;
  border: 1px solid #ddd;
  border-radius: 0 4px 4px 0;
  padding: 8px 12px;
  outline: none;
  font-size: 14px;
}


/* 历史记录区样式 */
.history-section {
  width: 100%;
}
.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 16px;
  font-weight: 500;
}
.clear-btn {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 20px;
}
.history-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}
.history-tag {
  background-color: #f5f5f5;
  border: 1px solid #ddd;
  border-radius: 16px;
  padding: 6px 16px;
  cursor: pointer;
  font-size: 14px;
}
.history-tag:hover {
  background-color: #e0e0e0;
}
/* 顶部搜索栏 */
.header {
  padding: 16px;
  background-color: #f8f9fa;
}

.search-box {
  display: flex;
  align-items: center;
  background-color: #e9ecef;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}
.search-box input {
  flex: 1;
  border: none;
  background-color: transparent;
  outline: none;
  font-size: 14px;
}
.search-btn {
  background-color: #4285f4;
  color: #fff;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  cursor: pointer;
}
</style>
