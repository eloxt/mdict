package replacer

import (
	"regexp"
	"strings"
)

var imageReg *regexp.Regexp

func init() {
	var err error
	imageReg, err = regexp.Compile(`src=["|'](\S+\.(png|jpg|gif|jpeg|svg))["|']`)
	if err != nil {
		panic(err)
	}
}

type ReplacerImage struct {
}

func (r *ReplacerImage) Replace(dictId string, html string) string {
	if html == "" || dictId == "" {
		return html
	}
	cache := make(map[string]bool)

	newHtml := html
	matchedGroup := imageReg.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 3 {
			continue
		}
		imgstr := matched[1]
		if _, ok := cache[imgstr]; ok {
			continue
		} else {
			cache[imgstr] = true
			newStr := "/api/resource/" + dictId + "/" + imgstr
			newHtml = strings.ReplaceAll(newHtml, imgstr, newStr)
		}
	}

	return newHtml
}
