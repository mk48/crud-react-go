/**
 * Browser tracing (see docs/tracing.md). Every API call made through
 * `apiFetch` is a span whose W3C `traceparent` header is sent along, so the
 * API's spans - and the operations it records - join the browser's trace.
 *
 * Deliberately not the auto-instrumentation of `fetch`: that would also put
 * trace headers on the browser's calls to Casdoor (token refresh), which
 * aren't ours to trace and would fail its CORS preflight.
 */
import { trace, type Tracer } from "@opentelemetry/api"
import { W3CTraceContextPropagator } from "@opentelemetry/core"
import { resourceFromAttributes } from "@opentelemetry/resources"
import {
  BatchSpanProcessor,
  ParentBasedSampler,
  TraceIdRatioBasedSampler,
  WebTracerProvider,
  type SpanProcessor,
} from "@opentelemetry/sdk-trace-web"
import {
  ATTR_SERVICE_NAME,
  ATTR_SERVICE_VERSION,
} from "@opentelemetry/semantic-conventions"

/**
 * Registers the tracer provider - call (and await) before the first API
 * call. Spans are always created, so every request carries a trace id, but
 * only exported when the API's /config.js names a collector
 * (`otelTracesUrl`, from WEB_OTEL_TRACES_URL). The exporter is loaded only
 * then, keeping it out of the main bundle.
 */
export async function initTelemetry(): Promise<void> {
  const config = window.__KFAMILY_CONFIG__
  if (!config) return // runtimeConfig() reports this when the app starts

  const spanProcessors: SpanProcessor[] = []
  if (config.otelTracesUrl) {
    const { OTLPTraceExporter } =
      await import("@opentelemetry/exporter-trace-otlp-http")
    spanProcessors.push(
      new BatchSpanProcessor(
        new OTLPTraceExporter({ url: config.otelTracesUrl })
      )
    )
  }

  const provider = new WebTracerProvider({
    resource: resourceFromAttributes({
      [ATTR_SERVICE_NAME]: "kfamily-web",
      [ATTR_SERVICE_VERSION]: config.version,
    }),
    // The browser starts each trace, so it makes the sampling decision; the
    // API follows it (parent-based).
    sampler: new ParentBasedSampler({
      root: new TraceIdRatioBasedSampler(config.otelSampleRatio ?? 1),
    }),
    spanProcessors,
  })
  provider.register({ propagator: new W3CTraceContextPropagator() })
}

export function tracer(): Tracer {
  return trace.getTracer("kfamily-web")
}

const idSegment =
  /^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}|\d+)$/i

/**
 * The span name for an API path - ids replaced by `:id` and the query
 * dropped (e.g. `/api/v1/samples/:id`), so spans group by endpoint instead
 * of one name per record.
 */
export function routeOf(path: string): string {
  return path
    .split("?")[0]
    .split("/")
    .map((segment) => (idSegment.test(segment) ? ":id" : segment))
    .join("/")
}
