package request

type FetchRequest struct {
	FetchRange string `json:"fetch_range"` // count of latest story to crawl
}
