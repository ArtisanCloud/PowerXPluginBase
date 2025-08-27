// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },

  // PowerX Plugin Configuration
  ssr: false, // 纯前端SPA模式
  
  // Development server configuration
  devServer: {
    port: 3036,
    host: '0.0.0.0'
  },
  
  // Base path for PowerX plugin integration
  app: {
    baseURL: process.env.NODE_ENV === 'production' 
      ? '/_p/com.powerx.plugins.scrum/admin/' 
      : '/',
    buildAssetsDir: '/assets/'
  },

  // Build configuration for plugin deployment
  nitro: {
    output: {
      dir: '.output',
      publicDir: '.output/public'
    }
  },

  modules: [
    '@nuxt/eslint',
    '@nuxt/content',
    '@nuxt/image',
    '@nuxt/scripts',
    '@nuxt/test-utils',
    '@nuxt/ui'
  ],

  // CSS framework configuration
  css: [
    // Add global styles here
  ],

  // Runtime config for API integration
  runtimeConfig: {
    public: {
      apiBaseUrl: process.env.NODE_ENV === 'production'
        ? '/_p/com.powerx.plugins.scrum/api/v1'
        : 'http://localhost:8086/v1'
    }
  }
})