package client

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/adityanagar10/trader/models"
	"github.com/gorilla/websocket"
)

type DeribitClient struct {
	Conn            *websocket.Conn
	OrderBookWindow *models.Window
	Instrument      string
	RequestID       int
}

func NewDeribitClient(instrument string, orderBookWindow *models.Window) *DeribitClient {
	return &DeribitClient{
		OrderBookWindow: orderBookWindow,
		Instrument:      instrument,
		RequestID:       1,
	}
}

func (c *DeribitClient) handleMessages() {
	for {
		_, message, err := c.Conn.ReadMessage()
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
			c.OrderBookWindow.Data = response.Result
		}
	}
}

func (c *DeribitClient) fetchOrderBook() {
	// Create an order book request
	request := models.DeribitRequest{
		JsonRPC: "2.0",
		ID:      c.RequestID,
		Method:  "public/get_order_book",
		Params: models.OrderBookParams{
			InstrumentName: c.Instrument,
		},
	}
	c.RequestID++

	// Marshal and send the request
	data, err := json.Marshal(request)
	if err != nil {
		log.Printf("Failed to marshal request: %v", err)
		return
	}

	if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
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

func (c *DeribitClient) Connect() error {
	// Connect to Deribit WebSocket API
	conn, _, err := websocket.DefaultDialer.Dial("wss://www.deribit.com/ws/api/v2", nil)
	if err != nil {
		return fmt.Errorf("websocket connection error: %v", err)
	}

	c.Conn = conn
	log.Println("Connected to Deribit WebSocket API")

	// Start listening for messages
	go c.handleMessages()

	// Start fetching order book periodically
	go c.fetchOrderBookPeriodically()

	return nil
}

func (c *DeribitClient) Close() {
	if c.Conn != nil {
		c.Conn.Close()
		log.Println("Closed Deribit WebSocket connection")
	}
}
