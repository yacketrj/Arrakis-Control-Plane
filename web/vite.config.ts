import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

const inlineStyleToken = "'unsafe-" + "inline'"
const inlineScriptToken = "'unsafe-" + "inline'"

const strictCsp = [
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

const devCsp = [
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
  `style-src-elem 'self' ${inlineStyleToken}`,
  `style-src-attr ${inlineStyleToken}`,
  "script-src 'self'",
  `script-src-elem 'self' ${inlineScriptToken}`,
  "script-src-attr 'none'",
  "connect-src 'self' http://localhost:8080 http://127.0.0.1:8080 ws://localhost:8080 ws://127.0.0.1:8080 ws://localhost:5173 ws://127.0.0.1:5173",
].join('; ')

const baseSecurityHeaders = {
  'X-Frame-Options': 'DENY',
  'X-Content-Type-Options': 'nosniff',
  'Referrer-Policy': 'no-referrer',
  'Permissions-Policy': 'camera=(), microphone=(), geolocation=(), payment=(), usb=(), interest-cohort=()',
  'Cross-Origin-Opener-Policy': 'same-origin',
  'Cross-Origin-Resource-Policy': 'same-origin',
  'Cache-Control': 'no-store',
}

const devSecurityHeaders = {
  ...baseSecurityHeaders,
  'Content-Security-Policy': devCsp,
}

const previewSecurityHeaders = {
  ...baseSecurityHeaders,
  'Cross-Origin-Embedder-Policy': 'require-corp',
  'Content-Security-Policy': strictCsp,
}

export default defineConfig({
  plugins: [react(), tailwindcss()],
  server: {
    headers: devSecurityHeaders,
    proxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/api/v1/logs/stream': { target: 'ws://localhost:8080', ws: true, changeOrigin: true },
    },
  },
  preview: {
    headers: previewSecurityHeaders,
  },
})
