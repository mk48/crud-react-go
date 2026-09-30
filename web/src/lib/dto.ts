export interface Result<T> {
  isSuccess: boolean
  error: string
  message: string
  result: T
}

export interface PaginationResult<T> {
  items: T[]
  pagination: {
    pageIndex: number
    resultsPerPage: number
    totalResults: number
  }
}

export interface IdEmail {
  id: string
  email: string
}

export interface IdName {
  id: string
  name: string
}

export interface AuditColumn {
  createdAt: string | null
  createdBy: IdEmail | null

  updatedAt: string | null
  updatedBy: IdEmail | null

  deletedAt: string | null
  deletedBy: IdEmail | null
}

// One audit_history row. GET /api/v1/audit-history/{id} returns
// `Result<AuditHistory<T>[]>` - an array of these, most recent first - not a
// single wrapper object. `data` is the full source row right after the change.
export interface AuditHistory<T> {
  id: string
  tableName: string
  sourceId: string
  action: "create" | "update" | "delete"
  changedBy: string
  createdAt: string
  data: T
}

// Matches the API's dto.ColumnMeta - GET /api/v1/samples/meta (and
// equivalent per-table meta endpoints) return `Result<ColumnMetaDataResponseModel[]>`.
export interface ColumnMetaDataResponseModel {
  name: string
  dataType: string
}
