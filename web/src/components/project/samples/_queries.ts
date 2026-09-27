import { queryOptions } from "@tanstack/react-query"
import type { ApiClient } from "@/hooks/use-api-client"
import { toQueryString } from "@/lib/api-client"
import type { PaginationResult, Result } from "@/lib/dto"
import type { SamplesDto, SamplesRequestDto } from "./types"

const apiPath = "/api/v1/samples"

export const samplesQueries = {
  get: (apiClient: ApiClient, id: string, includeDeleted: boolean = false) =>
    queryOptions({
      queryKey: ["samples", id, includeDeleted],
      queryFn: async () => {
        const r = await apiClient.get<Result<SamplesDto>>(
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
      queryKey: ["samples-list", searchText, pageIndex, recordsPerPage],
      queryFn: async () => {
        const response = await apiClient.get<
          Result<PaginationResult<SamplesDto>>
        >(`${apiPath}${toQueryString({ searchText, pageIndex, recordsPerPage })}`)
        return response.result
      },
    }),
}

export const samplesMutations = {
  create: (apiClient: ApiClient) => ({
    mutationFn: async (dataToServer: SamplesRequestDto) => {
      return apiClient.post(apiPath, { ...dataToServer })
    },
  }),

  update: (apiClient: ApiClient, id: string) => ({
    mutationFn: async (dataToServer: SamplesRequestDto) => {
      return apiClient.put(`${apiPath}/${id}`, { ...dataToServer })
    },
  }),

  delete: (apiClient: ApiClient, id: string) => ({
    mutationFn: async () => {
      return apiClient.delete(`${apiPath}/${id}`)
    },
  }),
}
