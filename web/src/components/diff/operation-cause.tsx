import {
  operationKindLabel,
  tableLabel,
} from "@/components/project/operations/labels"
import RecordLink from "@/components/project/operations/record-link"
import type { OperationSummary } from "@/lib/dto"
import { cn } from "@/lib/utils"
import { Link } from "@tanstack/react-router"
import { Workflow } from "lucide-react"
import { useTranslation } from "react-i18next"

interface props {
  operation: OperationSummary
  // The record whose history this change belongs to.
  sourceId: string
}

// Which operation (user action) caused an audit history change. A change to
// the operation's own target is direct ("via Edit sample"); anything else is
// highlighted as an indirect one, so users can tell e.g. an import or a
// side effect of an action on another record from an edit of their own.
const OperationCause: React.FC<props> = ({ operation, sourceId }) => {
  const { t } = useTranslation()

  const direct = operation.targetId === sourceId
  const prefix = direct
    ? t("audit-history.via")
    : operation.targetId
      ? t("audit-history.side-effect-of")
      : t("audit-history.part-of")

  return (
    <div
      className={cn(
        "mt-1 flex flex-wrap items-center gap-1 text-xs",
        direct ? "text-muted-foreground" : "font-medium text-amber-700"
      )}
    >
      {!direct && <Workflow className="size-3" />}
      <span>{prefix}</span>
      <Link
        to="/operations/$id"
        params={{ id: operation.id }}
        className="underline underline-offset-4"
      >
        {operationKindLabel(operation.kind)}
      </Link>
      {!direct && operation.targetTable && operation.targetId && (
        <>
          <span>·</span>
          <RecordLink tableName={operation.targetTable} id={operation.targetId}>
            {tableLabel(operation.targetTable)}
          </RecordLink>
        </>
      )}
    </div>
  )
}

export default OperationCause
