package checkers

import (
	"context"
	"encoding/json"
	"fmt"
	"go-instaloader/instaloader"
	"go-instaloader/models"
	"go-instaloader/models/response"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"sync"
)

func CheckProfileURL(talent *models.Talent, url string) (bool, string, error) {
	profile, err := instaloader.GetProfileNode(talent.Username)
	if err != nil {
		byt, _ := json.Marshal(talent)
		fwRedis.RedisQueue().RPush(context.Background(), models.RedisJobQueueKey+"_err", string(byt))
		return false, "", err
	}

	if profile == nil {
		rlog.Error("profile not found")
		return false, fmt.Sprintf("%s's profile not found", talent.Username), nil
	}

	if checkProfileUrl(profile, url) {
		return true, "", nil
	}

	return false, fmt.Sprintf("%s's profile doesn't contain the URL", talent.Username), nil
}

func checkProfileUrl(profile *response.ProfileNodeResponse, url string) bool {
	var isHasUrl bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	if profile.Data.BioLinks != nil {
		for _, bioLink := range profile.Data.BioLinks {
			wg.Add(1)

			go func(bioLink *models.BioLink) {
				defer wg.Done()

				bioUrl := bioLink.URL
				if CheckUrl(url, bioUrl) {
					mu.Lock()
					isHasUrl = true
					mu.Unlock()
					return
				}
			}(&bioLink)
		}
		wg.Wait()
	} else {
		rlog.Error("bio not found")
		return false
	}

	return isHasUrl
}
