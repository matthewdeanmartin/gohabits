package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	// One single internal import
	"github.com/yourname/habits/internal/app"
)

func main() {
	// 1. Load Config
	cfg, err := app.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Sheets Client
	ctx := context.Background()
	client, err := app.NewSheetClient(ctx, cfg.Auth.KeyPath, cfg.SpreadsheetID, cfg.SheetName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to sheets: %v\n", err)
		os.Exit(1)
	}

	// 3. Start TUI
	p := tea.NewProgram(app.NewModel(cfg, client))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
