import ErrorMessage from "@/components/ErrorMessage"
import ErrorReference from "@/components/ErrorReference"
import { useTranslation } from "react-i18next"

interface props {
  isError: boolean
  errorMessage: string
  // The mutation's error - its trace id is shown as a reference.
  error?: unknown
}

// The error message block repeated under every table's create/update form.
// Success needs no inline message: the form toasts and navigates to the list.
const MutationFeedback: React.FC<props> = ({
  isError,
  errorMessage,
  error,
}) => {
  const { t } = useTranslation()

  return (
    <div className="mx-auto min-w-96">
      {isError && (
        <ErrorMessage title={t("failed")}>
          {errorMessage}
          <ErrorReference error={error} />
        </ErrorMessage>
      )}
    </div>
  )
}

export default MutationFeedback
