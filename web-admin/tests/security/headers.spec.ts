import config from '../../nuxt.config'

describe('security headers', () => {
  it('enforces CSP and security headers on all routes', () => {
    const headers = config.nitro?.routeRules?.['/**']?.headers
    expect(headers).toBeTruthy()
    expect(headers).toMatchObject({
      'X-Frame-Options': 'SAMEORIGIN',
    })
    expect(headers?.['Content-Security-Policy']).toContain("default-src 'self'")
    expect(headers?.['Strict-Transport-Security']).toContain('max-age')
  })
})
