package conf

import (
	"encoding/csv"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func GenerateConfigFile() error {
	config := &yaml.Node{
		Kind: yaml.DocumentNode,
		Content: []*yaml.Node{
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Value: "email",
					},
					{
						Kind: yaml.MappingNode,
						Content: []*yaml.Node{

							{
								Kind:  yaml.ScalarNode,
								Value: "subject",
							},
							{
								Kind:        yaml.ScalarNode,
								Style:       yaml.DoubleQuotedStyle,
								Value:       "Secret Santa",
								LineComment: "# This is the subject of the secret santa email",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "address",
							},
							{
								Kind:        yaml.ScalarNode,
								Style:       yaml.DoubleQuotedStyle,
								Value:       "santa",
								LineComment: "# This along with the domain is used to create the email address of the sender",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "domain",
							},
							{
								Kind:        yaml.ScalarNode,
								Style:       yaml.DoubleQuotedStyle,
								LineComment: "# This is the domain of the email",
							},
							{
								Kind:  yaml.ScalarNode,
								Value: "sender",
							},
							{
								Kind: yaml.MappingNode,
								Content: []*yaml.Node{
									{
										Kind:  yaml.ScalarNode,
										Value: "name",
									},
									{
										Kind:        yaml.ScalarNode,
										Style:       yaml.DoubleQuotedStyle,
										Value:       "Santa Claus",
										LineComment: "# This is the name of the sender for use in the email body",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	y, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	err = os.WriteFile("config.yaml", y, 0644)
	if err != nil {
		return err
	}
	return nil
}

func GenerateParticipantsFile() error {
	_, err := os.Stat("./participants.csv")
	if err == nil {
		return fmt.Errorf("participants.csv already exists")
	}
	f, err := os.Create("participants.csv")
	if err != nil {
		return fmt.Errorf("error creating participants.csv: %v", err)
	}
	defer f.Close()
	content := [][]string{
		{"Name", "Email", "Partner", "Interests"},
		{"Barney", "barney@bedrock.com", "Betty", "Bowling, Jokes, Movies"},
		{"Fred", "fred@bedrock.com", "Wilma", "Bowling, Dinosaurs, Golf"},
		{"Wilma", "wilma@bedrock.com", "Fred", "Cooking, Gardening, Shopping"},
		{"Betty", "betty@bedrock.com", "Barney", "Reading, Music, Crafts"},
		{"Pebbles", "pebbles@bedrock.com", "", "Exploring, Drawing, Sports"},
		{"BamBam", "bambam@bedrock.com", "", "Rock Music, Cave Painting, Athletics"},
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()

	err = writer.WriteAll(content)
	if err != nil {
		return fmt.Errorf("error writing data to csv: %v", err)
	}
	return nil
}
