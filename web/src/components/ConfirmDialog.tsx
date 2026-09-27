import { useState } from "react"
import { useTranslation } from "react-i18next"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { Button } from "./ui/button"

interface props {
  open: boolean
  title: string
  children: React.ReactNode
  onConfirmed: () => void
  onCancel: () => void
}

const ConfirmDialog: React.FC<props> = ({
  open,
  title,
  children,
  onCancel,
  onConfirmed,
}) => {
  const { t } = useTranslation()
  const [isOpen, setIsOpen] = useState(open)

  const onConfirmClicked = () => {
    close()
    onConfirmed()
  }

  const onCancelClicked = () => {
    close()
    onCancel()
  }

  const close = () => {
    setIsOpen(false)
  }

  // -------------------------------------------- Render -----------------------------------
  if (!isOpen) {
    return <>{children}</>
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>
            {t("confirm-dialog.description")}
          </DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button onClick={onCancelClicked} variant="outline">
            {t("cancel")}
          </Button>
          <Button onClick={onConfirmClicked}>{t("continue")}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}

export default ConfirmDialog
