<template>
  <div>
    <div class="step">
      <h2>{{ t('step0.title') }}</h2>
      <div class="input-group">
        <label for="serverUrl">{{ t('step0.serverUrl') }}</label>
        <input id="serverUrl" v-model="serverUrl" :placeholder="t('step0.serverUrl')" />
      </div>
      <div class="input-group">
        <label for="appId">{{ t('step0.appId') }}</label>
        <input id="appId" v-model="appId" :placeholder="t('step0.appId')" />
      </div>
      <div class="input-group">
        <label for="secret">{{ t('step0.secret') }}</label>
        <input id="secret" v-model="secret" :placeholder="t('step0.secret')" type="password" />
      </div>
      <button @click="getToken">{{ t('step0.getToken') }}</button>
      <div>{{ tokenStatus }}</div>
      <div v-if="tokenInfo" class="token-display">
        <div class="token-info">
          <span class="token-label">{{ t('tokenInfo.title') }}：</span>
          <span class="token-value">{{ tokenInfo.token }}</span>
        </div>
        <div class="token-info">
          <span class="token-label">{{ t('tokenInfo.expiresDate') }}：</span>
          <span class="token-value">{{ tokenInfo.expiresDate }}</span>
        </div>
        <div class="token-info">
          <span class="token-label">{{ t('tokenInfo.ttl') }}：</span>
          <span class="token-value">{{ tokenInfo.ttl }}{{ t('seconds') }}</span>
        </div>
        <button @click="copyToken">{{ t('copyToken') }}</button>
      </div>
    </div>

    <div v-if="token" class="step">
      <h2>{{ t('step1.title') }}</h2>
      <button @click="connect">{{ t('step1.connect') }}</button>
      <div>{{ connectionStatus }}</div>
    </div>

    <div v-if="isConnected" class="step">
      <h2>{{ t('step2.title') }}</h2>
      <input v-model="product" :placeholder="t('step2.productInput')" />
      <button @click="requestQuote">{{ t('step2.requestQuote') }}</button>
    </div>

    <div v-if="isConnected" class="step">
      <h2>{{ t('messages') }}</h2>
      <div class="message-window" :class="{ 'fullscreen': isFullscreen }">
        <div class="message-content">
          <div v-for="(message, index) in messages" :key="index" class="message">
            {{ message }}
          </div>
        </div>
      </div>
      <button @click="toggleFullscreen" class="fullscreen-toggle">
        {{ isFullscreen ? t('exitFullscreen') : t('fullscreen') }}
      </button>
    </div>
  </div>
</template>

<script>
import { useI18n } from 'vue-i18n'

export default {
  name: 'QuoteService',
  setup() {
    const { t } = useI18n()
    return { t }
  },
  data() {
    return {
      serverUrl: 'wss://socket.zhycit.com/',
      appId: 'test_app_quote_01', // 更新默认值
      secret: '',
      token: '',
      tokenStatus: '',
      connectionStatus: '',
      isConnected: false,
      product: '',
      quoteResult: '',
      socket: null,
      CLIENT_ID: '/service/demo/quote_1',
      tokenInfo: null,
      messages: [],
      isFullscreen: false,
    }
  },
  methods: {
    async getToken() {
      console.log("getToken 函数被调用");
      if (!this.serverUrl || !this.appId || !this.secret) {
        alert("请填写所有必要的信息");
        return;
      }

      try {
        const response = await fetch(`/api/wecom/token?appid=${this.appId}&secret=${this.secret}`);
        const data = await response.json();
        this.handleTokenResponse(data);
      } catch (error) {
        console.error('获取 Token 失败:', error);
        this.tokenStatus = "获取 Token 失败，请检查输入信息";
      }
    },
    handleTokenResponse(response) {
      console.log("handleTokenResponse 被调用", response);
      if (response && response.token) {
        this.token = response.token;
        this.tokenStatus = "Token 获取成功";
        
        // 解析并格式化 token 信息
        const expiresDate = new Date(response.expires);
        this.tokenInfo = {
          token: response.token,
          expiresDate: this.formatDate(expiresDate),
          ttl: response.ttl
        };
      } else {
        console.error('获取 Token 失败:', response);
        this.tokenStatus = "获取 Token 失败，请检查输入信息";
        this.token = ''; // 清除之前的 token
        this.tokenInfo = null;
      }
    },
    formatDate(date) {
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, '0');
      const day = String(date.getDate()).padStart(2, '0');
      const hours = String(date.getHours()).padStart(2, '0');
      const minutes = String(date.getMinutes()).padStart(2, '0');
      const seconds = String(date.getSeconds()).padStart(2, '0');
      
      return `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`;
    },
    copyToken() {
      if (this.tokenInfo) {
        navigator.clipboard.writeText(this.tokenInfo.token).then(() => {
          alert('Token 已复制到剪贴板');
        }, (err) => {
          console.error('无法复制 Token: ', err);
          alert('复制 Token 失败，请手动复制');
        });
      }
    },
    connect() {
      console.log("connect 函数被调用");
      if (!this.token) {
        alert("请先获取 Token");
        return;
      }

      const secureUrl = this.serverUrl.replace('ws://', 'wss://');
      
      this.socket = new WebSocket(`${secureUrl}?clientID=${this.CLIENT_ID}&token=${this.token}`);

      this.socket.onopen = (event) => {
        console.log("WebSocket 连接已建立");
        this.connectionStatus = "已连接到服务器";
        this.isConnected = true;
      };

      this.socket.onmessage = (event) => {
        console.log("收到消息:", event.data);
        try {
          const response = JSON.parse(event.data);
          this.handleResponse(response);
        } catch (error) {
          console.error("解析消息失败:", error);
        }
      };

      this.socket.onerror = (error) => {
        console.error("WebSocket 错误:", error);
        this.connectionStatus = "连接错误，请重试";
        this.isConnected = false;
      };

      this.socket.onclose = (event) => {
        console.log("WebSocket 连接已关闭");
        this.connectionStatus = "连接已关闭";
        this.isConnected = false;
      };
    },
    requestQuote() {
      console.log("requestQuote 函数被调用");
      if (!this.product) {
        alert("请输入产品名称");
        return;
      }

      const message = {
        from: this.CLIENT_ID,
        to: "AI_CHATBOT",
        subject: "报价请求",
        content: `请为${this.product}提供报价`,
        type: "msg",
        createTime: Date.now()
      };

      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify(message));
      } else {
        alert("WebSocket 连接未建立，请先连接到服务器");
      }
    },
    handleResponse(response) {
      console.log("handleResponse 函数被调用", response);
      this.messages.push(JSON.stringify(response, null, 2));
    },
    toggleFullscreen() {
      this.isFullscreen = !this.isFullscreen;
    }
  }
}
</script>

<style scoped>
.step {
  margin-bottom: 25px;
  padding: 20px;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

.step h2 {
  color: #3498db;
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 1.2em;
}

.input-group {
  position: relative;
  margin-bottom: 15px;
}

.input-group label {
  position: absolute;
  top: -8px;
  left: 10px;
  font-size: 0.7em;
  color: #7f8c8d;
  background-color: #ffffff;
  padding: 0 5px;
}

input, button {
  margin: 8px 0;
  padding: 10px;
  width: 100%;
  box-sizing: border-box;
  border: 1px solid #bdc3c7;
  border-radius: 4px;
}

input {
  background-color: #ffffff;
  color: #2c3e50;
  padding-top: 15px; /* 为标签腾出空间 */
}

input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

button {
  background-color: #3498db;
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #2980b9;
}

.token-display {
  margin: 15px 0;
  padding: 15px;
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.token-info {
  display: flex;
  justify-content: flex-start; /* 改为左对齐 */
  margin-bottom: 8px;
  font-size: 0.9em;
  color: #6c757d;
}

.token-label {
  flex: 0 0 90px;
  text-align: right;
  padding-right: 10px;
  font-weight: 600;
}

.token-value {
  flex: 1;
  text-align: left;
  word-break: break-all;
  font-family: 'Courier New', Courier, monospace;
  padding-left: 10px; /* 添加左侧内边距 */
}

.token-display button {
  background-color: #6c757d;
  color: white;
  border: none;
  padding: 8px 15px;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  margin-top: 10px;
  font-size: 0.9em;
  width: auto;
}

.token-display button:hover {
  background-color: #5a6268;
}

.message-window {
  height: 300px;
  overflow-y: auto;
  border: 1px solid #ddd;
  padding: 10px;
  margin-top: 10px;
  background-color: #f8f9fa;
  transition: all 0.3s ease;
}

.message-window.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
}

.message {
  margin-bottom: 10px;
  padding: 5px;
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.1);
}

.fullscreen-toggle {
  margin-top: 10px;
  background-color: #34495e;
}

.fullscreen-toggle:hover {
  background-color: #2c3e50;
}
</style>
