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

// TODO: move into seperate file
func renderOrderBook(w *models.Window) {
	orderBook, ok := w.Data.(*models.OrderBookResult)
	if !ok || orderBook == nil {
		// Render placeholder if no data available
		rl.DrawTextEx(
			w.Font,
			"Loading order book...",
			rl.Vector2{X: w.Rect.X + w.Padding, Y: w.Rect.Y + 35},
			16,
			1,
			colors.ColorSubtext)
		return
	}

	// Use fixed-width columns for terminal style appearance
	numAsks := len(orderBook.Asks)
	if numAsks > 20 {
		numAsks = 20
	}

	numBids := len(orderBook.Bids)
	if numBids > 20 {
		numBids = 20
	}

	contentHeight := float32(45 + (numAsks+numBids)*20 + 30) // Header + rows + spread row
	w.MaxScroll = contentHeight - (w.Rect.Height - 30)
	if w.MaxScroll < 0 {
		w.MaxScroll = 0
	}

	startY := w.Rect.Y + 35 - w.ScrollPosition

	// Draw column headers with monospaced font (terminal style)
	headerColor := colors.ColorSubtext
	rowSpacing := float32(20) // Space between rows

	// Column positions - consistent with terminal style fixed-width columns
	priceX := w.Rect.X + w.Padding
	amountX := w.Rect.X + w.Padding + 120
	totalX := w.Rect.X + w.Padding + 220

	// Draw column headers
	rl.DrawTextEx(w.Font, "Price", rl.Vector2{X: priceX, Y: startY}, 16, 1, headerColor)
	rl.DrawTextEx(w.Font, "Amount", rl.Vector2{X: amountX, Y: startY}, 16, 1, headerColor)
	rl.DrawTextEx(w.Font, "Total", rl.Vector2{X: totalX, Y: startY}, 16, 1, headerColor)
	startY += rowSpacing + 5

	// Calculate max volume for visualization (subtle volume bars like in image 2)
	maxVolume := 0.0
	for _, ask := range orderBook.Asks {
		if len(ask) >= 2 && ask[1] > maxVolume {
			maxVolume = ask[1]
		}
	}
	for _, bid := range orderBook.Bids {
		if len(bid) >= 2 && bid[1] > maxVolume {
			maxVolume = bid[1]
		}
	}

	// Display asks (from lowest to highest) - RED
	totalAsks := 0.0

	for i := numAsks - 1; i >= 0; i-- {
		ask := orderBook.Asks[i]
		if len(ask) >= 2 {
			price := ask[0]
			amount := ask[1]
			totalAsks += amount

			// Draw subtle volume bar (similar to image 2)
			barWidth := (amount / maxVolume) * float64((w.Rect.Width - 50 - w.Padding*2))
			barRect := rl.Rectangle{
				X:      w.Rect.X + w.Rect.Width - float32(barWidth) - w.Padding,
				Y:      startY,
				Width:  float32(barWidth),
				Height: 16,
			}
			rl.DrawRectangleRec(barRect, rl.NewColor(229, 78, 103, 40)) // Very subtle red background

			// Draw text with monospaced font
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.2f", price), rl.Vector2{X: priceX, Y: startY}, 16, 1, colors.ColorRed)
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.4f", amount), rl.Vector2{X: amountX, Y: startY}, 16, 1, colors.ColorRed)
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.4f", totalAsks), rl.Vector2{X: totalX, Y: startY}, 16, 1, colors.ColorRed)

			startY += rowSpacing
		}
	}

	// Spread row - more subtle like in image 2
	spread := 0.0
	spreadPct := 0.0

	if len(orderBook.Asks) > 0 && len(orderBook.Bids) > 0 {
		// Calculate spread if we have both asks and bids
		bestAskPrice := orderBook.Asks[0][0]
		bestBidPrice := orderBook.Bids[0][0]
		spread = bestAskPrice - bestBidPrice
		spreadPct = (spread / bestBidPrice) * 100
	}

	// Draw spread info with a subtle background
	rl.DrawRectangle(
		int32(w.Rect.X+w.Padding),
		int32(startY),
		int32(w.Rect.Width-w.Padding*2),
		20,
		rl.NewColor(40, 44, 52, 100)) // Very subtle background

	rl.DrawTextEx(
		w.Font,
		fmt.Sprintf("Spread: %.2f (%.4f%%)", spread, spreadPct),
		rl.Vector2{X: w.Rect.X + w.Padding + 5, Y: startY + 2},
		16,
		1,
		colors.ColorText)
	startY += rowSpacing + 5

	// Draw bids (buy orders) - GREEN (similar to image 2)
	totalBids := 0.0

	for i := 0; i < numBids; i++ {
		bid := orderBook.Bids[i]
		if len(bid) >= 2 {
			price := bid[0]
			amount := bid[1]
			totalBids += amount

			// Draw subtle volume bar
			barWidth := (amount / maxVolume) * float64((w.Rect.Width - 50 - w.Padding*2))
			barRect := rl.Rectangle{
				X:      w.Rect.X + w.Rect.Width - float32(barWidth) - w.Padding,
				Y:      startY,
				Width:  float32(barWidth),
				Height: 16,
			}
			rl.DrawRectangleRec(barRect, rl.NewColor(75, 201, 155, 40)) // Very subtle green background

			// Draw text with monospaced font
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.2f", price), rl.Vector2{X: priceX, Y: startY}, 16, 1, colors.ColorGreen)
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.4f", amount), rl.Vector2{X: amountX, Y: startY}, 16, 1, colors.ColorGreen)
			rl.DrawTextEx(w.Font, fmt.Sprintf("%.4f", totalBids), rl.Vector2{X: totalX, Y: startY}, 16, 1, colors.ColorGreen)

			startY += rowSpacing
		}
	}
}

func renderRecentTrades(w *models.Window) {
	// Trade data could come from websocket similar to order book
	trades, ok := w.Data.([]models.Trade)
	if !ok || trades == nil {
		rl.DrawTextEx(
			w.Font,
			"No recent trades data...",
			rl.Vector2{X: w.Rect.X + w.Padding, Y: w.Rect.Y + 35},
			16,
			1,
			colors.ColorSubtext)
		return
	}

	startY := w.Rect.Y + 35 - w.ScrollPosition
	rowSpacing := float32(20)

	// Draw column headers
	priceX := w.Rect.X + w.Padding
	amountX := w.Rect.X + w.Padding + 120
	timeX := w.Rect.X + w.Padding + 220

	rl.DrawTextEx(w.Font, "Price", rl.Vector2{X: priceX, Y: startY}, 16, 1, colors.ColorSubtext)
	rl.DrawTextEx(w.Font, "Amount", rl.Vector2{X: amountX, Y: startY}, 16, 1, colors.ColorSubtext)
	rl.DrawTextEx(w.Font, "Time", rl.Vector2{X: timeX, Y: startY}, 16, 1, colors.ColorSubtext)
	startY += rowSpacing + 5

	// Display trades
	for _, trade := range trades {
		// Color based on trade direction
		textColor := colors.ColorGreen
		if trade.Direction == "sell" {
			textColor = colors.ColorRed
		}

		rl.DrawTextEx(w.Font, fmt.Sprintf("%.2f", trade.Price),
			rl.Vector2{X: priceX, Y: startY}, 16, 1, textColor)
		rl.DrawTextEx(w.Font, fmt.Sprintf("%.4f", trade.Amount),
			rl.Vector2{X: amountX, Y: startY}, 16, 1, textColor)
		rl.DrawTextEx(w.Font, trade.Timestamp,
			rl.Vector2{X: timeX, Y: startY}, 16, 1, colors.ColorSubtext)

		startY += rowSpacing
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1200, 800, "Go Trader")
	rl.SetTargetFPS(60)

	font := rl.LoadFont("JetBrainsMono-Regular.ttf")
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	// Create main windows with improved spacing
	orderBookWindow := NewWindow("deribit btcusdt - Orderbook", 10, 85, 580, 400, renderOrderBook, font)
	orderBookWindow.IsActive = true

	// Windows list
	windows := []*models.Window{
		orderBookWindow,
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
