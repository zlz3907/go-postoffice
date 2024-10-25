<template>
  <div class="sidebar-panel" :class="{ 'open': isOpen }">
    <div class="sidebar-header">
      <h2>{{ t('sidebar.title') }}</h2>
      <button @click="$emit('close')" class="close-button">×</button>
    </div>
    <div class="sidebar-content">
      <div class="tabs">
        <button 
          v-for="tab in tabs" 
          :key="tab.id" 
          @click="changeTab(tab.id)"
          :class="{ 'active': activeTab === tab.id }"
        >
          {{ t(`sidebar.tabs.${tab.id}`) }}
        </button>
      </div>
      <div class="tab-content">
        <div v-if="activeTab === 'logs'" class="log-entries" ref="logContainer">
          <div v-for="(log, index) in logs" :key="index" :class="['log-entry', log.type, log.direction]">
            <div class="log-header">
              <span class="log-timestamp">{{ formatTimestamp(log.timestamp) }}</span>
              <span class="log-type">
                {{ log.type.toUpperCase() }}
                <span class="direction-icon" :title="log.direction === 'sent' ? t('messageSent') : t('messageReceived')">
                  {{ log.direction === 'sent' ? '↗️' : '↘️' }}
                </span>
              </span>
            </div>
            <div class="log-content" :class="{ 'expanded': expandedLogs.includes(index) }">
              <pre v-if="log.type === 'message'"><code>{{ formatJSON(log.message) }}</code></pre>
              <div v-else class="log-message">{{ log.message }}</div>
            </div>
            <button v-if="shouldShowExpandButton(log)" @click="toggleLogExpansion(index)" class="expand-button">
              {{ expandedLogs.includes(index) ? t('collapse') : t('expand') }}
            </button>
          </div>
        </div>
        <div v-if="activeTab === 'msg-examples'" class="message-examples">
          <h3>{{ t('sidebar.tabs.msg-examples') }}</h3>
          <p>{{ t('msgExamplesDescription') }}</p>
          <div v-for="(example, index) in commandExamples" :key="index" class="example-item">
            <div class="example-header">
              <h4>{{ example.desc }}</h4>
              <button @click="fillCommandForm(example.example)" class="fill-form-btn">
                {{ t('fillCommandForm') }}
              </button>
            </div>
            <div class="code-block">
              <div class="code-header">
                <div class="code-language">JSON</div>
                <button @click="copyCode(JSON.stringify(example.example, null, 2))" class="copy-button" :title="t('copyCode')">
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                  </svg>
                </button>
              </div>
              <pre><code v-html="highlightCode(JSON.stringify(example.example, null, 2))"></code></pre>
            </div>
            <hr v-if="index !== commandExamples.length - 1" class="example-divider">
          </div>
        </div>
        <div v-if="activeTab === 'email-format'" class="message-examples">
          <h3>{{ t('sidebar.tabs.email-format') }}</h3>
          <p>{{ t('emailFormatDescription') }}</p>
          <div v-for="(example, index) in messageExamples" :key="index" class="example-item">
            <div class="example-header">
              <h4>{{ example.desc }}</h4>
            </div>
            <div class="code-block">
              <div class="code-header">
                <div class="code-language">JSON</div>
                <button @click="copyCode(JSON.stringify(example.example, null, 2))" class="copy-button" :title="t('copyCode')">
                  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                  </svg>
                </button>
              </div>
              <pre><code v-html="highlightCode(JSON.stringify(example.example, null, 2))"></code></pre>
              <!-- 添加分割线，但不在最后一个示例后添加 -->
              <hr v-if="index !== messageExamples.length - 1" class="example-divider">
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import 'highlight.js/styles/github.css'; // GitHub 主题

hljs.registerLanguage('json', json);

import commandExamples from '@/command-examples.json'
import messageExamples from '@/msg-examples.json'

export default {
  name: 'SidebarPanel',
  props: {
    isOpen: Boolean,
    logs: Array,
    messageExamples: Array,
    emailFormatExamples: Array,
    activeTab: String
  },
  emits: ['close', 'update:activeTab', 'fill-command-form', 'open-examples'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const localActiveTab = ref(props.activeTab || 'logs')
    const logContainer = ref(null)
    const expandedLogs = ref([])

    const tabs = [
      { id: 'logs', name: t('sidebar.tabs.logs') },
      { id: 'msg-examples', name: t('sidebar.tabs.msg-examples') },
      { id: 'email-format', name: t('sidebar.tabs.email-format') }
    ]

    watch(() => props.activeTab, (newTab) => {
      localActiveTab.value = newTab
    })

    const copyCode = (code) => {
      navigator.clipboard.writeText(code).then(() => {
        console.log('Code copied to clipboard')
      })
    }

    const highlightCode = (code) => {
      return hljs.highlight(code, { language: 'json' }).value
    }

    const fillCommandForm = (example) => {
      console.log('SidebarPanel: Filling command form with:', example);
      emit('fill-command-form', example)
    }

    const formatTimestamp = (timestamp) => {
      const date = new Date(timestamp);
      return date.toLocaleString(); // 使用本地化的日期和时间格式
    }

    const formatJSON = (message) => {
      try {
        const parsed = JSON.parse(message);
        return JSON.stringify(parsed, null, 2);
      } catch (e) {
        return message;
      }
    }

    const toggleLogExpansion = (index) => {
      const i = expandedLogs.value.indexOf(index)
      if (i > -1) {
        expandedLogs.value.splice(i, 1)
      } else {
        expandedLogs.value.push(index)
      }
    }

    const shouldShowExpandButton = (log) => {
      const formattedMessage = log.type === 'message' 
        ? formatJSON(log.message)
        : log.message;
      return formattedMessage.split('\n').length > 5 || formattedMessage.length > 200;
    }

    const changeTab = (tabId) => {
      emit('update:activeTab', tabId);
      if (tabId === 'msg-examples') {
        emit('open-examples');
      }
    };

    // 监听 logs 的变化
    watch(() => props.logs, () => {
      nextTick(() => {
        if (logContainer.value) {
          logContainer.value.scrollTop = logContainer.value.scrollHeight
        }
      })
    }, { deep: true })

    return {
      t,
      localActiveTab,
      tabs,
      copyCode,
      highlightCode,
      fillCommandForm,
      commandExamples,
      messageExamples,
      formatTimestamp,
      formatJSON,
      logContainer,
      expandedLogs,
      toggleLogExpansion,
      shouldShowExpandButton,
      changeTab,
    }
  }
}
</script>

<style scoped>
.sidebar-panel {
  position: fixed;
  top: 0;
  right: -500px; /* 增加宽度到 500px */
  width: 500px; /* 增加宽度到 500px */
  height: 100vh;
  background-color: #f6f8fa;
  box-shadow: -2px 0 5px rgba(0, 0, 0, 0.1);
  transition: right 0.3s ease;
  z-index: 1001;
  display: flex;
  flex-direction: column;
  text-align: left;
}

.sidebar-panel.open {
  right: 0;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  background-color: #f1f3f5;
  border-bottom: 1px solid #e1e4e8;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 1.2em;
  color: #24292e;
}

.close-button {
  background: none;
  border: none;
  font-size: 1.5em;
  cursor: pointer;
  color: #586069;
}

.sidebar-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  height: calc(100% - 50px); /* 假设头部高度为50px，根据实际情况调整 */
}

.tabs {
  display: flex;
  background-color: #f1f3f5;
  border-bottom: 1px solid #e1e4e8;
}

.tabs button {
  flex: 1;
  padding: 10px;
  border: none;
  background: none;
  cursor: pointer;
  font-size: 0.9em;
  color: #24292e;
  transition: background-color 0.2s, color 0.2s;
}

.tabs button.active {
  background-color: #fff;
  color: #0366d6;
  font-weight: bold;
  box-shadow: inset 0 -2px 0 #0366d6;
}

.tab-content {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
}

.log-entries {
  height: 100%;
  overflow-y: auto;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
  font-size: 0.9em;
  scrollbar-width: none; /* 对于 Firefox */
  -ms-overflow-style: none; /* 对于 Internet Explorer 和 Edge */
}

.log-entries::-webkit-scrollbar {
  display: none; /* 对于 Chrome, Safari 和 Opera */
}

.log-entry {
  margin-bottom: 15px;
  padding: 10px;
  border-radius: 4px;
  border: 1px solid #e1e4e8;
}

.log-entry.sent {
  background-color: #f1f8ff; /* 发送消息的淡蓝色背景 */
}

.log-entry.received {
  background-color: #f6f8fa; /* 接收消息的浅灰色背景 */
}

.log-entry.error {
  background-color: #ffe6e6; /* 错误消息的淡红色背景 */
}

.log-entry.info {
  background-color: #f6f8fa; /* 信息消息保持原来的背景色 */
}

.log-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 5px;
  font-size: 0.8em;
}

.log-timestamp {
  color: #586069;
}

.log-type {
  font-weight: bold;
  display: flex;
  align-items: center;
}

.direction-icon {
  margin-left: 5px;
  font-size: 1.2em;
}

.log-content {
  max-height: 100px;
  overflow: hidden;
  transition: max-height 0.3s ease-out;
}

.log-content.expanded {
  max-height: none;
}

.expand-button {
  background: none;
  border: none;
  color: #0366d6;
  cursor: pointer;
  padding: 5px 0;
  font-size: 0.9em;
  text-align: left;
  display: block;
  width: 100%;
}

.expand-button:hover {
  text-decoration: underline;
}

.log-content pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.log-entry.sent .log-type {
  color: #28a745;
}

.log-entry.received .log-type {
  color: #0366d6;
}

.log-entry.error .log-type {
  color: #d73a49;
}

.log-entry.info .log-type {
  color: #6f42c1;
}

.message-examples {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
}

.example-item {
  margin-bottom: 20px;
}

.example-item h3 {
  font-size: 1em;
  margin-bottom: 10px;
  color: #24292e;
}

.code-block {
  background-color: #f6f8fa;
  border-radius: 3px;
  margin: 10px 0;
  overflow: hidden;
  max-width: 100%; /* 确保代码块不会超出侧边栏 */
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #f1f3f5;
  padding: 5px 10px;
  border: 1px solid #e1e4e8;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
}

.code-language {
  color: #24292e;
  font-size: 0.8em;
}

.copy-button {
  background: transparent;
  border: none;
  color: #586069;
  cursor: pointer;
  padding: 2px;
  transition: color 0.3s;
  font-size: 0.8em;
  display: flex;
  align-items: center;
  justify-content: center;
}

.copy-button:hover {
  color: #0366d6;
}

.copy-button svg {
  width: 14px;
  height: 14px;
}

pre {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
  font-family: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 12px;
  line-height: 1.45;
  color: #24292e;
  white-space: pre-wrap; /* 允许长行换行 */
  word-break: break-word; /* 在单词内换行 */
}

code {
  font-family: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 12px;
  line-height: 1.45;
  word-wrap: normal;
}

.example-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.fill-form-btn {
  font-size: 0.8em;
  padding: 4px 8px;
  background-color: #f0f0f0;
  border: 1px solid #d0d0d0;
  border-radius: 3px;
  cursor: pointer;
  color: #333; /* 添加文本颜色 */
  transition: background-color 0.2s, color 0.2s; /* 添加过渡效果 */
}

.fill-form-btn:hover {
  background-color: #e0e0e0;
  color: #000; /* 悬停时文本颜色变深 */
}

.sidebar-item i {
  margin-right: 8px;
}

.message-examples h3 {
  font-size: 1.2em;
  margin-bottom: 10px;
  color: #24292e;
}

.message-examples p {
  margin-bottom: 20px;
  color: #586069;
}

.message-examples h4 {
  font-size: 1em;
  margin-bottom: 10px;
  color: #24292e;
}

.example-divider {
  margin: 20px 0;
  border: 0;
  border-top: 1px solid #e1e4e8;
}

/* 为 WebKit 浏览器（如 Chrome、Safari）设置滚动条样式 */
.log-entries::-webkit-scrollbar {
  width: 8px;
}

.log-entries::-webkit-scrollbar-track {
  background: #f1f1f1;
}

.log-entries::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.log-entries::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>

