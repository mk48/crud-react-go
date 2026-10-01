import { Link } from "@tanstack/react-router"
import type { ReactNode } from "react"

interface props {
  tableName: string
  id: string
  // Deleted records have no regular view page - link to their audit
  // history instead.
  deleted?: boolean
  children: ReactNode
}

// Links an audit_history/operation (tableName, id) pair to that record's
// page. One case per table, since TanStack Router's typed `to` needs route
// literals; tables without a page render children unlinked.
const RecordLink: React.FC<props> = ({ tableName, id, deleted, children }) => {
  const className = "text-primary underline-offset-4 hover:underline"

  switch (tableName) {
    case "sample_items":
      return deleted ? (
        <Link
          to="/samples/$id/audit-history"
          params={{ id }}
          className={className}
        >
          {children}
        </Link>
      ) : (
        <Link to="/samples/$id" params={{ id }} className={className}>
          {children}
        </Link>
      )
    case "sample_child_items":
      return deleted ? (
        <Link
          to="/sample-children/$id/audit-history"
          params={{ id }}
          className={className}
        >
          {children}
        </Link>
      ) : (
        <Link to="/sample-children/$id" params={{ id }} className={className}>
          {children}
        </Link>
      )
    case "user":
      return deleted ? (
        <Link
          to="/users/$id/audit-history"
          params={{ id }}
          className={className}
        >
          {children}
        </Link>
      ) : (
        <Link to="/users/$id" params={{ id }} className={className}>
          {children}
        </Link>
      )
    default:
      return <>{children}</>
  }
}

export default RecordLink
