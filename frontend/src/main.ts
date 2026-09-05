import { createApp } from 'vue'
import App from './App.vue'
import 'element-plus/dist/index.css'

// dev 环境下引入 mock 模块
if (import.meta.env.DEV) {
  await import('@/mock')
}

const app = createApp(App)
app.mount('#app')
