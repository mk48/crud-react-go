import { queryOptions } from "@tanstack/react-query"
import type { ApiClient } from "@/hooks/use-api-client"
import { toQueryString } from "@/lib/api-client"
import type { PaginationResult, Result } from "@/lib/dto"
import type { SampleChildrenDto, SampleChildrenRequestDto } from "./types"

const apiPath = "/api/v1/sample-children"

export const sampleChildrenQueries = {
  get: (apiClient: ApiClient, id: string, includeDeleted: boolean = false) =>
    queryOptions({
      queryKey: ["sample-children", id, includeDeleted],
      queryFn: async () => {
        const r = await apiClient.get<Result<SampleChildrenDto>>(
          `${apiPath}/${id}${toQueryString({ includeDeleted })}`
        )
        return r.result
      },
    }),

  list: (
    apiClient: ApiClient,
    searchText: string,
    pageIndex: number = 0,
    recordsPerPage: number = 50
  ) =>
    queryOptions({
      queryKey: ["sample-children-list", searchText, pageIndex, recordsPerPage],
      queryFn: async () => {
        const response = await apiClient.get<
          Result<PaginationResult<SampleChildrenDto>>
        >(
          `${apiPath}${toQueryString({ searchText, pageIndex, recordsPerPage })}`
        )
        return response.result
      },
    }),
}

export const sampleChildrenMutations = {
  create: (apiClient: ApiClient) => ({
    mutationFn: async (dataToServer: SampleChildrenRequestDto) => {
      return apiClient.post(apiPath, { ...dataToServer })
    },
  }),

  update: (apiClient: ApiClient, id: string) => ({
    mutationFn: async (dataToServer: SampleChildrenRequestDto) => {
      return apiClient.put(`${apiPath}/${id}`, { ...dataToServer })
    },
  }),

  delete: (apiClient: ApiClient, id: string) => ({
    mutationFn: async () => {
      return apiClient.delete(`${apiPath}/${id}`)
    },
  }),
}
