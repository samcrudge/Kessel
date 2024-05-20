/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ColorReset = "\033[0m"
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in using AWS SSO",
	Long:  `This command allows you to log in using AWS SSO and select a profile to authenticate.`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles, _, err := validateProfiles()
		if err != nil {
			fmt.Printf("%s%s%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

		if len(profiles) == 0 {
			fmt.Printf("%sNo valid profiles found%s\n", ColorRed, ColorReset)
			os.Exit(1)
		}

		// List profiles
		fmt.Printf("%sSelect a profile to log in:%s\n", ColorGreen, ColorReset)
		for i, profile := range profiles {
			fmt.Printf("%s%d. %s%s\n", ColorGreen, i+1, profile, ColorReset)
		}

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter the number of the profile you want to log in with: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("%sFailed to read input: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(profiles) {
			fmt.Printf("%sInvalid choice%s\n", ColorRed, ColorReset)
			os.Exit(1)
		}

		selectedProfile := profiles[choice-1]

		// Perform AWS SSO login
		loginCmd := exec.Command("aws", "sso", "login", "--profile", selectedProfile)
		loginCmd.Stdout = os.Stdout
		loginCmd.Stderr = os.Stderr
		if err := loginCmd.Run(); err != nil {
			fmt.Printf(
				"%sFailed to log in with profile %s: %v%s\n",
				ColorRed,
				selectedProfile,
				err,
				ColorReset,
			)
			os.Exit(1)
		}

		fmt.Printf(
			"%sSuccessfully logged in with profile %s%s\n",
			ColorGreen,
			selectedProfile,
			ColorReset,
		)

		// Get the AWS Console URL
		consoleUrl, err := getAwsConsoleUrl(selectedProfile)
		if err != nil {
			fmt.Printf("%sFailed to get AWS console URL: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}

		// Open the URL in the default browser
		fmt.Printf("%sOpening AWS console in your default browser...%s\n", ColorGreen, ColorReset)
		openBrowserCmd := exec.Command(
			"open",
			consoleUrl,
		) // Use "xdg-open" for Linux and "start" for Windows
		if err := openBrowserCmd.Start(); err != nil {
			fmt.Printf("%sFailed to open browser: %v%s\n", ColorRed, err, ColorReset)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
