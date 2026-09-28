/* eslint-disable react-refresh/only-export-components */
import * as React from "react"
import { redirectToSignin, signOutOfCasdoor } from "@/lib/casdoor"
import { queryClient, setAuthErrorHandler } from "@/lib/query-client"

const TOKEN_STORAGE_KEY = "kfamily.accessToken"
// When we last bounced the user to Casdoor because the API rejected their
// token - see onAuthError's loop guard.
const AUTH_REDIRECT_AT_KEY = "kfamily.authRedirectAt"
const AUTH_REDIRECT_COOLDOWN_MS = 30_000

interface AuthContextValue {
  token: string | null
  // Set when the API keeps rejecting fresh tokens, so RequireAuth shows an
  // error instead of redirecting to sign-in in a loop.
  authError: string | null
  setToken: (token: string) => void
  signOut: () => void
  retrySignin: () => void
}

const AuthContext = React.createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = React.useState<string | null>(() =>
    localStorage.getItem(TOKEN_STORAGE_KEY)
  )
  const [authError, setAuthError] = React.useState<string | null>(null)

  const setToken = React.useCallback((next: string) => {
    localStorage.setItem(TOKEN_STORAGE_KEY, next)
    setAuthError(null)
    setTokenState(next)
  }, [])

  // Clears this app's token and every cached response (so nothing of this
  // user's data survives for the next one), then ends the Casdoor session.
  const signOut = React.useCallback(() => {
    const current = localStorage.getItem(TOKEN_STORAGE_KEY)
    localStorage.removeItem(TOKEN_STORAGE_KEY)
    queryClient.clear()
    if (current) {
      // Navigates away - no state update, or RequireAuth would race it
      // with a sign-in redirect.
      signOutOfCasdoor(current)
    } else {
      setTokenState(null)
    }
  }, [])

  const retrySignin = React.useCallback(() => {
    sessionStorage.removeItem(AUTH_REDIRECT_AT_KEY)
    redirectToSignin()
  }, [])

  // Any 401 from the API (expired/revoked token): sign in again. If Casdoor
  // still has a session this is a quick bounce back to the same page.
  React.useEffect(() => {
    let redirecting = false

    return setAuthErrorHandler(() => {
      if (redirecting) return
      localStorage.removeItem(TOKEN_STORAGE_KEY)

      // A token that was *just* issued being rejected again means sign-in
      // itself is broken (e.g. API/Casdoor misconfiguration) - redirecting
      // would loop forever, so stop and show an error instead.
      const lastRedirect = Number(sessionStorage.getItem(AUTH_REDIRECT_AT_KEY))
      if (Date.now() - lastRedirect < AUTH_REDIRECT_COOLDOWN_MS) {
        queryClient.clear()
        setAuthError("The server keeps rejecting your sign-in.")
        setTokenState(null)
        return
      }

      redirecting = true
      sessionStorage.setItem(AUTH_REDIRECT_AT_KEY, String(Date.now()))
      redirectToSignin()
    })
  }, [])

  const value = React.useMemo(
    () => ({ token, authError, setToken, signOut, retrySignin }),
    [token, authError, setToken, signOut, retrySignin]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = React.useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
