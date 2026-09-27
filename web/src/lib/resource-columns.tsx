import AuditTime from "@/components/AuditTime"
import AuditUserName from "@/components/AuditUserName"
import ToggleSortColumnHeader from "@/components/ToggleSortColumnHeader"
import type { AuditColumn } from "@/lib/dto"
import type { AppTableFeatures } from "@/lib/table-features"
import type { ColumnHelper } from "@tanstack/react-table"
import i18next from "i18next"

// Shared actionBy/actionAt columns for any table's list screen - every
// table's rows carry the same audit trail. The row-action (edit/view) column
// is left to each table's own list-columns.tsx since it needs typed,
// per-table route literals (TanStack Router's `Link to` is not generic-safe
// across route bases).
export function buildAuditTrailColumns<TDto extends AuditColumn>(
  columnHelper: ColumnHelper<AppTableFeatures, TDto>
) {
  return [
    columnHelper.display({
      id: "actionBy",
      header: ({ column }) => (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("audit-history.action-by")}
        </ToggleSortColumnHeader>
      ),
      cell: ({ row }) => <AuditUserName auditModel={row.original} />,
    }),

    columnHelper.display({
      id: "actionAt",
      header: ({ column }) => (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("audit-history.action-at")}
        </ToggleSortColumnHeader>
      ),
      size: 150,
      cell: ({ row }) => <AuditTime auditModel={row.original} />,
    }),
  ]
}
