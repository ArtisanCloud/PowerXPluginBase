let withNuxt
try {
  withNuxt = (await import('./.nuxt/eslint.config.mjs')).default
} catch {
  // .nuxt 未生成时，降级为直通配置
  withNuxt = /** @param {import('eslint').Linter.FlatConfig[]|undefined} cfg */ (cfg) => (Array.isArray(cfg) ? cfg : [])
}

export default withNuxt(
  // Your custom configs here
)
