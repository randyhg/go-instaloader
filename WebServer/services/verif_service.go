package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-instaloader/WebServer/checkers"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"math/rand"
	"time"
)

var VerifService = new(verifService)

type verifService struct{}

func (v *verifService) VerifTalentService(storyLimit int, url string) error {
	if storyLimit <= 0 || url == "" {
		return errors.New("nothing to verif")
	}
	go v.startVerification(context.Background(), url, storyLimit, models.RedisJobQueueKey)
	return nil
}

func (v *verifService) startVerification(ctx context.Context, url string, storyLimit int, queueKey string) {
	for {
		q, err := fwRedis.RedisQueue().RPop(ctx, queueKey).Result()

		if errors.Is(err, redis.Nil) {
			//rlog.Error("no queue")
			rlog.Info("job finished!")
			break
		}

		if err != nil {
			rlog.Error(models.RedisJobQueueKey, "error getting queue", err)
			break
		}

		talent := v.parseTalentQueue(q)

		if talent == nil {
			rlog.Error("error parsing queue", q, err)
			i := time.Duration(randomInt())
			time.Sleep(i * time.Second)
			continue
		}

		// create new sheet service
		sheet := newSheetService()
		if sheet == nil {
			rlog.Error("sheet service is nil")
			break
		}

		isPass, resultMsg, err := v.CheckStoryAndProfile(talent, url, storyLimit)
		if err != nil {
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, err.Error())
			rlog.Info("job paused!")
			break
		}

		if !isPass {
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, resultMsg)
			i := time.Duration(randomInt())
			time.Sleep(i * time.Second)
			continue
		}

		talent.Status = models.StatusOk
		if err = TalentService.UpsertTalentData(talent); err != nil {
			rlog.Error(err)
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, "failed to store talent data to DB")
			i := time.Duration(randomInt())
			time.Sleep(i * time.Second)
			continue
		}

		remark := fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, url)
		sheet.UpdateTalentStatus(ctx, models.StatusOk, talent.Uuid, remark)

		//i := time.Duration(config.Instance.DelayWhenJobDoneInSeconds)
		i := time.Duration(randomInt())
		time.Sleep(i * time.Second)
	}
}

func (v *verifService) parseTalentQueue(s string) *models.Talent {
	var talent *models.Talent
	if err := json.Unmarshal([]byte(s), &talent); err != nil {
		return nil
	}
	return talent
}

func (v *verifService) CheckStoryAndProfile(talent *models.Talent, url string, storyLimit int) (bool, string, error) {
	var checkStoryResult, checkProfileResult bool
	var err error
	var resultMsg string

	// check story
	checkStoryResult, resultMsg, err = checkers.CheckStoryURL(talent, url, storyLimit)
	if err != nil {
		rlog.Error(fmt.Sprintf("checking %s's story node failed: %v", talent.Username, err))
		return false, "", err
	}

	if !checkStoryResult {
		return false, resultMsg, nil
	}

	// check profile
	checkProfileResult, resultMsg, err = checkers.CheckProfileURL(talent, url)
	if err != nil {
		rlog.Error(fmt.Sprintf("checking %s's profile node failed: %v", talent.Username, err))
		return false, "", err
	}

	if !checkProfileResult {
		return false, resultMsg, nil
	}

	// pass
	return true, "", nil

	//// determine the result
	//switch {
	//case isStoryHasUrl && isProfileHasUrl:
	//	return true, nil
	//case !isStoryHasUrl && isProfileHasUrl:
	//	return false, fmt.Errorf("%s's story does not contain the URL", talent.Username)
	//case isStoryHasUrl && !isProfileHasUrl:
	//	return false, fmt.Errorf("%s's profile does not contain the URL", talent.Username)
	//default:
	//	return false, fmt.Errorf("both %s's story and profile do not contain the URL", talent.Username)
	//}
}

func randomInt() int {
	rand.Seed(time.Now().UnixNano())
	arr := []int{6, 7, 8, 9, 10}
	randomIndex := rand.Intn(len(arr))
	randomValue := arr[randomIndex]
	return randomValue
}
