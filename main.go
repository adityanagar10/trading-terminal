package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/gorilla/websocket"
)

// Color theme based on the Market Monkey screenshot
var (
	colorBackground = rl.NewColor(21, 25, 30, 255)    // Dark background
	colorPanelBg    = rl.NewColor(28, 32, 38, 255)    // Slightly lighter panel background
	colorHeaderBg   = rl.NewColor(32, 36, 43, 255)    // Header background
	colorBorder     = rl.NewColor(45, 49, 55, 255)    // Border color
	colorText       = rl.NewColor(210, 210, 210, 255) // Main text color
	colorSubtext    = rl.NewColor(140, 145, 155, 255) // Subtitle/label text
	colorGreen      = rl.NewColor(75, 201, 155, 255)  // Green for positive/buy
	colorRed        = rl.NewColor(229, 78, 103, 255)  // Red for negative/sell
	colorHighlight  = rl.NewColor(86, 180, 233, 255)  // Highlight blue
	colorChartBg    = rl.NewColor(25, 29, 34, 255)    // Chart background
)

// Deribit API structures
type DeribitRequest struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type OrderBookParams struct {
	InstrumentName string `json:"instrument_name"`
}

type DeribitResponse struct {
	JsonRPC string           `json:"jsonrpc"`
	ID      int              `json:"id"`
	Result  *OrderBookResult `json:"result,omitempty"`
	Error   *DeribitError    `json:"error,omitempty"`
	UsIn    int64            `json:"usIn,omitempty"`
	UsOut   int64            `json:"usOut,omitempty"`
	UsDiff  int              `json:"usDiff,omitempty"`
	Testnet bool             `json:"testnet,omitempty"`
}

type DeribitError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type OrderBookStats struct {
	High           float64 `json:"high"`
	Low            float64 `json:"low"`
	PriceChange    float64 `json:"price_change"`
	Volume         float64 `json:"volume"`
	VolumeUSD      float64 `json:"volume_usd"`
	VolumeNotional float64 `json:"volume_notional"`
}

type OrderBookResult struct {
	Timestamp        int64          `json:"timestamp"`
	State            string         `json:"state"`
	Stats            OrderBookStats `json:"stats"`
	ChangeID         int64          `json:"change_id"`
	IndexPrice       float64        `json:"index_price"`
	InstrumentName   string         `json:"instrument_name"`
	Bids             [][]float64    `json:"bids"`
	Asks             [][]float64    `json:"asks"`
	LastPrice        float64        `json:"last_price"`
	SettlementPrice  float64        `json:"settlement_price"`
	MinPrice         float64        `json:"min_price"`
	MaxPrice         float64        `json:"max_price"`
	OpenInterest     float64        `json:"open_interest"`
	MarkPrice        float64        `json:"mark_price"`
	InterestValue    float64        `json:"interest_value"`
	BestAskPrice     float64        `json:"best_ask_price"`
	BestBidPrice     float64        `json:"best_bid_price"`
	EstDeliveryPrice float64        `json:"estimated_delivery_price"`
	BestAskAmount    float64        `json:"best_ask_amount"`
	BestBidAmount    float64        `json:"best_bid_amount"`
	CurrentFunding   float64        `json:"current_funding"`
	Funding8h        float64        `json:"funding_8h"`
}

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
	rl.DrawRectangleRec(w.rect, colorPanelBg)

	// Draw window border with different color if active
	borderColor := colorBorder
	if w.isActive {
		borderColor = colorHighlight
	}
	rl.DrawRectangleLinesEx(w.rect, 1, borderColor)

	// Draw header
	headerRect := rl.Rectangle{
		X:      w.rect.X,
		Y:      w.rect.Y,
		Width:  w.rect.Width,
		Height: 30,
	}
	rl.DrawRectangleRec(headerRect, colorHeaderBg)

	// Draw title text
	rl.DrawText(
		w.title,
		int32(w.rect.X+10),
		int32(w.rect.Y+8),
		18,
		colorText,
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
	orderBook, ok := w.data.(*OrderBookResult)
	if !ok || orderBook == nil {
		// Render placeholder if no data available
		rl.DrawText("Loading order book...",
			int32(w.rect.X+10),
			int32(w.rect.Y+40),
			18,
			colorSubtext)
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
	priceChangeColor := colorSubtext
	if orderBook.Stats.PriceChange < 0 {
		priceChangeColor = colorRed
	} else if orderBook.Stats.PriceChange > 0 {
		priceChangeColor = colorGreen
	}

	// Last price row
	rl.DrawText("Last:",
		int32(w.rect.X+10),
		int32(startY),
		16,
		colorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f", orderBook.LastPrice),
		int32(w.rect.X+70),
		int32(startY),
		16,
		colorText)

	// Mark price
	rl.DrawText("Mark:",
		int32(w.rect.X+150),
		int32(startY),
		16,
		colorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f", orderBook.MarkPrice),
		int32(w.rect.X+210),
		int32(startY),
		16,
		colorText)
	startY += 22

	// 24h change row
	rl.DrawText("24h:",
		int32(w.rect.X+10),
		int32(startY),
		16,
		colorSubtext)
	rl.DrawText(fmt.Sprintf("%.2f%%", orderBook.Stats.PriceChange),
		int32(w.rect.X+70),
		int32(startY),
		16,
		priceChangeColor)

	// Funding
	fundingColor := colorSubtext
	if orderBook.Funding8h < 0 {
		fundingColor = colorGreen
	} else if orderBook.Funding8h > 0 {
		fundingColor = colorRed
	}

	rl.DrawText("Funding:",
		int32(w.rect.X+150),
		int32(startY),
		16,
		colorSubtext)
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
		colorBorder)
	startY += 15

	// Draw order book headers
	rl.DrawText("Price", int32(w.rect.X+10), int32(startY), 16, colorSubtext)
	rl.DrawText("Amount", int32(w.rect.X+120), int32(startY), 16, colorSubtext)
	rl.DrawText("Total", int32(w.rect.X+220), int32(startY), 16, colorSubtext)
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
				colorRed,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", amount),
				int32(w.rect.X+120),
				int32(startY),
				16,
				colorRed,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", totalAsks),
				int32(w.rect.X+220),
				int32(startY),
				16,
				colorRed,
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
		colorText,
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
				colorGreen,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", amount),
				int32(w.rect.X+120),
				int32(startY),
				16,
				colorGreen,
			)
			rl.DrawText(
				fmt.Sprintf("%.0f", totalBids),
				int32(w.rect.X+220),
				int32(startY),
				16,
				colorGreen,
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
		var response DeribitResponse
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
	request := DeribitRequest{
		JsonRPC: "2.0",
		ID:      c.requestID,
		Method:  "public/get_order_book",
		Params: OrderBookParams{
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
	rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), 30, colorHeaderBg)
	rl.DrawLine(0, 30, int32(rl.GetScreenWidth()), 30, colorBorder)

	// Draw navigation buttons

	// Draw title on the right
	title := "Go Trader"
	titleWidth := rl.MeasureText(title, 18)
	rl.DrawText(title, int32(rl.GetScreenWidth())-int32(titleWidth)-20, 8, 18, colorText)
}

func drawStatusBar() {
	// Draw status bar at the bottom
	statusBarHeight := int32(25)
	rl.DrawRectangle(
		0,
		int32(rl.GetScreenHeight())-statusBarHeight,
		int32(rl.GetScreenWidth()),
		statusBarHeight,
		colorHeaderBg,
	)

	// Draw status text
	statusText := "Connected to Deribit | v0.1.0"
	rl.DrawText(
		statusText,
		10,
		int32(rl.GetScreenHeight())-statusBarHeight+5,
		16,
		colorSubtext,
	)

	// Draw timestamp on the right
	timeText := time.Now().Format("15:04:05")
	timeWidth := rl.MeasureText(timeText, 16)
	rl.DrawText(
		timeText,
		int32(rl.GetScreenWidth())-int32(timeWidth)-10,
		int32(rl.GetScreenHeight())-statusBarHeight+5,
		16,
		colorSubtext,
	)
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(1280, 800, "Go Trader")
	rl.SetTargetFPS(60)

	// Set custom font for a more modern look
	// You'll need to download and include a suitable monospaced font
	// For example, "JetBrainsMono-Regular.ttf" would be good for this
	font := rl.LoadFont("JetBrainsMono-Regular.ttf")
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	// Create order book window
	orderBookWindow := NewWindow("BTC-PERPETUAL", 350, 50, 300, 600, renderOrderBook)
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
		rl.ClearBackground(colorBackground)

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
