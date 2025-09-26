package replacer

import (
	"regexp"
	"strings"
)

type ReplacerCss struct {
}

var CSS_REG *regexp.Regexp

func init() {
	var err error
	CSS_REG, err = regexp.Compile(`href=\"(\S+\.css)\"`)
	if err != nil {
		panic(err)
	}
}

func (r *ReplacerCss) Replace(dictId string, html string) string {

	if html == "" || dictId == "" {
		return html
	}

	newHtml := html
	matchedGroup := CSS_REG.FindAllStringSubmatch(html, -1)
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
