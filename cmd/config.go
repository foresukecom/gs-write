package cmd

import (
	"fmt"
	"gs-write/pkg/config"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage gs-write configuration / gs-writeの設定を管理",
	Long: `Manage gs-write configuration settings.
gs-writeの設定を管理します。

Available settings / 利用可能な設定:
  freeze.rows       - Number of rows to freeze / 固定する行数 (default: not set)
  freeze.cols       - Number of columns to freeze / 固定する列数 (default: not set)
  filter.header_row - Header row for basic filter / フィルタのヘッダー行 (default: not set)

Examples / 使用例:
  gs-write config list
  gs-write config get freeze.rows
  gs-write config set freeze.rows 1
  gs-write config set filter.header_row 1
  gs-write config unset freeze.rows`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration settings / すべての設定を表示",
	Long: `Display all current configuration settings.
現在の設定をすべて表示します。`,
	RunE: runConfigList,
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value / 設定値を取得",
	Long: `Get the value of a specific configuration setting.
指定した設定の値を取得します。

Available keys / 利用可能なキー:
  freeze.rows       - Number of rows to freeze / 固定する行数
  freeze.cols       - Number of columns to freeze / 固定する列数
  filter.header_row - Header row for basic filter / フィルタのヘッダー行`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigGet,
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value / 設定値を変更",
	Long: `Set the value of a specific configuration setting.
指定した設定の値を変更します。

Available keys / 利用可能なキー:
  freeze.rows       - Number of rows to freeze / 固定する行数 (must be non-negative integer / 非負の整数)
  freeze.cols       - Number of columns to freeze / 固定する列数 (must be non-negative integer / 非負の整数)
  filter.header_row - Header row for basic filter / フィルタのヘッダー行 (must be non-negative integer / 非負の整数)

Examples / 使用例:
  gs-write config set freeze.rows 1
  gs-write config set freeze.cols 2
  gs-write config set filter.header_row 1`,
	Args: cobra.ExactArgs(2),
	RunE: runConfigSet,
}

var configUnsetCmd = &cobra.Command{
	Use:   "unset <key>",
	Short: "Unset a configuration value / 設定値を削除",
	Long: `Remove a configuration setting, reverting to default behavior.
設定を削除し、デフォルト動作に戻します。

Available keys / 利用可能なキー:
  freeze.rows       - Number of rows to freeze / 固定する行数
  freeze.cols       - Number of columns to freeze / 固定する列数
  filter.header_row - Header row for basic filter / フィルタのヘッダー行

Examples / 使用例:
  gs-write config unset freeze.rows
  gs-write config unset freeze.cols
  gs-write config unset filter.header_row`,
	Args: cobra.ExactArgs(1),
	RunE: runConfigUnset,
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configUnsetCmd)
}

func runConfigList(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current configuration:")

	// Always show the effective value (user configured or default)
	rows, _ := cfg.GetFreezeRows()
	fmt.Printf("  freeze.rows = %d\n", rows)

	cols, _ := cfg.GetFreezeCols()
	fmt.Printf("  freeze.cols = %d\n", cols)

	headerRow, _ := cfg.GetFilterHeaderRow()
	fmt.Printf("  filter.header_row = %d\n", headerRow)

	return nil
}

func runConfigGet(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "freeze.rows":
		// Always return the effective value (user configured or default)
		rows, _ := cfg.GetFreezeRows()
		fmt.Println(rows)
	case "freeze.cols":
		// Always return the effective value (user configured or default)
		cols, _ := cfg.GetFreezeCols()
		fmt.Println(cols)
	case "filter.header_row":
		// Always return the effective value (user configured or default)
		headerRow, _ := cfg.GetFilterHeaderRow()
		fmt.Println(headerRow)
	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	valueStr := args[1]

	cfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "freeze.rows":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid value for freeze.rows: must be an integer")
		}
		if value < 0 {
			return fmt.Errorf("invalid value for freeze.rows: must be non-negative (got: %d)", value)
		}
		cfg.SetFreezeRows(value)
		fmt.Printf("Set freeze.rows = %d\n", value)

	case "freeze.cols":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid value for freeze.cols: must be an integer")
		}
		if value < 0 {
			return fmt.Errorf("invalid value for freeze.cols: must be non-negative (got: %d)", value)
		}
		cfg.SetFreezeCols(value)
		fmt.Printf("Set freeze.cols = %d\n", value)

	case "filter.header_row":
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return fmt.Errorf("invalid value for filter.header_row: must be an integer")
		}
		if value < 0 {
			return fmt.Errorf("invalid value for filter.header_row: must be non-negative (got: %d)", value)
		}
		cfg.SetFilterHeaderRow(value)
		fmt.Printf("Set filter.header_row = %d\n", value)

	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	if err := config.SaveUserConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath, _ := config.GetConfigPath()
	fmt.Printf("Configuration saved to: %s\n", configPath)

	return nil
}

func runConfigUnset(cmd *cobra.Command, args []string) error {
	key := args[0]

	cfg, err := config.LoadUserConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	switch key {
	case "freeze.rows":
		cfg.UnsetFreezeRows()
		fmt.Println("Unset freeze.rows")

	case "freeze.cols":
		cfg.UnsetFreezeCols()
		fmt.Println("Unset freeze.cols")

	case "filter.header_row":
		cfg.UnsetFilterHeaderRow()
		fmt.Println("Unset filter.header_row")

	default:
		return fmt.Errorf("unknown configuration key: %s", key)
	}

	if err := config.SaveUserConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath, _ := config.GetConfigPath()
	fmt.Printf("Configuration saved to: %s\n", configPath)

	return nil
}

// normalizeKey converts various key formats to dot notation
func normalizeKey(key string) string {
	// Replace underscores and hyphens with dots
	key = strings.ReplaceAll(key, "_", ".")
	key = strings.ReplaceAll(key, "-", ".")
	return key
}
