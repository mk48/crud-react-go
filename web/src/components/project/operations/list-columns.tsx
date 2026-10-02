import DisplayTime from "@/components/DisplayTime"
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import ToggleSortColumnHeader from "@/components/ToggleSortColumnHeader"
import type { AppTableFeatures } from "@/lib/table-features"
import { Link } from "@tanstack/react-router"
import { createColumnHelper } from "@tanstack/react-table"
import i18next from "i18next"
import { View } from "lucide-react"
import ClientBadge from "./client-badge"
import { operationKindLabel, tableLabel } from "./labels"
import type { OperationDto } from "./types"

const columnHelper = createColumnHelper<AppTableFeatures, OperationDto>()

export const columns = [
  // ----------------- Column: Operation ---------------------------
  columnHelper.accessor("kind", {
    header: ({ column }) => (
      <ToggleSortColumnHeader column={column}>
        {i18next.t("operation.operation")}
      </ToggleSortColumnHeader>
    ),
    cell: (info) => (
      <Link
        to="/operations/$id"
        params={{ id: info.row.original.id }}
        className="font-medium underline-offset-4 hover:underline"
      >
        {operationKindLabel(info.getValue())}
      </Link>
    ),
  }),
  // ----------------- Column: Target ---------------------------
  // Not linked here - whether the record has been deleted since (which
  // changes where it can be viewed) is only loaded on the detail page.
  columnHelper.accessor("targetTable", {
    header: ({ column }) => (
      <ToggleSortColumnHeader column={column}>
        {i18next.t("operation.target")}
      </ToggleSortColumnHeader>
    ),
    cell: ({ row }) =>
      row.original.targetTable && row.original.targetId ? (
        <span>
          {tableLabel(row.original.targetTable)}{" "}
          <span className="font-mono text-xs text-muted-foreground">
            {row.original.targetId.slice(0, 8)}
          </span>
        </span>
      ) : (
        "-"
      ),
  }),
  // ----------------- Column: Source ---------------------------
  columnHelper.accessor("client", {
    header: ({ column }) => (
      <ToggleSortColumnHeader column={column}>
        {i18next.t("operation.source")}
      </ToggleSortColumnHeader>
    ),
    cell: (info) => <ClientBadge client={info.getValue()} />,
  }),
  // ----------------- Column: Changes ---------------------------
  columnHelper.accessor("changeCount", {
    header: () => i18next.t("operation.changes"),
    cell: (info) => info.getValue(),
  }),
  // ----------------- Column: Performed by ---------------------------
  columnHelper.display({
    id: "performedBy",
    header: () => i18next.t("operation.performed-by"),
    cell: ({ row }) => row.original.performedBy.email,
  }),
  // ----------------- Column: Performed at ---------------------------
  columnHelper.accessor("createdAt", {
    header: ({ column }) => (
      <ToggleSortColumnHeader column={column}>
        {i18next.t("operation.performed-at")}
      </ToggleSortColumnHeader>
    ),
    size: 150,
    cell: (info) => <DisplayTime time={info.getValue()} />,
  }),

  // ----------------- Column: Action icons ---------------------------
  columnHelper.display({
    id: "actions",
    cell: ({ row }) => (
      <div className="flex justify-end">
        <Tooltip>
          <TooltipTrigger>
            <Link to="/operations/$id" params={{ id: row.original.id }}>
              <View className="size-4" />
            </Link>
          </TooltipTrigger>
          <TooltipContent>{i18next.t("view")}</TooltipContent>
        </Tooltip>
      </div>
    ),
  }),
]
