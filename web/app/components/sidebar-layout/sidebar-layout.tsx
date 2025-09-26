import { Outlet } from "react-router";
import AppSidebar from "../sidebar/app-sidebar";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "@/components/ui/sidebar";

export default function SidebarLayout() {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <main>
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
