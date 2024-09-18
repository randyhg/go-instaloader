package request

type FetchRequest struct {
	StoryLimit int `json:"story_limit"` // count of latest story to crawl
}
