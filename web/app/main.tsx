import { useEffect, useState } from "react";
import { useDictStore } from "./hooks/use-dict-store";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { DropdownMenu, DropdownMenuCheckboxItem, DropdownMenuContent, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { ChevronsUpDown } from "lucide-react";

export interface Dictionary {
    id: string
    name: string
}

interface lookupResult {
    html: string;
}

export default function Main() {
    const selectedWord = useDictStore((state) => state.word);
    const [dictList, setDictlist] = useState<Dictionary[]>([])
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(false)
    const [selectedDict, setSelectedDict] = useState<Dictionary | null>(null)
    const [wordData, setWordData] = useState<lookupResult | null>(null)

    // Fetch dictionary list on component mount
    useEffect(() => {
        const fetchDicts = async () => {
            try {
                const response = await fetch('/api/dictionaries')
                if (!response.ok) {
                    throw new Error('Network response was not ok')
                }
                const data = await response.json()
                setDictlist(data)
                if (data.length > 0) {
                    setSelectedDict(data[0])
                    useDictStore.getState().setDictionary(data[0])
                }
            } catch (error) {
                setError(true)
            } finally {
                setLoading(false)
            }
        }

        fetchDicts()
    }, [])

    // Log selected word changes
    useEffect(() => {
        const fetchHtml = async () => {
            try {
                const response = await fetch(`/api/lookup/${selectedWord}?dict=${useDictStore.getState().dictionary?.id}`);
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                const data: lookupResult = await response.json();
                setWordData(data);
                const iframe = document.getElementById('definition-iframe') as HTMLIFrameElement | null;
                if (iframe && data) {
                    iframe.srcdoc = data.html;
                }
            } catch (error) {
            } finally {
                setLoading(false);
            }
        };

        if (selectedDict && selectedWord) {
            fetchHtml();
        }
    }, [selectedWord]);

    useEffect(() => {
        window.onmessage = (event) => {
            if (event.data.evtype === '_INNER_FRAME_MSG_EVTP_ENTRY_JUMP') {
                let word = event.data.word;
                let dict_id = event.data.dict_id;
                useDictStore.getState().setWord(word);
                if (dict_id) {
                    const dict = dictList.find(d => d.id === dict_id);
                    if (dict) {
                        setSelectedDict(dict);
                        useDictStore.getState().setDictionary(dict);
                    }
                }
            }
        };
    }, [dictList]);

    return (
        <>
            <header className="flex h-16 shrink-0 items-center gap-2">
                <div className="flex items-center gap-2 px-4">
                    <SidebarTrigger className="-ml-1" />
                    <Separator
                        orientation="vertical"
                        className="mr-2 data-[orientation=vertical]:h-4"
                    />
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <div>
                                {loading ? <span className="text-muted-foreground text-sm">Loading...</span> :
                                    <span className="p-2 rounded-lg text-foreground text-sm cursor-pointer hover:bg-accent">
                                        {selectedDict?.name}
                                        <ChevronsUpDown className="inline-block size-4 ml-2 mb-0.5" />
                                    </span>}
                                {error && <span className="text-red-500">Failed to load dictionary</span>}
                            </div>
                        </DropdownMenuTrigger>
                        {!loading && !error && dictList.length > 0 && (
                            <DropdownMenuContent>
                                {dictList.map((dict) => (
                                    <DropdownMenuCheckboxItem
                                        key={dict.name}
                                        checked={dict === selectedDict}
                                        onCheckedChange={(checked) => {
                                            if (checked) {
                                                setSelectedDict(dict)
                                                useDictStore.getState().setDictionary(dict)
                                            }
                                        }}
                                    >
                                        {dict.name}
                                    </DropdownMenuCheckboxItem>
                                ))}
                            </DropdownMenuContent>
                        )}
                    </DropdownMenu>
                </div>
            </header>

            <div className="w-full">
                <iframe id="definition-iframe"
                    className={`px-4 ${selectedWord ? "" : "hidden"}`} style={{ width: '100%', height: 'calc(100vh - 6rem)', border: 'none' }}
                ></iframe>
                {!selectedWord && (
                    <div id="iframe-container" className="pt-4 flex h-full flex-col items-center justify-center text-center text-muted-foreground">
                        <div>
                            <h1 className="mb-4 text-2xl font-bold">Select a word</h1>
                            <p className="max-w-md px-4">
                                Please select a word from the search results to see its definition and details.
                            </p>
                        </div>
                    </div>
                )}
            </div>
        </>
    );
}
