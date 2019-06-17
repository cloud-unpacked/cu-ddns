package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops cu-ddns from running automatically",
	Run: func(cmd *cobra.Command, args []string) {

		// If we're running within a Snap, delete "ready" file
		if snapDir := os.Getenv("SNAP_DATA"); snapDir != "" {
			err := os.Remove(snapDir + "/cu-ddns.active")
			if err != nil {
				fmt.Print("Failed to stop within snap.")
				return
			}
		} else {
			cronPath := "/etc/cron.hourly/cu-ddns"
			err := os.Remove(cronPath)
			if err != nil {
				fmt.Println("cu-ddns isn't set to run automatically.")
			}
		}

		fmt.Println("cu-ddns will no longer run automatically.")

	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
