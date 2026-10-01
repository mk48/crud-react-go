import AdvancedQueryBuilder from "@/components/AdvancedQueryBuilder"
import ErrorMessage from "@/components/ErrorMessage"
import PageLoadingIcon from "@/components/PageLoadingIcon"
import Pagination from "@/components/Pagination"
import { Button } from "@/components/ui/button"
import { Checkbox } from "@/components/ui/checkbox"
import { Input } from "@/components/ui/input"
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useApiClient } from "@/hooks/use-api-client"
import { toQueryString } from "@/lib/api-client"
import type { PaginationResult, Result } from "@/lib/dto"
import { tableFeaturesConfig } from "@/lib/table-features"
import type { AppTableFeatures } from "@/lib/table-features"
import { useQuery } from "@tanstack/react-query"
import type { ColumnDef, RowData, SortingState } from "@tanstack/react-table"
import { functionalUpdate, useTable } from "@tanstack/react-table"
import { useNavigate, useSearch } from "@tanstack/react-router"
import {
  DEFAULT_PAGE_SIZE,
  sortingFromParam,
  sortingToApi,
  sortingToParam,
  validateListSearch,
  type ListSearch,
} from "@/lib/list-search"
import { Loader, Search, X } from "lucide-react"
import type { ReactNode } from "react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { useDebounceCallback } from "usehooks-ts"

export interface ResourceListProps<TDto extends RowData> {
  // Every table's backend exposes GET {apiPath}, GET {apiPath}/query and
  // GET {apiPath}/meta with the same shapes (see sample's init.go/handler.go
  // for the pattern this relies on) - this is also used as the react-query
  // cache key prefix.
  apiPath: string
  // TValue differs per column, so it can't be a single concrete type here.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  columns: ColumnDef<AppTableFeatures, TDto, any>[]
  defaultSorting?: SortingState
  searchPlaceholder: string
  // Rendered top-right of the search bar - e.g. a "create new" link/button.
  // Table-specific because it's a route the generic list can't know about.
  headerActions?: ReactNode
  // Off for tables without soft delete (e.g. operations).
  showIncludeDeleted?: boolean
}

const DEFAULT_SORTING: SortingState = [{ id: "createdAt", desc: true }]

// Generic search/advanced-query/sort/paginate list screen for any table that
// follows the samples list/query/meta backend convention. Table-specific
// pieces (columns, row actions, search placeholder) come in as props; a
// per-table `list.tsx` just supplies those and renders this.
export default function ResourceList<TDto extends RowData>({
  apiPath,
  columns,
  defaultSorting = DEFAULT_SORTING,
  searchPlaceholder,
  headerActions,
  showIncludeDeleted = true,
}: ResourceListProps<TDto>) {
  const { t } = useTranslation()
  const apiClient = useApiClient()
  const navigate = useNavigate()

  // Paging/sorting/search live in the URL (see lib/list-search.ts) - every
  // list route declares validateSearch: validateListSearch.
  const search = useSearch({ strict: false }) as ListSearch
  const sorting = sortingFromParam(search.sort, defaultSorting)
  const pagination = {
    pageIndex: search.page ?? 0,
    pageSize: search.size ?? DEFAULT_PAGE_SIZE,
  }
  const debouncedSearchValue = search.q ?? ""
  const includeDeletedRecords = search.deleted ?? false

  // `replace` for keystroke-driven changes, so Back isn't flooded with them.
  const setSearch = (patch: ListSearch, replace = false) =>
    void navigate({
      to: ".",
      search: (prev: Record<string, unknown>) =>
        validateListSearch({ ...prev, ...patch }),
      replace,
    })

  // The search box updates the URL (debounced) - but Back/Forward can also
  // change `q`, so resync the box when it does.
  const [searchValue, setSearchValue] = useState(debouncedSearchValue)
  const [syncedSearchValue, setSyncedSearchValue] =
    useState(debouncedSearchValue)
  if (syncedSearchValue !== debouncedSearchValue) {
    setSyncedSearchValue(debouncedSearchValue)
    setSearchValue(debouncedSearchValue)
  }

  // advanced query
  const [queryType, setQueryType] = useState("search") // search or advancedQuery
  const [whereCondition, setWhereCondition] = useState("")
  const [whereConditionParameters, setWhereConditionParameters] = useState("")

  // ------------------------- Query: search -----------------------------------
  const { data, isLoading, isError } = useQuery({
    queryKey: [
      apiPath,
      "list",
      sorting,
      pagination.pageIndex,
      pagination.pageSize,
      debouncedSearchValue,
      includeDeletedRecords,
    ],
    queryFn: async () => {
      const response = await apiClient.get<Result<PaginationResult<TDto>>>(
        `${apiPath}${toQueryString({
          pageIndex: pagination.pageIndex,
          recordsPerPage: pagination.pageSize,
          searchText: debouncedSearchValue,
          sortBy: sortingToApi(sorting),
          includeDeleted: includeDeletedRecords,
        })}`
      )
      return response.result
    },
    enabled: queryType === "search",
  })

  // ------------------------- Query: advanced -----------------------------------
  const {
    data: advancedQueryData,
    isLoading: advancedQueryIsLoading,
    isError: advancedQueryIsError,
  } = useQuery({
    queryKey: [
      apiPath,
      "query",
      pagination.pageIndex,
      pagination.pageSize,
      whereCondition,
      whereConditionParameters,
      sorting,
    ],
    queryFn: async () => {
      const response = await apiClient.get<Result<PaginationResult<TDto>>>(
        `${apiPath}/query${toQueryString({
          pageIndex: pagination.pageIndex,
          recordsPerPage: pagination.pageSize,
          whereCondition: whereCondition,
          whereConditionParametersJson: whereConditionParameters,
          sortBy: sortingToApi(sorting),
        })}`
      )
      return response.result
    },
    enabled: queryType === "advancedQuery" && whereCondition !== "",
  })

  const table = useTable({
    features: tableFeaturesConfig,
    data:
      (queryType == "search" ? data?.items : advancedQueryData?.items) || [],
    columns,
    manualPagination: true, // turn off client-side pagination
    rowCount:
      queryType == "search"
        ? data?.pagination.totalResults
        : advancedQueryData?.pagination.totalResults,

    manualSorting: true,
    onSortingChange: (updater) =>
      setSearch({
        sort: sortingToParam(
          functionalUpdate(updater, sorting),
          defaultSorting
        ),
        page: 0,
      }),

    // Client-side filtering was never registered as a feature - search is
    // handled entirely server-side via the queries above, not react-table's
    // column filtering - so there's no `manualFiltering` option in v9 to set.

    onPaginationChange: (updater) => {
      const next = functionalUpdate(updater, pagination)
      setSearch({ page: next.pageIndex, size: next.pageSize })
    },
    state: {
      sorting,
      pagination,
    },
  })

  const debounceSearch = useDebounceCallback((str: string) => {
    setSearch({ q: str, page: 0 }, true)
  }, 500)

  const searchInputHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
    const enteredText = event.target.value
    setSearchValue(enteredText)
    debounceSearch(enteredText)
  }

  const clearSearchText = () => {
    debounceSearch.cancel()
    setSearchValue("")
    setSearch({ q: undefined, page: 0 })
  }

  const searchIncludeDeletedChange = (include: boolean) => {
    setSearch({ deleted: include, page: 0 })
  }

  // Search and advanced query results are paged separately - start each
  // from the first page instead of carrying the other's page over.
  const onQueryTypeChange = (value: string) => {
    setQueryType(value)
    setSearch({ page: 0 })
  }

  const onAdvancedQuery = (where: string, parameters: string) => {
    setWhereCondition(where)
    setWhereConditionParameters(parameters)
    setSearch({ page: 0 })
  }

  // ----------------- Render
  return (
    <div>
      <Tabs
        defaultValue="search"
        value={queryType}
        onValueChange={onQueryTypeChange}
        className="mb-4"
      >
        <TabsList>
          <TabsTrigger value="search">{t("search")}</TabsTrigger>
          <TabsTrigger value="advancedQuery">
            {t("advancedQuery.advanced-query")}
          </TabsTrigger>
        </TabsList>

        {/* ------------------------- Search -------------------------*/}
        <TabsContent keepMounted value="search" className="data-hidden:hidden">
          <div className="flex items-center justify-between">
            <div className="flex flex-1 items-center space-x-2">
              {/* ------------------------- Search box -------------------------*/}
              <div className="relative w-full">
                {isLoading ? (
                  <Loader className="absolute top-2.5 left-2 size-4 animate-spin" />
                ) : (
                  <Search className="absolute top-2.5 left-2 h-4 w-4 text-muted-foreground" />
                )}
                <Input
                  placeholder={searchPlaceholder}
                  value={searchValue}
                  onChange={searchInputHandler}
                  className="w-full pl-8"
                />
                <Button
                  className="absolute top-2.5 right-2 h-4 w-4 text-muted-foreground"
                  variant="ghost"
                  onClick={clearSearchText}
                >
                  <X />
                </Button>
              </div>

              {/* ------------------------- Search by: Include deleted records -------------------------*/}
              {showIncludeDeleted && (
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="includeDeleted"
                    checked={includeDeletedRecords}
                    onCheckedChange={searchIncludeDeletedChange}
                  />
                  <label
                    htmlFor="includeDeleted"
                    className="text-sm leading-none font-medium peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    {t("include-deleted-records")}
                  </label>
                </div>
              )}
            </div>
            {headerActions}
          </div>
          {isError && (
            <ErrorMessage title="Error!" className="mt-4 mb-4">
              {t("err-loading-data")}
            </ErrorMessage>
          )}
        </TabsContent>
        <TabsContent value="advancedQuery">
          <AdvancedQueryBuilder
            isBusy={advancedQueryIsLoading}
            columnMetaDataUrl={`${apiPath}/meta`}
            onSearch={onAdvancedQuery}
          />
          {advancedQueryIsError && (
            <ErrorMessage title="Error!" className="mt-4 mb-4">
              {t("err-loading-data")}
            </ErrorMessage>
          )}
        </TabsContent>
      </Tabs>

      <div className="rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id}>
                    {header.isPlaceholder ? null : (
                      <table.FlexRender header={header} />
                    )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getAllCells().map((cell) => (
                    <TableCell key={cell.id}>
                      <table.FlexRender cell={cell} />
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center"
                >
                  {isLoading ? <PageLoadingIcon /> : t("no-data")}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <Pagination table={table} />
    </div>
  )
}
