package cmd

import (
	"github.com/spf13/cobra"
	"go-instaloader/WebSocket/socket_app"
)

var socketCmd = &cobra.Command{
	Use:   "socket",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		socket_app.SocketStart()
	},
}

func init() {
	rootCmd.AddCommand(socketCmd)
}
