import type { AuditColumn } from "@/lib/dto"
import { FilePlus, Pencil, Trash2 } from "lucide-react"
import { useTranslation } from "react-i18next"

interface props {
  auditModel: AuditColumn
}

const AuditUserName: React.FC<props> = ({ auditModel }) => {
  const { t } = useTranslation()

  // --------- Delete -----------------
  if (auditModel.deletedAt) {
    return (
      <div className="flex items-center gap-2 text-red-400">
        <Trash2 size="0.9rem" />
        {auditModel.deletedBy?.email}
      </div>
    )
  }

  // --------- Update -----------------
  if (auditModel.updatedAt) {
    return (
      <div className="flex items-center gap-2">
        <Pencil size="0.9rem" />
        {auditModel.updatedBy?.email}
      </div>
    )
  }

  // --------- Create -----------------
  if (auditModel.createdAt) {
    return (
      <div className="flex items-center gap-2 text-green-700">
        <FilePlus size="0.9rem" />
        {auditModel.createdBy?.email}
      </div>
    )
  }

  return <>{t("audit-history.no-actions")}</>
}

export default AuditUserName
