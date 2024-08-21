package cmd

import (
	"encoding/json"
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"go-instaloader/WebServer/caches"
	"go-instaloader/models"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"gorm.io/gorm/schema"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "A brief description of your command",
	Long:  "",
	Run:   taskStart,
}

func init() {
	rootCmd.AddCommand(taskCmd)
}

func taskStart(cmd *cobra.Command, args []string) {
	go func() {
		// running it at startup
		getTalentData()
		createMonthTable()
	}()

	c := cron.New()
	err := c.AddFunc("0 5 * * *", func() { // everyday at 05:00 AM
		createMonthTable()
	})
	if err != nil {
		rlog.Fatal("Error adding cron job:", err)
	}

	err = c.AddFunc("5 * * * *", func() { // every 5 minutes
		getTalentData()
	})
	if err != nil {
		rlog.Fatal("Error adding cron job:", err)
	}

	c.Start()
	rlog.Info("task started...")
	select {}
}

func createMonthTable() {
	createNextMothTable(models.Talent{})
}

func createNextMothTable(dst schema.Tabler) {
	tableName := myDb.GetMonthTableName(dst)
	if err := myDb.CreateMonthTable(myDb.GetDb(), dst, tableName); err != nil {
		rlog.Error("create month table", tableName, "error:", err.Error())
	}

	nextTableName := myDb.GetMonthTableNameNext(dst)
	if err := myDb.CreateMonthTable(myDb.GetDb(), dst, nextTableName); err != nil {
		rlog.Error("create next month table", nextTableName, "error:", err.Error())
	}
}

func getTalentData() {
	tableName := myDb.GetMonthTableName(models.Talent{})
	var talents []models.Talent
	err := myDb.GetDb().Table(tableName).Order("created_at DESC").Find(&talents).Error
	if err != nil {
		rlog.Error(err)
		return
	}

	byt, err := json.Marshal(&talents)
	if err != nil {
		rlog.Error(err)
		return
	}
	caches.TalentCache.SetAllTalents(tableName, string(byt))
}
