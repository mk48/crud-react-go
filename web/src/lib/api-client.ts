/**
 * Thin fetch wrapper around the KFamily API (same origin as this app).
 * Auth (token attachment, refresh, sign-in redirects) lives in `useApiClient`;
 * this module only knows how to talk HTTP once it has a token.
 *
 * Every call is traced (see lib/telemetry.ts): it's a span, its
 * `traceparent` header continues the trace into the API, and a failed
 * call's ApiError carries the trace id for the user to quote.
 */
import {
  context,
  propagation,
  SpanKind,
  SpanStatusCode,
  trace,
} from "@opentelemetry/api"
import {
  ATTR_HTTP_REQUEST_METHOD,
  ATTR_HTTP_RESPONSE_STATUS_CODE,
  ATTR_HTTP_ROUTE,
  ATTR_URL_PATH,
} from "@opentelemetry/semantic-conventions"
import { routeOf, tracer } from "./telemetry"

export class ApiError extends Error {
  readonly status: number
  readonly body: unknown
  // The request's trace id - shown in error messages, so a report can be
  // matched to its trace, log lines and operations.
  readonly traceId?: string

  constructor(
    message: string,
    status: number,
    body?: unknown,
    traceId?: string
  ) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.body = body
    this.traceId = traceId
  }
}

/** The access token was missing, expired, or rejected by the server. */
export class ApiAuthError extends ApiError {
  constructor(
    message: string,
    status: number,
    body?: unknown,
    traceId?: string
  ) {
    super(message, status, body, traceId)
    this.name = "ApiAuthError"
  }
}

export interface ApiFetchOptions extends Omit<RequestInit, "body"> {
  body?: unknown
  token?: string
}

// FormData bodies (file uploads) must be sent as-is, letting the browser set
// its own multipart/form-data boundary - JSON.stringify-ing one just turns
// it into the useless string "[object FormData]".
function isFormData(body: unknown): body is FormData {
  return typeof FormData !== "undefined" && body instanceof FormData
}

// How this app identifies itself (see the API's util.ClientVersionHeader).
// Recorded on each operation as a hint; the app itself is identified by
// the access token, not by these.
function clientHeaders(): Record<string, string> {
  return {
    "X-Client-Platform": "web",
    "X-Client-Version": window.__KFAMILY_CONFIG__?.version ?? "unknown",
  }
}

export async function apiFetch<T>(
  path: string,
  { token, headers, body, ...init }: ApiFetchOptions = {}
): Promise<T> {
  const bodyIsFormData = isFormData(body)
  const method = (init.method ?? "GET").toUpperCase()
  const route = routeOf(path)

  const span = tracer().startSpan(`${method} ${route}`, {
    kind: SpanKind.CLIENT,
    attributes: {
      [ATTR_HTTP_REQUEST_METHOD]: method,
      [ATTR_HTTP_ROUTE]: route,
      [ATTR_URL_PATH]: path.split("?")[0],
    },
  })
  const traceHeaders: Record<string, string> = {}
  propagation.inject(trace.setSpan(context.active(), span), traceHeaders)

  // The API echoes the id in X-Trace-Id; fall back to our own span's id
  // when the request never got an answer.
  const spanContext = span.spanContext()
  const fallbackTraceId = trace.isSpanContextValid(spanContext)
    ? spanContext.traceId
    : undefined

  try {
    // Same origin: the API serves this app (and Vite proxies /api in dev).
    const response = await fetch(path, {
      ...init,
      headers: {
        Accept: "application/json",
        ...(body !== undefined && !bodyIsFormData
          ? { "Content-Type": "application/json" }
          : {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...clientHeaders(),
        ...traceHeaders,
        ...headers,
      },
      body: bodyIsFormData
        ? body
        : body !== undefined
          ? JSON.stringify(body)
          : undefined,
    })
    span.setAttribute(ATTR_HTTP_RESPONSE_STATUS_CODE, response.status)
    const traceId = response.headers.get("X-Trace-Id") ?? fallbackTraceId

    if (response.status === 401) {
      span.setStatus({ code: SpanStatusCode.ERROR })
      const detail = await response.text().catch(() => undefined)
      throw new ApiAuthError(
        "The API server rejected the access token",
        response.status,
        detail,
        traceId
      )
    }

    if (!response.ok) {
      span.setStatus({ code: SpanStatusCode.ERROR })
      const detail = await response.text().catch(() => undefined)
      throw new ApiError(
        `API request failed: ${response.status} ${response.statusText}`,
        response.status,
        detail,
        traceId
      )
    }

    if (response.status === 204) {
      return undefined as T
    }

    return (await response.json()) as T
  } catch (error) {
    if (!(error instanceof ApiError)) {
      // Network failure, aborted request or an unreadable body.
      span.recordException(error as Error)
      span.setStatus({ code: SpanStatusCode.ERROR })
      throw new ApiError(
        (error as Error).message || "API request failed",
        0,
        undefined,
        fallbackTraceId
      )
    }
    throw error
  } finally {
    span.end()
  }
}

export function toQueryString(
  params: Record<string, string | number | boolean | undefined>
): string {
  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined && value !== "") {
      search.set(key, String(value))
    }
  }
  const qs = search.toString()
  return qs ? `?${qs}` : ""
}
