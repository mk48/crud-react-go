import type { AuditColumn } from "@/lib/dto"
import { useTranslation } from "react-i18next"
import DisplayTime from "./DisplayTime"

interface props {
  auditModel: AuditColumn
}

const AuditTime: React.FC<props> = ({ auditModel }) => {
  const { t } = useTranslation()

  // --------- Delete -----------------
  if (auditModel.deletedAt) {
    return <DisplayTime time={auditModel.deletedAt} />
  }

  // --------- Update -----------------
  if (auditModel.updatedAt) {
    return <DisplayTime time={auditModel.updatedAt} />
  }

  // --------- Create -----------------
  if (auditModel.createdAt) {
    return <DisplayTime time={auditModel.createdAt} />
  }

  return <>{t("audit-history.no-actions")}</>
}

export default AuditTime
