package main

import (
	"fmt"
	"log"
	"time"

	"github.com/adityanagar10/trader/client"
	"github.com/adityanagar10/trader/components"
	colors "github.com/adityanagar10/trader/constants"
	"github.com/adityanagar10/trader/models"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func NewWindow(title string, x, y, width, height float32, contentFunc func(*models.Window), font rl.Font) *models.Window {
	return &models.Window{
		Title:            title,
		Rect:             rl.NewRectangle(x, y, width, height),
		Content:          contentFunc,
		IsActive:         false,
		ScrollPosition:   0,
		IsResizing:       false,
		ResizeDir:        0,
		ResizeHandleSize: 10,
		Padding:          12,
		Font:             font,
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1200, 800, "Go Trader")
	rl.SetTargetFPS(60)

	font := rl.LoadFont("JetBrainsMono-Regular.ttf")
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	orderBookWindow := NewWindow("deribit btcusdt - Orderbook", 510, 85, 580, 400, components.RenderOrderBook, font)
	orderBookWindow.IsActive = true

	intrumentGraph := NewWindow("derbit btcusdt - Graph", 10, 85, 500, 500, components.RenderIntrumentGraph, font)

	// Windows list
	windows := []*models.Window{
		orderBookWindow,
		intrumentGraph,
	}

	// Create dropdown for instrument selection
	instrumentDropdown := components.NewDropdown(
		10, 50, 200,
		[]string{"BTC-PERPETUAL", "ETH-PERPETUAL", "SOL-PERPETUAL", "XRP-PERPETUAL"},
		"Instrument",
		font)

	// Create Deribit client and connect
	deribitClient := client.NewDeribitClient("BTC-PERPETUAL", orderBookWindow)
	err := deribitClient.Connect()
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}
	defer deribitClient.Close()

	// Set handler for instrument change
	instrumentDropdown.SetOnChangeHandler(func(idx int) {
		selectedInstrument := instrumentDropdown.GetSelectedOption()

		// TODO: move to constants
		var symbol string
		switch selectedInstrument {
		case "BTC-PERPETUAL":
			symbol = "btcusdt"
		case "ETH-PERPETUAL":
			symbol = "ethusdt"
		case "SOL-PERPETUAL":
			symbol = "solusdt"
		case "XRP-PERPETUAL":
			symbol = "xrpusdt"
		}

		orderBookWindow.Title = fmt.Sprintf("deribit %s - Orderbook", symbol)

		// Connect to appropriate WebSocket for selected instrument
		if deribitClient != nil {
			deribitClient.Instrument = selectedInstrument
			// Clear data while loading
			orderBookWindow.Data = nil
		}
		fmt.Printf("Switched to instrument: %s\n", selectedInstrument)
	})

	// Main loop
	for !rl.WindowShouldClose() {
		// Update
		for _, win := range windows {
			win.Update(windows)
		}
		instrumentDropdown.Update()

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(colors.ColorBackground)

		rl.DrawText("Go Trader", 10, 10, 20, colors.ColorText)

		// Draw instrument dropdown
		instrumentDropdown.Draw()

		// Draw all windows
		for _, win := range windows {
			win.Draw()
		}

		statusY := rl.GetScreenHeight() - 20
		rl.DrawText("Connected to Deribit", 10, int32(statusY), 16, colors.ColorSubtext)

		// Draw timestamp on the right
		timeText := time.Now().Format("15:04:05")
		timeWidth := rl.MeasureText(timeText, 16)
		rl.DrawText(timeText, int32(rl.GetScreenWidth())-int32(timeWidth)-10, int32(statusY), 16, colors.ColorSubtext)

		rl.EndDrawing()
	}

	rl.UnloadFont(font)
	rl.CloseWindow()
}
