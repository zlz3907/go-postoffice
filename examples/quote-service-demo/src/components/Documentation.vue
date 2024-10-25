<template>
  <div class="documentation">
    <h2>{{ t('documentation.title') }}</h2>
    <div v-for="(step, index) in steps" :key="index" class="doc-step">
      <h3>{{ t(`documentation.steps[${index}].title`) }}</h3>
      <p>{{ t(`documentation.steps[${index}].description`) }}</p>
      <h4>{{ t('documentation.requestExample') }}</h4>
      <div class="code-block">
        <div class="code-header">
          <div class="code-language">{{ detectLanguage(codeExamples[index].request) }}</div>
          <button @click="copyCode(codeExamples[index].request)" class="copy-button" :title="t('copyCode')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
          </button>
        </div>
        <pre><code v-html="highlightCode(codeExamples[index].request, detectLanguage(codeExamples[index].request))"></code></pre>
      </div>
      <h4>{{ t('documentation.responseExample') }}</h4>
      <div class="code-block">
        <div class="code-header">
          <div class="code-language">{{ detectLanguage(codeExamples[index].response) }}</div>
          <button @click="copyCode(codeExamples[index].response)" class="copy-button" :title="t('copyCode')">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
          </button>
        </div>
        <pre><code v-html="highlightCode(codeExamples[index].response, detectLanguage(codeExamples[index].response))"></code></pre>
      </div>
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n'
import hljs from 'highlight.js/lib/core';
import json from 'highlight.js/lib/languages/json';
import 'highlight.js/styles/github.css'; // 更改为 GitHub 主题

export default {
  name: 'Documentation',
  setup() {
    const { t } = useI18n()

    const highlightCode = (code, language) => {
      if (language === 'http' || language === 'plaintext') {
        // 对于 HTTP 和 plaintext，我们使用自定义的高亮逻辑
        return customHighlight(code, language);
      }
      return hljs.highlight(code, { language: language || 'plaintext' }).value
    }

    const customHighlight = (code, language) => {
      if (language === 'http') {
        // 为 HTTP 请求添加简单的语法高亮
        return code.replace(/(GET|POST|PUT|DELETE|PATCH)/, '<span class="http-method">$1</span>')
                   .replace(/(https?:\/\/[^\s]+)/, '<span class="http-url">$1</span>');
      } else {
        // 为 plaintext 添加基本的 HTML 转义
        return code.replace(/&/g, '&amp;')
                   .replace(/</g, '&lt;')
                   .replace(/>/g, '&gt;')
                   .replace(/"/g, '&quot;')
                   .replace(/'/g, '&#039;');
      }
    }

    const detectLanguage = (content) => {
      if (content.startsWith('{') || content.startsWith('[')) {
        return 'json'
      } else if (content.startsWith('GET') || content.startsWith('POST')) {
        return 'http'
      }
      return 'plaintext'
    }

    const copyCode = (code) => {
      navigator.clipboard.writeText(code).then(() => {
        console.log('Code copied to clipboard')
      })
    }

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

    return { t, highlightCode, detectLanguage, copyCode, languages }
  },
  data() {
    return {
      steps: [0, 1, 2],
      codeExamples: [
        {
          request: 'GET https://ai.zhycit.com/wecom/token?appid=<APPID>&secret=<SECRET>',
          response: `{
  "expires": 1732256675685,
  "token": "af71********6922",
  "ttl": 2589694
}`
        },
        {
          request: 'WebSocket连接: wss://socket.zhycit.com/?clientID=<CLIENT_ID>&token=<TOKEN>',
          response: '连接成功后，服务器不会立即返回消息。'
        },
        {
          request: `{
  "from": "/service/demo/quote_1",
  "to": "AI_CHATBOT",
  "subject": "报价请求",
  "content": "请为<产品名称>提供报价",
  "type": "msg",
  "createTime": 1623456789
}`,
          response: `{
  "from": "AI_CHATBOT",
  "to": "/service/demo/quote_1",
  "subject": "报价结果",
  "content": "<产品名称>的报价为：XXX元",
  "type": "msg",
  "createTime": 1623456790
}`
        }
      ]
    }
  }
}
</script>

<style scoped>
.documentation {
  background-color: #ffffff;
  color: #24292e;
  padding: 20px;
  border-radius: 8px;
  text-align: left;
}
.doc-step {
  margin-bottom: 30px;
}
h2, h3, h4 {
  color: #24292e;
}
p {
  margin-left: 0;
}
.code-block {
  background-color: #f6f8fa;
  border-radius: 6px;
  margin: 10px 0;
  overflow: hidden;
}
.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #f1f3f5;
  padding: 5px 10px;
  border-bottom: 1px solid #e1e4e8;
}
.code-language {
  color: #24292e;
  font-size: 0.8em;
}
pre {
  margin: 0;
  padding: 16px;
  overflow-x: auto;
}
code {
  font-family: SFMono-Regular, Consolas, 'Liberation Mono', Menlo, monospace;
  font-size: 14px;
  line-height: 1.5;
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

/* GitHub 主题的语法高亮样式会由 highlight.js 的 CSS 文件提供 */
</style>
