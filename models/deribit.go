package models

type DeribitRequest struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
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
