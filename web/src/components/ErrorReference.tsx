import { ApiError } from "@/lib/api-client"
import { useTranslation } from "react-i18next"

interface props {
  error: unknown
}

// The failed request's trace id, for the user to quote when reporting it -
// it finds the request's trace, log lines and operations (see
// docs/tracing.md). Renders nothing for errors without one.
const ErrorReference: React.FC<props> = ({ error }) => {
  const { t } = useTranslation()

  if (!(error instanceof ApiError) || !error.traceId) {
    return null
  }

  return (
    <div className="mt-1 text-xs text-muted-foreground">
      {t("error-reference")}{" "}
      <span className="font-mono select-all">{error.traceId}</span>
    </div>
  )
}

export default ErrorReference
