import { Outlet, createRootRoute } from "@tanstack/react-router"

export const Route = createRootRoute({
  component: () => <Outlet />,
  notFoundComponent: NotFound,
})

function NotFound() {
  return (
    <main className="container mx-auto p-4 pt-16">
      <h1>Page not found</h1>
      <p>The page you're looking for doesn't exist.</p>
    </main>
  )
}
