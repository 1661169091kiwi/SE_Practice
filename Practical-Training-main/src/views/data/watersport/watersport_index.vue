 <template>
    <div class="soccer_index">
       <div class="data_index">
      <div class="date_chance">
        <van-dropdown-menu active-color="#3775F6">
          <van-dropdown-item
            v-model="value1"
            :options="option1"
            background-color="#f5f5f5"
            @change="updateData"
          />
        </van-dropdown-menu>
      </div>
      <div class="my_data_opction">
        <ul @click="handleClickoption">
          <router-link to="/data/watersport/dragonboat">
            <li :class="{ active: currentoptionindex == 1 }" data-index="1">龙舟</li>
          </router-link>
          <router-link to="/data/watersport/kayak">
            <li :class="{ active: currentoptionindex == 2 }" data-index="2">皮艇</li>
          </router-link>
          <router-link to="/data/watersport/paddleboard">
            <li :class="{ active: currentoptionindex == 3 }" data-index="3">桨板</li>
          </router-link>
          <router-link to="/data/watersport/schedule">
            <li :class="{ active: currentoptionindex == 4 }" data-index="4">赛程</li>
          </router-link>
        </ul>
      </div>
    </div>
    <router-view></router-view>
    </div>

  </template>
  <script setup lang="ts">
  import { ref } from 'vue'
  const currentoptionindex = ref(1)
// 点击数据类型导航栏触发事件
const handleClickoption = (e) => {
  // currentindex.value = index.dataset.index;
  if (e.target.tagName == 'LI') {
    currentoptionindex.value = e.target.dataset.index
    //更新数据
    updateData()
  }
}
const value1 = ref(2025)
const option1 = [
  {
    text: '2025',
    value: 2025,
  },
  {
    text: '2024',
    value: 2024,
  },
  {
    text: '2023',
    value: 2023,
  },
  {
    text: '2022',
    value: 2022,
  },
  {
    text: '2021',
    value: 2021,
  },
  {
    text: '2020',
    value: 2020,
  },
  {
    text: '2019',
    value: 2019,
  },
]
// 更新数据
const updateData = () => {
  // 这里可以根据 currentindex.value 和 currentoptionindex.value 来更新数据
  // console.log('更新数据', currentindex.value, currentoptionindex.value)
}
  </script>
  <style scoped>
    /* 数据类型导航栏 */
.data_index {
  width: 100%;
  height: 12.0773vw;
  background-color: #f5f5f5;
  /* padding:4.8309vw; */
}
.date_chance {
  width: 16%;
  display: inline-block;
}
.my_data_opction {
  width: 84%;
  display: inline-block;
  padding: 0 0 0 4.8309vw;
}
.my_data_opction ul {
  display: flex;
  justify-content: start;
  gap: 4.8309vw;
}
.my_data_opction ul li {
  font-size: 4.8309vw;
  color: #333;
  flex: auto;
  text-align: center;
}
.my_data_opction ul .active {
  position: relative;
  color: #3775f6;
  font-weight: bold;
  font-size: 5.314vw;
  /* border-bottom: .4831vw solid #000; */
}
.my_data_opction ul .active::after {
  content: '';
  display: inline-block;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  bottom: -2px;
  width: 60%;
  height: .4831vw;
  background-color: #3775f6;
}
  </style>
