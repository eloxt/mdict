package replacer

import (
	"fmt"
	"regexp"
	"strings"
)

var ENTRY_REG *regexp.Regexp

func init() {
	var err error
	ENTRY_REG, err = regexp.Compile(`href=\"entry://([\w#_ -]+)\"`)
	if err != nil {
		panic(err)
	}
}

type ReplacerEntry struct {
}

func (r *ReplacerEntry) Replace(dictId string, html string) string {

	if html == "" || dictId == "" {
		return html
	}

	newHtml := html
	matchedGroup := ENTRY_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[0]
		oldWord := strings.TrimRight(matched[0], "\"")
		oldWord = strings.TrimPrefix(oldWord, "href=\"entry://")

		newStr := fmt.Sprintf("href=\"javascript:_entry_jump('%s', '%s');\"", oldWord, dictId)
		fmt.Printf("old %s => new %s\n", oldStr, newStr)

		newHtml = strings.ReplaceAll(newHtml, oldStr, newStr)
	}

	return newHtml
}
