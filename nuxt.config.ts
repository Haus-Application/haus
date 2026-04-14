export default defineNuxtConfig({
  srcDir: 'frontend/',
  devtools: { enabled: true },
  ssr: false,
  compatibilityDate: '2025-01-01',
  css: ['~/assets/css/global.css'],
  // Vite handles the dev server when ssr: false, so its proxy config is what
  // actually forwards /api and /api/ws to the Go backend on :8080. The
  // nitro.devProxy block is kept for production-style previews.
  vite: {
    server: {
      proxy: {
        // Vite's HTTP proxy forwards REST + SSE to the Go backend. WebSocket
        // (/api/ws) is NOT proxied here — the Vite/Nuxt dev server doesn't
        // reliably pass WS upgrades through. useWebSocket.ts connects
        // directly to :8080 in dev instead.
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true,
          // Long-lived SSE streams (/api/scan/stream) close client-side when
          // the browser disconnects, which bubbles up as ECONNRESET in
          // http-proxy. Without a handler, Vite escalates it to an
          // unhandledRejection and Nuxt restarts mid-scan. Swallow it.
          configure: (proxy: any) => {
            proxy.on('error', (err: any) => {
              if (err?.code !== 'ECONNRESET') {
                console.warn('[vite proxy]', err?.message ?? err)
              }
            })
          },
        },
      },
    },
  },
  nitro: {
    devProxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
  app: {
    head: {
      title: 'Haus',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1, viewport-fit=cover' },
        { name: 'description', content: 'Haus — Smart Home Discovery' },
        { name: 'apple-mobile-web-app-capable', content: 'yes' },
        { name: 'theme-color', content: '#0D0D0F' },
      ],
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        { rel: 'stylesheet', href: 'https://fonts.googleapis.com/css2?family=DM+Sans:wght@400;500;600;700&family=Inter:wght@400;500;600&display=swap' },
      ],
    },
  },
})
