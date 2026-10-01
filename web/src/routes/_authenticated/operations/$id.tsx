import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { operationQueries } from "@/components/project/operations/_queries"
import { operationKindLabel } from "@/components/project/operations/labels"
import OperationView from "@/components/project/operations/view"
import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute("/_authenticated/operations/$id")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()
  const apiClient = useApiClient()

  const { data } = useQuery(operationQueries.get(apiClient, id))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("operation.page-title"), href: "/operations" },
          { label: data ? operationKindLabel(data.kind) : "⏳" },
        ]}
      />
      <div className="p-4">
        <QueryBoundary query={operationQueries.get(apiClient, id)}>
          {(operation) => <OperationView data={operation} />}
        </QueryBoundary>
      </div>
    </>
  )
}
