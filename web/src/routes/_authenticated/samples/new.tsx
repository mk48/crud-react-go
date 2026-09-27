import { createFileRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import SamplesNewForm from "@/components/project/samples/form-new"

export const Route = createFileRoute("/_authenticated/samples/new")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("samples.page-title"), href: "/samples" },
          { label: t("create") },
        ]}
      />
      <div className="mx-auto w-96 p-4">
        <SamplesNewForm />
      </div>
    </>
  )
}
