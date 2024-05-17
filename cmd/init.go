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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalise AWS profiles",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		usr, err := user.Current()
		if err != nil {
			fmt.Printf("Failed to read file: %v", err)
			os.Exit(1)
		}
		credFilePath := filepath.Join(usr.HomeDir, ".aws", "config")
		credFile, err := ini.Load(credFilePath)
		if err != nil {
			fmt.Printf("Failed to read file: %v", err)
			os.Exit(1)
		}

		var profiles []string

		for _, sections := range credFile.Sections() {
			sectionName := sections.Name()
			if sectionName != "DEFAULT" && sectionName != "default" {
				if strings.HasPrefix(sectionName, "profile") {
					profileName := strings.TrimPrefix(sectionName, "profile ")
					profiles = append(profiles, profileName)
				} else {
					profiles = append(profiles, sectionName)
				}
			}
		}
		if len(profiles) == 0 {
			fmt.Println("Failed to find profiles")
		} else {
			fmt.Println("Success, Profiles found:", profiles)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
