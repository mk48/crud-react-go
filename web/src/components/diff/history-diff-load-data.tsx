import HistoryDiff from "@/components/diff/history-diff"
import ErrorMessage from "@/components/ErrorMessage"
import PageLoadingIcon from "@/components/PageLoadingIcon"
import { useApiClient } from "@/hooks/use-api-client"
import type { AuditSnapshotBase } from "@/lib/audit"
import type { AuditHistory, Result } from "@/lib/dto"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

interface props {
  id: string
}

const AuditHistoryDataLoad: React.FC<props> = ({ id }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  // ------------------------- Query: get audit history -----------------------------------
  const { data, isLoading, isError } = useQuery({
    queryKey: ["audit-history", id],
    queryFn: async () => {
      const response = await apiClient.get<
        Result<AuditHistory<AuditSnapshotBase>[]>
      >(`/api/v1/audit-history/${id}`)
      return response.result
    },
  })

  // ---------------------------------- Render ----------------------------------
  if (isLoading) {
    return <PageLoadingIcon />
  }

  if (isError) {
    return <ErrorMessage title="Error!">{t("err-loading-data")}</ErrorMessage>
  }

  return (
    <div className="mt-8 p-4">
      <h1 className="text-2xl">{t("audit-history.audit-history-records")}</h1>
      <hr />
      <HistoryDiff data={data!} />
    </div>
  )
}

export default AuditHistoryDataLoad
