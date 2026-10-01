import DisplayTime from "@/components/DisplayTime"
import SnapshotDiff from "@/components/diff/snapshot-diff"
import { Table, TableBody, TableCell, TableRow } from "@/components/ui/table"
import { cn } from "@/lib/utils"
import type { ReactNode } from "react"
import { useTranslation } from "react-i18next"
import { operationKindLabel, tableLabel } from "./labels"
import RecordLink from "./record-link"
import type { OperationChangeDto, OperationDetailDto } from "./types"

interface props {
  data: OperationDetailDto
}

const actionClass: Record<OperationChangeDto["action"], string> = {
  create: "text-green-700",
  update: "",
  delete: "text-red-400",
}

const OperationView: React.FC<props> = ({ data }) => {
  const { t } = useTranslation()

  // The target's deleted state is only known if it's among the changes
  // (it nearly always is).
  const targetDeleted = data.changes.some(
    (c) => c.sourceId === data.targetId && c.sourceDeleted
  )
  const metadata = Object.entries(data.metadata ?? {})

  const fields: { label: string; value: ReactNode }[] = [
    { label: t("operation.operation"), value: operationKindLabel(data.kind) },
    { label: t("operation.performed-by"), value: data.performedBy.email },
    {
      label: t("operation.performed-at"),
      value: <DisplayTime time={data.createdAt} />,
    },
    {
      label: t("operation.target"),
      value:
        data.targetTable && data.targetId ? (
          <RecordLink
            tableName={data.targetTable}
            id={data.targetId}
            deleted={targetDeleted}
          >
            {tableLabel(data.targetTable)}{" "}
            <span className="font-mono text-xs">{data.targetId}</span>
          </RecordLink>
        ) : (
          "-"
        ),
    },
    ...(metadata.length > 0
      ? [
          {
            label: t("operation.details"),
            value: (
              <dl className="grid grid-cols-[auto_1fr] gap-x-4">
                {metadata.map(([key, value]) => (
                  <div key={key} className="contents">
                    <dt className="text-muted-foreground">{key}</dt>
                    <dd className="break-all">
                      {typeof value === "string"
                        ? value
                        : JSON.stringify(value)}
                    </dd>
                  </div>
                ))}
              </dl>
            ),
          },
        ]
      : []),
    { label: t("operation.changes"), value: data.changeCount },
  ]

  return (
    <>
      {/* ----------------- Summary -----------------*/}
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
            </dl>
          </div>
        </div>
      </div>

      {/* ----------------- Changes -----------------*/}
      <div className="mt-8 p-4">
        <h1 className="text-2xl">{t("operation.changes")}</h1>
        <hr />
        {data.changes.length < data.changeCount && (
          <p className="mt-2 text-sm text-muted-foreground">
            {t("operation.changes-truncated", {
              shown: data.changes.length,
              total: data.changeCount,
            })}
          </p>
        )}
        <Table>
          <TableBody>
            {data.changes.map((change) => (
              <TableRow key={change.id}>
                <TableCell className="w-[600px]">
                  <SnapshotDiff
                    oldData={change.previousData}
                    newData={change.data}
                  />
                </TableCell>
                <TableCell className="align-top text-sm">
                  <div>
                    <RecordLink
                      tableName={change.tableName}
                      id={change.sourceId}
                      deleted={change.sourceDeleted}
                    >
                      {tableLabel(change.tableName)}{" "}
                      <span className="font-mono text-xs">
                        {change.sourceId.slice(0, 8)}
                      </span>
                    </RecordLink>
                  </div>
                  <div
                    className={cn("font-medium", actionClass[change.action])}
                  >
                    {t(`operation.action.${change.action}`)}
                  </div>
                  {data.targetId && change.sourceId !== data.targetId && (
                    <div className="text-xs font-medium text-amber-700">
                      {t("operation.side-effect")}
                    </div>
                  )}
                  <DisplayTime time={change.createdAt} />
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  )
}

export default OperationView
