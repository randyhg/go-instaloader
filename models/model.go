package models

const (
	StatusPending   = 0
	StatusOnProcess = 1
	StatusOk        = 2
	StatusFail      = 3
)

const DefaultStoryLimit = 3

const RedisJobQueueKey = "talent_queue"

type Talent struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Url      string `json:"url"`
	Status   int    `json:"status"`
}
