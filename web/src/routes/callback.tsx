import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useEffect, useRef, useState } from "react"
import { useAuth } from "@/lib/auth-context"
import { completeSignin } from "@/lib/casdoor"

export const Route = createFileRoute("/callback")({
  component: CallbackPage,
})

function getCodeAndState() {
  const params = new URLSearchParams(window.location.search)
  return { code: params.get("code"), state: params.get("state") }
}

function CallbackPage() {
  const navigate = useNavigate()
  const { setToken } = useAuth()
  const [error, setError] = useState<string | null>(() => {
    const { code, state } = getCodeAndState()
    return code && state ? null : "Casdoor redirect is missing code/state."
  })
  const ran = useRef(false)

  useEffect(() => {
    if (ran.current || error) return
    ran.current = true

    const { code, state } = getCodeAndState()
    // error above already guards against either being missing
    completeSignin(code!, state!)
      .then((token) => {
        setToken(token)
        void navigate({ to: "/", replace: true })
      })
      .catch((err: unknown) => {
        setError(err instanceof Error ? err.message : "Sign-in failed.")
      })
  }, [navigate, setToken, error])

  return (
    <div className="flex min-h-svh items-center justify-center">
      <p className="text-sm text-muted-foreground">
        {error ?? "Signing you in…"}
      </p>
    </div>
  )
}
