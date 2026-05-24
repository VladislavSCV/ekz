import { resolve } from 'path'
import { defineConfig } from 'vite'

export default defineConfig({
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        register: resolve(__dirname, 'register.html'),
        booking: resolve(__dirname, 'booking.html'),
        cabinet: resolve(__dirname, 'cabinet.html'),
        admin: resolve(__dirname, 'admin.html'),
      },
    },
  },
})
