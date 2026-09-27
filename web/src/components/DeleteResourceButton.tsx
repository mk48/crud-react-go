import ConfirmDialog from "@/components/ConfirmDialog"
import { Button } from "@/components/ui/button"
import { useMutation } from "@tanstack/react-query"
import type { UseMutationOptions } from "@tanstack/react-query"
import { Loader2, Trash2 } from "lucide-react"
import { useState } from "react"
import { useTranslation } from "react-i18next"
import { toast } from "sonner"

interface props {
  mutation: UseMutationOptions<unknown, Error, void>
  id: string
  name: string
  // Expects `<translationNamespace>.delete-confirm-title`,
  // `.deleted-successfully` and `.delete-failed` keys, each interpolating
  // `{{Name}}` - see samples' translation.json entries for the pattern.
  translationNamespace: string
  onSuccess: () => void
}

// Generic delete button for any table's detail/list row: confirm dialog +
// delete mutation + success/error toast.
const DeleteResourceButton: React.FC<props> = ({
  mutation: mutationOptions,
  id,
  name,
  translationNamespace,
  onSuccess: deletedSuccessfully,
}) => {
  const { t } = useTranslation()
  const [open, setOpen] = useState(false)

  // -------------------------- Mutations --------------------------------
  const mutation = useMutation({
    ...mutationOptions,
    onSuccess: async () => {
      toast.success(
        t(`${translationNamespace}.deleted-successfully`, { Name: name })
      )
      deletedSuccessfully()
    },
    onError: () => {
      toast.error(t(`${translationNamespace}.delete-failed`, { Name: name }))
    },
  })

  const onConfirmedToDelete = () => {
    setOpen(false)
    mutation.mutate()
  }

  const onOpenClick = () => {
    setOpen(true)
  }

  const onClosed = () => {
    setOpen(false)
  }

  return (
    <ConfirmDialog
      title={t(`${translationNamespace}.delete-confirm-title`, {
        Name: name,
      })}
      key={`${translationNamespace}-delete-${id}-${open}`} // included `open` state value so that it will re-render the confirm dialog
      open={open}
      onConfirmed={onConfirmedToDelete}
      onCancel={onClosed}
    >
      <Button
        variant="ghost"
        onClick={onOpenClick}
        disabled={mutation.isPending}
      >
        {mutation.isPending ? (
          <Loader2 className="animate-spin" />
        ) : (
          <Trash2 className="size-4" />
        )}
        {t("delete")}
      </Button>
    </ConfirmDialog>
  )
}

export default DeleteResourceButton
