import {
  columnSizingFeature,
  createPaginatedRowModel,
  createSortedRowModel,
  rowPaginationFeature,
  rowSortingFeature,
  tableFeatures,
} from "@tanstack/react-table"

// Shared feature set for every ResourceList table (see ResourceList.tsx).
// All of them are server-driven (manualSorting/manualPagination) and only
// ever use sorting + pagination + a fixed column `size` (see
// resource-columns.tsx's actionAt column) - no column filtering, grouping,
// pinning, interactive resizing, or row selection anywhere in either app -
// so this is the only `tableFeatures()` config needed, reused via
// `typeof tableFeaturesConfig` by every list-columns.tsx column helper and
// by Pagination/ToggleSortColumnHeader's prop types.
export const tableFeaturesConfig = tableFeatures({
  rowSortingFeature,
  rowPaginationFeature,
  columnSizingFeature,
  sortedRowModel: createSortedRowModel(),
  paginatedRowModel: createPaginatedRowModel(),
})

export type AppTableFeatures = typeof tableFeaturesConfig
