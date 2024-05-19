/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

// validateCmd represents the init command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate AWS profiles",
	Long:  `This command validates your AWS local config to ensure that the tool will function as intended.`,
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("%sFailed to get current user: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
		credFilePath := filepath.Join(usr.HomeDir, ".aws", "config")
		credFile, err := ini.Load(credFilePath)
		if err != nil {
			fmt.Printf("%sFailed to read file: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

		var profiles []string
		requiredKeys := []string{
			"sso_account_id",
			"sso_start_url",
			"sso_region",
			"sso_role_name",
			"region",
		}
		missingKeys := make(map[string][]string)

		for _, section := range credFile.Sections() {
			sectionName := section.Name()
			if sectionName != ini.DefaultSection && sectionName != "default" {
				if strings.HasPrefix(sectionName, "profile") {
					sectionName = strings.TrimPrefix(sectionName, "profile ")
				}

				HasMissingKeys := false
				for _, key := range requiredKeys {
					if !section.HasKey(key) {
						missingKeys[sectionName] = append(missingKeys[sectionName], key)
						HasMissingKeys = true
					}
				}

				if !HasMissingKeys {
					profiles = append(profiles, sectionName)
				}
			}
		}

		if len(profiles) == 0 {
			fmt.Printf("%sFailed to find profiles%s\n", ColorRed, ColorReset)
		} else {
			fmt.Printf("%sValid configuration for %d profiles:%s\n", ColorGreen, len(profiles), ColorReset)
			for _, profile := range profiles {
				fmt.Printf("%s- %s%s\n", ColorGreen, profile, ColorReset)
			}
		}

		if len(missingKeys) == 0 {
			fmt.Printf("%sAll profiles have the required keys.%s\n", ColorGreen, ColorReset)
		} else {
			fmt.Printf("%sFailed to validate configuration, ensure profiles are configured correctly with the following values:%s\n", ColorRed, ColorReset)
			for profile, keys := range missingKeys {
				fmt.Printf("%sProfile [%s] is missing keys: %v%s\n", ColorRed, profile, keys, ColorReset)
			}
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
