import { BookA, ChevronsUpDown } from "lucide-react"
import { SidebarMenu, SidebarMenuButton, SidebarMenuItem } from "../ui/sidebar"



export function AppHeader() {
    return (
        <SidebarMenu>
            <SidebarMenuItem>
                <div className="flex items-center gap-2 px-3 py-1" >
                    <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-9 items-center justify-center rounded-lg">
                        <BookA className="size-5" />
                    </div>
                    <div className="flex flex-col gap-0.5 leading-none">
                        <span className="font-medium text-sm">Dictionary</span>
                    </div>
                </div>
            </SidebarMenuItem>
        </SidebarMenu>
    )
}
