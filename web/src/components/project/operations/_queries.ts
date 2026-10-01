import { queryOptions } from "@tanstack/react-query"
import type { ApiClient } from "@/hooks/use-api-client"
import type { Result } from "@/lib/dto"
import type { OperationDetailDto } from "./types"

// Also the list's react-query key prefix (see ResourceList) - every
// resource's mutations invalidate it, since each one records an operation.
export const operationsApiPath = "/api/v1/operations"

// Operations are read-only - they're only ever written by the API as part
// of the action they record, so there are no mutations here.
export const operationQueries = {
  get: (apiClient: ApiClient, id: string) =>
    queryOptions({
      queryKey: [operationsApiPath, "detail", id],
      queryFn: async () => {
        const r = await apiClient.get<Result<OperationDetailDto>>(
          `${operationsApiPath}/${id}`
        )
        return r.result
      },
    }),
}
