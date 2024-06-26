package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-instaloader/WebServer/checkers"
	"go-instaloader/config"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
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

		isPass, err := v.CheckStoryAndProfile(talent, url, storyLimit, isRetry)
		if err != nil || !isPass {
			SheetService.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, err.Error())
			continue
		}

		remark := fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, url)
		SheetService.UpdateTalentStatus(ctx, models.StatusOk, talent.Uuid, remark)

		// todo: store to DB

		// Use goroutine for concurrent processing
		//go func(talent *models.Talent) {
		//	defer func() {
		//		if r := recover(); r != nil {
		//			rlog.Error("recovered from panic in goroutine:", r)
		//		}
		//	}()
		//
		//	isPass, err := v.CheckStoryAndProfile(talent, url, storyLimit)
		//	if err != nil || !isPass {
		//		SheetService.UpdateTalentStatus(ctx, models.StatusFail, talent.Uuid, err.Error())
		//		return
		//	}
		//
		//	remark := fmt.Sprintf("both of %s's story and profile contain %s url", talent.Username, url)
		//	SheetService.UpdateTalentStatus(ctx, models.StatusOk, talent.Uuid, remark)
		//
		//	// todo: store to DB
		//
		//}(talent)

		// Sleep outside of the goroutine to avoid delaying fetching the next job
		i := time.Duration(config.Instance.DelayWhenJobDoneInSeconds)
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
