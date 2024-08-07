package instaloader

import (
	"encoding/json"
	"fmt"
	"go-instaloader/models/response"
	"go-instaloader/utils/rlog"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	GetStoryNodeURL   = "http://127.0.0.1:8090/api/profile/get_story_node"
	GetProfileNodeURL = "http://127.0.0.1:8090/api/profile/get_profile_node"
)

func GetStoryNode(username string, intLimit int) (*response.StoryNodeResponse, error) {
	limit := strconv.Itoa(intLimit)

	params := url.Values{}
	params.Add("username", username)
	params.Add("limit", limit)

	fullURL := fmt.Sprintf("%s?%s", GetStoryNodeURL, params.Encode())
	rlog.Debugf("URL: %s", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		rlog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rlog.Error(err)
		return nil, err
	}

	var storyNode *response.StoryNodeResponse
	if err = json.Unmarshal(body, &storyNode); err != nil {
		rlog.Error(err)
		err = fmt.Errorf("please renew session first")
		return nil, err
	}

	return storyNode, nil

	//respData, ok := storyNode["data"]
	//var stories []*models.Node
	//if ok {
	//	dataBytes, _ := json.Marshal(respData)
	//
	//	if err = json.Unmarshal(dataBytes, &stories); err != nil {
	//		rlog.Error(err)
	//		return nil, err
	//	}
	//}
	//
	//return stories, nil
}

func GetProfileNode(username string) (*response.ProfileNodeResponse, error) {
	params := url.Values{}
	params.Add("username", username)

	fullURL := fmt.Sprintf("%s?%s", GetProfileNodeURL, params.Encode())
	rlog.Debug("URL:", fullURL)

	resp, err := http.Get(fullURL)
	if err != nil {
		rlog.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		rlog.Error(err)
		return nil, err
	}

	var profileNode *response.ProfileNodeResponse
	if err = json.Unmarshal(body, &profileNode); err != nil {
		rlog.Error(err)
		return nil, err
	}

	return profileNode, nil
}
