import AuditTime from "@/components/AuditTime"
import AuditUserName from "@/components/AuditUserName"
import type { AuditColumn } from "@/lib/dto"
import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"

export interface ResourceDetailField {
  label: string
  value: ReactNode
}

interface props {
  fields: ResourceDetailField[]
  // Every table's Dto extends AuditColumn (see lib/dto.ts), so the table's
  // own data object can be passed straight through here.
  audit: AuditColumn
}

// Generic detail/view screen for any table: a label/value row per data
// field, plus the two audit rows every table shares.
const ResourceDetailView: React.FC<props> = ({ fields, audit }) => {
  const { t } = useTranslation()

  return (
    <div className="flex justify-start">
      <div className="overflow-hidden rounded-lg border">
        <div className="border-t border-gray-200 px-4 py-5 sm:p-0">
          <dl className="sm:divide-y sm:divide-gray-200">
            {fields.map(({ label, value }) => (
              <div
                key={label}
                className="py-1 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-2"
              >
                <dt className="text-sm font-medium text-gray-500">{label}</dt>
                <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                  {value}
                </dd>
              </div>
            ))}

            {/* ----------------- Action By -----------------*/}
            <div className="py-1 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-2">
              <dt className="text-sm font-medium text-gray-500">
                {t("audit-history.action-by")}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                <AuditUserName auditModel={audit} />
              </dd>
            </div>

            {/* ----------------- Action At -----------------*/}
            <div className="py-1 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-2">
              <dt className="text-sm font-medium text-gray-500">
                {t("audit-history.action-at")}
              </dt>
              <dd className="mt-1 text-sm text-gray-900 sm:col-span-2 sm:mt-0">
                <AuditTime auditModel={audit} />
              </dd>
            </div>
          </dl>
        </div>
      </div>
    </div>
  )
}

export default ResourceDetailView
