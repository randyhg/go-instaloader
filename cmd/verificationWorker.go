/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"go-instaloader/config"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
	"golang.org/x/net/context"
	"os"
	"regexp"
	"time"
)

// verificationWorkerCmd represents the verificationWorker command
var verificationWorkerCmd = &cobra.Command{
	Use:   "verificationWorker",
	Short: "Verification Worker",
	Long:  "Data Verification Process",
	Run:   verificationStart,
}

func init() {
	rootCmd.AddCommand(verificationWorkerCmd)
}

func verificationStart(cmd *cobra.Command, args []string) {
	client, err := GetHttpClient(config.Instance.CredentialPath)
	if err != nil {
		rlog.Fatal("Unable to get http client:", err.Error())
	}
	ctx := context.Background()

	fmt.Println("================================ Verification Process Started ================================")
	for {
		q, err := fwRedis.RedisQueue().RPop(ctx, models.RedisJobQueueKey).Result()

		if errors.Is(err, redis.Nil) {
			rlog.Debug(models.RedisJobQueueKey, ":no queue")
			i := time.Duration(config.Instance.DelayWhenNoJobInSeconds)
			time.Sleep(i * time.Second)
			continue
		}

		if err != nil {
			rlog.Error(models.RedisJobQueueKey, "Error getting queue", err)
			i := time.Duration(config.Instance.DelayWhenErrorInSeconds)
			time.Sleep(i * time.Second)
			continue
		}

		talent := parseTalentQueue(q)
		if talent == nil {
			rlog.Error("Error parsing queue", err)
			i := time.Duration(config.Instance.DelayWhenErrorInSeconds)
			time.Sleep(i * time.Second)
			continue
		}

		isPass, err := CheckStoryAndProfile(talent)
		if err != nil || !isPass {
			UpdateTalentStatus(client, ctx, models.StatusFail, talent.Uuid, err.Error())
			i := time.Duration(config.Instance.DelayWhenErrorInSeconds)
			time.Sleep(i * time.Second)
			continue
		}

		UpdateTalentStatus(client, ctx, models.StatusOk, talent.Uuid, "talent pass")
		i := time.Duration(config.Instance.DelayWhenJobDoneInSeconds)
		time.Sleep(i * time.Second)
	}
}

func parseTalentQueue(s string) *models.Talent {
	var talent *models.Talent
	if err := json.Unmarshal([]byte(s), &talent); err != nil {
		return nil
	}
	return talent
}

func CheckStoryAndProfile(talent *models.Talent) (bool, error) {
	var isStoryHasUrl, isProfileHasUrl bool
	var err error

	// check story
	isStoryHasUrl, err = CheckStoryURL(talent)
	if err != nil {
		rlog.Error(fmt.Sprintf("checking %s story node failed: %v", talent.Username, err))
		return false, err
	}

	// check profile
	isProfileHasUrl, err = CheckProfileURL(talent)
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

func testSaveToFile(stories []*models.StoryNode, talent *models.Talent) {
	jsonData, err := json.MarshalIndent(stories, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s.json", talent.Username))
	if err != nil {
		rlog.Error(err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		rlog.Error(err)
	}
}

func checkUrl(url string) bool {
	re := regexp.MustCompile(TheURL)
	if re.MatchString(url) {
		return true
	}
	return false
}
