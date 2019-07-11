package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/cloudflare/cloudflare-go"
	"github.com/linode/linodego"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the verify command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Determine IP addresses and update DNS if needed",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {

		// If we're running within a Snap, this is called by a SystemD Timer.
		// We should only continue if "ready" file exists.
		if snapDir := os.Getenv("SNAP_DATA"); snapDir != "" {
			_, err := os.Stat(snapDir + "/cu-ddns.active")
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Print("cu-ddns isn't set to run. Exiting.")
					return
				} else {
					log.Fatal("Something went wrong checking for snap exists file.")
				}
			}
		}

		sourceURL := "http://4.myip.cloudunpacked.com"
		apiToken := viper.GetString("providertoken")
		theDomain := viper.GetString("domainname")
		theHostname := viper.GetString("arecord")
		domainID := 0
		recordID := 0
		dIP := ""

		resp, err := http.Get(sourceURL)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		dIP = string(body)

		// Below starts provider specific code. In the future, this is an area
		// where separating out into functions and likely interfaces will be
		// useful.

		if viper.GetString("provider") == "linode" {

			tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})

			oauth2Client := &http.Client{
				Transport: &oauth2.Transport{
					Source: tokenSource,
				},
			}

			linodeClient := linodego.NewClient(oauth2Client)
			//linodeClient.SetDebug(true)

			domains, err := linodeClient.ListDomains(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}

			for _, domain := range domains {

				if domain.Domain == theDomain {
					domainID = domain.ID
					break
				}
			}

			records, err := linodeClient.ListDomainRecords(context.Background(), domainID, nil)
			if err != nil {
				fmt.Println("It was this.")
				log.Fatal(err)
			}

			for _, record := range records {

				if record.Name == theHostname {
					recordID = record.ID
					break
				}
			}

			_, err = linodeClient.UpdateDomainRecord(context.Background(), domainID, recordID, linodego.DomainRecordUpdateOptions{Target: dIP})
		} else if viper.GetString("provider") == "cloudflare" {

			api, err := cloudflare.New(apiToken, viper.GetString("providerEmail"))
			if err != nil {
				log.Fatal(err)
			}

			zoneID, err := api.ZoneIDByName(theDomain)
			if err != nil {
				log.Fatal(err)
			}

			recFilter := cloudflare.DNSRecord{Name: theHostname + "." + theDomain}
			records, err := api.DNSRecords(zoneID, recFilter)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("This record is: " + records[0].Name)
			records[0].Content = dIP
			err = api.UpdateDNSRecord(zoneID, records[0].ID, records[0])
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("Provider not supported.")
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
