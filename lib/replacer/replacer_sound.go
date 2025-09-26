package replacer

import (
	"fmt"
	"regexp"
	"strings"
)

var SOUND_REG *regexp.Regexp

func init() {
	var err error
	// <a class="fayin" href="sound://us/pond__us_1.mp3"><span class="phon-us">pɑːnd</span><img src="us_pron.png?dict_id=f234356c227f82a54afdaa3514de188a&amp;d=0" class="fayin"></a>
	SOUND_REG, err = regexp.Compile(`href=[\"|'](sound://\S+\.mp3)[\"|']`)
	if err != nil {
		panic(err)
	}
}

type ReplacerSound struct {
}

func (r *ReplacerSound) Replace(dictId string, html string) string {

	if html == "" || dictId == "" {
		return html
	}

	newHtml := html
	matchedGroup := SOUND_REG.FindAllStringSubmatch(html, -1)
	for _, matched := range matchedGroup {
		if len(matched) != 2 {
			continue
		}
		oldStr := matched[1]
		// sound://us/pond__us_1.mp3?dict_id=1236&d=0
		oldStr = strings.TrimPrefix(oldStr, "sound://")
		soundURL := "/api/resource/" + dictId + "/" + oldStr
		newStr := fmt.Sprintf("javascript:_play_sound('%s')", soundURL)
		newHtml = strings.ReplaceAll(newHtml, "sound://"+oldStr, newStr)
	}

	return newHtml
}
