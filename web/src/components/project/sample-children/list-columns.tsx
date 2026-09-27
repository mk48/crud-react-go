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
import type { SampleChildrenDto } from "./types"

const columnHelper = createColumnHelper<AppTableFeatures, SampleChildrenDto>()

export const columns = [
  // ----------------- Column: Name ---------------------------
  columnHelper.accessor("name", {
    header: ({ column }) => {
      return (
        <ToggleSortColumnHeader column={column}>
          {i18next.t("sampleChildren.name")}
        </ToggleSortColumnHeader>
      )
    },
    cell: (info) => info.getValue(),
  }),
  // ----------------- Column: Sample item (parent) ---------------------------
  // Not sortable: the sort keys the API accepts are this table's own columns
  // (see util.BuildSortableColumns), and sorting by sample_item_id would
  // order by UUID, not by the name shown here.
  columnHelper.accessor("sampleItem", {
    header: () => i18next.t("sampleChildren.sample-item"),
    cell: (info) => (
      <Link
        className="hover:underline"
        to="/samples/$id"
        params={{ id: info.getValue().id }}
      >
        {info.getValue().name}
      </Link>
    ),
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
              <Link
                to={"/sample-children/$id/edit"}
                params={{ id: row.original.id }}
              >
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
                  ? "/sample-children/$id"
                  : "/sample-children/$id/audit-history"
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
