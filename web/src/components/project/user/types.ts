import type { AuditColumn } from "@/lib/dto"

export interface UserDto extends AuditColumn {
  id: string
  sub: string
  email: string
  name: string | null
  isAdmin: boolean
  // A service account (batch job, system) - never signs in.
  isService: boolean
}

// Users have no Create - they're provisioned automatically by the auth
// middleware on first sign-in. Only name and isAdmin can be edited.
export interface UserRequestDto {
  name: string | null
  isAdmin: boolean
}
