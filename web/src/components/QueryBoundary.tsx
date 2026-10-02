import ErrorMessage from "@/components/ErrorMessage"
import ErrorReference from "@/components/ErrorReference"
import type { QueryKey, UseQueryOptions } from "@tanstack/react-query"
import { useQuery } from "@tanstack/react-query"
import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { Spinner } from "./ui/spinner"

interface props<TData, TQueryKey extends QueryKey> {
  // queryOptions() return values are keyed with a concrete (non-readonly)
  // tuple type, so TQueryKey has to be inferred per call site rather than
  // defaulted - a fixed default (e.g. QueryKey) rejects those concrete keys.
  query: UseQueryOptions<TData, Error, TData, TQueryKey>
  children: (data: TData) => ReactNode
  // Table-cell-sized usages (e.g. load-name.tsx) need a compact spinner/error
  // instead of the full-page fallbacks below.
  loadingFallback?: ReactNode
  errorFallback?: ReactNode
}

// Wraps the isLoading/isError/render-data branching that every "load a
// single record then show it" component (load-view, load-name,
// form-update, ...) repeats for every table.
export default function QueryBoundary<TData, TQueryKey extends QueryKey>({
  query,
  children,
  loadingFallback,
  errorFallback,
}: props<TData, TQueryKey>) {
  const { t } = useTranslation()
  const { data, isLoading, isError, error } = useQuery(query)

  if (isLoading) {
    return <>{loadingFallback ?? <Spinner />}</>
  }

  if (isError || data === undefined) {
    return (
      <>
        {errorFallback ?? (
          <ErrorMessage title="Error!">
            {t("err-loading-data")}
            <ErrorReference error={error} />
          </ErrorMessage>
        )}
      </>
    )
  }

  return <>{children(data)}</>
}
