import ErrorMessage from "@/components/ErrorMessage"
import ErrorReference from "@/components/ErrorReference"
import SuccessMessage from "@/components/SuccessMessage"
import { useTranslation } from "react-i18next"

interface props {
  isSuccess: boolean
  isError: boolean
  successMessage: string
  errorMessage: string
  // The mutation's error - its trace id is shown as a reference.
  error?: unknown
}

// The success/error toast-adjacent message block repeated under every
// table's create/update form.
const MutationFeedback: React.FC<props> = ({
  isSuccess,
  isError,
  successMessage,
  errorMessage,
  error,
}) => {
  const { t } = useTranslation()

  return (
    <div className="mx-auto min-w-96">
      {isSuccess && (
        <SuccessMessage title={t("success")}>{successMessage}</SuccessMessage>
      )}
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
