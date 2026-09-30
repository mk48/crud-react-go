// Settings the API serves at /config.js (loaded by index.html before the
// app) from its own environment - so one build runs in every environment,
// instead of values being baked into the bundle at build time. In dev, Vite
// proxies /config.js to the API (see vite.config.ts).
export interface RuntimeConfig {
  casdoorEndpoint: string
  casdoorClientId: string
}

declare global {
  interface Window {
    __KFAMILY_CONFIG__?: RuntimeConfig
  }
}

export function runtimeConfig(): RuntimeConfig {
  const config = window.__KFAMILY_CONFIG__
  if (!config) {
    throw new Error(
      "App configuration (/config.js) didn't load - is the API server running?"
    )
  }
  return config
}
