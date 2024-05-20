/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// validateCmd represents the init command

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate AWS profiles",
	Long:  `This command validates your AWS local config to ensure that the tool will function as intended.`,

	Run: func(cmd *cobra.Command, args []string) {
		profiles, missingKeys, err := validateProfiles()
		if err != nil {
			fmt.Printf("%s%s%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

		if len(profiles) == 0 {
			fmt.Printf("%sFailed to find valid profiles%s\n", ColorRed, ColorReset)
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
