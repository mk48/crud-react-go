import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { sampleChildrenQueries } from "./_queries"
import SampleChildrenView from "./view"

interface props {
  id: string
}

const LoadAndViewSampleChildren: React.FC<props> = ({ id }) => {
  const apiClient = useApiClient()

  return (
    <QueryBoundary query={sampleChildrenQueries.get(apiClient, id)}>
      {(data) => <SampleChildrenView data={data} />}
    </QueryBoundary>
  )
}

export default LoadAndViewSampleChildren
