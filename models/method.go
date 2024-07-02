package models

import (
	"fmt"
	"strings"
)

func (t *Talent) AddStoryUrls(links []string) {
	joinedStr := fmt.Sprintf("[%s]", strings.Join(links, ", "))
	t.StoryImgUrl = joinedStr
}

func (t *Talent) GetStoryUrls() []string {
	trimmed := strings.Trim(t.StoryImgUrl, "[]")
	return strings.Split(trimmed, ", ")
}
