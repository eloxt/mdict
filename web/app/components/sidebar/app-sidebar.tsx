import { Sidebar, SidebarContent, SidebarFooter, SidebarGroup, SidebarHeader } from "@/components/ui/sidebar";
import { AppHeader } from "./app-header";
import { SearchForm } from "./search-form";

export default function AppSidebar() {
    return (
        <Sidebar variant="inset">
            <SidebarHeader>
                <AppHeader />
            </SidebarHeader>
            <SidebarContent>
                <SearchForm />
            </SidebarContent>
            <SidebarFooter />
        </Sidebar>
    );
}