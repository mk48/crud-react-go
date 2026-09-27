import { queryOptions } from "@tanstack/react-query"
import type { ApiClient } from "@/hooks/use-api-client"
import { toQueryString } from "@/lib/api-client"
import type { PaginationResult, Result } from "@/lib/dto"
import type { UserDto, UserRequestDto } from "./types"
import { MINUTE } from "@/lib/constants"

const apiPath = "/api/v1/users"

export const userQueries = {
  me: (apiClient: ApiClient) =>
    queryOptions({
      queryKey: ["users-me"],
      queryFn: async () => {
        const r = await apiClient.get<Result<UserDto>>(`${apiPath}/me`)
        return r.result
      },
    }),

  get: (apiClient: ApiClient, id: string, includeDeleted: boolean = false) =>
    queryOptions({
      queryKey: ["users", id, includeDeleted],
      queryFn: async () => {
        const r = await apiClient.get<Result<UserDto>>(
          `${apiPath}/${id}${toQueryString({ includeDeleted })}`
        )
        return r.result
      },
      staleTime: 30 * MINUTE, // /api/v1/users/:id changes rarely once loaded
    }),

  list: (
    apiClient: ApiClient,
    searchText: string,
    pageIndex: number = 0,
    recordsPerPage: number = 50
  ) =>
    queryOptions({
      queryKey: ["users-list", searchText, pageIndex, recordsPerPage],
      queryFn: async () => {
        const response = await apiClient.get<Result<PaginationResult<UserDto>>>(
          `${apiPath}${toQueryString({ searchText, pageIndex, recordsPerPage })}`
        )
        return response.result
      },
    }),
}

export const userMutations = {
  create: (apiClient: ApiClient) => ({
    mutationFn: async (dataToServer: UserRequestDto) => {
      return apiClient.post(apiPath, { ...dataToServer })
    },
  }),

  update: (apiClient: ApiClient, id: string) => ({
    mutationFn: async (dataToServer: UserRequestDto) => {
      return apiClient.put(`${apiPath}/${id}`, { ...dataToServer })
    },
  }),

  delete: (apiClient: ApiClient, id: string) => ({
    mutationFn: async () => {
      return apiClient.delete(`${apiPath}/${id}`)
    },
  }),
}
