<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import QuoteService from './components/QuoteService.vue'
import Documentation from './components/Documentation.vue'
import SidebarPanel from './components/SidebarPanel.vue'
import commandExamples from './command-examples.json'
import messageExamples from './msg-examples.json'

const { t, locale } = useI18n()
const isSidebarOpen = ref(false)
const isConnected = ref(false)
const isHeartbeating = ref(false)
const activeTab = ref('logs')
const isSocketOpen = ref(false)

const languages = [
  { code: 'zh', name: '中文', flag: '🇨🇳' },
  { code: 'ug', name: 'ئۇيغۇرچە', flag: '🇨🇳' },
  { code: 'bo', name: 'བོད་ཡིག', flag: '🇨🇳' },
  { code: 'de', name: 'DE', flag: '🇩🇪' },
  { code: 'en', name: 'EN', flag: '🇺🇸' },
  { code: 'es', name: 'ES', flag: '🇪🇸' },
  { code: 'fr', name: 'FR', flag: '🇫🇷' },
  { code: 'lo', name: 'ລາວ', flag: '🇱🇦' },
  { code: 'sn', name: 'SN', flag: '🇿🇼' },
  { code: 'ja', name: '日本語', flag: '🇯🇵' }
]

function changeLanguage(lang) {
  locale.value = lang
}

function toggleSidebar() {
  isSidebarOpen.value = !isSidebarOpen.value
}

function updateConnectionStatus(status, heartbeat) {
  console.log('Connection status updated:', status, heartbeat);
  isConnected.value = status
  isHeartbeating.value = heartbeat
  isSocketOpen.value = status // 添加这行来更新 isSocketOpen
}

const socketStatusIcon = computed(() => {
  console.log('Socket status:', isConnected.value, isHeartbeating.value);
  if (isConnected.value) {
    return isHeartbeating.value ? '🟡' : '🟢'
  }
  return '🔴'
})

const socketStatusClass = computed(() => {
  if (isConnected.value) {
    return isHeartbeating.value ? 'connected heartbeat' : 'connected'
  }
  return 'disconnected'
})

const logs = ref([])

function addLog(log) {
  console.log('Log added:', log);
  logs.value.push(log)
}

function openSidebar(tab) {
  console.log('Opening sidebar with tab:', tab);
  isSidebarOpen.value = true;
  activeTab.value = tab;
}

const quoteServiceRef = ref(null)

function fillCommandForm(example) {
  console.log('App: Filling command form with:', example);
  if (quoteServiceRef.value) {
    quoteServiceRef.value.fillCommandForm(example);
  } else {
    console.error('QuoteService component reference not found');
  }
}

onMounted(() => {
  // 在这里添加任何需要在组件挂载时执行的代码
})

onUnmounted(() => {
  // 在这里添加任何需要在组件卸载时执行的清理代码
})
</script>

<template>
  <div class="app-container">
    <nav class="top-nav">
      <div class="nav-content">
        <div class="language-selector">
          <button 
            v-for="lang in languages" 
            :key="lang.code" 
            @click="changeLanguage(lang.code)" 
            :class="{ active: locale === lang.code }"
          >
            <span class="flag">{{ lang.flag }}</span> {{ lang.name }}
          </button>
        </div>
        <div class="status-and-sidebar">
          <!-- 更新 socket 状态图标 -->
          <span 
            class="socket-status-icon" 
            :class="socketStatusClass"
            :title="isConnected ? (isHeartbeating ? 'Connected (Heartbeat)' : 'Connected') : 'Disconnected'"
          >
            {{ socketStatusIcon }}
          </span>
          <!-- 侧边栏切换按钮 -->
          <button @click="toggleSidebar" class="sidebar-toggle">
            ☰
          </button>
        </div>
      </div>
    </nav>

    <div class="container">
      <div class="content">
        <div class="interactive-section">
          <QuoteService 
            ref="quoteServiceRef"
            @connection-status-change="updateConnectionStatus"
            @add-log="addLog"
            @open-sidebar="openSidebar"
          />
        </div>
        <div class="documentation-section">
          <Documentation />
        </div>
      </div>
    </div>
    <footer>
      <div class="footer-content">
        <p>&copy; {{ new Date().getFullYear() }} jyiai.com. All rights reserved.</p>
        <p>
          Contents licensed under 
          <a href="http://creativecommons.org/licenses/by-nc-nd/4.0/" target="_blank" rel="noopener noreferrer">
            CC BY-NC-ND 4.0
          </a> 
          with attribution required
        </p>
      </div>
    </footer>

    <!-- 添加侧边栏组件 -->
    <SidebarPanel 
      :isOpen="isSidebarOpen" 
      :logs="logs"
      :commandExamples="commandExamples"
      :activeTab="activeTab"
      @close="isSidebarOpen = false"
      @update:activeTab="activeTab = $event"
      @fill-command-form="fillCommandForm"
    />
  </div>
</template>

<style>
@import url('https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@400;500;700&display=swap');

body {
  font-family: 'Noto Sans SC', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  line-height: 1.6;
  margin: 0;
  padding: 0;
  background-color: #f5f5f5;
  color: #333;
}

.app-container {
  padding-top: 60px; /* 为顶部导航栏留出空间 */
  padding-bottom: 60px; /* 为状态栏留出空间 */
}

.top-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background-color: #333;
  color: white;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  z-index: 1000;
}

.nav-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.nav-content h1 {
  margin: 0;
  font-size: 1.5em;
  color: white;
}

.language-selector {
  display: flex;
  gap: 10px;
}

.language-selector button {
  background-color: transparent;
  border: 1px solid #555;
  border-radius: 4px;
  color: white;
  padding: 5px 10px;
  cursor: pointer;
  transition: background-color 0.3s;
  display: flex;
  align-items: center;
}

.language-selector button:hover,
.language-selector button.active {
  background-color: #555;
}

.flag {
  margin-right: 5px;
}

.container {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  box-sizing: border-box;
}

.content {
  display: flex;
  gap: 30px;
}

.interactive-section {
  flex: 1;
  max-width: 38.2%; /* 黄金分割比例 */
  min-width: 300px; /* 设置最小宽度 */
}

.documentation-section {
  flex: 1.618; /* 黄金分割比例 */
  min-width: 400px; /* 设置最小宽度 */
}

footer {
  margin-top: 40px;
  padding: 20px 0;
  border-top: 1px solid #e0e0e0;
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  text-align: center;
  font-size: 0.9em;
  color: #666;
}

footer p {
  margin: 5px 0;
}

footer a {
  color: #3498db;
  text-decoration: none;
}

footer a:hover {
  text-decoration: underline;
}

/* 添加媒体查询以处理小屏幕 */
@media (max-width: 1200px) {
  .content {
    flex-direction: column;
  }
  .interactive-section,
  .documentation-section {
    max-width: 100%;
    min-width: auto;
  }
  .nav-content {
    flex-direction: column;
    align-items: flex-start;
  }
  .language-selector {
    margin-top: 10px;
  }
}

/* 修改侧边栏切换按钮样式 */
.sidebar-toggle {
  background-color: transparent;
  border: none;
  color: white;
  font-size: 1.5em;
  cursor: pointer;
  padding: 0;
  width: 40px;
  height: 40px;
  display: flex;
  justify-content: center; /* 保留水平居中 */
  transition: background-color 0.3s;
  border-radius: 4px;
}

.sidebar-toggle:hover {
  background-color: #555;
}

.nav-content {
  justify-content: space-between;
  align-items: center;
}

.status-and-sidebar {
  display: flex;
  align-items: center;
}

.socket-status-icon {
  font-size: 1.2em;
  margin-right: 10px;
  transition: all 0.3s ease;
}

.socket-status-icon.disconnected {
  opacity: 0.5;
}

.socket-status-icon.connected {
  opacity: 1;
}

@keyframes heartbeat {
  0% { transform: scale(1); }
  50% { transform: scale(1.2); }
  100% { transform: scale(1); }
}

.socket-status-icon.heartbeat {
  animation: heartbeat 1s infinite;
}

/* 其他样式保持不变 */
</style>
