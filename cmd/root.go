/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go-instaloader/WebServer/app"
	"go-instaloader/config"
	"go-instaloader/utils/fwRedis"
	"go-instaloader/utils/myDb"
	"go-instaloader/utils/rlog"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-instaloader",
	Short: "Web Server",
	Long:  "Go Instaloader Web Server",
	Run:   start,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	config.Init()
	fwRedis.RedisInit()
	myDb.DBInit()
	rlog.Info("configs initiated successfully!!")
}

func start(cmd *cobra.Command, args []string) {
	app.IrisInit()
	app.IrisStart()
}
