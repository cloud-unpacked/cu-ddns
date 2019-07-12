package cmd

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Setup hostname and provider credentials",
	Long: `Setup cu-ddns with the information and credentials it needs to
operate.

This will configure:
  - the domain name
  - A/AAAA record
  - provider credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Let's configure cu-ddns.\n")

		provider := ""
		q1 := &survey.Select{
			Message: "Please choose a DNS provider:",
			Options: []string{"Linode", "Cloudflare"},
		}
		survey.AskOne(q1, &provider, survey.WithValidator(survey.Required))
		provider = strings.ToLower(provider)
		viper.Set("provider", provider)

		promptString := ""
		switch provider {
		case "linode":
			promptString = "Linode Personal Access Token"
			break
		case "cloudflare":
			promptString = "Cloudflare API key"
			break
		}

		apiToken := ""
		q2 := &survey.Password{
			Message: "Please paste a " + promptString + " for authentication:",
			// in the future, add validators depending on the provider's
			// specific key.
			// For the record, Linode APIv4 tokens:
			//    64chars
			//    /[a-f0-9]{64}/
		}
		survey.AskOne(q2, &apiToken, survey.WithValidator(survey.Required))
		viper.Set("providerToken", apiToken)

		if provider == "cloudflare" {

			providerEmail := ""
			q3 := &survey.Input{
				Message: "Enter your Cloudflare email address:",
			}
			survey.AskOne(q3, &providerEmail, survey.WithValidator(survey.Required))
			viper.Set("providerEmail", providerEmail)
		}

		domainName := ""
		q4 := &survey.Input{
			Message: "Enter the domain name portion of the hostname to configure. For example, if the hostname is home.example.com, enter example.com:",
		}
		survey.AskOne(q4, &domainName, survey.WithValidator(survey.Required))
		viper.Set("domainname", domainName)

		arecord := ""
		q5 := &survey.Input{
			Message: "Enter the A/AAAA portion of the hostname. For example, if the hostname is home.example.com, enter home:",
		}
		survey.AskOne(q5, &arecord, survey.WithValidator(survey.Required))
		viper.Set("arecord", arecord)

		// Below is neccessary due to an unmerged PR in Viper
		// https://github.com/spf13/viper/pull/450
		//err := viper.SafeWriteConfig()
		err := viper.WriteConfig()
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
}
