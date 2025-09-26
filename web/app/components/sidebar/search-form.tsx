import { X } from "lucide-react"

import {
    SidebarGroup,
    SidebarGroupContent,
    SidebarInput,
} from "@/components/ui/sidebar"
import { useEffect, useState } from "react"
import { useDictStore } from "~/hooks/use-dict-store"

export function SearchForm() {
    const [searchTerm, setSearchTerm] = useState("")
    const [searchResults, setSearchResults] = useState<string[]>([])
    const [selectedWord, setSelectedWord] = useState<string | null>(null)

    useEffect(() => {
        if (!useDictStore.getState().dictionary) {
            return
        }
        if (!searchTerm) {
            setSearchResults([])
            return
        }
        const fetchSearchResults = async () => {
            const response = await fetch(`/api/suggest/${searchTerm}?dict=${useDictStore.getState().dictionary?.id}`)
            const data = await response.json()
            setSearchResults(data)
        }

        if (searchTerm) {
            fetchSearchResults()
        } else {
            setSearchResults([])
        }
    }, [searchTerm])

    return (
        <SidebarGroup>
            <SidebarGroupContent className="relative">
                <SidebarInput
                    className={`${selectedWord ? 'pr-10' : ''}`}
                    type="text"
                    value={searchTerm}
                    id="search"
                    placeholder="Search..."
                    onChange={(event) => {
                        setSearchTerm(event.target.value)
                    }}
                />
                {selectedWord && (
                    <div
                        className="absolute right-0 h-8 w-8 cursor-pointer grid place-items-center top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
                        onClick={() => {
                            setSearchTerm("")
                            setSearchResults([])
                            setSelectedWord(null)
                            useDictStore.getState().setWord("")
                        }}
                    >
                        <X className="size-3" />
                    </div>
                )}
            </SidebarGroupContent>

            {searchResults.length > 0 && (
                <SidebarGroupContent className="mt-2 overflow-y-auto">
                    {searchResults.map((result, index) => (
                        <div key={index} className={`p-2 ${selectedWord === result ? "bg-accent-foreground text-accent" : "hover:bg-accent"}  rounded-md cursor-pointer`} onClick={() => {
                            setSelectedWord(result)
                            useDictStore.getState().setWord(result)
                        }}>
                            {result}
                        </div>
                    ))}
                </SidebarGroupContent>
            )}
        </SidebarGroup>
    )
}
