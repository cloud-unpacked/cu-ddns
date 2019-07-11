package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh/terminal"

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

		for {

			fmt.Println("Please choose a provider:\n")
			fmt.Println("[1] Linode DNS")
			fmt.Println("[2] Cloudflare DNS\n")
			fmt.Println("Enter the number of the provider [1]: ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error")
			}

			input = strings.TrimSpace(input)

			choice, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Not a valid option.")
				continue
			}

			// Basic, almost useless validation
			choiceValid := false
			switch choice {
			case 1:
				choiceValid = true
				viper.Set("provider", "linode")
				break
			case 2:
				choiceValid = true
				viper.Set("provider", "cloudflare")
				break
			default:
				fmt.Println("Not a valid option.")
			}

			if choiceValid {
				break
			}
		}

		for {

			promptString := ""
			switch viper.GetString("provider") {
			case "linode":
				promptString = "Linode Personal Access Token"
				break
			case "cloudflare":
				promptString = "Cloudflare API key"
				break
			}

			fmt.Println("Please paste a " + promptString)
			fmt.Print("for authentication: ")

			passwdBytes, err := terminal.ReadPassword(0)
			if err != nil {
				fmt.Println("Error")
			}

			// Disabling below for now as Cloudflare keys are different
			// lengths from Linode keys.

			// Basic, almost useless validation
			// For the record, Linode APIv4 tokens:
			//    64chars
			//    /[a-f0-9]{64}/
			//if len(string(passwdBytes)) == 64 {
			if len(string(passwdBytes)) >= 10 {
				viper.Set("providerToken", string(passwdBytes))
				break
			} else {
				fmt.Println("The token is invalid. Try again.\n")
			}
		}

		for {

			if viper.GetString("provider") != "cloudflare" {
				// Exit this loop. Only Cloudflare needs an email address.
				break
			}

			fmt.Print("\nEnter your Cloudflare email address: ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error")
			}

			input = strings.TrimSpace(input)

			if len(input) >= 5 {
				viper.Set("providerEmail", input)
				break
			} else {
				fmt.Println("The email is invalid. Try again.\n")
				continue
			}

		}

		for {

			fmt.Println("Enter the domain name portion of the hostname")
			fmt.Println("to configure. For example, if the hostname is")
			fmt.Print("home.example.com, enter example.com: ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error")
			}

			input = strings.TrimSpace(input)

			if len(input) == 0 {
				fmt.Println("Error: Please provide a domain name.")
				continue
			}

			viper.Set("domainname", input)
			break
		}

		for {

			fmt.Println("Enter the A/AAAA portion of the hostname. For")
			fmt.Println("example, if the hostname is home.example.com,")
			fmt.Println("enter home: ")

			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error")
			}

			input = strings.TrimSpace(input)

			if len(input) == 0 {
				fmt.Println("Error: Please provide an A record.")
				continue
			}

			viper.Set("arecord", input)
			break
		}

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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configureCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configureCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
