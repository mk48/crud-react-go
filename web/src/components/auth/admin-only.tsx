import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import ErrorMessage from "@/components/ErrorMessage"
import { useIsAdmin } from "@/hooks/use-current-user"

interface props {
  children: ReactNode
  // Rendered for non-admins - nothing by default (hide a button/icon). Pass
  // `page` to show an explanation instead, for whole admin-only pages.
  fallback?: ReactNode | "page"
}

// Renders children only for admins - see useIsAdmin.
export default function AdminOnly({ children, fallback = null }: props) {
  const { t } = useTranslation()
  const isAdmin = useIsAdmin()

  if (isAdmin) {
    return <>{children}</>
  }

  if (fallback === "page") {
    return (
      <div className="p-4">
        <ErrorMessage title={t("admin-only.title")}>
          {t("admin-only.description")}
        </ErrorMessage>
      </div>
    )
  }

  return <>{fallback}</>
}
