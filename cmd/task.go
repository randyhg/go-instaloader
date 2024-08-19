package cmd

import (
	"github.com/robfig/cron"
	"github.com/spf13/cobra"
	"go-instaloader/models"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"gorm.io/gorm/schema"
	"log"
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
	c := cron.New()
	err := c.AddFunc("0 0 * * *", func() { // 1 day
		createMonthTable()
	})
	if err != nil {
		log.Fatal("Error adding cron job:", err)
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
