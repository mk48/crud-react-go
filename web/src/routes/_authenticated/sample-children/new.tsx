import { createFileRoute, useNavigate } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import AdminOnly from "@/components/auth/admin-only"
import SampleChildrenNewForm from "@/components/project/sample-children/form-new"

export const Route = createFileRoute("/_authenticated/sample-children/new")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const navigate = useNavigate()

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("sampleChildren.page-title"), href: "/sample-children" },
          { label: t("create") },
        ]}
      />
      <AdminOnly fallback="page">
        <div className="mx-auto w-96 p-4">
          <SampleChildrenNewForm
            onCreated={() => navigate({ to: "/sample-children" })}
          />
        </div>
      </AdminOnly>
    </>
  )
}
