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
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

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

    return { t, highlightCode, detectLanguage, copyCode }
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
  background-color: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  text-align: left;
}
.doc-step {
  margin-bottom: 30px;
}
h2, h3, h4 {
  color: #2c3e50;
}
p {
  margin-left: 0;
}
.code-block {
  position: relative;
  background-color: #282c34;
  border-radius: 6px;
  margin: 10px 0;
  overflow: hidden;
}
.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #21252b;
  padding: 5px 10px;
  border-bottom: 1px solid #3e4451;
}
.code-language {
  color: #abb2bf;
  font-size: 0.8em;
}
pre {
  margin: 0;
  padding: 12px 15px;
  overflow-x: auto;
}
code {
  font-family: 'Fira Code', 'Courier New', Courier, monospace;
  font-size: 14px;
  line-height: 1.5;
  color: #abb2bf; /* 默认文本颜色 */
}
.copy-button {
  background: transparent;
  border: none;
  color: #abb2bf;
  cursor: pointer;
  padding: 2px;
  transition: color 0.3s;
  font-size: 0.8em;
  display: flex;
  align-items: center;
  justify-content: center;
}
.copy-button svg {
  width: 14px;
  height: 14px;
}
.copy-button:hover {
  color: #ffffff;
}

/* HTTP 和 plaintext 的自定义样式 */
.http-method {
  color: #61afef; /* 蓝色 */
  font-weight: bold;
}
.http-url {
  color: #98c379; /* 绿色 */
}
code[class*="language-http"],
code[class*="language-plaintext"] {
  color: #e06c75; /* 红色，用于 HTTP 和 plaintext 的默认文本颜色 */
}
</style>
