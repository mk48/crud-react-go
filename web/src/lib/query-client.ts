import {
  MutationCache,
  QueryCache,
  QueryClient,
  type QueryKey,
} from "@tanstack/react-query"
import { ApiAuthError } from "@/lib/api-client"

declare module "@tanstack/react-query" {
  interface Register {
    mutationMeta: {
      // Query key prefixes to invalidate once the mutation succeeds - e.g.
      // [apiPath] refreshes every list/detail/option query of that resource.
      invalidates?: QueryKey[]
    }
  }
}

// AuthProvider registers what to do when any query/mutation gets a 401
// (expired or revoked token) - see setAuthErrorHandler.
let authErrorHandler: (() => void) | undefined

export function setAuthErrorHandler(handler: () => void) {
  authErrorHandler = handler
  return () => {
    if (authErrorHandler === handler) authErrorHandler = undefined
  }
}

function onError(error: Error) {
  if (error instanceof ApiAuthError) authErrorHandler?.()
}

export const queryClient: QueryClient = new QueryClient({
  queryCache: new QueryCache({ onError }),
  mutationCache: new MutationCache({
    onError,
    onSuccess: async (_data, _variables, _context, mutation) => {
      await Promise.all(
        (mutation.meta?.invalidates ?? []).map((queryKey) =>
          queryClient.invalidateQueries({ queryKey })
        )
      )
    },
  }),
  defaultOptions: {
    queries: {
      retry: false,
      refetchOnWindowFocus: false,
    },
  },
})
