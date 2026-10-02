import MutationFeedback from "@/components/MutationFeedback"
import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { useMutation } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { sampleChildrenMutations, sampleChildrenQueries } from "./_queries"
import type { SampleChildrenFormSchema } from "./form"
import SampleChildrenForm from "./form"
import type { SampleChildrenRequestDto } from "./types"

interface props {
  id: string
  onUpdated?: () => void
}

const SampleChildrenUpdateForm: React.FC<props> = ({ id, onUpdated }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  // ------------------------- Mutations  -----------------------------------
  const mutation = useMutation({
    ...sampleChildrenMutations.update(apiClient, id),
    onSuccess: () => {
      toast.success(t("update-success"))
      onUpdated?.()
    },
  })

  const onSubmit = (submitData: SampleChildrenFormSchema) => {
    const dataToServer: SampleChildrenRequestDto = {
      sampleItemId: submitData.sampleItemId,
      name: submitData.name,
    }
    mutation.mutate(dataToServer)
  }

  // ---------------------------------- Render ----------------------------------
  return (
    <QueryBoundary query={sampleChildrenQueries.get(apiClient, id)}>
      {(data) => (
        <>
          <SampleChildrenForm
            defaultValues={{
              sampleItemId: data.sampleItem.id,
              sampleItemName: data.sampleItem.name,
              name: data.name,
            }}
            submitButtonText={t("update")}
            onSubmit={onSubmit}
            isBusy={mutation.isPending}
          />
          <MutationFeedback
            isError={mutation.isError}
            error={mutation.error}
            errorMessage={t("update-failed")}
          />
        </>
      )}
    </QueryBoundary>
  )
}

export default SampleChildrenUpdateForm
