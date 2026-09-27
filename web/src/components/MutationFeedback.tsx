import ErrorMessage from "@/components/ErrorMessage"
import SuccessMessage from "@/components/SuccessMessage"
import { useTranslation } from "react-i18next"

interface props {
  isSuccess: boolean
  isError: boolean
  successMessage: string
  errorMessage: string
}

// The success/error toast-adjacent message block repeated under every
// table's create/update form.
const MutationFeedback: React.FC<props> = ({
  isSuccess,
  isError,
  successMessage,
  errorMessage,
}) => {
  const { t } = useTranslation()

  return (
    <div className="mx-auto min-w-96">
      {isSuccess && (
        <SuccessMessage title={t("success")}>{successMessage}</SuccessMessage>
      )}
      {isError && (
        <ErrorMessage title={t("failed")}>{errorMessage}</ErrorMessage>
      )}
    </div>
  )
}

export default MutationFeedback
