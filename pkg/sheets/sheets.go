package sheets

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Client wraps the Google Sheets API client
type Client struct {
	service *sheets.Service
}

// NewClient creates a new Sheets client
func NewClient(ctx context.Context, config *oauth2.Config, token *oauth2.Token) (*Client, error) {
	httpClient := config.Client(ctx, token)

	service, err := sheets.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &Client{service: service}, nil
}

// CreateSpreadsheet creates a new spreadsheet with the given title and data
func (c *Client) CreateSpreadsheet(ctx context.Context, title string, data [][]string) (string, error) {
	// If no title is provided, generate one from timestamp
	if title == "" {
		title = generateDefaultTitle()
	}

	// Create a new spreadsheet
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: title,
		},
		Sheets: []*sheets.Sheet{
			{
				Properties: &sheets.SheetProperties{
					Title: "Sheet1",
				},
			},
		},
	}

	resp, err := c.service.Spreadsheets.Create(spreadsheet).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to create spreadsheet: %w", err)
	}

	spreadsheetID := resp.SpreadsheetId

	// Write data to the spreadsheet
	if len(data) > 0 {
		if err := c.writeData(ctx, spreadsheetID, "Sheet1", data); err != nil {
			return "", fmt.Errorf("failed to write data: %w", err)
		}
	}

	// Return the spreadsheet URL
	url := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/edit", spreadsheetID)
	return url, nil
}

// writeData writes data to the specified sheet
func (c *Client) writeData(ctx context.Context, spreadsheetID, sheetName string, data [][]string) error {
	// Convert [][]string to [][]interface{} for the API
	var values [][]interface{}
	for _, row := range data {
		interfaceRow := make([]interface{}, len(row))
		for i, cell := range row {
			interfaceRow[i] = cell
		}
		values = append(values, interfaceRow)
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	rangeStr := fmt.Sprintf("%s!A1", sheetName)
	_, err := c.service.Spreadsheets.Values.Update(
		spreadsheetID,
		rangeStr,
		valueRange,
	).ValueInputOption("RAW").Context(ctx).Do()

	if err != nil {
		return err
	}

	return nil
}

// generateDefaultTitle generates a default title using the current timestamp
func generateDefaultTitle() string {
	now := time.Now()
	return now.Format("20060102150405") + "+gs"
}
