package app

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetClient struct {
	srv           *sheets.Service
	spreadsheetID string
	sheetName     string
	headerMap     map[string]int
}

func NewSheetClient(ctx context.Context, credentialsPath, spreadsheetID, sheetName string) (*SheetClient, error) {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("read credentials: %w", err)
	}

	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(data))
	if err != nil {
		return nil, fmt.Errorf("init sheets service: %w", err)
	}

	return &SheetClient{
		srv:           srv,
		spreadsheetID: spreadsheetID,
		sheetName:     sheetName,
	}, nil
}

func (c *SheetClient) FetchHeaders() ([]string, error) {
	readRange := fmt.Sprintf("%s!1:1", c.sheetName)
	resp, err := c.srv.Spreadsheets.Values.Get(c.spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("sheet %s is empty", c.sheetName)
	}

	headers := make([]string, 0)
	c.headerMap = make(map[string]int)

	for i, v := range resp.Values[0] {
		h := fmt.Sprintf("%v", v)
		headers = append(headers, h)
		c.headerMap[h] = i
	}

	return headers, nil
}

func (c *SheetClient) FetchExistingDates() (map[string]bool, error) {
	dateColIdx, ok := c.headerMap["date"]
	if !ok {
		return nil, fmt.Errorf("required column 'date' missing in sheet")
	}

	colLetter := string(rune('A' + dateColIdx))
	readRange := fmt.Sprintf("%s!%s2:%s", c.sheetName, colLetter, colLetter)

	resp, err := c.srv.Spreadsheets.Values.Get(c.spreadsheetID, readRange).Do()
	if err != nil {
		return nil, err
	}

	exists := make(map[string]bool)
	for _, row := range resp.Values {
		if len(row) > 0 {
			dateStr := fmt.Sprintf("%v", row[0])
			exists[dateStr] = true
		}
	}
	return exists, nil
}

func (c *SheetClient) AppendRow(data map[string]interface{}) error {
	maxIdx := 0
	for _, idx := range c.headerMap {
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	row := make([]interface{}, maxIdx+1)
	for header, val := range data {
		if idx, ok := c.headerMap[header]; ok {
			row[idx] = val
		}
	}

	vr := &sheets.ValueRange{
		Values: [][]interface{}{row},
	}

	_, err := c.srv.Spreadsheets.Values.Append(c.spreadsheetID, c.sheetName, vr).
		ValueInputOption("USER_ENTERED").Do()

	return err
}
