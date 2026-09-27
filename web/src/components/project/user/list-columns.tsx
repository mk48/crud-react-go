import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip"
import ToggleSortColumnHeader from "@/components/ToggleSortColumnHeader"
import { buildAuditTrailColumns } from "@/lib/resource-columns"
import type { AppTableFeatures } from "@/lib/table-features"
import { Link } from "@tanstack/react-router"
import { createColumnHelper } from "@tanstack/react-table"
import i18next from "i18next"
import { Pencil, View } from "lucide-react"
import type { UserDto } from "./types"

const columnHelper = createColumnHelper<AppTableFeatures, UserDto>()

export const columns = [
  // ----------------- Column: Name ---------------------------
  columnHelper.accessor("name", {
    header: ({ column }) => {
      return (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("user.name")}
        </ToggleSortColumnHeader>
      )
    },
    cell: (info) => info.getValue() || "-",
  }),
  // ----------------- Column: Email ---------------------------
  columnHelper.accessor("email", {
    header: ({ column }) => {
      return (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("user.email")}
        </ToggleSortColumnHeader>
      )
    },
    cell: (info) => info.getValue(),
  }),
  // ----------------- Column: Is Admin ---------------------------
  columnHelper.accessor("isAdmin", {
    header: ({ column }) => {
      return (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("user.is-admin")}
        </ToggleSortColumnHeader>
      )
    },
    cell: (info) => (info.getValue() ? i18next.t("yes") : i18next.t("no")),
  }),

  // ----------------- Columns: Action by / Action at ---------------------------
  ...buildAuditTrailColumns(columnHelper),

  // ----------------- Column: Action icons ---------------------------
  columnHelper.display({
    id: "actions",
    cell: ({ row }) => (
      <div className="flex justify-end gap-x-8">
        {/* Edit is not allowed for deleted rows */}
        {row.original.deletedBy == null && (
          <Tooltip>
            <TooltipTrigger>
              <Link to={"/users/$id/edit"} params={{ id: row.original.id }}>
                <Pencil className="size-4" />
              </Link>
            </TooltipTrigger>
            <TooltipContent>{i18next.t("edit")}</TooltipContent>
          </Tooltip>
        )}
        <Tooltip>
          <TooltipTrigger>
            <Link
              to={
                row.original.deletedBy == null
                  ? "/users/$id"
                  : "/users/$id/audit-history"
              }
              params={{ id: row.original.id }}
            >
              <View className="size-4" />
            </Link>
          </TooltipTrigger>
          <TooltipContent>{i18next.t("view")}</TooltipContent>
        </Tooltip>
      </div>
    ),
  }),
]
