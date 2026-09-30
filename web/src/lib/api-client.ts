/**
 * Thin fetch wrapper around the KFamily API (same origin as this app).
 * Auth (token attachment, refresh, sign-in redirects) lives in `useApiClient`;
 * this module only knows how to talk HTTP once it has a token.
 */

export class ApiError extends Error {
  readonly status: number
  readonly body: unknown

  constructor(message: string, status: number, body?: unknown) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.body = body
  }
}

/** The access token was missing, expired, or rejected by the server. */
export class ApiAuthError extends ApiError {
  constructor(message: string, status: number, body?: unknown) {
    super(message, status, body)
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

export async function apiFetch<T>(
  path: string,
  { token, headers, body, ...init }: ApiFetchOptions = {}
): Promise<T> {
  const bodyIsFormData = isFormData(body)

  // Same origin: the API serves this app (and Vite proxies /api in dev).
  const response = await fetch(path, {
    ...init,
    headers: {
      Accept: "application/json",
      ...(body !== undefined && !bodyIsFormData
        ? { "Content-Type": "application/json" }
        : {}),
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
    body: bodyIsFormData
      ? body
      : body !== undefined
        ? JSON.stringify(body)
        : undefined,
  })

  if (response.status === 401) {
    const detail = await response.text().catch(() => undefined)
    throw new ApiAuthError(
      "The API server rejected the access token",
      response.status,
      detail
    )
  }

  if (!response.ok) {
    const detail = await response.text().catch(() => undefined)
    throw new ApiError(
      `API request failed: ${response.status} ${response.statusText}`,
      response.status,
      detail
    )
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
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
