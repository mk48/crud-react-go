import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useEffect, useRef, useState } from "react"
import { Button } from "@/components/ui/button"
import { useAuth } from "@/lib/auth-context"
import {
  SIGNED_OUT_STATE,
  completeSignin,
  redirectToSignin,
} from "@/lib/casdoor"

export const Route = createFileRoute("/callback")({
  component: CallbackPage,
})

function getCodeAndState() {
  const params = new URLSearchParams(window.location.search)
  return { code: params.get("code"), state: params.get("state") }
}

function CallbackPage() {
  const { code, state } = getCodeAndState()

  // Casdoor sends the browser back here after signOutOfCasdoor() too.
  if (!code && state === SIGNED_OUT_STATE) {
    return <SignedOut />
  }

  return <CompleteSignin />
}

function SignedOut() {
  return (
    <div className="flex min-h-svh items-center justify-center">
      <div className="flex flex-col items-center gap-4 text-sm">
        <p className="text-muted-foreground">You have been signed out.</p>
        <Button variant="outline" onClick={() => redirectToSignin("/")}>
          Sign in
        </Button>
      </div>
    </div>
  )
}

function CompleteSignin() {
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
      .then(({ token, returnTo }) => {
        setToken(token)
        void navigate({ href: returnTo, replace: true })
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
