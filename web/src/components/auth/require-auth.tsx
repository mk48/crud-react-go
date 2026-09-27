import { useEffect } from "react"
import { useQuery } from "@tanstack/react-query"
import { useAuth } from "@/lib/auth-context"
import { redirectToSignin } from "@/lib/casdoor"
import { useApiClient } from "@/hooks/use-api-client"
import { userQueries } from "@/components/project/user/_queries"
import { Button } from "@/components/ui/button"

export function RequireAuth({ children }: { children: React.ReactNode }) {
  const { token, signOut } = useAuth()
  const apiClient = useApiClient()

  useEffect(() => {
    if (!token) {
      redirectToSignin()
    }
  }, [token])

  const meQuery = useQuery({
    ...userQueries.me(apiClient),
    enabled: !!token,
    retry: false,
  })

  if (!token || meQuery.isPending) {
    return (
      <div className="flex min-h-svh items-center justify-center">
        <p className="text-sm text-muted-foreground">
          {token ? "Loading…" : "Redirecting to sign-in…"}
        </p>
      </div>
    )
  }

  if (meQuery.isError) {
    return (
      <div className="flex min-h-svh items-center justify-center p-6">
        <div className="flex max-w-md min-w-0 flex-col gap-4 text-sm leading-loose">
          <div>
            <h1 className="font-medium">Couldn't verify your session</h1>
            <p>Please sign in again.</p>
          </div>

          <Button className="w-fit" variant="outline" onClick={signOut}>
            Sign out
          </Button>
        </div>
      </div>
    )
  }

  return <>{children}</>
}
