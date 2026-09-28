import type { SortingState } from "@tanstack/react-table"

// URL search params of every ResourceList page (see ResourceList.tsx), so
// paging/sorting/search survive a refresh, work with Back/Forward and can be
// shared as a link. Default values are left out of the URL.
export interface ListSearch {
  page?: number // 0-based page index
  size?: number // rows per page
  q?: string // search text
  sort?: string // "<columnId>:asc|desc", or "none" for unsorted
  deleted?: boolean // include soft-deleted rows
}

export const PAGE_SIZES = [10, 20, 30, 40, 50]
export const DEFAULT_PAGE_SIZE = PAGE_SIZES[0]

/** Route `validateSearch` for list pages - drops anything malformed. */
export function validateListSearch(
  search: Record<string, unknown>
): ListSearch {
  const page = Number(search.page)
  const size = Number(search.size)
  // The router JSON-parses values, so a numeric search text arrives as a number.
  const q =
    typeof search.q === "string" || typeof search.q === "number"
      ? String(search.q)
      : ""
  const sort = typeof search.sort === "string" ? search.sort : ""

  return {
    page: Number.isInteger(page) && page > 0 ? page : undefined,
    size:
      PAGE_SIZES.includes(size) && size !== DEFAULT_PAGE_SIZE
        ? size
        : undefined,
    q: q || undefined,
    sort: sort === "none" || /^\w+:(asc|desc)$/.test(sort) ? sort : undefined,
    deleted: search.deleted === true || search.deleted === "true" || undefined,
  }
}

export function sortingFromParam(
  sort: string | undefined,
  defaultSorting: SortingState
): SortingState {
  if (!sort) return defaultSorting
  if (sort === "none") return []
  const [id, direction] = sort.split(":")
  return [{ id, desc: direction === "desc" }]
}

/** Sort param for the URL - undefined when it's the page's default. */
export function sortingToParam(
  sorting: SortingState,
  defaultSorting: SortingState
): string | undefined {
  const param = sortingToApi(sorting) || "none"
  return param === (sortingToApi(defaultSorting) || "none") ? undefined : param
}

/** `sortBy` value the API expects ("" = its default order). */
export function sortingToApi(sorting: SortingState): string {
  return sorting.length >= 1
    ? `${sorting[0].id}:${sorting[0].desc ? "desc" : "asc"}`
    : ""
}
