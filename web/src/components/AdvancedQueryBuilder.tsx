import { formatQuery, QueryBuilder } from "react-querybuilder"
import type {
  Field,
  RuleGroupArray,
  RuleGroupType,
  RuleType,
  ValueEditorType,
} from "react-querybuilder"
import { useTranslation } from "react-i18next"
import { useState } from "react"
import { Button } from "./ui/button"
import "react-querybuilder/dist/query-builder.css"
import { Card, CardContent } from "@/components/ui/card"
import { ArrowBigLeftDash, Loader2 } from "lucide-react"
import { Separator } from "./ui/separator"
import { useQuery } from "@tanstack/react-query"
import PageLoadingIcon from "./PageLoadingIcon"
import ErrorMessage from "./ErrorMessage"
import { useApiClient } from "@/hooks/use-api-client"
import type { ColumnMetaDataResponseModel, Result } from "@/lib/dto"

interface props {
  columnMetaDataUrl: string
  isBusy: boolean
  defaultRules?: RuleGroupArray<RuleGroupType<RuleType, string>, RuleType>
  onSearch: (whereCondition: string, whereConditionParameters: string) => void
}

// Keyed by Postgres's information_schema.columns `data_type` values (e.g.
// "boolean"), which is what the backend's meta endpoint reports - not C#/JSON
// type names.
const postgresDataTypeToValueEditor: { [key: string]: ValueEditorType } = {
  boolean: "checkbox",
}

const AdvancedQueryBuilder: React.FC<props> = ({
  columnMetaDataUrl,
  isBusy,
  defaultRules,
  onSearch,
}) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()
  const [dirty, setDirty] = useState(false)
  const [query, setQuery] = useState<RuleGroupType>({
    combinator: "and",
    rules: defaultRules ?? [{ field: "id", operator: "notNull", value: "" }],
  })

  // ------------------------- Query: get column meta data -----------------------------------
  const { data, isLoading, isError } = useQuery({
    queryKey: ["column-meta-data", columnMetaDataUrl],
    queryFn: async () => {
      const response = await apiClient.get<
        Result<ColumnMetaDataResponseModel[]>
      >(columnMetaDataUrl)
      return response.result
    },
    select: (columns) => {
      const fields: Field[] = columns.map((col) => ({
        name: col.name,
        label: col.name,
        valueEditorType: postgresDataTypeToValueEditor[col.dataType],
      }))
      return fields
    },
  })

  const onQueryChanged = (
    q: RuleGroupType<RuleType<string, string, unknown, string>, string>
  ) => {
    setQuery(q)
    setDirty(true)
  }

  const onSearchClicked = () => {
    // paramPrefix ":p" (not the postgresql preset's default "$") so the
    // fragment binds the same way as this app's other named-parameter SQL
    // (see the API's util.AdvancedQuery use of sqlx's :name/BindNamed) - the backend
    // maps each ":pN" back to whereConditionParameters[N-1] by position.
    const q = formatQuery(query, {
      format: "parameterized",
      preset: "postgresql",
      paramPrefix: ":p",
    })

    const where = q.sql
    const paramString = JSON.stringify(q.params)

    setDirty(false)
    onSearch(where, paramString)
  }

  // ---------------------------------- Render ----------------------------------
  if (isLoading) {
    return <PageLoadingIcon />
  }

  if (isError) {
    return <ErrorMessage title="Error!">{t("err-loading-data")}</ErrorMessage>
  }

  return (
    <Card>
      <CardContent>
        <div className="mt-4">
          <QueryBuilder
            fields={data}
            query={query}
            onQueryChange={onQueryChanged}
          />
        </div>

        <div className="mt-4 flex h-5 items-center space-x-4 text-sm">
          <div>
            <Button onClick={onSearchClicked} disabled={isBusy}>
              {isBusy && <Loader2 className="animate-spin" />}
              {t("search")}
            </Button>
          </div>
          <Separator orientation="vertical" />

          {dirty && (
            <div className="flex animate-bounce">
              <ArrowBigLeftDash />
              <div className="text-sm text-muted-foreground">
                {t("advancedQuery.clik-search")}
              </div>
            </div>
          )}
          <div className="text-sm text-muted-foreground">
            {formatQuery(query, "natural_language")}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}

export default AdvancedQueryBuilder
