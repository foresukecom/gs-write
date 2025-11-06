package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"gs-write/pkg/auth"
	"gs-write/pkg/config"
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
	// freezeRows is the number of rows to freeze (nil means not set via CLI)
	freezeRowsFlag *int
	// freezeCols is the number of columns to freeze (nil means not set via CLI)
	freezeColsFlag *int
	// filterHeaderRow is the header row for basic filter (nil means not set via CLI)
	filterHeaderRowFlag *int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gs-write",
	Short: "Write stdin to a new Google Spreadsheet / 標準入力を新しいGoogleスプレッドシートに書き込む",
	Long: `gs-write is a simple CLI tool that writes standard input to a new Google Spreadsheet.
標準入力を新しいGoogleスプレッドシートに書き込むシンプルなCLIツールです。

It is designed to work with pipes (|) based on UNIX philosophy.
UNIX哲学に基づき、パイプ(|)で他のコマンドと連携することを前提に設計されています。

Examples / 使用例:
  ls -l | gs-write
  cat report.csv | gs-write --title "Monthly Report"
  cat data.csv | gs-write --freeze-rows 1 --freeze-cols 0
  cat data.csv | gs-write --filter-header-row 1
  ps aux | gs-write --title "Processes" --freeze-rows 1 --filter-header-row 1`,
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

	// Use pointer flags to distinguish between "not set" and "set to 0"
	freezeRowsFlag = rootCmd.Flags().Int("freeze-rows", -1, "Number of rows to freeze (overrides config file)")
	freezeColsFlag = rootCmd.Flags().Int("freeze-cols", -1, "Number of columns to freeze (overrides config file)")
	filterHeaderRowFlag = rootCmd.Flags().Int("filter-header-row", -1, "Header row for basic filter (overrides config file)")

	// Disable completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// Add subcommands
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
}

func runRoot(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load user config
	userConfig, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine freeze parameters with priority: CLI > config > default
	freezeRows := resolveFreezeRows(cmd, userConfig)
	freezeCols := resolveFreezeCols(cmd, userConfig)
	filterHeaderRow := resolveFilterHeaderRow(cmd, userConfig)

	// Validate parameters
	if freezeRows < 0 || freezeCols < 0 {
		return fmt.Errorf("freeze-rows and freeze-cols must be non-negative (got: rows=%d, cols=%d)", freezeRows, freezeCols)
	}
	if filterHeaderRow < 0 {
		return fmt.Errorf("filter-header-row must be non-negative (got: %d)", filterHeaderRow)
	}

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
	url, err := client.CreateSpreadsheet(ctx, title, data, freezeRows, freezeCols, filterHeaderRow)
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

// resolveFreezeRows determines the freeze rows value with priority: CLI > config > default
func resolveFreezeRows(cmd *cobra.Command, userConfig *config.UserConfig) int {
	// Check if CLI flag was explicitly set
	if cmd.Flags().Changed("freeze-rows") {
		return *freezeRowsFlag
	}

	// Check if config has a value
	if rows, ok := userConfig.GetFreezeRows(); ok {
		return rows
	}

	// Return default value
	return 0
}

// resolveFreezeCols determines the freeze cols value with priority: CLI > config > default
func resolveFreezeCols(cmd *cobra.Command, userConfig *config.UserConfig) int {
	// Check if CLI flag was explicitly set
	if cmd.Flags().Changed("freeze-cols") {
		return *freezeColsFlag
	}

	// Check if config has a value
	if cols, ok := userConfig.GetFreezeCols(); ok {
		return cols
	}

	// Return default value
	return 0
}

// resolveFilterHeaderRow determines the filter header row value with priority: CLI > config > default
func resolveFilterHeaderRow(cmd *cobra.Command, userConfig *config.UserConfig) int {
	// Check if CLI flag was explicitly set
	if cmd.Flags().Changed("filter-header-row") {
		return *filterHeaderRowFlag
	}

	// Check if config has a value
	if headerRow, ok := userConfig.GetFilterHeaderRow(); ok {
		return headerRow
	}

	// Return default value (0 means no filter)
	return 0
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
