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
import { useTable } from "@tanstack/react-table"
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
}

const sortByParam = (sorting: SortingState): string =>
  sorting.length >= 1 ? `${sorting[0].id}:${sorting[0].desc ? "desc" : "asc"}` : ""

// Generic search/advanced-query/sort/paginate list screen for any table that
// follows the samples list/query/meta backend convention. Table-specific
// pieces (columns, row actions, search placeholder) come in as props; a
// per-table `list.tsx` just supplies those and renders this.
export default function ResourceList<TDto extends RowData>({
  apiPath,
  columns,
  defaultSorting = [{ id: "createdAt", desc: true }],
  searchPlaceholder,
  headerActions,
}: ResourceListProps<TDto>) {
  const { t } = useTranslation()
  const apiClient = useApiClient()
  const [sorting, setSorting] = useState<SortingState>(defaultSorting)
  const [pagination, setPagination] = useState({
    pageIndex: 0, // initial page index
    pageSize: 10, // default page size
  })

  // default search
  const [searchValue, setSearchValue] = useState("")
  const [debouncedSearchValue, setDebouncedSearchValue] = useState("")
  const [includeDeletedRecords, setIncludeDeletedFiles] = useState(false)

  // advanced query
  const [queryType, setQueryType] = useState("search") // search or advancedQuery
  const [whereCondition, setWhereCondition] = useState("")
  const [whereConditionParameters, setWhereConditionParameters] = useState("")

  // ------------------------- Query: search -----------------------------------
  const { data, isLoading, isError } = useQuery({
    queryKey: [
      apiPath,
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
          sortBy: sortByParam(sorting),
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
      "advanced-query",
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
          sortBy: sortByParam(sorting),
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
    onSortingChange: setSorting,

    // Client-side filtering was never registered as a feature - search is
    // handled entirely server-side via the queries above, not react-table's
    // column filtering - so there's no `manualFiltering` option in v9 to set.

    onPaginationChange: setPagination, // update the pagination state when internal APIs mutate the pagination state
    state: {
      sorting,
      pagination,
    },
  })

  const debounceSearch = useDebounceCallback((str: string) => {
    setDebouncedSearchValue(str)
    table.firstPage()
  }, 500)

  const searchInputHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
    const enteredText = event.target.value
    setSearchValue(enteredText)
    debounceSearch(enteredText)
  }

  const clearSearchText = () => {
    setSearchValue("")
    setDebouncedSearchValue("")
    table.firstPage()
  }

  const searchIncludeDeletedChange = (include: boolean) => {
    setIncludeDeletedFiles(include)
    table.firstPage()
  }

  const onAdvancedQuery = (where: string, parameters: string) => {
    setWhereCondition(where)
    setWhereConditionParameters(parameters)
  }

  // ----------------- Render
  return (
    <div>
      <Tabs
        defaultValue="search"
        value={queryType}
        onValueChange={setQueryType}
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
              <div className="flex items-center space-x-2">
                <Checkbox
                  id="includeDeleted"
                  defaultChecked={false}
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
