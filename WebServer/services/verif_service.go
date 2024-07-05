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
	go v.startVerification(context.Background(), url, storyLimit, models.RedisJobQueueKey, false)
	return nil
}

func (v *verifService) RetryFailedVerificationService(storyLimit int, url string) error {
	if storyLimit <= 0 || url == "" {
		return errors.New("nothing to verif")
	}
	go v.startVerification(context.Background(), url, storyLimit, models.RedisJobQueueKey+"_err", true)
	return nil
}

func (v *verifService) startVerification(ctx context.Context, url string, storyLimit int, queueKey string, isRetry bool) {
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
			continue
		}

		// create new sheet service
		sheet := newSheetService()

		isPass, err := v.CheckStoryAndProfile(talent, url, storyLimit, isRetry)
		if err != nil || !isPass {
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, err.Error())
			continue
		}

		talent.Status = models.StatusOk
		if err = TalentService.UpsertTalentData(talent); err != nil {
			rlog.Error(err)
			sheet.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, "failed to store talent data to DB")
			continue
		}

		remark := fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, url)
		sheet.UpdateTalentStatus(ctx, models.StatusOk, talent.Uuid, remark)

		//i := time.Duration(config.Instance.DelayWhenJobDoneInSeconds)
		rand.Seed(time.Now().UnixNano())
		i := time.Duration(rand.Intn(10) + 1)
		rlog.Info(i, "))))))))))))))))))))")
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

func (v *verifService) CheckStoryAndProfile(talent *models.Talent, url string, storyLimit int, isRetry bool) (bool, error) {
	var isStoryHasUrl, isProfileHasUrl bool
	var err error

	// check story
	isStoryHasUrl, err = checkers.CheckStoryURL(talent, url, storyLimit, isRetry)
	if err != nil {
		rlog.Error(fmt.Sprintf("checking %s story node failed: %v", talent.Username, err))
		return false, err
	}

	// check profile
	isProfileHasUrl, err = checkers.CheckProfileURL(talent, url, isRetry)
	if err != nil {
		rlog.Error(fmt.Sprintf("checking %s profile node failed: %v", talent.Username, err))
		return false, err
	}

	// determine the result
	switch {
	case isStoryHasUrl && isProfileHasUrl:
		return true, nil
	case !isStoryHasUrl && isProfileHasUrl:
		return false, fmt.Errorf("%s's story does not contain the URL", talent.Username)
	case isStoryHasUrl && !isProfileHasUrl:
		return false, fmt.Errorf("%s's profile does not contain the URL", talent.Username)
	default:
		return false, fmt.Errorf("both %s's story and profile do not contain the URL", talent.Username)
	}
}

//for {
//q, err := fwRedis.RedisQueue().RPop(ctx, models.RedisJobQueueKey).Result()
//
//if errors.Is(err, redis.Nil) {
//rlog.Error("no queue")
//break
//}
//
//if err != nil {
//rlog.Error(models.RedisJobQueueKey, "error getting queue", err)
//break
//}
//
//talent := v.parseTalentQueue(q)
//
//if talent == nil {
//rlog.Error("error parsing queue", q, err)
//continue
//}
//
//isPass, err := v.CheckStoryAndProfile(talent, url)
//if err != nil || !isPass {
//SheetService.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, err.Error())
//continue
//}
//
//remark := fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, url)
//SheetService.UpdateTalentStatus(ctx, models.StatusOk, talent.Uuid, remark)
//
//// todo: store to DB
//
//i := time.Duration(config.Instance.DelayWhenJobDoneInSeconds)
//time.Sleep(i * time.Second)
//}
