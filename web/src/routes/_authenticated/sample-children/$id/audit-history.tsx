import { createFileRoute } from "@tanstack/react-router"
import { useTranslation } from "react-i18next"

import { PageHeader } from "@/components/layout/page-header"
import AuditHistoryDataLoad from "@/components/diff/history-diff-load-data"

// This page is to show only the Audit history, mainly for deleted records.
// The audit history is already shown in the View page, but the view page
// won't work for deleted records, so this is a dedicated page for it.
export const Route = createFileRoute(
  "/_authenticated/sample-children/$id/audit-history"
)({
  component: RouteComponent,
})

function RouteComponent() {
  const { t } = useTranslation()
  const { id } = Route.useParams()

  return (
    <>
      <PageHeader
        breadcrumbs={[
          { label: t("sampleChildren.page-title"), href: "/sample-children" },
          { label: `(${t("deleted")})` },
        ]}
      />
      <div className="p-4">
        <AuditHistoryDataLoad id={id} />
      </div>
    </>
  )
}
