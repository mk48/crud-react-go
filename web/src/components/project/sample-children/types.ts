import type { AuditColumn, IdName } from "@/lib/dto"
import type { AuditSnapshotBase } from "@/lib/audit"

export interface SampleChildrenDto extends AuditColumn {
  id: string
  // The parent sample item, joined in by the API - always present, since
  // sample_item_id is a required foreign key.
  sampleItem: IdName
  name: string
}

export interface SampleChildrenRequestDto {
  sampleItemId: string
  name: string
}

// Shape of `AuditHistory<T>["data"]` for a sample child item: the
// table-specific data columns on top of the common audit fields.
export interface SampleChildrenAuditSnapshot extends AuditSnapshotBase {
  sample_item_id?: string
  name?: string
}
