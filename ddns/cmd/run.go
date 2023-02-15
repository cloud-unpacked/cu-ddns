package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/cloudflare/cloudflare-go"
	"github.com/digitalocean/godo"
	"github.com/linode/linodego"
	log "github.com/sirupsen/logrus"
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

		modeDryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Error("Failed to retrieve dry-run flag value.")
		}

		if modeDryRun {

			fmt.Printf("The current public IP address is: %v", dIP)
			return
		}

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

			if recordID == 0 {
				log.Fatal("Record not found.")
			}

			_, err = linodeClient.UpdateDomainRecord(context.Background(), domainID, recordID, linodego.DomainRecordUpdateOptions{Target: dIP})
		} else if viper.GetString("provider") == "cloudflare" {

			var api *cloudflare.API

			// backwards compatibility support for the old Cloudflare key/email combo
			// will remove somewhere down the line, likely before v1.0 release
			if viper.IsSet("providerEmail") {

				api, err = cloudflare.New(apiToken, viper.GetString("providerEmail"))
				if err != nil {
					log.Fatal(err)
				}
			} else {

				api, err = cloudflare.NewWithAPIToken(apiToken)
				if err != nil {
					log.Fatal(err)
				}
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
		} else if viper.GetString("provider") == "digitalocean" {

			tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiToken})
			oauthClient := oauth2.NewClient(context.Background(), tokenSource)

			doClient := godo.NewClient(oauthClient)
			ctx := context.TODO()

			opt := &godo.ListOptions{
				Page:    1,
				PerPage: 25,
			}

			records, _, err := doClient.Domains.Records(ctx, theDomain, opt)
			if err != nil {
				log.Fatal(err)
			}

			for _, record := range records {

				if record.Name == theHostname {
					recordID = record.ID
					break
				}
			}

			if recordID == 0 {
				log.Fatal("Record not found.")
			}

			editRequest := &godo.DomainRecordEditRequest{
				Data: dIP,
			}

			_, _, err = doClient.Domains.EditRecord(ctx, theDomain, recordID, editRequest)
		} else {
			log.Fatal("Provider not supported.")
		}

		log.Info("Successfully updated " + viper.GetString("provider") + " the IP: " + dIP)
	},
}

func init() {

	rootCmd.AddCommand(runCmd)
	runCmd.Flags().Bool("dry-run", false, "Pretend to run but don't update provider")
}
