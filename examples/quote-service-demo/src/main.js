import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import i18n from './i18n'
import 'highlight.js/styles/atom-one-dark.css' // 更改为 atom-one-dark 主题
import 'bootstrap/dist/css/bootstrap.min.css'
// 添加全局样式
const style = document.createElement('style')
style.textContent = `
  @import url('https://fonts.googleapis.com/css2?family=Fira+Code&display=swap');
  
  code {
    font-family: 'Fira Code', monospace;
  }
`
document.head.appendChild(style)

const app = createApp(App)
app.use(i18n)
app.mount('#app')
