import type { Column, RowData } from "@tanstack/react-table"
import { ArrowUpDown, MoveUp, MoveDown } from "lucide-react"
import type { AppTableFeatures } from "@/lib/table-features"

interface props<TData extends RowData, TValue> {
  children: React.ReactNode
  column: Column<AppTableFeatures, TData, TValue>
}

// Generic over TData/TValue (rather than pinned to `any`) because v9's
// Column type is invariant in TData - a Column<Features, SomeDto, string>
// is no longer structurally assignable to Column<Features, any, any>, so
// each call site's column must be able to instantiate this generically.
function ToggleSortColumnHeader<TData extends RowData, TValue>({
  column,
  children,
}: props<TData, TValue>) {
  const sortIcon = column.getIsSorted() ? (
    column.getIsSorted() === "asc" ? (
      <MoveDown className="ml-2 size-4 flex-none" />
    ) : (
      <MoveUp className="ml-2 size-4 flex-none" />
    )
  ) : (
    <ArrowUpDown className="ml-2 size-4 flex-none" />
  )

  return (
    <div
      className="flex cursor-pointer items-center p-2 break-words whitespace-normal hover:bg-gray-100"
      // variant="ghost"
      onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
    >
      {children}
      {sortIcon}
    </div>
  )
}

export default ToggleSortColumnHeader
