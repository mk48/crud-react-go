import type { AuditColumn, IdEmail } from "./dto"

// Base shape of `AuditHistory<T>["data"]` as recorded by the server for any
// table (see the API's util.RecordAudit): the full row as JSON, keyed by its
// snake_case column names, with by-fields as bare user ids (no email join).
// Per-table snapshot types should extend this with their own data columns.
export interface AuditSnapshotBase {
  id?: string
  created_at?: string
  created_by?: string
  updated_at?: string
  updated_by?: string
  deleted_at?: string
  deleted_by?: string
}

// AuditUserName/AuditTime expect an AuditColumn (camelCase, by-fields as
// IdEmail). The server only records a bare user id per audit snapshot, not
// an email, so email falls back to the id here.
export function toAuditColumn(snapshot: AuditSnapshotBase): AuditColumn {
  const idEmail = (userId?: string): IdEmail | null =>
    userId ? { id: userId, email: userId } : null

  return {
    createdAt: snapshot.created_at ?? null,
    createdBy: idEmail(snapshot.created_by),
    updatedAt: snapshot.updated_at ?? null,
    updatedBy: idEmail(snapshot.updated_by),
    deletedAt: snapshot.deleted_at ?? null,
    deletedBy: idEmail(snapshot.deleted_by),
  }
}
