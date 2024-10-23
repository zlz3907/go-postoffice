<script setup>
import { useI18n } from 'vue-i18n'
import QuoteService from './components/QuoteService.vue'
import Documentation from './components/Documentation.vue'

const { t, locale } = useI18n()

const languages = [
  { code: 'zh', name: '中文', flag: '🇨🇳' },
  { code: 'ug', name: 'ئۇيغۇرچە', flag: '🇨🇳' },
  { code: 'bo', name: 'བོད་ཡིག', flag: '🇨🇳' },
  { code: 'de', name: 'DE', flag: '🇩🇪' },
  { code: 'en', name: 'EN', flag: '🇺🇸' },
  { code: 'es', name: 'ES', flag: '🇪🇸' },
  { code: 'fr', name: 'FR', flag: '🇫🇷' },
  { code: 'lo', name: 'ລາວ', flag: '🇱🇦' },
  { code: 'sn', name: 'SN', flag: '🇿🇼' }
]

function changeLanguage(lang) {
  locale.value = lang
}
</script>

<template>
  <div class="app-wrapper">
    <div class="container">
      <div class="language-selector">
        <button v-for="lang in languages" :key="lang.code" @click="changeLanguage(lang.code)" :class="{ active: locale === lang.code }">
          <span class="flag">{{ lang.flag }}</span> {{ lang.name }}
        </button>
      </div>
      <h1>{{ t('title') }}</h1>
      <div class="content">
        <div class="interactive-section">
          <QuoteService />
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

.app-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.container {
  flex: 1;
  width: 100%;
  max-width: 1400px;
  margin: 0 auto;
  padding: 20px;
  box-sizing: border-box;
}

h1 {
  color: #2c3e50;
  text-align: center;
  margin: 40px 0 20px;
  font-size: 2em;
  font-weight: 500;
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

.language-selector {
  display: flex;
  justify-content: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 30px;
}

.language-selector button {
  background-color: #ffffff;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  padding: 5px 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  font-size: 0.9em;
  color: #333;
  min-width: 80px;
  justify-content: center;
}

.language-selector button:hover {
  background-color: #f0f0f0;
}

.language-selector button.active {
  background-color: #3498db;
  color: white;
  border-color: #3498db;
}

.flag {
  margin-right: 5px;
  font-size: 1.2em;
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
}
</style>
