package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-instaloader/WebServer/checkers"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

var VerifService = new(verifService)

type verifService struct{}

func (v *verifService) VerifTalentService(storyLimit int) error {
	if storyLimit <= 0 {
		return errors.New("nothing to verif")
	}
	go v.startVerification(context.Background(), storyLimit, models.RedisJobQueueKey)
	return nil
}

func (v *verifService) startVerification(ctx context.Context, storyLimit int, queueKey string) {
	for {
		q, err := fwRedis.RedisStore().RPop(ctx, queueKey).Result()

		if errors.Is(err, redis.Nil) {
			rlog.Info("job finished!")
			break
		}

		if err != nil {
			rlog.Error(models.RedisJobQueueKey, "error getting queue", err)
			ErrorHandler(err)
			break
		}

		// parse talent
		talent, err := v.parseTalentQueue(q)
		if err != nil {
			rlog.Error("error parsing queue", q, err)
			ErrorHandler(err)
			i := time.Duration(randomInt())
			time.Sleep(i * time.Second)
			continue
		}

		// create new sheet service
		sheet := newSheetService()
		if sheet == nil {
			rlog.Error("sheet service is nil")
			fwRedis.RedisStore().RPush(context.Background(), models.RedisJobQueueKey, q)
			ErrorHandler(errors.New("sheet service is nil"))
			break
		}

		// check talent story and profile
		isPass, resultMsg, err := v.CheckStoryAndProfile(talent, talent.Url, storyLimit)
		if err != nil {
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Username, err.Error())
			fwRedis.RedisStore().RPush(context.Background(), models.RedisJobQueueKey, q)
			ErrorHandler(err)
			rlog.Info("job paused!")
			break
		}

		var remark string
		if isPass { // verification pass
			talent.Status = models.StatusOk
			remark = fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, talent.Url)
		} else { // verification failed
			talent.Status = models.StatusFail
			remark = resultMsg
		}

		// store to DB
		if err = TalentService.UpsertTalentData(talent); err != nil {
			rlog.Error(err)
			ErrorHandler(err)
			fwRedis.RedisStore().RPush(context.Background(), models.RedisJobQueueKey, q)
			remark = fmt.Sprint("failed to store talent data to DB")
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Username, remark)
			i := time.Duration(randomInt())
			time.Sleep(i * time.Second)
			continue
		}

		sheet.UpdateTalentStatus(ctx, talent.Status, talent.Username, remark)
		i := time.Duration(randomInt())
		time.Sleep(i * time.Second)
	}
}

func (v *verifService) parseTalentQueue(s string) (*models.Talent, error) {
	var talent *models.Talent
	if err := json.Unmarshal([]byte(s), &talent); err != nil {
		rlog.Error(err)
		return nil, err
	}
	return talent, nil
}

func (v *verifService) CheckStoryAndProfile(talent *models.Talent, url string, storyLimit int) (bool, string, error) {
	var checkStoryResult, checkProfileResult bool
	var err error
	var resultMsg string

	// check profile
	checkProfileResult, resultMsg, err = checkers.CheckProfileURL(talent, url)
	if err != nil {
		rlog.Errorf("checking %s's profile node failed: %v", talent.Username, err)
		return false, "", err
	}

	if !checkProfileResult {
		return false, resultMsg, nil
	}

	// check story
	checkStoryResult, resultMsg, err = checkers.CheckStoryURL(talent, url, storyLimit)
	if err != nil {
		rlog.Errorf("checking %s's story node failed: %v", talent.Username, err)
		return false, "", err
	}

	if !checkStoryResult {
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
