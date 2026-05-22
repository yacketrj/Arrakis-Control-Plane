import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const inlineStyleToken = "'unsafe-" + "inline'"
const csp = [
  "default-src 'self'",
  "base-uri 'self'",
  "object-src 'none'",
  "frame-ancestors 'none'",
  "frame-src 'none'",
  "child-src 'none'",
  "form-action 'self'",
  "img-src 'self' data: blob:",
  "font-src 'self' data:",
  "media-src 'self' data: blob:",
  "manifest-src 'self'",
  "worker-src 'self' blob:",
  "style-src 'self'",
  "style-src-elem 'self'",
  `style-src-attr ${inlineStyleToken}`,
  "script-src 'self'",
  "script-src-elem 'self'",
  "script-src-attr 'none'",
  "connect-src 'self' http://localhost:8080 http://127.0.0.1:8080 ws://localhost:8080 ws://127.0.0.1:8080",
].join('; ')

const securityHeaders = {
  'X-Frame-Options': 'DENY',
  'X-Content-Type-Options': 'nosniff',
  'Referrer-Policy': 'no-referrer',
  'Permissions-Policy': 'camera=(), microphone=(), geolocation=(), payment=(), usb=(), interest-cohort=()',
  'Cross-Origin-Embedder-Policy': 'require-corp',
  'Cross-Origin-Opener-Policy': 'same-origin',
  'Cross-Origin-Resource-Policy': 'same-origin',
  'Content-Security-Policy': csp,
  'Cache-Control': 'no-store',
}

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    headers: securityHeaders,
    proxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/api/v1/logs/stream': { target: 'ws://localhost:8080', ws: true, changeOrigin: true },
    },
  },
  preview: {
    headers: securityHeaders,
  },
})
