import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite' // Bu satır şart

export default defineConfig({
  plugins: [react(), tailwindcss()], // tailwindcss() buraya eklenmeli
  server: {
    host: true,
    port: 5173,
    allowedHosts: ["furkan.local"],
    watch: {
      usePolling: true
    }
  }
})