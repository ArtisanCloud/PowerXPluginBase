// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  devtools: { enabled: true },

  // PowerX Plugin Configuration
  ssr: false, // 纯前端SPA模式
  srcDir: "app", // Nuxt 4 规范：所有源码在app目录下

  // Development server configuration
  devServer: {
    port: 3036,
    host: "0.0.0.0",
  },

  // Modules configuration
  modules: [
    "@nuxt/ui",
    "@nuxt/icon",
    "@nuxtjs/i18n",
    "@pinia/nuxt",
    "@nuxtjs/color-mode",
    "@nuxt/eslint",
    "@nuxt/content",
    "@nuxt/image",
    "@nuxt/scripts",
    "@nuxt/test-utils",
  ],

  nitro: {
    experimental: {
      websocket: true, // ✅ 开启 Nitro 原生 WS
    },
  },

  // CSS configuration
  css: ["~/assets/css/main.css", "@/assets/scss/main.scss"],

  // UI configuration
  ui: {
    fonts: false,
  },

  // Color mode configuration
  colorMode: {
    preference: "system", // default value of $colorMode.preference
    fallback: "light", // fallback value if not system preference found
    hid: "nuxt-color-mode-script",
    globalName: "__NUXT_COLOR_MODE__",
    componentName: "ColorScheme",
    classPrefix: "",
    classSuffix: "",
    storageKey: "nuxt-color-mode",
  },

  // Base path for PowerX plugin integration
  app: {
    baseURL:
      process.env.NODE_ENV === "production"
        ? "/_p/com.powerx.plugins.note/admin/"
        : "/",
    buildAssetsDir: "/assets/",
  },

  // Nitro build configuration for plugin deployment
  nitro: {
    output: {
      dir: ".output",
      publicDir: ".output/public",
    },
  },

  // Runtime config for API integration
  runtimeConfig: {
    public: {
      apiBaseUrl:
        process.env.NODE_ENV === "production"
          ? "/_p/com.powerx.plugin.base/api/v1"
          : "/api/v1", // 改为相对路径，交给 vite 代理，避免 CORS
    },
  },

  // Vite 配置 - 开发环境代理
  vite: {
    server: {
      proxy: {
        "/api": {
          target: "http://localhost:8086",
          changeOrigin: true,
          ws: true,
          rewrite: (p: string) => p.replace(/^\/api/, ""),
        },
        "/ws": {
          target: "ws://127.0.0.1:4000", // 修改为你的 WebSocket 服务地址
          changeOrigin: true,
          ws: true, // 启用 WebSocket 代理
        },
      },
    },
  },

  // Internationalization configuration
  i18n: {
    defaultLocale: "zh",
    locales: [
      { code: "zh", name: "简体中文", file: "zh.json" },
      { code: "en", name: "English", file: "en.json" },
    ],
    langDir: "locales",
    strategy: "no_prefix",
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_redirected",
      alwaysRedirect: false,
      fallbackLocale: "zh",
    },
  },
});
