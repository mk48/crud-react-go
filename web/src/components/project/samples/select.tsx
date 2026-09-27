import ResourceCombobox from "@/components/ResourceCombobox"
import { useTranslation } from "react-i18next"
import { samplesQueries } from "./_queries"
import type { SamplesDto } from "./types"

interface props {
  id: string
  name: string
  onSelect: (id: string, name: string) => void
}

const SamplesSelect: React.FC<props> = ({ id, name, onSelect }) => {
  const { t } = useTranslation()

  return (
    <ResourceCombobox<SamplesDto>
      id={id}
      name={name}
      onSelect={onSelect}
      queries={samplesQueries}
      getLabel={(item) => item.name}
      placeholder={t("samples.select-placeholder")}
    />
  )
}

export default SamplesSelect
