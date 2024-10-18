import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueSetupExtend from 'vite-plugin-vue-setup-extend'
import path from 'path'

export default defineConfig({
  plugins: [
    vue(),
    VueSetupExtend(),
  ],
  // resolve: {
  //   alias: {
  //     '@': fileURLToPath(new URL('./src', import.meta.url)),
  //     'vue$': path.resolve(__dirname, 'node_modules/vue/dist/vue.runtime.esm-bundler.js'), 
  //   }
  // },
  server: {
    proxy: {
      '/user': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
})
