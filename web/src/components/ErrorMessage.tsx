import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert"
import { AlertCircle, RefreshCw } from "lucide-react"
import { Button } from "./ui/button"
import { useTranslation } from "react-i18next"

interface props {
  children: React.ReactNode
  title?: string
  className?: string
  retry?: () => void
}

const ErrorMessage: React.FC<props> = ({
  title,
  children,
  retry,
  className,
}) => {
  const { t } = useTranslation()

  return (
    <Alert variant="destructive" className={className}>
      <AlertCircle className="h-4 w-4" />
      <AlertTitle>{title || t("error")}</AlertTitle>
      <AlertDescription>
        {children}
        {retry && (
          <div className="flex w-full justify-end">
            <Button onClick={retry} variant="outline">
              <RefreshCw />
              {t("try-again")}
            </Button>
          </div>
        )}
      </AlertDescription>
    </Alert>
  )
}

export default ErrorMessage
