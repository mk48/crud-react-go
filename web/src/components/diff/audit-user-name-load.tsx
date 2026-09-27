import AuditUserName from "@/components/AuditUserName"
import QueryBoundary from "@/components/QueryBoundary"
import { userQueries } from "@/components/project/user/_queries"
import { useApiClient } from "@/hooks/use-api-client"
import type { AuditColumn } from "@/lib/dto"
import { Loader2 } from "lucide-react"

interface props {
  auditModel: AuditColumn
}

// toAuditColumn (lib/audit.ts) falls back to the bare user id for email,
// since raw audit_history snapshots only record an id, not an email join.
// This resolves that id to the real email before handing off to
// AuditUserName.
const AuditUserNameLoad: React.FC<props> = ({ auditModel }) => {
  const apiClient = useApiClient()

  const activeField = auditModel.deletedBy
    ? "deletedBy"
    : auditModel.updatedBy
      ? "updatedBy"
      : auditModel.createdBy
        ? "createdBy"
        : null

  if (!activeField) {
    return <AuditUserName auditModel={auditModel} />
  }

  const userId = auditModel[activeField]!.id

  return (
    <QueryBoundary
      query={userQueries.get(apiClient, userId)}
      loadingFallback={<Loader2 className="size-4 animate-spin" />}
      errorFallback={<AuditUserName auditModel={auditModel} />}
    >
      {(user) => (
        <AuditUserName
          auditModel={{
            ...auditModel,
            [activeField]: { id: user.id, email: user.email },
          }}
        />
      )}
    </QueryBoundary>
  )
}

export default AuditUserNameLoad
