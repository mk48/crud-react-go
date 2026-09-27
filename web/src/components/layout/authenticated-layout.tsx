import { Outlet } from "@tanstack/react-router"

import { AppSidebar } from "@/components/side-bar/app-sidebar"
import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"

export function AuthenticatedLayout() {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <Outlet />
      </SidebarInset>
    </SidebarProvider>
  )
}
