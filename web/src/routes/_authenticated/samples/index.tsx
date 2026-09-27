import { createFileRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import SamplesList from "@/components/project/samples/list"

export const Route = createFileRoute("/_authenticated/samples/")({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()

  return (
    <>
      <PageHeader breadcrumbs={[{ label: t("samples.page-title") }]} />
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <SamplesList />
      </div>
    </>
  )
}
