/* eslint-disable react-refresh/only-export-components */
import * as React from "react"

const TOKEN_STORAGE_KEY = "kfamily.accessToken"

interface AuthContextValue {
  token: string | null
  setToken: (token: string) => void
  signOut: () => void
}

const AuthContext = React.createContext<AuthContextValue | undefined>(
  undefined
)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [token, setTokenState] = React.useState<string | null>(() =>
    localStorage.getItem(TOKEN_STORAGE_KEY)
  )

  const setToken = React.useCallback((next: string) => {
    localStorage.setItem(TOKEN_STORAGE_KEY, next)
    setTokenState(next)
  }, [])

  const signOut = React.useCallback(() => {
    localStorage.removeItem(TOKEN_STORAGE_KEY)
    setTokenState(null)
  }, [])

  const value = React.useMemo(
    () => ({ token, setToken, signOut }),
    [token, setToken, signOut]
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
