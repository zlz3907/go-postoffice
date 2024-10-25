<template>
  <div class="quote-service container">
    <h1>{{ t('title') }}</h1>
    
    <div class="step">
      <h2>{{ t('step0.title') }}</h2>
      <div class="form-group">
        <label for="tokenUrl">{{ t('step0.tokenUrl') }}</label>
        <input id="tokenUrl" v-model="tokenUrl" disabled class="form-control" />
      </div>
      <div class="form-group">
        <label for="appId">{{ t('step0.appId') }}</label>
        <input id="appId" v-model="appId" :placeholder="t('step0.appId')" class="form-control" @input="updateTokenUrl" />
      </div>
      <div class="form-group">
        <label for="secret">{{ t('step0.secret') }}</label>
        <input id="secret" v-model="secret" :placeholder="t('step0.secret')" type="password" class="form-control" @input="updateTokenUrl" />
      </div>
      <button @click="getToken" class="btn btn-primary">{{ t('step0.getToken') }}</button>
      <div class="mt-3">{{ tokenStatus }}</div>
      <div v-if="tokenInfo" class="mt-3 token-info">
        <div><strong>{{ t('tokenInfo.title') }}：</strong> {{ tokenInfo.token }}</div>
        <div><strong>{{ t('tokenInfo.expiresDate') }}：</strong> {{ tokenInfo.expiresDate }}</div>
        <div><strong>{{ t('tokenInfo.ttl') }}：</strong> {{ tokenInfo.ttl }}{{ t('seconds') }}</div>
      </div>
    </div>

    <div v-if="token" class="step">
      <h2>{{ t('step1.title') }}</h2>
      <div class="form-group">
        <label for="serverUrl">{{ t('step1.serverUrl') }}</label>
        <input id="serverUrl" v-model="serverUrl" :placeholder="t('step1.serverUrl')" class="form-control" />
      </div>
      <button @click="connect" class="btn btn-primary">{{ t('step1.connect') }}</button>
      <div class="mt-3">{{ connectionStatus }}</div>
    </div>

    <div v-if="isConnected" class="step">
      <div class="step-header">
        <h2>{{ t('step2.title') }}</h2>
        <button @click="openExamplesSidebar" class="btn btn-sm btn-light show-examples-btn">
          {{ t('step2.showExamples') }}
        </button>
      </div>
      <!-- 必填项 -->
      <div class="form-group">
        <label for="from">{{ t('step2.from') }}</label>
        <input id="from" v-model="message.from" :placeholder="t('step2.fromPlaceholder')" class="form-control" required />
      </div>
      <div class="form-group">
        <label for="to">{{ t('step2.to') }}</label>
        <input id="to" v-model="message.to" :placeholder="t('step2.toPlaceholder')" class="form-control" required />
      </div>
      <div class="form-group">
        <label for="subject">{{ t('step2.subject') }}</label>
        <input id="subject" v-model="message.subject" :placeholder="t('step2.subjectPlaceholder')" class="form-control" required />
      </div>
      <div class="form-group">
        <label for="content">{{ t('step2.content') }}</label>
        <textarea 
          id="content" 
          v-model="formattedContent" 
          :placeholder="t('step2.contentPlaceholder')" 
          class="form-control" 
          required
          :style="{ fontFamily: isContentObject ? 'monospace' : 'inherit' }"
        ></textarea>
      </div>
      <div class="form-group">
        <label for="type">{{ t('step2.type') }}</label>
        <select id="type" v-model="message.type" class="form-control" required>
          <option value="UPDATE_CONTACTS">更新联系人列表</option>
          <option value="UPDATE_BOT_CONFIG">更新机器人配置</option>
          <option value="SEND_REPLY_MSG">发送回复消息</option>
        </select>
      </div>
      <div class="form-group">
        <label for="attachments">{{ t('step2.attachments') }}</label>
        <input id="attachments" v-model="message.attachments" :placeholder="t('step2.attachmentsPlaceholder')" class="form-control" />
      </div>

      <!-- 可选项展开/折叠链接 -->
      <a @click="toggleOptionalFields" class="toggle-optional">
        {{ showOptionalFields ? t('step2.hideOptional') : t('step2.showOptional') }}
        <span class="arrow-icon">{{ showOptionalFields ? '▲' : '▼' }}</span>
      </a>

      <!-- 可选项 -->
      <div v-if="showOptionalFields" class="optional-fields">
        <div class="form-group">
          <label for="cc">{{ t('step2.cc') }}</label>
          <input id="cc" v-model="message.cc" :placeholder="t('step2.ccPlaceholder')" class="form-control" />
        </div>
        <div class="form-group">
          <label for="contentType">{{ t('step2.contentType') }}</label>
          <input id="contentType" v-model="message.contentType" type="number" :placeholder="t('step2.contentTypePlaceholder')" class="form-control" />
        </div>
        <div class="form-group">
          <label for="charset">{{ t('step2.charset') }}</label>
          <input id="charset" v-model="message.charset" :placeholder="t('step2.charsetPlaceholder')" class="form-control" />
        </div>
        <div class="form-group">
          <label for="level">{{ t('step2.level') }}</label>
          <input id="level" v-model="message.level" type="number" :placeholder="t('step2.levelPlaceholder')" class="form-control" />
        </div>
        <div class="form-group">
          <label for="tags">{{ t('step2.tags') }}</label>
          <input id="tags" v-model="message.tags" :placeholder="t('step2.tagsPlaceholder')" class="form-control" />
        </div>
      </div>

      <!-- 发送指令按钮 -->
      <div class="send-message-container">
        <button @click="sendMessage" class="btn btn-primary">{{ t('step2.sendCommand') }}</button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, reactive, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import commandExamples from '../command-examples.json'

export default {
  name: 'QuoteService',
  props: {
    // ... props 保持不变 ...
  },
  emits: ['update:connectionStatus', 'addLog', 'open-sidebar', 'connection-status-change', 'add-log'],
  setup(props, { emit }) {
    const { t, locale } = useI18n()
    const showOptionalFields = ref(false)
    const message = reactive({
      from: '',
      to: '',
      subject: '',
      content: '', // 将 content 初始化为空字符串
      type: 'msg',
      cc: '',
      contentType: null,
      charset: '',
      level: 0,
      tags: [],
      attachments: '',
    })

    const isContentObject = ref(false)

    const formattedContent = computed({
      get: () => {
        if (isContentObject.value) {
          return JSON.stringify(message.content, null, 2)
        }
        return message.content
      },
      set: (value) => {
        try {
          const parsed = JSON.parse(value)
          message.content = parsed
          isContentObject.value = true
        } catch (e) {
          message.content = value
          isContentObject.value = false
        }
      }
    })

    const token = ref('')
    const serverUrl = ref('wss://socket.zhycit.com/')
    const appId = ref('test_app_quote_01')
    const secret = ref('')
    const tokenUrl = ref('')
    const tokenStatus = ref('')
    const connectionStatus = ref(false)
    const isConnected = ref(false)
    const socket = ref(null)
    const CLIENT_ID = '/service/erp/company/001'
    const tokenInfo = ref(null)
    const isHeartbeating = ref(false)
    let heartbeatInterval = null

    const toggleOptionalFields = () => {
      showOptionalFields.value = !showOptionalFields.value
    }

    const openExamplesSidebar = () => {
      console.log('Opening examples sidebar'); // 添加日志
      emit('open-sidebar', 'msg-examples');
    }

    const updateTokenUrl = () => {
      tokenUrl.value = `/wecom/token?appid=${appId.value}&secret=${secret.value}`;
    }

    const getToken = async () => {
      tokenStatus.value = t('step0.gettingToken');
      try {
        const response = await fetch(tokenUrl.value);
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const data = await response.json();
        if (data.token) {
          token.value = data.token;
          tokenInfo.value = {
            token: data.token,
            expiresDate: new Date(data.expires).toLocaleString(),
            ttl: data.ttl
          };
          tokenStatus.value = t('step0.tokenSuccess');
        } else {
          throw new Error('Token not found in response');
        }
      } catch (error) {
        console.error('Error fetching token:', error);
        tokenStatus.value = t('step0.tokenError');
      }
    }

    const connect = () => {
      if (!token.value) {
        alert(t('step1.tokenRequired'));
        return;
      }

      const secureUrl = serverUrl.value.replace('ws://', 'wss://');
      socket.value = new WebSocket(`${secureUrl}?clientID=${CLIENT_ID}&token=${token.value}`);

      socket.value.onopen = (event) => {
        console.log("WebSocket 连接已建立");
        connectionStatus.value = t('step1.connected');
        isConnected.value = true;
        emit('connection-status-change', true, false);
        startHeartbeat();
        addLog(t('step1.connected'), 'success');
      };

      socket.value.onmessage = (event) => {
        console.log("收到消息:", event.data);
        addLog(event.data, 'message', 'received');
      };

      socket.value.onerror = (error) => {
        console.error("WebSocket 错误:", error);
        connectionStatus.value = t('step1.connectionError');
        isConnected.value = false;
      };

      socket.value.onclose = (event) => {
        console.log("WebSocket 连接已关闭");
        connectionStatus.value = t('step1.connectionClosed');
        isConnected.value = false;
        emit('connection-status-change', false, false);
        stopHeartbeat();
      };
    };

    const startHeartbeat = () => {
      console.log("开始心跳");
      heartbeatInterval = setInterval(() => {
        sendHeartbeat();
      }, 25000);
      sendHeartbeat();
    }

    const stopHeartbeat = () => {
      console.log("停止心跳");
      if (heartbeatInterval) {
        clearInterval(heartbeatInterval);
        heartbeatInterval = null;
      }
      isHeartbeating.value = false;
      emit('connection-status-change', isConnected.value, false);
    }

    const sendHeartbeat = () => {
      console.log("发送心跳");
      if (socket.value && socket.value.readyState === WebSocket.OPEN) {
        const heartbeatMessage = {
          from: CLIENT_ID,
          to: "AI_CHATBOT",
          subject: "Heartbeat",
          content: "Keep-Alive",
          type: "heartbeat",
          createTime: Date.now()
        };
        socket.value.send(JSON.stringify(heartbeatMessage));
        isHeartbeating.value = true;
        emit('connection-status-change', true, true);
        
        setTimeout(() => {
          isHeartbeating.value = false;
          emit('connection-status-change', true, false);
        }, 1000);
      }
    }

    const updateMessageContent = () => {
      try {
        message.content = JSON.parse(formattedContent.value)
      } catch (error) {
        console.error('Invalid JSON:', error)
        // 可以在这里添加错误处理，比如显示错误消息给用户
      }
    }

    const sendMessage = () => {
      if (socket.value && socket.value.readyState === WebSocket.OPEN) {
        // 移除所有空值或未定义的字段
        const cleanMessage = Object.fromEntries(
          Object.entries(message).filter(([_, v]) => v != null && v !== '')
        )
        
        // ��需要额外处理 content 字段，因为已经是正确的格式了
        
        socket.value.send(JSON.stringify(cleanMessage))
        console.log("指令已发送:", cleanMessage)
        addLog(JSON.stringify(cleanMessage), 'message', 'sent')
        // 重置表单
        Object.keys(message).forEach(key => {
          if (key === 'content') {
            message[key] = isContentObject.value ? {} : ''
          } else {
            message[key] = ''
          }
        })
        isContentObject.value = false
        message.from = CLIENT_ID
        message.type = 'msg'
      } else {
        console.error("WebSocket is not connected")
        alert(t('step2.connectionError'))
        addLog(t('step2.connectionError'), 'error')
      }
    }

    const addLog = (message, type, direction = 'system') => {
      const timestamp = new Date().toISOString() // 使用 ISO 格式的时间戳
      emit('add-log', { timestamp, message, type, direction })
    };

    const fillCommandForm = (example) => {
      Object.keys(message).forEach(key => {
        if (example.hasOwnProperty(key)) {
          if (key === 'content') {
            if (typeof example[key] === 'object') {
              message[key] = example[key] // 直接赋值对象，不转换为字符串
              isContentObject.value = true
            } else {
              message[key] = example[key]
              isContentObject.value = false
            }
          } else {
            message[key] = example[key]
          }
        }
      })
    }

    onMounted(() => {
      updateTokenUrl();
      locale.value = 'zh'
    })

    onUnmounted(() => {
      stopHeartbeat();
      if (socket.value) {
        socket.value.close();
      }
    })

    return {
      t,
      locale,
      showOptionalFields,
      message,
      toggleOptionalFields,
      openExamplesSidebar,
      commandExamples,
      fillCommandForm,
      connect,
      addLog,
      token,
      serverUrl,
      appId,
      secret,
      tokenUrl,
      tokenStatus,
      connectionStatus,
      isConnected,
      socket,
      CLIENT_ID,
      tokenInfo,
      isHeartbeating,
      updateTokenUrl,
      getToken,
      sendMessage,
      formattedContent,
      updateMessageContent,
      stopHeartbeat, // 添加这行，将 stopHeartbeat 函数暴露给组件实例
      isContentObject, // 添加这行
    }
  },
  data() {
    return {
      product: '',
      quoteResult: '',
    }
  },
  created() {
    this.updateTokenUrl();
  },
}
</script>

<style scoped>
.quote-service {
  max-width: 600px;
  margin: 0 auto;
  padding: 10px;
}

h1 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 1rem;
  text-align: center;
}

h2 {
  font-size: 1.2rem;
  color: #444;
  margin-top: 0;
  margin-bottom: 1rem;
}

.step {
  margin-bottom: 1.5rem;
  background-color: #fff;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 1rem;
  position: relative; /* 添加这行 */
}

.step-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.step-header h2 {
  margin: 0;
}

.show-examples-btn {
  font-size: 0.75rem;
  padding: 0.2rem 0.4rem;
  line-height: 1;
}

.form-group {
  position: relative;
  margin-bottom: 1rem;
}

.form-group label {
  position: absolute;
  top: -0.5rem;
  left: 0.5rem;
  padding: 0 0.25rem;
  background-color: #fff;
  color: #666;
  font-size: 0.75rem;
  font-weight: 500;
  z-index: 1;
}

.form-control {
  display: block;
  width: 100%;
  padding: 0.5rem 0.75rem;
  font-size: 0.9rem;
  line-height: 1.5;
  color: #495057;
  background-color: #fff;
  background-clip: padding-box;
  border: 1px solid #ced4da;
  border-radius: 0.25rem;
  transition: border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

/* 为 tokenUrl 和 secret 文本框添加淡背景色 */
#tokenUrl,
#secret {
  background-color: #f8f9fa; /* 设置为淡灰色 */
}

.form-control:focus {
  border-color: #80bdff;
  outline: 0;
  box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
}

select.form-control {
  appearance: none;
  -webkit-appearance: none;
  -moz-appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='8' height='8' viewBox='0 0 8 8'%3E%3Cpath fill='%23343a40' d='M2.5 0L1 1.5 3.5 4 1 6.5 2.5 8l4-4-4-4z' transform='rotate(90 4 4)'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 8px 10px;
  padding-right: 1.5rem;
}

.btn {
  display: inline-block;
  font-weight: 400;
  text-align: center;
  vertical-align: middle;
  user-select: none;
  border: 1px solid transparent;
  padding: 0.375rem 0.75rem;
  font-size: 0.9rem;
  line-height: 1.5;
  border-radius: 0.25rem;
  transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
}

.btn-primary {
  color: #fff;
  background-color: #007bff;
  border-color: #007bff;
}

.btn-primary:hover {
  background-color: #0056b3;
  border-color: #0056b3;
}

.btn-secondary {
  color: #fff;
  background-color: #6c757d;
  border-color: #6c757d;
}

.btn-secondary:hover {
  background-color: #5a6268;
  border-color: #545b62;
}

.toggle-optional {
  display: inline-block;
  color: #6c757d;
  cursor: pointer;
  margin-bottom: 1rem;
  text-decoration: none;
  font-size: 0.9rem;
}

.toggle-optional:hover {
  text-decoration: underline;
}

.arrow-icon {
  margin-left: 5px;
  font-size: 0.8rem;
}

.optional-fields {
  background-color: #f8f9fa;
  padding: 0.75rem;
  border-radius: 0.25rem;
  margin-top: 0.5rem;
  margin-bottom: 1rem;
}

.token-info {
  background-color: #f8f9fa;
  padding: 0.75rem;
  border-radius: 0.25rem;
  margin-top: 0.75rem;
  font-size: 0.9rem;
  text-align: left; /* 添加这行来确保文本左对 */
}

.token-info div {
  margin-bottom: 0.25rem;
  display: flex; /* 使用 flex 布局 */
  align-items: flex-start; /* 顶部对齐 */
}

.token-info strong {
  min-width: 80px; /* 为标签设置最小宽度 */
  margin-right: 10px; /* 在标签和值之间添加一些间距 */
}

.send-message-container {
  margin-top: 1rem;
  text-align: center;
}

.send-message-container .btn {
  width: 100%;
  max-width: 200px;
}

.mr-3 {
  margin-right: 1rem;
}

textarea {
  white-space: pre-wrap;
  min-height: 100px;
}
</style>







