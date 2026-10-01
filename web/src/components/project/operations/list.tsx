import ResourceList from "@/components/ResourceList"
import { useTranslation } from "react-i18next"
import { operationsApiPath } from "./_queries"
import { columns } from "./list-columns"
import type { OperationDto } from "./types"

export default function OperationList() {
  const { t } = useTranslation()

  return (
    <ResourceList<OperationDto>
      apiPath={operationsApiPath}
      columns={columns}
      searchPlaceholder={t("operation.search-placeholder")}
      showIncludeDeleted={false}
    />
  )
}
