package replacer

import (
	"regexp"
	"strings"
)

var JS_REG *regexp.Regexp

func init() {
	var err error
	JS_REG, err = regexp.Compile(`src=\"(\S+\.js)\"`)
	if err != nil {
		panic(err)
	}
}

type ReplacerJs struct {
}

func (r *ReplacerJs) Replace(dictId string, html string) string {
	if html == "" || dictId == "" {
		return html
	}

	newHtml := html
	matchedGroup := JS_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[1]
		newStr := "/api/iframe/" + dictId + "/" + oldStr
		newHtml = strings.ReplaceAll(newHtml, oldStr, newStr)
	}

	return newHtml
}
