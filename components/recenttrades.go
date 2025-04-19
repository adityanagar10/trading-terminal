package components

import (
	"fmt"

	colors "github.com/adityanagar10/trader/constants"
	"github.com/adityanagar10/trader/models"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func RenderRecentTrades(w *models.Window) {
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
