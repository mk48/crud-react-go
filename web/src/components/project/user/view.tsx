import ResourceDetailView from "@/components/ResourceDetailView"
import { useTranslation } from "react-i18next"
import type { UserDto } from "./types"

interface props {
  data: UserDto
}

const UserView: React.FC<props> = ({ data }) => {
  const { t } = useTranslation()

  return (
    <ResourceDetailView
      audit={data}
      fields={[
        { label: t("user.name"), value: data.name || "-" },
        { label: t("user.email"), value: data.email },
        { label: t("user.is-admin"), value: data.isAdmin ? t("yes") : t("no") },
      ]}
    />
  )
}

export default UserView
