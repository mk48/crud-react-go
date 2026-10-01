import type { AuditSnapshotBase } from "@/lib/audit"
import { create } from "jsondiffpatch"
import { format } from "jsondiffpatch/formatters/html"
import "jsondiffpatch/formatters/styles/html.css"

interface props {
  // undefined/null for a create - every initial value shows as added
  oldData?: AuditSnapshotBase | null
  newData: AuditSnapshotBase
}

const jsondiffpatch = create({})

// Diff of two audit_history snapshots of the same record, limited to its own
// data columns - the audit columns are shown separately next to it.
const SnapshotDiff: React.FC<props> = ({ oldData, newData }) => {
  const delta = jsondiffpatch.diff(
    removeAuditCols(oldData ?? {}),
    removeAuditCols(newData)
  )
  const htmlDiff = format(delta)

  return (
    <div
      className={`json-diff-container jsondiffpatch-unchanged-hidden overflow-x-auto bg-gray-100`}
    >
      <div dangerouslySetInnerHTML={{ __html: htmlDiff || "" }}></div>
    </div>
  )
}

export default SnapshotDiff

const auditCols: ReadonlySet<string> = new Set<keyof AuditSnapshotBase>([
  "id",
  "created_at",
  "created_by",
  "updated_at",
  "updated_by",
  "deleted_at",
  "deleted_by",
])

const removeAuditCols = (data: AuditSnapshotBase) =>
  Object.fromEntries(Object.entries(data).filter(([k]) => !auditCols.has(k)))
