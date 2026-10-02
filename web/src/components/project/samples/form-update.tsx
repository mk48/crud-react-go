import MutationFeedback from "@/components/MutationFeedback"
import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { useMutation } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { samplesMutations, samplesQueries } from "./_queries"
import type { SamplesFormSchema } from "./form"
import SamplesForm from "./form"
import type { SamplesRequestDto } from "./types"

interface props {
  id: string
  onUpdated?: () => void
}

const SamplesUpdateForm: React.FC<props> = ({ id, onUpdated }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  // ------------------------- Mutations  -----------------------------------
  const mutation = useMutation({
    ...samplesMutations.update(apiClient, id),
    onSuccess: () => {
      toast.success(t("update-success"))
      onUpdated?.()
    },
  })

  const onSubmit = (submitData: SamplesFormSchema) => {
    const dataToServer: SamplesRequestDto = {
      name: submitData.name,
      description: submitData.description.trim() || null,
    }
    mutation.mutate(dataToServer)
  }

  // ---------------------------------- Render ----------------------------------
  return (
    <QueryBoundary query={samplesQueries.get(apiClient, id)}>
      {(data) => (
        <>
          <SamplesForm
            defaultValues={{
              name: data.name,
              description: data.description ?? "",
            }}
            submitButtonText={t("update")}
            onSubmit={onSubmit}
            isBusy={mutation.isPending}
          />
          <MutationFeedback
            isSuccess={mutation.isSuccess}
            isError={mutation.isError}
            error={mutation.error}
            successMessage={t("update-success")}
            errorMessage={t("update-failed")}
          />
        </>
      )}
    </QueryBoundary>
  )
}

export default SamplesUpdateForm
