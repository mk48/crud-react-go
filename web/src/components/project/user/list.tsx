import ResourceList from "@/components/ResourceList"
import { useTranslation } from "react-i18next"
import { columns } from "./list-columns"
import type { UserDto } from "./types"

export default function UserList() {
  const { t } = useTranslation()

  return (
    <ResourceList<UserDto>
      apiPath="/api/v1/users"
      columns={columns}
      searchPlaceholder={t("user.search-by-name-or-email")}
    />
  )
}
