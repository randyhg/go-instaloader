package response

import "go-instaloader/models"

type StoryNodeResponse struct {
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Data    []struct {
		Node *models.StoryNode `json:"node"`
	} `json:"data"`
}

type ProfileNodeResponse struct {
	Code    int                 `json:"code"`
	Success bool                `json:"success"`
	Msg     string              `json:"msg"`
	Data    *models.ProfileNode `json:"data"`
}
