import { createFileRoute, Link, useNavigate } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { Pencil } from "lucide-react"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { userMutations, userQueries } from "@/components/project/user/_queries"
import AuditHistoryDataLoad from "@/components/diff/history-diff-load-data"
import DeleteResourceButton from "@/components/DeleteResourceButton"
import LoadAndViewUser from "@/components/project/user/load-view"
import { Button } from "@/components/ui/button"
import AdminOnly from "@/components/auth/admin-only"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute("/_authenticated/users/$id/")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()
  const navigate = useNavigate()
  const apiClient = useApiClient()

  // ------------------------- Query: Get User -----------------------------------
  const { data } = useQuery(userQueries.get(apiClient, id))
  // The API refuses self-deletion (it would lock you out), so don't offer it.
  const { data: me } = useQuery(userQueries.me(apiClient))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("user.page-title"), href: "/users" },
          { label: data?.name || data?.email || "⏳" },
        ]}
      />
      <div className="p-4">
        <AdminOnly>
          <div className="mb-4 flex justify-between">
            <Button
              render={<Link to="/users/$id/edit" params={{ id }} />}
              nativeButton={false}
            >
              <Pencil className="mr-2 size-4" />
              {t("edit")}
            </Button>

            {me && me.id !== id && (
              <DeleteResourceButton
                mutation={userMutations.delete(apiClient, id)}
                id={id}
                name={data?.name || data?.email || ""}
                translationNamespace="user"
                onSuccess={() => navigate({ to: "/users" })}
              />
            )}
          </div>
        </AdminOnly>

        <LoadAndViewUser id={id} />
        <AuditHistoryDataLoad id={id} />
      </div>
    </>
  )
}
