import ResourceDetailView from "@/components/ResourceDetailView"
import { useTranslation } from "react-i18next"
import type { SamplesDto } from "./types"

interface props {
  data: SamplesDto
}

const SamplesView: React.FC<props> = ({ data }) => {
  const { t } = useTranslation()

  return (
    <ResourceDetailView
      audit={data}
      fields={[
        { label: t("samples.name"), value: data.name },
        { label: t("samples.description"), value: data.description },
      ]}
    />
  )
}

export default SamplesView
