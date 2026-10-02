import MutationFeedback from "@/components/MutationFeedback"
import type { SamplesFormSchema } from "@/components/project/samples/form"
import SamplesForm from "@/components/project/samples/form"
import type { SamplesRequestDto } from "@/components/project/samples/types"
import { useApiClient } from "@/hooks/use-api-client"
import { useMutation } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { samplesMutations } from "./_queries"

interface props {
  onCreated?: () => void
}

const SamplesNewForm: React.FC<props> = ({ onCreated }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  // Mutations
  const mutation = useMutation({
    ...samplesMutations.create(apiClient),
    onSuccess: () => {
      toast.success(t("create-success"))
      onCreated?.()
    },
  })

  const onSubmit = (data: SamplesFormSchema) => {
    const dataToServer: SamplesRequestDto = {
      name: data.name,
      description: data.description.trim() || null,
    }
    mutation.mutate(dataToServer)
  }

  return (
    <>
      <SamplesForm
        defaultValues={{
          name: "",
          description: "",
        }}
        submitButtonText={t("create")}
        onSubmit={onSubmit}
        isBusy={mutation.isPending}
      />
      <MutationFeedback
        isError={mutation.isError}
        error={mutation.error}
        errorMessage={t("create-failed")}
      />
    </>
  )
}

export default SamplesNewForm
