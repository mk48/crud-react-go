import MutationFeedback from "@/components/MutationFeedback"
import type { SampleChildrenFormSchema } from "@/components/project/sample-children/form"
import SampleChildrenForm from "@/components/project/sample-children/form"
import type { SampleChildrenRequestDto } from "@/components/project/sample-children/types"
import { useApiClient } from "@/hooks/use-api-client"
import { useMutation } from "@tanstack/react-query"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { sampleChildrenMutations } from "./_queries"

interface props {
  onCreated?: () => void
}

const SampleChildrenNewForm: React.FC<props> = ({ onCreated }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()
  const [keyUpdate, setKeyUpdate] = useState(0)

  // Mutations
  const mutation = useMutation({
    ...sampleChildrenMutations.create(apiClient),
    onSuccess: () => {
      toast.success(t("create-success"))
      setKeyUpdate((k) => k + 1)
      onCreated?.()
    },
  })

  const onSubmit = (data: SampleChildrenFormSchema) => {
    const dataToServer: SampleChildrenRequestDto = {
      sampleItemId: data.sampleItemId,
      name: data.name,
    }
    mutation.mutate(dataToServer)
  }

  return (
    <>
      <SampleChildrenForm
        key={keyUpdate}
        defaultValues={{
          sampleItemId: "",
          sampleItemName: "",
          name: "",
        }}
        submitButtonText={t("create")}
        onSubmit={onSubmit}
        isBusy={mutation.isPending}
      />
      <MutationFeedback
        isSuccess={mutation.isSuccess}
        isError={mutation.isError}
        successMessage={t("create-success")}
        errorMessage={t("create-failed")}
      />
    </>
  )
}

export default SampleChildrenNewForm
