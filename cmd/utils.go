package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

var requiredKeys = []string{
	"sso_account_id",
	"sso_start_url",
	"sso_region",
	"sso_role_name",
	"region",
}

func validateProfiles() ([]string, map[string][]string, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current user: %v", err)
	}
	credFilePath := filepath.Join(usr.HomeDir, ".aws", "config")
	credFile, err := ini.Load(credFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file: %v", err)
	}

	var profiles []string
	missingKeys := make(map[string][]string)

	for _, section := range credFile.Sections() {
		rawSectionName := section.Name()
		if rawSectionName != ini.DefaultSection && rawSectionName != "default" {
			sectionName := rawSectionName
			if strings.HasPrefix(rawSectionName, "profile ") {
				sectionName = strings.TrimPrefix(rawSectionName, "profile ")
			}

			hasMissingKeys := false
			for _, key := range requiredKeys {
				if !section.HasKey(key) {
					missingKeys[sectionName] = append(missingKeys[sectionName], key)
					hasMissingKeys = true
				}
			}

			if !hasMissingKeys {
				profiles = append(profiles, sectionName)
			}
		}
	}

	return profiles, missingKeys, nil
}

func getAwsConsoleUrl(profile string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %v", err)
	}
	cacheDir := filepath.Join(usr.HomeDir, ".aws", "sso", "cache")

	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return "", fmt.Errorf("failed to read SSO cache directory: %v", err)
	}

	var cachedCredentials map[string]interface{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			data, err := ioutil.ReadFile(filepath.Join(cacheDir, file.Name()))
			if err != nil {
				return "", fmt.Errorf("failed to read SSO cache file: %v", err)
			}

			var creds map[string]interface{}
			if err := json.Unmarshal(data, &creds); err != nil {
				return "", fmt.Errorf("failed to parse SSO cache file: %v", err)
			}

			if startUrl, ok := creds["startUrl"].(string); ok &&
				strings.Contains(startUrl, profile) {
				cachedCredentials = creds
				break
			}
		}
	}

	if cachedCredentials == nil {
		return "", fmt.Errorf("no cached credentials found for profile %s", profile)
	}

	accessKeyId := cachedCredentials["accessKeyId"].(string)
	secretAccessKey := cachedCredentials["secretAccessKey"].(string)
	sessionToken := cachedCredentials["sessionToken"].(string)

	// Generate a sign-in token
	signinTokenUrl := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=getSigninToken&Session=%s",
		url.QueryEscape(
			fmt.Sprintf(
				`{"sessionId":"%s","sessionKey":"%s","sessionToken":"%s"}`,
				accessKeyId,
				secretAccessKey,
				sessionToken,
			),
		),
	)

	getTokenCmd := exec.Command("curl", "-s", signinTokenUrl)
	tokenOutput, err := getTokenCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get sign-in token: %v", err)
	}

	var tokenResponse map[string]string
	if err := json.Unmarshal(tokenOutput, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to parse sign-in token response: %v", err)
	}

	signinToken := tokenResponse["SigninToken"]
	// Construct the console URL
	consoleUrl := fmt.Sprintf(
		"https://signin.aws.amazon.com/federation?Action=login&Issuer=Example.org&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&SigninToken=%s",
		signinToken,
	)
	return consoleUrl, nil
}
