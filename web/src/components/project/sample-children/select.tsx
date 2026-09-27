import ResourceCombobox from "@/components/ResourceCombobox"
import { useTranslation } from "react-i18next"
import { sampleChildrenQueries } from "./_queries"
import type { SampleChildrenDto } from "./types"

interface props {
  id: string
  name: string
  onSelect: (id: string, name: string) => void
}

const SampleChildrenSelect: React.FC<props> = ({ id, name, onSelect }) => {
  const { t } = useTranslation()

  return (
    <ResourceCombobox<SampleChildrenDto>
      id={id}
      name={name}
      onSelect={onSelect}
      queries={sampleChildrenQueries}
      getLabel={(item) => item.name}
      placeholder={t("sampleChildren.select-placeholder")}
    />
  )
}

export default SampleChildrenSelect
