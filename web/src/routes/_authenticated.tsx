import { createFileRoute } from "@tanstack/react-router"
import { RequireAuth } from "@/components/auth/require-auth"
import { AuthenticatedLayout } from "@/components/layout/authenticated-layout"

export const Route = createFileRoute("/_authenticated")({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <RequireAuth>
      <AuthenticatedLayout />
    </RequireAuth>
  )
}
