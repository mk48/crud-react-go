import AuditTime from "@/components/AuditTime"
import { TableCell, TableRow } from "@/components/ui/table"
import AuditUserNameLoad from "./audit-user-name-load"
import { toAuditColumn } from "@/lib/audit"
import type { AuditSnapshotBase } from "@/lib/audit"
import type { AuditHistory } from "@/lib/dto"
import { create } from "jsondiffpatch"
import { format } from "jsondiffpatch/formatters/html"
import "jsondiffpatch/formatters/styles/html.css"

interface props {
  // undefined for the create entry
  oldData?: AuditHistory<AuditSnapshotBase>
  newData: AuditHistory<AuditSnapshotBase>
}

const jsondiffpatch = create({})

const DiffRowPatch: React.FC<props> = ({ oldData, newData }) => {
  const delta = jsondiffpatch.diff(
    removeAuditCols(oldData?.data ?? {}),
    removeAuditCols(newData.data)
  )
  const htmlDiff = format(delta)

  return (
    <TableRow>
      <TableCell className="w-[600px]">
        <div
          className={`json-diff-container jsondiffpatch-unchanged-hidden overflow-x-auto bg-gray-100`}
        >
          <div dangerouslySetInnerHTML={{ __html: htmlDiff || "" }}></div>
        </div>
      </TableCell>
      <TableCell>
        <div className="text-sm">
          <AuditUserNameLoad auditModel={toAuditColumn(newData.data)} />
        </div>
        <AuditTime auditModel={toAuditColumn(newData.data)} />
      </TableCell>
    </TableRow>
  )
}

export default DiffRowPatch

const auditCols: ReadonlySet<string> = new Set<keyof AuditSnapshotBase>([
  "id",
  "created_at",
  "created_by",
  "updated_at",
  "updated_by",
  "deleted_at",
  "deleted_by",
])

// Only the table's own data columns are diffed - the audit columns are shown
// separately next to the diff.
const removeAuditCols = (data: AuditSnapshotBase) =>
  Object.fromEntries(Object.entries(data).filter(([k]) => !auditCols.has(k)))
