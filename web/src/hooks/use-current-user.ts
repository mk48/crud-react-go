import { useQuery } from "@tanstack/react-query"
import { userQueries } from "@/components/project/user/_queries"
import { useApiClient } from "@/hooks/use-api-client"

/** The signed-in user (already loaded and cached by RequireAuth). */
export function useCurrentUser() {
  const apiClient = useApiClient()
  return useQuery(userQueries.me(apiClient)).data
}

/**
 * Whether the signed-in user may create/edit/delete. Only a UI hint - the
 * API enforces it (AdminMiddleware) regardless.
 */
export function useIsAdmin(): boolean {
  return useCurrentUser()?.isAdmin ?? false
}
