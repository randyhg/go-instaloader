package checkers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-instaloader/instaloader"
	"go-instaloader/models"
	"go-instaloader/models/response"
	"go-instaloader/utils/fwRedis"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func CheckStoryURL(talent *models.Talent, url string, storyLimit int, isRetry bool) (bool, string, error) {
	stories, err := instaloader.GetStoryNode(talent.Username, storyLimit)
	if err != nil {
		byt, _ := json.Marshal(talent)
		fwRedis.RedisQueue().LPush(context.Background(), models.RedisJobQueueKey, string(byt))
		return false, "", err
	}

	if len(stories.Data) == 0 {
		return false, fmt.Sprintf("%s doesn't has a story right now", talent.Username), nil
	}

	if checkStoryUrl(stories, talent, url) {
		return true, "", nil
	}

	return false, fmt.Sprintf("%s's story does not contain the URL", talent.Username), nil
}

func checkStoryUrl(stories *response.StoryNodeResponse, talent *models.Talent, url string) bool {
	var isHasUrl bool
	var mu sync.Mutex
	var wg sync.WaitGroup
	var storyUrls []string

	for i, story := range stories.Data {
		wg.Add(1)

		storyUrls = append(storyUrls, story.Node.DisplayURL)
		go func(story *models.StoryNode) {
			defer wg.Done()

			// download story screenshot
			downloadStoryImg(story, talent, i+1)
			if story.IphoneStruct.StoryLinkStickers != nil {
				for _, storyLink := range story.IphoneStruct.StoryLinkStickers {
					storyUrl := storyLink.StoryLink.URL

					if CheckUrl(url, storyUrl) {
						mu.Lock()
						isHasUrl = true
						mu.Unlock()
						return
					}
				}
			}
		}(story.Node)
	}
	talent.AddStoryUrls(storyUrls)

	wg.Wait()
	return isHasUrl
}

func downloadStoryImg(story *models.StoryNode, talent *models.Talent, storyNum int) {
	fileName := fmt.Sprintf("%d.png", storyNum)

	// cristiano/1.png
	outputFilePath := filepath.Join("stories", talent.Username, fileName)

	downloadFile(story.DisplayURL, outputFilePath)
}

func downloadFile(url string, outputPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status code: %v", resp.StatusCode)
	}

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create a new file
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Copy the reader to the file
	l, err := io.Copy(f, resp.Body)
	if err != nil {
		return err
	}
	if l == 0 {
		err := os.Remove(outputPath)
		if err != nil {
			return err
		}
		return errors.New("file length after copy is 0")
	}

	f.Seek(0, 0)
	return nil
}
