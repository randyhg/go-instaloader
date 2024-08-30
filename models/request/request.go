package request

type FetchRequest struct {
	FetchRange string `json:"fetch_range"` // count of latest story to crawl
}

type VerifyRequest struct {
	StoryLimit int    `json:"story_limit"`
	Url        string `json:"url"`
}
