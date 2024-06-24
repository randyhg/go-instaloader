/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"github.com/spf13/cobra"
	"go-instaloader/config"
	"go-instaloader/models"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/rlog"
)

// fetchWorkerCmd represents the fetchWorker command
var fetchWorkerCmd = &cobra.Command{
	Use:   "fetchWorker",
	Short: "Fetch Worker",
	Long:  "Fetch Data Process",
	Run:   fetchStart,
}

func init() {
	rootCmd.AddCommand(fetchWorkerCmd)
}

func fetchStart(cmd *cobra.Command, args []string) {
	client, err := GetHttpClient(config.Instance.CredentialPath)
	if err != nil {
		rlog.Fatal("Unable to get http client:", err.Error())
	}
	ctx := context.Background()

	talents, err := GetTalents(client, ctx)
	if err != nil {
		rlog.Error("unable to get talents:", err.Error())
	}
	for _, talent := range talents {
		if len(talent.Uuid) > 0 {
			byt, err := json.Marshal(&talent)
			if err != nil {
				rlog.Error(err)
			} else {
				fwRedis.RedisQueue().LPush(ctx, models.RedisJobQueueKey, string(byt))
			}
		}
	}

	rlog.Info("fetch-worker finished")

}
