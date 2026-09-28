import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { useAuth } from "@/lib/auth-context"
import { redirectToSignin } from "@/lib/casdoor"
import { useApiClient } from "@/hooks/use-api-client"
import { userQueries } from "@/components/project/user/_queries"
import { Button } from "@/components/ui/button"

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { token, authError, signOut, retrySignin } = useAuth()
  const apiClient = useApiClient()

  useEffect(() => {
    if (!token && !authError) {
      redirectToSignin()
    }
  }, [token, authError])

  const meQuery = useQuery({
    ...userQueries.me(apiClient),
    enabled: !!token,
    retry: false,
  })

  if (authError) {
    return (
      <CenteredMessage title="Couldn't sign you in" detail={authError}>
        <Button className="w-fit" variant="outline" onClick={retrySignin}>
          Try again
        </Button>
      </CenteredMessage>
    )
  }

  if (!token || meQuery.isPending) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <p className="text-sm text-muted-foreground">
          {token ? "Loading…" : "Redirecting to sign-in…"}
        </p>
      </div>
    )
  }

  // 401s are handled globally (redirect to sign-in - see AuthProvider); this
  // covers the rest, e.g. 403 for a deleted account or the API being down.
  if (meQuery.isError) {
    return (
      <CenteredMessage
        title="Couldn't verify your session"
        detail="Please sign in again."
      >
        <Button className="w-fit" variant="outline" onClick={signOut}>
          Sign out
        </Button>
      </CenteredMessage>
    )
  }

  return <>{children}</>
}

function CenteredMessage({
  title,
  detail,
  children,
}: {
  title: string
  detail: string
  children: React.ReactNode
}) {
  return (
    <div className="flex min-h-svh items-center justify-center p-6">
      <div className="flex max-w-md min-w-0 flex-col gap-4 text-sm leading-loose">
        <div>
          <h1 className="font-medium">{title}</h1>
          <p>{detail}</p>
        </div>
        {children}
      </div>
    </div>
  )
}
