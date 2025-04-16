package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	colors "github.com/adityanagar10/trader/constants"
	models "github.com/adityanagar10/trader/models"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gorilla/websocket"
)

// Window struct for UI
type Window struct {
	title      string
	rect       rl.Rectangle
	isDragging bool
	dragOffset rl.Vector2
	content    func(*Window) // Function to render window content
	data       interface{}   // Generic data field for window content
	isActive   bool          // Whether this window is currently active/focused
}

func NewWindow(title string, x, y, width, height float32, contentFunc func(*Window)) *Window {
	return &Window{
		title:    title,
		rect:     rl.NewRectangle(x, y, width, height),
		content:  contentFunc,
		isActive: false,
	}
}

func (w *Window) Update(windows []*Window) {
	mousePos := rl.GetMousePosition()

	// Header/title bar area
	headerRect := rl.Rectangle{
		X:      w.rect.X,
		Y:      w.rect.Y,
		Width:  w.rect.Width,
		Height: 30,
	}

	// Check if clicked on this window to make it active
	if rl.CheckCollisionPointRec(mousePos, w.rect) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		// Make this window active and all others inactive
		for _, win := range windows {
			win.isActive = false
		}
		w.isActive = true
	}

	// Dragging logic
	if rl.CheckCollisionPointRec(mousePos, headerRect) {
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			w.isDragging = true
			w.dragOffset = rl.Vector2{
				X: mousePos.X - w.rect.X,
				Y: mousePos.Y - w.rect.Y,
			}
			// Make this window active and all others inactive
			for _, win := range windows {
				win.isActive = false
			}
			w.isActive = true
		}
	}

	if w.isDragging {
		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			w.rect.X = mousePos.X - w.dragOffset.X
			w.rect.Y = mousePos.Y - w.dragOffset.Y
		} else {
			w.isDragging = false
		}
	}
}

func (w *Window) Draw() {
	// Draw window background
	rl.DrawRectangleRec(w.rect, colors.ColorPanelBg)

	// Draw window border with different color if active
	borderColor := colors.ColorBorder
	if w.isActive {
		borderColor = colors.ColorHighlight
	}
	rl.DrawRectangleLinesEx(w.rect, 1, borderColor)

	// Draw header
	headerRect := rl.Rectangle{
		X:      w.rect.X,
		Y:      w.rect.Y,
		Width:  w.rect.Width,
		Height: 30,
	}
	rl.DrawRectangleRec(headerRect, colors.ColorHeaderBg)
	rl.DrawRectangleLinesEx(w.rect, 1, borderColor)

	// Draw header bottom border line
	rl.DrawLine(
		int32(headerRect.X),
		int32(headerRect.Y+headerRect.Height),
		int32(headerRect.X+headerRect.Width),
		int32(headerRect.Y+headerRect.Height),
		borderColor)

	// Draw title text
	rl.DrawText(
		w.title,
		int32(w.rect.X+10),
		int32(w.rect.Y+8),
		18,
		colors.ColorText,
	)

	// Draw content
	if w.content != nil {
		// Set clipping to ensure content stays within window
		rl.BeginScissorMode(
			int32(w.rect.X),
			int32(w.rect.Y+30),
			int32(w.rect.Width),
			int32(w.rect.Height-30),
		)
		w.content(w)
		rl.EndScissorMode()
	}
}

// Render order book with real data in the Market Monkey style
func renderOrderBook(w *Window) {
	orderBook, ok := w.data.(*models.OrderBookResult)
	if !ok || orderBook == nil {
		// Render placeholder if no data available
		rl.DrawText("Loading order book...",
			int32(w.rect.X+10),
			int32(w.rect.Y+40),
			18,
			colors.ColorSubtext)
		return
	}

	startY := w.rect.Y + 40

	// // Draw instrument name
	// instrumentText := fmt.Sprintf("%s", orderBook.InstrumentName)
	// rl.DrawText(instrumentText,
	// 	int32(w.rect.X+10),
	// 	int32(startY),
	// 	18,
	// 	colorHighlight)
	startY += 30

	// Draw price information in a more compact format
	priceChangeColor := colors.ColorSubtext
	if orderBook.Stats.PriceChange < 0 {
		priceChangeColor = colors.ColorRed
	} else if orderBook.Stats.PriceChange > 0 {
		priceChangeColor = colors.ColorGreen
	}

	// Last price row
	rl.DrawText("Last:",
		int32(w.rect.X+10),
		int32(startY),
		16,
		colors.ColorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f", orderBook.LastPrice),
		int32(w.rect.X+70),
		int32(startY),
		16,
		colors.ColorText)

	// Mark price
	rl.DrawText("Mark:",
		int32(w.rect.X+150),
		int32(startY),
		16,
		colors.ColorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f", orderBook.MarkPrice),
		int32(w.rect.X+210),
		int32(startY),
		16,
		colors.ColorText)
	startY += 22

	// 24h change row
	rl.DrawText("24h:",
		int32(w.rect.X+10),
		int32(startY),
		16,
		colors.ColorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f%%", orderBook.Stats.PriceChange),
		int32(w.rect.X+70),
		int32(startY),
		16,
		priceChangeColor)

	// Funding
	fundingColor := colors.ColorSubtext
	if orderBook.Funding8h < 0 {
		fundingColor = colors.ColorGreen
	} else if orderBook.Funding8h > 0 {
		fundingColor = colors.ColorRed
	}

	rl.DrawText("Funding:",
		int32(w.rect.X+150),
		int32(startY),
		16,
		colors.ColorSubtext)
	rl.DrawText(fmt.Sprintf("%.4f%%", orderBook.Funding8h*100),
		int32(w.rect.X+210),
		int32(startY),
		16,
		fundingColor)
	startY += 30

	// Divider line
	rl.DrawLine(
		int32(w.rect.X+10),
		int32(startY),
		int32(w.rect.X+w.rect.Width-10),
		int32(startY),
		colors.ColorBorder)
	startY += 15

	// Draw order book headers
	rl.DrawText("Price", int32(w.rect.X+10), int32(startY), 16, colors.ColorSubtext)
	rl.DrawText("Amount", int32(w.rect.X+120), int32(startY), 16, colors.ColorSubtext)
	rl.DrawText("Total", int32(w.rect.X+220), int32(startY), 16, colors.ColorSubtext)
	startY += 25

	// Calculate max volume for visualization
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

	// Display asks (from lowest to highest)
	maxDisplayItems := 12
	totalAsks := 0.0

	numAsks := len(orderBook.Asks)
	if numAsks > maxDisplayItems {
		numAsks = maxDisplayItems
	}

	for i := numAsks - 1; i >= 0; i-- {
		ask := orderBook.Asks[i]
		if len(ask) >= 2 {
			price := ask[0]
			amount := ask[1]
			totalAsks += amount

			// Draw volume visualization bar
			barWidth := (amount / maxVolume) * 150
			barRect := rl.Rectangle{
				X:      w.rect.X + w.rect.Width - 10 - float32(barWidth),
				Y:      startY - 2,
				Width:  float32(barWidth),
				Height: 18,
			}
			rl.DrawRectangleRec(barRect, rl.NewColor(229, 78, 103, 50)) // Semi-transparent red

			rl.DrawText(
				fmt.Sprintf("%.2f", price),
				int32(w.rect.X+10),
				int32(startY),
				16,
				colors.ColorRed,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", amount),
				int32(w.rect.X+120),
				int32(startY),
				16,
				colors.ColorRed,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", totalAsks),
				int32(w.rect.X+220),
				int32(startY),
				16,
				colors.ColorRed,
			)
			startY += 20
		}
	}

	// Spread row
	spread := orderBook.BestAskPrice - orderBook.BestBidPrice
	spreadPct := (spread / orderBook.BestBidPrice) * 100

	rl.DrawRectangle(
		int32(w.rect.X+10),
		int32(startY),
		int32(w.rect.Width-20),
		25,
		rl.NewColor(40, 44, 52, 255))

	rl.DrawText(
		fmt.Sprintf("Spread: %.2f (%.4f%%)", spread, spreadPct),
		int32(w.rect.X+10),
		int32(startY+4),
		16,
		colors.ColorText,
	)
	startY += 30

	// Draw bids (buy orders) - green
	totalBids := 0.0

	// Display bids (from highest to lowest)
	numBids := len(orderBook.Bids)
	if numBids > maxDisplayItems {
		numBids = maxDisplayItems
	}

	for i := 0; i < numBids; i++ {
		bid := orderBook.Bids[i]
		if len(bid) >= 2 {
			price := bid[0]
			amount := bid[1]
			totalBids += amount

			// Draw volume visualization bar
			barWidth := (amount / maxVolume) * 150
			barRect := rl.Rectangle{
				X:      w.rect.X + w.rect.Width - 10 - float32(barWidth),
				Y:      startY - 2,
				Width:  float32(barWidth),
				Height: 18,
			}
			rl.DrawRectangleRec(barRect, rl.NewColor(75, 201, 155, 50)) // Semi-transparent green

			rl.DrawText(
				fmt.Sprintf("%.2f", price),
				int32(w.rect.X+10),
				int32(startY),
				16,
				colors.ColorGreen,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", amount),
				int32(w.rect.X+120),
				int32(startY),
				16,
				colors.ColorGreen,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", totalBids),
				int32(w.rect.X+220),
				int32(startY),
				16,
				colors.ColorGreen,
			)
			startY += 20
		}
	}
}

// DeribitClient handles WebSocket communication with Deribit
type DeribitClient struct {
	conn            *websocket.Conn
	orderBookWindow *Window
	instrument      string
	requestID       int
}

func NewDeribitClient(instrument string, orderBookWindow *Window) *DeribitClient {
	return &DeribitClient{
		orderBookWindow: orderBookWindow,
		instrument:      instrument,
		requestID:       1,
	}
}

func (c *DeribitClient) Connect() error {
	// Connect to Deribit WebSocket API
	conn, _, err := websocket.DefaultDialer.Dial("wss://www.deribit.com/ws/api/v2", nil)
	if err != nil {
		return fmt.Errorf("websocket connection error: %v", err)
	}

	c.conn = conn
	log.Println("Connected to Deribit WebSocket API")

	// Start listening for messages
	go c.handleMessages()

	// Start fetching order book periodically
	go c.fetchOrderBookPeriodically()

	return nil
}

func (c *DeribitClient) Close() {
	if c.conn != nil {
		c.conn.Close()
		log.Println("Closed Deribit WebSocket connection")
	}
}

func (c *DeribitClient) handleMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			return
		}

		// Parse the response
		var response models.DeribitResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Failed to unmarshal response: %v", err)
			continue
		}

		// Handle errors
		if response.Error != nil {
			log.Printf("Deribit API error: %d - %s", response.Error.Code, response.Error.Message)
			continue
		}

		// Process order book responses
		if response.Result != nil {
			// Update window data directly with the complete order book result
			c.orderBookWindow.data = response.Result
		}
	}
}

func (c *DeribitClient) fetchOrderBook() {
	// Create an order book request
	request := models.DeribitRequest{
		JsonRPC: "2.0",
		ID:      c.requestID,
		Method:  "public/get_order_book",
		Params: models.OrderBookParams{
			InstrumentName: c.instrument,
		},
	}
	c.requestID++

	// Marshal and send the request
	data, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to send request: %v", err)
		return
	}
}

func (c *DeribitClient) fetchOrderBookPeriodically() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		c.fetchOrderBook()
	}
}

// UI components
func drawTopBar() {
	// Draw top navigation bar
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), 30, colors.ColorHeaderBg)
	rl.DrawLine(0, 30, int32(rl.GetScreenWidth()), 30, colors.ColorBorder)

	// Draw navigation buttons
	navItems := []string{"Orderbooks"}
	xPos := 20

	for _, item := range navItems {
		rl.DrawText(item, int32(xPos), 8, 18, colors.ColorText)
		xPos += int(rl.MeasureText(item, 18) + 30)
	}

	// Draw title on the right
	title := "Go Trader"
	titleWidth := rl.MeasureText(title, 18)
	rl.DrawText(title, int32(rl.GetScreenWidth())-int32(titleWidth)-20, 8, 18, colors.ColorText)
}

func drawStatusBar() {
	// Draw status bar at the bottom
	statusBarHeight := int32(25)
	rl.DrawRectangle(
		0,
		int32(rl.GetScreenHeight())-statusBarHeight,
		int32(rl.GetScreenWidth()),
		statusBarHeight,
		colors.ColorHeaderBg,
	)

	// Draw status text
	statusText := "Connected to Deribit | v0.1.0"
	rl.DrawText(
		statusText,
		10,
		int32(rl.GetScreenHeight())-statusBarHeight+5,
		16,
		colors.ColorSubtext,
	)

	// Draw timestamp on the right
	timeText := time.Now().Format("15:04:05")
	timeWidth := rl.MeasureText(timeText, 16)
	rl.DrawText(
		timeText,
		int32(rl.GetScreenWidth())-int32(timeWidth)-10,
		int32(rl.GetScreenHeight())-statusBarHeight+5,
		16,
		colors.ColorSubtext,
	)
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(500, 500, "Go Trader")
	rl.SetTargetFPS(60)

	// Set custom font for a more modern look
	// You'll need to download and include a suitable monospaced font
	// For example, "JetBrainsMono-Regular.ttf" would be good for this
	font := rl.LoadFont("JetBrainsMono-Regular.ttf")
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	// Create order book window
	orderBookWindow := NewWindow("Orderbook: BTC-PERPETUAL", 350, 50, 300, 600, renderOrderBook)
	orderBookWindow.isActive = true

	// Create empty windows for chart and depth
	// chartWindow := NewWindow("BTC-PERPETUAL Chart", 50, 50, 280, 400, nil)
	// depthWindow := NewWindow("BTC-PERPETUAL Depth", 680, 50, 280, 400, nil)

	// Create Deribit client
	client := NewDeribitClient("BTC-PERPETUAL", orderBookWindow)
	err := client.Connect()
	if err != nil {
		log.Printf("Failed to connect: %v", err)
	}
	defer client.Close()

	// Windows list
	windows := []*Window{
		orderBookWindow,
	}

	for !rl.WindowShouldClose() {
		// Update
		for _, win := range windows {
			win.Update(windows)
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(colors.ColorBackground)

		// Draw UI elements
		drawTopBar()
		drawStatusBar()

		// Draw all windows
		for _, win := range windows {
			win.Draw()
		}

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
