import * as React from "react"
import { FlaskConical, ListTree, Trees, Users } from "lucide-react"
import { Link } from "@tanstack/react-router"
import { useQuery } from "@tanstack/react-query"
import { useTranslation } from "react-i18next"

import { NavMain } from "@/components/side-bar/nav-main"
import { NavUser } from "@/components/side-bar/nav-user"
import { userQueries } from "@/components/project/user/_queries"
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"
import { useApiClient } from "@/hooks/use-api-client"

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const { t } = useTranslation()
  const apiClient = useApiClient()

  // Already loaded (and cached) by RequireAuth before this renders.
  const meQuery = useQuery(userQueries.me(apiClient))

  const user = {
    name: meQuery.data?.name ?? "",
    email: meQuery.data?.email ?? "",
  }

  const navMain = [
    {
      title: t("samples.page-title"),
      url: "/samples",
      icon: FlaskConical,
    },
    {
      title: t("sampleChildren.page-title"),
      url: "/sample-children",
      icon: ListTree,
    },
    {
      title: t("user.page-title"),
      url: "/users",
      icon: Users,
    },
  ]

  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" render={<Link to="/" />}>
              <div className="flex aspect-square size-8 shrink-0 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <Trees className="size-4" />
              </div>
              <span className="font-heading font-semibold">
                {t("app.title")}
              </span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
