package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Run cu-ddns hourly",
	Long:  `This command configures cu-ddns to run hourly via Cron.`,
	Run: func(cmd *cobra.Command, args []string) {

		// If we're running within a Snap, use SystemD Timer and "ready"
		// file to determine when to run
		if snapDir := os.Getenv("SNAP_DATA"); snapDir != "" {
			_, err := os.Create(snapDir + "/cu-ddns.active")
			if err != nil {
				fmt.Print("Failed to start within snap.")
				return
			}
		} else {
			// otherwise we use Cron
			execPath := "/usr/bin/cu-ddns"
			cronPath := "/etc/cron.hourly/cu-ddns"
			os.Symlink(execPath, cronPath)
		}

		fmt.Println("cu-ddns has been set to run hourly.")
		log.Info("Hourly cron enabled.")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
