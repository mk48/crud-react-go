/**
 * Minimal Casdoor OAuth2 authorization-code flow, hand-rolled instead of
 * pulling in casdoor-js-sdk/casdoor-react-sdk: those expect a peer
 * dependency we don't otherwise need and haven't been updated for React 19.
 *
 * Flow:
 * 1. redirectToSignin() sends the browser to Casdoor's sign-in page.
 * 2. Casdoor redirects back to /callback with `code` and `state`.
 * 3. completeSignin(code, state) verifies `state` and POSTs `code` to the
 *    API server's /api/v1/auth/signin, which holds the client secret needed
 *    to exchange it for an access token.
 * 4. That access token (a Casdoor-issued JWT) is used as the Bearer token
 *    for every other API call - see lib/auth-context.tsx.
 */
import { apiFetch } from "@/lib/api-client"
import type { Result } from "@/lib/dto"

const STATE_STORAGE_KEY = "kfamily.casdoor.state"

function requireEnv(name: string, value: string | undefined): string {
  if (!value) {
    throw new Error(`${name} is not set. Add it to your .env file (see .env.example).`)
  }
  return value
}

function getCasdoorEndpoint(): string {
  return requireEnv(
    "VITE_CASDOOR_ENDPOINT",
    import.meta.env.VITE_CASDOOR_ENDPOINT
  )
}

function getCasdoorClientId(): string {
  return requireEnv(
    "VITE_CASDOOR_CLIENT_ID",
    import.meta.env.VITE_CASDOOR_CLIENT_ID
  )
}

function getRedirectUri(): string {
  return `${window.location.origin}/callback`
}

/** Redirects the browser to Casdoor's sign-in page. */
export function redirectToSignin() {
  const state = crypto.randomUUID()
  sessionStorage.setItem(STATE_STORAGE_KEY, state)

  const url = new URL("/login/oauth/authorize", getCasdoorEndpoint())
  url.searchParams.set("client_id", getCasdoorClientId())
  url.searchParams.set("response_type", "code")
  url.searchParams.set("redirect_uri", getRedirectUri())
  url.searchParams.set("scope", "profile email")
  url.searchParams.set("state", state)

  window.location.assign(url.toString())
}

export interface SigninResult {
  accessToken: string
  refreshToken: string
  tokenType: string
  expiresIn: number
}

/**
 * Completes the sign-in flow started by redirectToSignin(): checks the
 * returned `state` against the one we generated (CSRF protection), then
 * exchanges `code` for an access token via the API server.
 */
export async function completeSignin(
  code: string,
  state: string
): Promise<string> {
  const expectedState = sessionStorage.getItem(STATE_STORAGE_KEY)
  sessionStorage.removeItem(STATE_STORAGE_KEY)
  if (!expectedState || expectedState !== state) {
    throw new Error("Sign-in state mismatch. Please try signing in again.")
  }

  const response = await apiFetch<Result<SigninResult>>("/api/v1/auth/signin", {
    method: "POST",
    body: { code, state },
  })

  if (!response.isSuccess) {
    throw new Error(response.message || "Sign-in failed.")
  }

  return response.result.accessToken
}
