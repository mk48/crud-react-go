import { createFileRoute } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { userQueries } from "@/components/project/user/_queries"
import UserUpdateForm from "@/components/project/user/form-update"
import { useApiClient } from "@/hooks/use-api-client"

export const Route = createFileRoute("/_authenticated/users/$id/edit")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()
  const apiClient = useApiClient()

  // ------------------------- Query: Get User -----------------------------------
  const { data } = useQuery(userQueries.get(apiClient, id))

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("user.page-title"), href: "/users" },
          { label: data?.name || data?.email || "⏳", href: `/users/${id}` },
          { label: t("edit") },
        ]}
      />
      <div className="p-4">
        <UserUpdateForm id={id} />
      </div>
    </>
  )
}
