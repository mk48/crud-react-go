import ResourceDetailView from "@/components/ResourceDetailView"
import LoadAndShowSamplesName from "@/components/project/samples/load-name"
import { useTranslation } from "react-i18next"
import type { SampleChildrenDto } from "./types"

interface props {
  data: SampleChildrenDto
}

const SampleChildrenView: React.FC<props> = ({ data }) => {
  const { t } = useTranslation()

  return (
    <ResourceDetailView
      audit={data}
      fields={[
        {
          label: t("sampleChildren.sample-item"),
          // load-name (rather than data.sampleItem.name) so a since-deleted
          // parent is flagged and links to its audit history instead.
          value: <LoadAndShowSamplesName id={data.sampleItem.id} />,
        },
        { label: t("sampleChildren.name"), value: data.name },
      ]}
    />
  )
}

export default SampleChildrenView
