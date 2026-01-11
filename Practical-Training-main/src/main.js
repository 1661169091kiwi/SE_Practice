import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import './utils/vant-ui.js'
import './styles/common.css'

import GameDisplay from './component/GameDisplay.vue'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')
app.component('GameDisplay', GameDisplay)
