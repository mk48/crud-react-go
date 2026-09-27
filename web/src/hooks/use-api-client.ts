import * as React from "react"
import { apiFetch, type ApiFetchOptions } from "@/lib/api-client"
import { useAuth } from "@/lib/auth-context"

export interface ApiClient {
  get<T>(path: string): Promise<T>
  post<T>(path: string, body?: unknown): Promise<T>
  put<T>(path: string, body?: unknown): Promise<T>
  delete<T>(path: string): Promise<T>
}

/** A fetch client bound to the current signed-in user's access token. */
export function useApiClient(): ApiClient {
  const { token } = useAuth()

  return React.useMemo<ApiClient>(() => {
    const request = <T,>(path: string, options?: ApiFetchOptions) =>
      apiFetch<T>(path, { ...options, token: token ?? undefined })

    return {
      get: (path) => request(path),
      post: (path, body) => request(path, { method: "POST", body }),
      put: (path, body) => request(path, { method: "PUT", body }),
      delete: (path) => request(path, { method: "DELETE" }),
    }
  }, [token])
}
