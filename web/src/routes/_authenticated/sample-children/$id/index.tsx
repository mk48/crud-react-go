import { createFileRoute, Link, useNavigate } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { Pencil } from "lucide-react"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import {
  sampleChildrenMutations,
  sampleChildrenQueries,
} from "@/components/project/sample-children/_queries"
import AuditHistoryDataLoad from "@/components/diff/history-diff-load-data"
import DeleteResourceButton from "@/components/DeleteResourceButton"
import LoadAndViewSampleChildren from "@/components/project/sample-children/load-view"
import { Button } from "@/components/ui/button"
import AdminOnly from "@/components/auth/admin-only"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute("/_authenticated/sample-children/$id/")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const apiClient = useApiClient()

  // ------------------------- Query: Get SampleChildren -----------------------------------
  const { data } = useQuery(sampleChildrenQueries.get(apiClient, id))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("sampleChildren.page-title"), href: "/sample-children" },
          { label: data?.name || "⏳" },
        ]}
      />
      <div className="p-4">
        <AdminOnly>
          <div className="mb-4 flex justify-between">
            <Button
              render={<Link to="/sample-children/$id/edit" params={{ id }} />}
              nativeButton={false}
            >
              <Pencil className="mr-2 size-4" />
              {t("edit")}
            </Button>

            <DeleteResourceButton
              mutation={sampleChildrenMutations.delete(apiClient, id)}
              id={id}
              name={data?.name || ""}
              translationNamespace="sampleChildren"
              onSuccess={() => navigate({ to: "/sample-children" })}
            />
          </div>
        </AdminOnly>

        <LoadAndViewSampleChildren id={id} />
        <AuditHistoryDataLoad id={id} />
      </div>
    </>
  )
}
