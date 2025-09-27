// https://nuxt.com/docs/api/configuration/nuxt-config
const pluginId = "com.powerx.plugins.base";
const pluginAdminBase = `/_p/${pluginId}/admin/`;
const pluginApiBase = `/_p/${pluginId}/api/v1`;

// 标识当前是否被 PowerX 反代（建议在 PowerX lifecycle 启动插件进程时注入 POWERX_PROXY=1）
const INSIDE_POWERX = process.env.POWERX_PROXY === '1'

export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",

  // ✅ 被 PowerX 反代时关闭 DevTools，避免注入 __up 等开发端点
  devtools: { enabled: !INSIDE_POWERX },

  // ✅ 你选择的是纯前端 SPA（可以，Nuxt 会用 Nitro 提供 history fallback）
  ssr: false,

  // 规范：源码目录
  srcDir: "app",

  // 开发服务器（仅直连本地开发用；被 PowerX 反代时不生效）
  devServer: {
    port: 3036,
    host: "0.0.0.0",
  },

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

  css: ["~/assets/css/main.css", "@/assets/scss/main.scss"],

  ui: { fonts: false },

  colorMode: {
    preference: "system",
    fallback: "light",
    hid: "nuxt-color-mode-script",
    globalName: "__NUXT_COLOR_MODE__",
    componentName: "ColorScheme",
    classPrefix: "",
    classSuffix: "",
    storageKey: "nuxt-color-mode",
  },

  // ✅ 基础路径
  // - 直连开发（INSIDE_POWERX=false）：保持 "/"，便于 http://127.0.0.1:3036 直接访问
  // - 被 PowerX 反代（INSIDE_POWERX=true）：必须固定到 /_p/<id>/admin/，避免跳回根 "/"
  app: {
    baseURL: INSIDE_POWERX ? pluginAdminBase : "/",
    buildAssetsDir: "/assets/",
  },

  // ✅ Nitro（即使 ssr:false 也会生成 Node 服务，负责 history fallback 与静态资源）
  nitro: {
    preset: "node-server",
    serveStatic: true,
    experimental: { websocket: true },

    // ✅ 统一下发允许被同源 iframe 嵌入（双保险；反代里也会改写响应头）
    routeRules: {
      "/**": {
        headers: {
          "X-Frame-Options": "SAMEORIGIN",
          "Content-Security-Policy": "frame-ancestors 'self'",
        },
      },
    },

    // 产出目录（保持你原来）
    output: {
      dir: ".output",
      publicDir: ".output/public",
    },
  },

  // ✅ API 基址
  // - 直连开发：走本地 Vite 代理（见下）
  // - 被 PowerX 反代：走插件 API 前缀 /_p/<id>/api/...
  runtimeConfig: {
    public: {
      apiBaseUrl: INSIDE_POWERX ? pluginApiBase : "/api/v1",
    },
  },

  // ✅ Vite 代理（仅直连开发用；被 PowerX 反代时不会触发）
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
          target: "ws://127.0.0.1:4000",
          changeOrigin: true,
          ws: true,
        },
      },
    },
  },

  // ✅ i18n 调整
  // - 仍保留 prefix_except_default（/en/...）
  // - 但被 PowerX 反代时，Nuxt 生成链接会相对 baseURL，不会跑到根 "/"
  i18n: {
    defaultLocale: "zh",
    strategy: "prefix_except_default",
    locales: [
      { code: "zh", name: "简体中文", file: "zh.json" },
      { code: "en", name: "English", file: "en.json" },
    ],
    langDir: "locales",
    pages: {
      "_p/com.powerx.plugins.base/admin/plugins/base/[...internal]": {
        zh: "/_p/com.powerx.plugins.base/admin/plugins/base/:internal(.*)*",
        en: "/en/_p/com.powerx.plugins.base/admin/plugins/base/:internal(.*)*",
      },
    },
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_redirected",
      alwaysRedirect: false,
      fallbackLocale: "zh",
      // 避免因为语言检测而把你从 pluginAdminBase 重定向到根路径
      redirectOn: "no_prefix",
    },
    // 可选：SEO base，用于拼绝对 URL（不影响前端路由）
    // baseUrl: INSIDE_POWERX ? pluginAdminBase : "/",
  },
});
