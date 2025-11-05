package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"gs-write/pkg/auth"
	"gs-write/pkg/sheets"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// cfgFile は設定ファイルのパスを保持するグローバル変数です
	cfgFile string
	// title is the title of the spreadsheet
	title string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gs-write",
	Short: "Write stdin to a new Google Spreadsheet",
	Long: `gs-write is a simple CLI tool that writes standard input to a new Google Spreadsheet.
It is designed to work with pipes (|) based on UNIX philosophy.

Examples:
  ls -l | gs-write
  cat report.csv | gs-write --title "Monthly Report"
  ps aux | gs-write`,
	RunE: runRoot,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Add flags
	rootCmd.Flags().StringVar(&title, "title", "", "Title of the spreadsheet (default: auto-generated from timestamp)")

	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(versionCmd)
}

func runRoot(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load authentication config
	oauthConfig, token, err := auth.GetClient(ctx)
	if err != nil {
		return err
	}

	// Create Sheets client
	client, err := sheets.NewClient(ctx, oauthConfig, token)
	if err != nil {
		return err
	}

	// Read CSV data from stdin
	data, err := readCSVFromStdin()
	if err != nil {
		return fmt.Errorf("failed to read CSV from stdin: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("no data provided")
	}

	// Create spreadsheet
	url, err := client.CreateSpreadsheet(ctx, title, data)
	if err != nil {
		return err
	}

	// Output the URL
	fmt.Println(url)

	return nil
}

// readCSVFromStdin reads CSV data from standard input
func readCSVFromStdin() ([][]string, error) {
	reader := csv.NewReader(os.Stdin)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

// initViper reads in config file and ENV variables if set.
func initViper() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in current directory and $HOME directory
		viper.AddConfigPath(".") // プロジェクトルート
		viper.SetConfigName("config") // config.toml, config.json, config.yaml... を探す (拡張子なし)
		viper.SetConfigType("toml")   // TOML形式であることを明示的に指定
	}

	viper.SetDefault("debug", false) // デバッグモードのデフォルト

	// 設定ファイルを読み込む
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		// 設定ファイルが見つからない場合やその他の読み込みエラー
		// ConfigFileNotFoundError の場合は警告のみ、それ以外は Fatal
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore
		} else {
			// Config file was found but another error was produced
			fmt.Fprintf(os.Stderr, "Error reading config file: %s\n", err)
		}
	}
}
