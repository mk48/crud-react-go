import i18next from "i18next"

// Operation kinds are dotted "<resource>.<action>" keys (e.g.
// "sample.create"), which nest naturally under operation.kind.* in
// translation.json. Unknown kinds fall back to the raw key, so a new
// backend kind still shows up before it's translated.
export const operationKindLabel = (kind: string) =>
  i18next.t(`operation.kind.${kind}`, { defaultValue: kind })

// audit_history/operation table names, e.g. "sample_items" -> "Sample".
export const tableLabel = (tableName: string) =>
  i18next.t(`operation.table.${tableName}`, { defaultValue: tableName })
