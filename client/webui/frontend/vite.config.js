import {defineConfig} from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
    base: '/extras/gostc/',
    plugins: [
        vue(),
    ],
    server: {
        proxy: {
            '/api': {
                target: 'http://127.0.0.1:18080',
                changeOrigin: true
            }
        },
        host: '0.0.0.0'
    }
})
