import { Table, TableBody } from "@/components/ui/table"
import type { AuditSnapshotBase } from "@/lib/audit"
import type { AuditHistory } from "@/lib/dto"
import DiffRowPatch from "./diff-row-patch"

interface props {
  data: AuditHistory<AuditSnapshotBase>[]
}

const HistoryDiff: React.FC<props> = ({ data }) => {
  return (
    <Table>
      <TableBody>
        {/* Newest first; the oldest (create) entry is diffed against nothing,
            so it shows every initial value. */}
        {data.map((item, index) => (
          <DiffRowPatch
            key={item.id}
            newData={item}
            oldData={data[index + 1]}
          />
        ))}
      </TableBody>
    </Table>
  )
}

export default HistoryDiff
