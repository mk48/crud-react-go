// Settings the API serves at /config.js (loaded by index.html before the
// app) from its own environment - so one build runs in every environment,
// instead of values being baked into the bundle at build time. In dev, Vite
// proxies /config.js to the API (see vite.config.ts).
export interface RuntimeConfig {
  casdoorEndpoint: string
  casdoorClientId: string
  // The API's version - this app is embedded in it, so it's ours too.
  version: string
  // Browser tracing (see lib/telemetry.ts): the OTLP/HTTP traces URL to
  // export spans to (none: spans are only propagated, not exported) and the
  // share of traces to sample.
  otelTracesUrl?: string
  otelSampleRatio?: number
  // Link to a trace in the tracing UI, with {traceId} as placeholder.
  traceUrlTemplate?: string
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

/** The tracing UI's page for traceId, if TRACE_URL_TEMPLATE is set. */
export function traceUrl(traceId: string): string | undefined {
  return window.__KFAMILY_CONFIG__?.traceUrlTemplate?.replace(
    "{traceId}",
    encodeURIComponent(traceId)
  )
}
