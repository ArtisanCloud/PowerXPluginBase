let withNuxt
try {
  withNuxt = (await import('./.nuxt/eslint.config.mjs')).default
} catch {
  // .nuxt 未生成时，降级为直通配置
  withNuxt = /** @param {import('eslint').Linter.FlatConfig[]|undefined} cfg */ (cfg) => (Array.isArray(cfg) ? cfg : [])
}

export default withNuxt(
  {
    ignores: [
      ".output/**",
    ],
  },
  {
    files: ["**/*.{js,jsx,ts,tsx,vue}"],
    rules: {
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-unused-vars": ["warn", { argsIgnorePattern: "^_", varsIgnorePattern: "^_" }],
      "@typescript-eslint/unified-signatures": "off",
      "no-empty": "off",
      "no-console": "off",
      "no-useless-escape": "off",
      "unicorn/no-array-reduce": "off",
      "vue/attributes-order": "off",
      "vue/html-self-closing": "off",
      "vue/require-default-prop": "off",
    },
  },
)
