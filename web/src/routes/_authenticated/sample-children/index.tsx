import { createFileRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { validateListSearch } from "@/lib/list-search"
import SampleChildrenList from "@/components/project/sample-children/list"

export const Route = createFileRoute("/_authenticated/sample-children/")({
  component: RouteComponent,
  validateSearch: validateListSearch,
})

function RouteComponent() {
  const { t } = useTranslation()

  return (
    <>
      <PageHeader breadcrumbs={[{ label: t("sampleChildren.page-title") }]} />
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <SampleChildrenList />
      </div>
    </>
  )
}
