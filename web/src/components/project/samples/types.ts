import type { AuditColumn } from "@/lib/dto"
import type { AuditSnapshotBase } from "@/lib/audit"

export interface SamplesDto extends AuditColumn {
  id: string
  name: string
  description: string | null
}

export interface SamplesRequestDto {
  name: string
  description: string | null
}

// Shape of `AuditHistory<T>["data"]` for a sample item: the sample-specific
// data columns on top of the common audit fields - not the SamplesDto shape.
export interface SamplesAuditSnapshot extends AuditSnapshotBase {
  name?: string
  description?: string | null
}
