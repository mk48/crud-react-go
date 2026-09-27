import QueryBoundary from "@/components/QueryBoundary"
import { useApiClient } from "@/hooks/use-api-client"
import { userQueries } from "./_queries"
import UserView from "./view"

interface props {
  id: string
}

const LoadAndViewUser: React.FC<props> = ({ id }) => {
  const apiClient = useApiClient()

  return (
    <QueryBoundary query={userQueries.get(apiClient, id)}>
      {(data) => <UserView data={data} />}
    </QueryBoundary>
  )
}

export default LoadAndViewUser
