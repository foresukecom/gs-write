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
func (c *Client) CreateSpreadsheet(ctx context.Context, title string, data [][]string, freezeRows, freezeCols, filterHeaderRow int) (string, error) {
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
	sheetID := resp.Sheets[0].Properties.SheetId

	// Write data to the spreadsheet
	if len(data) > 0 {
		if err := c.writeData(ctx, spreadsheetID, "Sheet1", data); err != nil {
			return "", fmt.Errorf("failed to write data: %w", err)
		}
	}

	// Apply freeze panes if specified
	if freezeRows > 0 || freezeCols > 0 {
		if err := c.setFreezePanes(ctx, spreadsheetID, sheetID, freezeRows, freezeCols); err != nil {
			return "", fmt.Errorf("failed to set freeze panes: %w", err)
		}
	}

	// Apply basic filter if specified
	if filterHeaderRow > 0 {
		if err := c.setBasicFilter(ctx, spreadsheetID, sheetID, filterHeaderRow, len(data), len(data[0])); err != nil {
			return "", fmt.Errorf("failed to set basic filter: %w", err)
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

// setFreezePanes sets frozen rows and columns for the sheet
func (c *Client) setFreezePanes(ctx context.Context, spreadsheetID string, sheetID int64, freezeRows, freezeCols int) error {
	gridProperties := &sheets.GridProperties{}

	if freezeRows > 0 {
		gridProperties.FrozenRowCount = int64(freezeRows)
	}

	if freezeCols > 0 {
		gridProperties.FrozenColumnCount = int64(freezeCols)
	}

	requests := []*sheets.Request{
		{
			UpdateSheetProperties: &sheets.UpdateSheetPropertiesRequest{
				Properties: &sheets.SheetProperties{
					SheetId:        sheetID,
					GridProperties: gridProperties,
				},
				Fields: "gridProperties.frozenRowCount,gridProperties.frozenColumnCount",
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err := c.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

// setBasicFilter sets a basic filter for the sheet
func (c *Client) setBasicFilter(ctx context.Context, spreadsheetID string, sheetID int64, headerRow, numRows, numCols int) error {
	// Basic filter range starts from the header row (0-indexed)
	// and spans all columns and rows from header to end
	filterRange := &sheets.GridRange{
		SheetId:          sheetID,
		StartRowIndex:    int64(headerRow - 1), // Convert to 0-indexed
		EndRowIndex:      int64(numRows),
		StartColumnIndex: 0,
		EndColumnIndex:   int64(numCols),
	}

	requests := []*sheets.Request{
		{
			SetBasicFilter: &sheets.SetBasicFilterRequest{
				Filter: &sheets.BasicFilter{
					Range: filterRange,
				},
			},
		},
	}

	batchUpdateRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}

	_, err := c.service.Spreadsheets.BatchUpdate(spreadsheetID, batchUpdateRequest).Context(ctx).Do()
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
