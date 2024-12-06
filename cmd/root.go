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

func getConfigurationFiles(cmd *cobra.Command) (string, string, error) {
	// Get participants Path
	participantsPath, err := cmd.Flags().GetString("participants")
	if err != nil || participantsPath == "" {
		participantsPath = "./participants.csv"
	}
	// Get config path
	configPath, err := cmd.Flags().GetString("config")
	if err != nil || configPath == "" {
		configPath = "./config.yaml"
	}

	return configPath, participantsPath, nil
}

func checkConfigFiles(configPath, participantsPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file does not exist at %s", configPath)
	}
	if _, err := os.Stat(participantsPath); os.IsNotExist(err) {
		return fmt.Errorf("participants file does not exist at %s", participantsPath)
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "secret-santa",
	Short: "Email secret santa messages to a group",
	Long: `Emails a list of interests to a random recipient from a list.
	If a partner is defined, a person will not get their partner.`,
	Run: func(cmd *cobra.Command, args []string) {

		configPath, participantsPath, err := getConfigurationFiles(cmd)
		if err != nil {
			fmt.Printf("error getting configuration files: %v\n", err)
			os.Exit(1)
		}

		// Get generate config file flag
		generateConfigFile, err := cmd.Flags().GetBool("generate-config")
		if err != nil {
			fmt.Printf("error retrieving generate-config flag\n")
			os.Exit(1)
		}
		// Get generate participants file flag
		generateParticipantsFile, err := cmd.Flags().GetBool("generate-participants")
		if err != nil {
			fmt.Printf("error retrieving generate-participants flag: %v", err)
			os.Exit(1)
		}

		// Generate config files if flags are set and exit the program
		if generateConfigFile || generateParticipantsFile {
			err := conf.GenerateConfigFiles(configPath, generateConfigFile, participantsPath, generateParticipantsFile)
			if err != nil {
				fmt.Printf("error generating config files: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		err = checkConfigFiles(configPath, participantsPath)
		if err != nil {
			fmt.Printf("error checking config files: %v\n", err)
			os.Exit(1)
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
			apiKey := viper.GetString("mailgun.apikey")
			if apiKey == "" {
				fmt.Printf("Please set a mailgun api key in the config file\n")
				os.Exit(1)
			}

			sender.Emailer = mgmailer.NewMailgunEmailer(emailDomain, apiKey)
		}
		err = sender.Send(participantsPath)
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
	rootCmd.Flags().BoolP("dry-run", "d", false, "dry-run will print a list rather than emailing people")
	rootCmd.Flags().StringP("participants", "p", "", "a csv file with participants (required)")
	rootCmd.Flags().StringP("email-template", "e", "", "a go template file for the email body")
	rootCmd.Flags().StringP("config", "c", "", "A configuration file for the application (required)")
	rootCmd.Flags().BoolP("generate-config", "", false, "generate a config file. Note that this will overwrite an existing config file, and the application will not run. Can be used with the --config flag to specify a path and name")
	rootCmd.Flags().BoolP("generate-participants", "", false, "generate a participants file. Note that this will overwrite an existing participants file of the same name, and the application will not run. Can be used with the --participants flag to specify a path and name")

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
