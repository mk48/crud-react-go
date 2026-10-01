import type { AuditSnapshotBase } from "@/lib/audit"
import type { IdEmail } from "@/lib/dto"

// One operation - a user action that wrote data (see the API's
// util.RunOperation). Matches the API's operation.Dto.
export interface OperationDto {
  id: string
  kind: string
  performedBy: IdEmail
  // The record the user acted on directly; null for operations with no
  // single target (e.g. an import).
  targetTable: string | null
  targetId: string | null
  metadata: Record<string, unknown>
  createdAt: string
  changeCount: number
}

// One record change (audit_history row) an operation caused, with the
// record's previous version (null for a create) to diff against.
export interface OperationChangeDto {
  id: string
  tableName: string
  sourceId: string
  action: "create" | "update" | "delete"
  data: AuditSnapshotBase
  previousData: AuditSnapshotBase | null
  // Deleted since - its regular view page no longer shows it.
  sourceDeleted: boolean
  createdAt: string
}

// GET /api/v1/operations/{id} - at most 1000 changes; compare with
// changeCount.
export interface OperationDetailDto extends OperationDto {
  changes: OperationChangeDto[]
}
