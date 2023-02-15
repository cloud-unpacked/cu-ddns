package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cu-ddns",
	Short: "A dynamic DNS client for Linode",
	Long: `The Cloud Unpacked Dynamic DNS tool (cu-ddns) sets up Linode as a DDNS 
provider. Useful for networks with a changing IP address such as a home 
network.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		logFile, err := os.OpenFile("/var/log/cu-ddns.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Error("Error opening log file. Logging to stderr instead.")
		} else {
			log.SetOutput(logFile)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// set config defaults
	viper.SetDefault("version", "0.1.0")
	viper.SetDefault("ipv4", true)
	viper.SetDefault("ipv6", true)
	viper.SetDefault("provider", "")
	viper.SetDefault("providerToken", "")
	viper.SetDefault("domainName", "")
	viper.SetDefault("aRecord", "")

	// If we're running within a Snap, change where we store the config
	if snapDir := os.Getenv("SNAP_DATA"); snapDir != "" {
		viper.SetConfigFile(snapDir + "/cu-ddns.yml")
	} else {
		viper.SetConfigFile("/etc/cu-ddns.yml")
	}

	viper.ReadInConfig()

	viper.AutomaticEnv() // read in environment variables that match
}
