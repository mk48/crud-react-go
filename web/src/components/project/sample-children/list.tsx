import AdminOnly from "@/components/auth/admin-only"
import ResourceList from "@/components/ResourceList"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { columns } from "./list-columns"
import type { SampleChildrenDto } from "./types"

export default function SampleChildrenList() {
  const { t } = useTranslation()

  return (
    <ResourceList<SampleChildrenDto>
      apiPath="/api/v1/sample-children"
      columns={columns}
      searchPlaceholder={t("sampleChildren.search-by-name")}
      headerActions={
        <AdminOnly>
          <Button
            className="ml-2"
            render={<Link to="/sample-children/new" />}
            nativeButton={false}
          >
            <Plus className="mr-2 h-4 w-4" />
            {t("sampleChildren.create-new-sample-children")}
          </Button>
        </AdminOnly>
      }
    />
  )
}
