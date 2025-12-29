import { ActionPanel, Action, Keyboard, Detail } from "@raycast/api";
import { useEffect, useState } from "react";

export default function Actions(props: { word: string }) {

    const [markdown, setMarkdown] = useState<string>("");

    useEffect(() => {
        async function fetchDetail() {
            try {
                const response = await fetch(`https://dict.eloxt.cn/api/lookup/${props.word}?dict=f80a82cc5a241b775d3f5e41416beb697f70ad5e`);
                const data = await response.json() as { markdown: string };
                const markdown = data.markdown;
                setMarkdown(markdown);
            } catch (error) {
                console.error(error);
            }
        }
        fetchDetail();
    }, [props.word]);

    return (
        <ActionPanel>
            <ActionPanel.Section>
                <Action.Push title="Show Details" target={<Detail markdown={markdown} />} />
            </ActionPanel.Section>
        </ActionPanel>
    );
}