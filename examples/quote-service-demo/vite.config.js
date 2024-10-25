import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/wecom': {
        target: 'https://ai.zhycit.com',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/wecom/, '/wecom')
      }
    }
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  }
})
