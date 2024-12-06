package cmd

import (
	"fmt"
	"os"

	"github.com/dcmcand/go-secret-santa/package/conf"
	csvLoader "github.com/dcmcand/go-secret-santa/package/csvparticipantloader"
	fakeMailer "github.com/dcmcand/go-secret-santa/package/fakemailer"
	"github.com/dcmcand/go-secret-santa/package/mgmailer"
	"github.com/dcmcand/go-secret-santa/package/send"
	"github.com/dcmcand/go-secret-santa/package/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "secret-santa",
	Short: "Email secret santa messages to a group",
	Long: `Emails a list of interests to a random recipient from a list.
	If a partner is defined, a person will not get their partner.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get participants
		participants, err := cmd.Flags().GetString("participants")
		if err != nil || participants == "" {
			fmt.Printf("No path set for participants\n")
			_, err := os.Stat("./participants.csv")
			if err != nil {
				fmt.Printf("No participants file found. Generating Skeleton at ./participants.csv\n")
				// err = conf.GenerateParticipantsFile()
				// if err != nil {
				// 	fmt.Printf("error generating participants file: %v\n", err)
				// 	os.Exit(1)
				// }
				os.Exit(1)
			}
			participants = "./participants.csv"
		}

		// Get config
		configPath, err := cmd.Flags().GetString("config")
		if err != nil || configPath == "" {
			fmt.Printf("No Config path set\n")
			_, err := os.Stat("./config.yaml")
			if err != nil {
				fmt.Println("No config file found. Generating Skeleton at ./config.yaml")
				err = conf.GenerateConfigFile()
				if err != nil {
					fmt.Printf("error generating config file: %v\n", err)
				}
				os.Exit(1)
			}
			fmt.Printf("found ./config.yaml\n")
			configPath = "./config.yaml"
		}
		initConfig(configPath)

		// Setup Email
		subject := viper.GetString("email.subject")
		if subject == "" {
			subject = "Secret Santa Assignment"
		}
		domain := viper.GetString("email.domain")
		if domain == "" {
			fmt.Printf("Please set a domain to send email from")
			os.Exit(1)
		}
		senderEmail := viper.GetString("email.sender.address")
		if senderEmail == "" {
			senderEmail = fmt.Sprintf("santa@%s", domain)
		}
		senderName := viper.GetString("email.sender.name")
		if senderName == "" {
			senderName = "Santa Claus"
		}
		emailDomain := viper.GetString("email.domain")
		if emailDomain == "" {
			fmt.Printf("Please set an email domain in the config file\n")
			os.Exit(1)
		}
		var emailTemplate *send.Email
		e, err := cmd.Flags().GetString("email-template")
		if err != nil || e == "" {
			emailTemplate, err = template.GetDefaultTemplate(subject, senderName, senderEmail)
			if err != nil {
				fmt.Printf("error getting default template: %v\n", err)
				os.Exit(1)
			}
		} else {
			emailTemplate, err = template.GetTemplate(e, subject, senderName, senderEmail)
			if err != nil {
				fmt.Printf("error getting template: %v\n", err)
				os.Exit(1)
			}
		}

		// Send Emails
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			fmt.Printf("error retrieving dry-run flag\n")
			os.Exit(1)
		}

		loader := csvLoader.Loader{}
		sender := send.Sender{
			ParticipantLoader: &loader,
			EmailTemplate:     emailTemplate,
		}
		if dryRun {
			sender.Emailer = &fakeMailer.Mailer{}
		} else {
			sender.Emailer = mgmailer.NewMailgunEmailer(emailDomain, "abc123")
		}
		err = sender.Send(participants)
		if err != nil {
			fmt.Printf("error sending emails: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("dry-run", "t", false, "dry-run will print a list rather than emailing people")
	rootCmd.Flags().StringP("participants", "p", "", "a csv file with participants (required)")
	rootCmd.Flags().StringP("email-template", "e", "", "a go template file for the email body")
	rootCmd.Flags().StringP("config", "c", "", "A configuration file for the application (required)")
}

func initConfig(configPath string) {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}
	viper.ReadInConfig()

	viper.AutomaticEnv() // read in environment variables that match
}
