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
        {data.map((item, index) => {
          if (index === data.length - 1) return null
          return (
            <DiffRowPatch
              key={item.id}
              newData={item}
              oldData={data[index + 1]}
            />
          )
        })}
      </TableBody>
    </Table>
  )
}

export default HistoryDiff
