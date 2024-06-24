package cmd

import (
	"fmt"
	"go-instaloader/models"
	"go-instaloader/models/response"
	"go-instaloader/utils/rlog"
	"sync"
)

func CheckProfileURL(talent *models.Talent) (bool, error) {
	profile, err := GetProfileNode(talent.Username)
	if err != nil {
		rlog.Error(err)
		return false, err
	}

	if profile == nil {
		rlog.Error("profile not found")
		return false, err
	}

	if checkProfileUrl(profile) {
		return true, nil
	}

	return false, fmt.Errorf("%s's profile doesn't contain the URL", talent.Username)
}

func checkProfileUrl(profile *response.ProfileNodeResponse) bool {
	var isHasUrl bool
	var mu sync.Mutex
	var wg sync.WaitGroup

	if profile.Data.BioLinks != nil {
		for _, bioLink := range profile.Data.BioLinks {
			wg.Add(1)

			go func(bioLink *models.BioLink) {
				defer wg.Done()

				bioUrl := bioLink.URL
				if checkUrl(bioUrl) {
					rlog.Info(bioUrl)
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
