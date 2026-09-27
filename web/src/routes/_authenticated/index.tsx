import { createFileRoute, Link } from "@tanstack/react-router"

import { PageHeader } from "@/components/layout/page-header"

export const Route = createFileRoute("/_authenticated/")({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <PageHeader breadcrumbs={[{ label: "Home" }]} />
      <div className="flex flex-col gap-4 p-4 pt-0">
        <h1 className="text-xl font-semibold">Welcome to KFamily</h1>
        <p className="text-sm text-muted-foreground">
          This is a fresh start - auth (Casdoor) and one sample CRUD resource
          are wired up so the rest of the family-tree app can be built on top.
        </p>
        <Link to="/samples" className="text-sm font-medium underline">
          Go to Samples →
        </Link>
      </div>
    </>
  )
}
