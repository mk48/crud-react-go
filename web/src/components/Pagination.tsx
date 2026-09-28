import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Button } from "@/components/ui/button"
import {
  ChevronLeft,
  ChevronRightIcon,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react"
import type { RowData, Table } from "@tanstack/react-table"
import { useTranslation } from "react-i18next"
import type { AppTableFeatures } from "@/lib/table-features"
import { PAGE_SIZES } from "@/lib/list-search"

interface props<TData extends RowData> {
  table: Table<AppTableFeatures, TData>
}

// Generic over TData (rather than pinned to `any`) because v9's Table type
// is invariant in TData - a Table<Features, SomeDto> is no longer
// structurally assignable to Table<Features, any>, so each call site's
// table must be able to instantiate this generically.
function Pagination<TData extends RowData>({ table }: props<TData>) {
  const { t } = useTranslation()

  return (
    <div className="mt-2 flex items-center justify-between px-2">
      <div className="flex-1 text-sm text-muted-foreground">
        {t("pagination.total-result")} {table.getRowCount()}
      </div>
      <div className="flex items-center space-x-6 lg:space-x-8">
        <div className="flex items-center space-x-2">
          <p className="text-sm font-medium">
            {t("pagination.result-per-page")}
          </p>
          <Select
            value={table.store.state.pagination.pageSize.toString()}
            onValueChange={(value: string | null) =>
              table.setPageSize(Number(value))
            }
          >
            <SelectTrigger className="h-8 w-[70px]">
              <SelectValue
                placeholder={table.store.state.pagination.pageSize}
              />
            </SelectTrigger>
            <SelectContent side="top">
              {PAGE_SIZES.map((pageSize) => (
                <SelectItem key={pageSize} value={pageSize.toString()}>
                  {pageSize}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        <div className="flex w-[100px] items-center justify-center text-sm font-medium">
          {t("pagination.page-n-of-total", {
            currentPage: table.store.state.pagination.pageIndex + 1,
            totalPage: table.getPageCount(),
          })}
        </div>
        <div className="flex items-center space-x-2">
          {/* -------------- << --------------- */}
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => table.firstPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to first page</span>
            <ChevronsLeft className="size-4" />
          </Button>

          {/* -------------- < --------------- */}
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            <span className="sr-only">Go to previous page</span>
            <ChevronLeft className="size-4" />
          </Button>

          {/* -------------- > --------------- */}
          <Button
            variant="outline"
            className="h-8 w-8 p-0"
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to next page</span>
            <ChevronRightIcon className="size-4" />
          </Button>

          {/* -------------- >> --------------- */}
          <Button
            variant="outline"
            className="hidden h-8 w-8 p-0 lg:flex"
            onClick={() => table.lastPage()}
            disabled={!table.getCanNextPage()}
          >
            <span className="sr-only">Go to last page</span>
            <ChevronsRight className="size-4" />
          </Button>
        </div>
      </div>
    </div>
  )
}

export default Pagination
