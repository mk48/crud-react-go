import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { Link } from "@tanstack/react-router"
import { CircleX, Loader2, Link as LinkIcon } from "lucide-react"
import { useTranslation } from "react-i18next"
import { samplesQueries } from "./_queries"

interface props {
  id: string
}

const LoadAndShowSamplesName: React.FC<props> = ({ id }) => {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  return (
    <QueryBoundary
      query={samplesQueries.get(apiClient, id, true)}
      loadingFallback={<Loader2 className="size-4 animate-spin" />}
      errorFallback={
        <div className="text-red-400">{t("err-loading-data")}</div>
      }
    >
      {(data) =>
        data.deletedBy ? (
          <Link
            className="text-red-400 hover:underline"
            to="/samples/$id/audit-history"
            params={{ id: data.id }}
          >
            <CircleX className="inline-block size-4" /> {data.name}{" "}
            <LinkIcon className="inline-block size-4" />
          </Link>
        ) : (
          <Link
            className="hover:underline"
            to="/samples/$id"
            params={{ id: data.id }}
          >
            {data.name}{" "}
            <LinkIcon className="inline-block size-4 text-muted-foreground" />
          </Link>
        )
      }
    </QueryBoundary>
  )
}

export default LoadAndShowSamplesName
