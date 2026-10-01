import AuditTime from "@/components/AuditTime"
import { TableCell, TableRow } from "@/components/ui/table"
import AuditUserNameLoad from "./audit-user-name-load"
import OperationCause from "./operation-cause"
import SnapshotDiff from "./snapshot-diff"
import { toAuditColumn } from "@/lib/audit"
import type { AuditSnapshotBase } from "@/lib/audit"
import type { AuditHistory } from "@/lib/dto"

interface props {
  // undefined for the create entry
  oldData?: AuditHistory<AuditSnapshotBase>
  newData: AuditHistory<AuditSnapshotBase>
}

const DiffRowPatch: React.FC<props> = ({ oldData, newData }) => {
  return (
    <TableRow>
      <TableCell className="w-[600px]">
        <SnapshotDiff oldData={oldData?.data} newData={newData.data} />
      </TableCell>
      <TableCell>
        <div className="text-sm">
          <AuditUserNameLoad auditModel={toAuditColumn(newData.data)} />
        </div>
        <AuditTime auditModel={toAuditColumn(newData.data)} />
        <OperationCause
          operation={newData.operation}
          sourceId={newData.sourceId}
        />
      </TableCell>
    </TableRow>
  )
}

export default DiffRowPatch
