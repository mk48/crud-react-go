import { createFileRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import { validateListSearch } from "@/lib/list-search"
import OperationList from "@/components/project/operations/list"

export const Route = createFileRoute("/_authenticated/operations/")({
  component: RouteComponent,
  validateSearch: validateListSearch,
})

function RouteComponent() {
  const { t } = useTranslation()

  return (
    <>
      <PageHeader breadcrumbs={[{ label: t("operation.page-title") }]} />
      <div className="flex flex-1 flex-col gap-4 p-4 pt-0">
        <OperationList />
      </div>
    </>
  )
}
