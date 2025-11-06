package cmd

import (
	"bufio"
	"context"
	"fmt"
	"gs-write/pkg/auth"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	credentialsFile string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Google Sheets API",
	Long: `Authenticate with Google Sheets API using OAuth 2.0.

You can provide credentials in two ways:
1. Interactively paste the credentials JSON
2. Provide a credentials file using --credentials flag

Example:
  gs-write auth
  gs-write auth --credentials ./credentials.json`,
	RunE: runAuth,
}

func init() {
	authCmd.Flags().StringVar(&credentialsFile, "credentials", "", "Path to credentials.json file")
}

func runAuth(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Get credentials JSON
	var credentialsJSON []byte
	var err error

	if credentialsFile != "" {
		// Read from file
		credentialsJSON, err = os.ReadFile(credentialsFile)
		if err != nil {
			return fmt.Errorf("failed to read credentials file: %w", err)
		}
	} else {
		// Read from stdin interactively
		fmt.Println("Please paste your credentials JSON (press Ctrl+D when done):")
		credentialsJSON, err = readCredentialsFromStdin()
		if err != nil {
			return fmt.Errorf("failed to read credentials: %w", err)
		}
	}

	// Parse credentials
	oauthConfig, err := auth.ParseCredentials(credentialsJSON)
	if err != nil {
		return err
	}

	// Get authorization URL
	authURL := auth.GetAuthURL(oauthConfig)
	fmt.Printf("\nPlease visit the following URL to authorize this application:\n%s\n\n", authURL)

	// Wait for authorization code
	fmt.Print("Enter the authorization code: ")
	reader := bufio.NewReader(os.Stdin)
	code, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read authorization code: %w", err)
	}
	code = strings.TrimSpace(code)

	// Exchange code for token
	token, err := auth.ExchangeToken(ctx, oauthConfig, code)
	if err != nil {
		return err
	}

	// Save authentication config
	config := &auth.AuthConfig{
		Credentials: oauthConfig,
		Token:       token,
	}

	if err := auth.SaveAuthConfig(config); err != nil {
		return err
	}

	authPath, _ := auth.GetAuthPath()
	fmt.Printf("\nAuthentication successful!\nAuthentication saved to: %s\n", authPath)

	return nil
}

// readCredentialsFromStdin reads multi-line JSON input from stdin
func readCredentialsFromStdin() ([]byte, error) {
	var builder strings.Builder
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		builder.WriteString(line)
	}

	input := strings.TrimSpace(builder.String())
	if input == "" {
		return nil, fmt.Errorf("no credentials provided")
	}

	return []byte(input), nil
}
