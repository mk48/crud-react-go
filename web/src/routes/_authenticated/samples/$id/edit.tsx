import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { samplesQueries } from "@/components/project/samples/_queries"
import SamplesUpdateForm from "@/components/project/samples/form-update"
import AdminOnly from "@/components/auth/admin-only"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute("/_authenticated/samples/$id/edit")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const navigate = useNavigate()
  const { id } = Route.useParams()
  const apiClient = useApiClient()

  // ------------------------- Query: Get Samples -----------------------------------
  const { data } = useQuery(samplesQueries.get(apiClient, id))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("samples.page-title"), href: "/samples" },
          { label: data?.name || "⏳", href: `/samples/${id}` },
          { label: t("edit") },
        ]}
      />
      <AdminOnly fallback="page">
        <div className="p-4">
          <SamplesUpdateForm
            id={id}
            onUpdated={() => navigate({ to: "/samples" })}
          />
        </div>
      </AdminOnly>
    </>
  )
}
