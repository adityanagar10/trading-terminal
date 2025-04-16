package models

type OrderBookParams struct {
	InstrumentName string `json:"instrument_name"`
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
