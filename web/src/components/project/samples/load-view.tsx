import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { samplesQueries } from "./_queries"
import SamplesView from "./view"

interface props {
  id: string
}

const LoadAndViewSamples: React.FC<props> = ({ id }) => {
  const apiClient = useApiClient()

  return (
    <QueryBoundary query={samplesQueries.get(apiClient, id)}>
      {(data) => <SamplesView data={data} />}
    </QueryBoundary>
  )
}

export default LoadAndViewSamples
