import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { sampleChildrenQueries } from "@/components/project/sample-children/_queries"
import SampleChildrenUpdateForm from "@/components/project/sample-children/form-update"
import AdminOnly from "@/components/auth/admin-only"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute(
  "/_authenticated/sample-children/$id/edit"
)({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()
  const apiClient = useApiClient()

  // ------------------------- Query: Get SampleChildren -----------------------------------
  const { data } = useQuery(sampleChildrenQueries.get(apiClient, id))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("sampleChildren.page-title"), href: "/sample-children" },
          { label: data?.name || "⏳", href: `/sample-children/${id}` },
          { label: t("edit") },
        ]}
      />
      <AdminOnly fallback="page">
        <div className="p-4">
          <SampleChildrenUpdateForm id={id} />
        </div>
      </AdminOnly>
    </>
  )
}
