// web-admin/nuxt.config.ts
const pluginId = "com.powerx.plugins.base";
const pluginAdminBase = `/_p/${pluginId}/admin/`;
const pluginApiBase   = `/_p/${pluginId}/api/v1`;
const localApiBase    = process.env.NUXT_PUBLIC_API_BASE || "http://localhost:8086/api/v1";

// Host（被 PowerX 反代）用 1；独立/开发用 0
const INSIDE_POWERX = process.env.POWERX_PROXY === '1';
const rawBridgeDebug = process.env.NUXT_PUBLIC_BRIDGE_DEBUG ?? process.env.BRIDGE_DEBUG;
const BRIDGE_DEBUG = rawBridgeDebug !== undefined
  ? /^(1|true)$/i.test(rawBridgeDebug)
  : !INSIDE_POWERX;

export default defineNuxtConfig({
  compatibilityDate: "2025-07-15",
  ssr: false,
  srcDir: "app",
  devtools: { enabled: !INSIDE_POWERX },

  app: {
    baseURL: INSIDE_POWERX ? pluginAdminBase : "/",
    buildAssetsDir: "/assets/",
    head: {
      meta: [
        { name: "referrer", content: "no-referrer" },
        { httpEquiv: "X-Content-Type-Options", content: "nosniff" },
        { name: "permissions-policy", content: "camera=(), microphone=(), geolocation=()" },
      ],
    },
  },

  css: ["~/assets/css/main.css", "@/assets/scss/main.scss"],

  // Tailwind v4 正确写法（避免你遇到的 @tailwindcss/postcss 报错）
  postcss: {
    plugins: {
      "@tailwindcss/postcss": {},
      autoprefixer: {},
    },
  },

  modules: [
    "@nuxt/ui",
    "@nuxt/icon",
    "@nuxtjs/i18n",
    "@pinia/nuxt",
    "@nuxtjs/color-mode",
    "@nuxt/content",
    "@nuxt/image",
  ],

  colorMode: {
    preference: "system",
    fallback: "light",
    storageKey: "powerx-color-mode",
  },

  i18n: {
    defaultLocale: "zh",
    // strategy: "prefix_except_default", // zh 无前缀，en 带 /en
    strategy: "no_prefix",
    locales: [
      { code: "zh", name: "简体中文", file: "zh.json" },
      { code: "en", name: "English", file: "en.json" },
    ],
    langDir: "locales",
    detectBrowserLanguage: INSIDE_POWERX ? false: {
      useCookie: true,
      cookieKey: "px_lang",
      alwaysRedirect: true,
      fallbackLocale: "zh",
      redirectOn: "root", // 仅在根路径做检测，避免插件地址被自动加语言前缀
    },
  },

  runtimeConfig: {
    public: {
      apiBaseUrl: INSIDE_POWERX ? pluginApiBase : localApiBase,
      insidePowerX: INSIDE_POWERX,
      pluginAdminBase,
      bridgeDebug: BRIDGE_DEBUG,
    },
  },

  nitro: {
    preset: "node-server",
    serveStatic: true,
    experimental: { websocket: true },
    routeRules: {
      "/**": {
        headers: {
          "X-Frame-Options": "SAMEORIGIN",
          "Content-Security-Policy": `default-src 'self'; img-src 'self' data:; style-src 'self' 'unsafe-inline'; script-src 'self'; connect-src 'self' ${INSIDE_POWERX ? pluginApiBase : localApiBase}; font-src 'self' data:; frame-ancestors 'self';`,
          "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
          "Referrer-Policy": "no-referrer",
        },
      },
    },
    output: {
      dir: ".output",
      publicDir: ".output/public",
    },
  },

  vite: {
    server: {
      proxy: INSIDE_POWERX ? {} : {
        "/api": {
          target: "http://localhost:8086",
          changeOrigin: true,
          ws: true,
          // rewrite: (p: string) => p.replace(/^\/api/, ""),
        },
        "/ws": {
          target: "ws://127.0.0.1:4000",
          changeOrigin: true,
          ws: true,
        },
      },
    },
  },

  devServer: { port: 3036, host: "0.0.0.0" },

  ui: { fonts: false },
});
