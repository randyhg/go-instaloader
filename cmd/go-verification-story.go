package cmd

import (
	"errors"
	"fmt"
	"go-instaloader/models"
	"go-instaloader/models/response"
	"go-instaloader/utils/rlog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

const TheURL = "youtube.com"

func CheckStoryURL(talent *models.Talent) (bool, error) {
	stories, err := GetStoryNode(talent.Username, models.DefaultStoryLimit)
	if err != nil {
		rlog.Error(err)
		return false, err
	}

	if len(stories.Data) == 0 {
		err = fmt.Errorf("%s doesn't has a story right now", talent.Username)
		rlog.Error(err)
		return false, err
	}

	if checkStoryUrl(stories, talent) {
		return true, nil
	}

	return false, fmt.Errorf("%s's story doesn't contain the URL", talent.Username)
}

func checkStoryUrl(stories *response.StoryNodeResponse, talent *models.Talent) bool {
	var isHasUrl bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, story := range stories.Data {
		wg.Add(1)

		go func(story *models.StoryNode) {
			defer wg.Done()

			// download story screenshot
			downloadStoryImg(story, talent, i+1)

			if story.IphoneStruct.StoryLinkStickers != nil {
				for _, storyLink := range story.IphoneStruct.StoryLinkStickers {
					storyUrl := storyLink.StoryLink.URL

					if checkUrl(storyUrl) {
						rlog.Info(storyUrl)
						mu.Lock()
						isHasUrl = true
						mu.Unlock()
						return
					}
				}
			}
		}(story.Node)
	}

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
