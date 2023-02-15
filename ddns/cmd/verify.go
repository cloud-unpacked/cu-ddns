package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/linode/linodego"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Confirms the installation status of cu-ddns",
	Long: `Verify confirms that everything is working with cu-ddns.
	It checks the config file, API token, and service status.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Checking your installation...\n\n")

		_, err := os.Stat("/etc/cu-ddns.yml")
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("1. Configuration file not found.")
			} else {
				fmt.Println("Something is wrong with the config file.")
			}
		} else {
			fmt.Println("1. Configuration file found.")
		}

		// Basic, almost useless validation
		// For the record, Linode APIv4 tokens:
		//    64chars
		//    /[a-f0-9]{64}/
		apiToken := viper.GetString("providertoken")
		if len(apiToken) == 64 {
			fmt.Println("2. Provider token conforms to spec.")
		} else {
			fmt.Println("2. Provider token does not conform to spec.")
			return
		}

		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})

		oauth2Client := &http.Client{
			Transport: &oauth2.Transport{
				Source: tokenSource,
			},
		}

		linodeClient := linodego.NewClient(oauth2Client)

		_, err = linodeClient.ListTokens(context.Background(), nil)

		// We just want to make sure we can make this call. The results don't matter.
		_, err = linodeClient.ListDomains(context.Background(), nil)
		if err != nil {
			fmt.Println("3. The provider token is invalid.")
		} else {
			fmt.Println("3. The provider token is valid.")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
