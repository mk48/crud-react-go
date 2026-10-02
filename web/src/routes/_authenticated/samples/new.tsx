import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import AdminOnly from "@/components/auth/admin-only"
import SamplesNewForm from "@/components/project/samples/form-new"

export const Route = createFileRoute("/_authenticated/samples/new")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const navigate = useNavigate()

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("samples.page-title"), href: "/samples" },
          { label: t("create") },
        ]}
      />
      <AdminOnly fallback="page">
        <div className="mx-auto w-96 p-4">
          <SamplesNewForm onCreated={() => navigate({ to: "/samples" })} />
        </div>
      </AdminOnly>
    </>
  )
}
