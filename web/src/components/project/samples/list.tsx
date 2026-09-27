import ResourceList from "@/components/ResourceList"
import { Button } from "@/components/ui/button"
import { Link } from "@tanstack/react-router"
import { Plus } from "lucide-react"
import { useTranslation } from "react-i18next"
import { columns } from "./list-columns"
import type { SamplesDto } from "./types"

export default function SamplesList() {
  const { t } = useTranslation()

  return (
    <ResourceList<SamplesDto>
      apiPath="/api/v1/samples"
      columns={columns}
      searchPlaceholder={t("samples.search-by-name")}
      headerActions={
        <Button
          className="ml-2"
          render={<Link to="/samples/new" />}
          nativeButton={false}
        >
          <Plus className="mr-2 h-4 w-4" />
          {t("samples.create-new-samples")}
        </Button>
      }
    />
  )
}
