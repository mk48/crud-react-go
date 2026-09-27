import MutationFeedback from "@/components/MutationFeedback"
import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"
import { userMutations, userQueries } from "./_queries"
import type { UserFormSchema } from "./form"
import UserForm from "./form"
import type { UserRequestDto } from "./types"

interface props {
  id: string
  onUpdated?: () => void
}

const UserUpdateForm: React.FC<props> = ({ id, onUpdated }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()
  const queryClient = useQueryClient()

  // ------------------------- Mutations  -----------------------------------
  const mutation = useMutation({
    ...userMutations.update(apiClient, id),
    onSuccess: () => {
      toast.success(t("update-success"))
      queryClient.invalidateQueries({ queryKey: ["users", id] })
      // Editing your own record changes the name shown in the sidebar.
      queryClient.invalidateQueries({ queryKey: ["users-me"] })
      onUpdated?.()
    },
  })

  const onSubmit = (submitData: UserFormSchema) => {
    const dataToServer: UserRequestDto = {
      name: submitData.name || null,
      isAdmin: submitData.isAdmin,
    }
    mutation.mutate(dataToServer)
  }

  // ---------------------------------- Render ----------------------------------
  return (
    <QueryBoundary query={userQueries.get(apiClient, id)}>
      {(data) => (
        <>
          <UserForm
            defaultValues={{
              name: data.name || "",
              isAdmin: data.isAdmin,
            }}
            submitButtonText={t("update")}
            onSubmit={onSubmit}
            isBusy={mutation.isPending}
          />
          <MutationFeedback
            isSuccess={mutation.isSuccess}
            isError={mutation.isError}
            successMessage={t("update-success")}
            errorMessage={t("update-failed")}
          />
        </>
      )}
    </QueryBoundary>
  )
}

export default UserUpdateForm
