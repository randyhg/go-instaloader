package models

import (
	"encoding/json"
	"fmt"
	"go-instaloader/utils/rlog"
	"strings"
)

func (t *Talent) AddStoryUrls(links []string) {
	joinedStr := fmt.Sprintf("[%s]", strings.Join(links, ", "))
	t.StoryImgUrl = joinedStr
}

func (t *Talent) GetStoryUrls() string {
	if t.StoryImgUrl == "" {
		return ""
	}
	trimmed := strings.Trim(t.StoryImgUrl, "[]")
	urls := strings.Split(trimmed, ", ")
	byt, err := json.Marshal(&urls)
	if err != nil {
		rlog.Error(err)
		return ""
	}
	return string(byt)
}

func (t *Talent) AddStoryPaths(paths []string) {
	joinedStr := fmt.Sprintf("[%s]", strings.Join(paths, ", "))
	t.StoryImgPath = joinedStr
}

func (t *Talent) GetStoryPaths() string {
	if t.StoryImgPath == "" {
		return ""
	}
	trimmed := strings.Trim(t.StoryImgPath, "[]")
	paths := strings.Split(trimmed, ", ")
	byt, err := json.Marshal(&paths)
	if err != nil {
		rlog.Error(err)
		return ""
	}
	return string(byt)
}
