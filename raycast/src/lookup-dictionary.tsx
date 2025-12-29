import { List } from "@raycast/api";
import { useEffect, useState } from "react";
import Actions from "./actions";

export default function Command() {
  const [searchText, setSearchText] = useState("");
  const [result, setResult] = useState<string[]>([]);

  useEffect(() => {
    async function fetchList() {
      try {
        const response = await fetch("https://dict.eloxt.cn/api/suggest/" + searchText + "?dict=f80a82cc5a241b775d3f5e41416beb697f70ad5e");
        const result = await response.json() as string[];
        setResult(result);
      } catch (error) {
        console.error(error);
      }
    }
    fetchList();
  }, [searchText]);

  return (
    <List
      onSearchTextChange={setSearchText}
    >
      {result.map((item) => (
        <List.Item
        key={item}
          title={item}
          actions={
            <Actions word={item} />
          }
        />
      ))}
    </List>
  );
}
