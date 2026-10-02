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

// The app an operation came from (operation.client): "web", "mobile",
// "admin", "system" or "batch:<job>". Keys are looked up literally - the ":"
// of a batch name isn't an i18next namespace separator here.
export const clientLabel = (client: string) => {
  const [kind, job] = client.split(/:(.*)/s)
  if (kind === "batch" && job) {
    return i18next.t("operation.client.batch", { job })
  }
  return i18next.t(`operation.client.${client}`, {
    defaultValue: client,
    nsSeparator: false,
  })
}
