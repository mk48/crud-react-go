import type { AuditColumn, IdEmail } from "./dto"

// Base shape of `AuditHistory<T>["data"]` as recorded by the server for any
// table (see the API's util.RecordAudit): snake_case keys - the literal Go
// map keys, not JSON-tagged struct fields - and by-fields are bare user ids
// (no email join). Only the fields relevant to that one action are present -
// e.g. a delete snapshot has no other columns at all. Per-table snapshot
// types should extend this with their own optional data columns.
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
