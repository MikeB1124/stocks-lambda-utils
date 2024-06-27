package stockslambdautils

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// All strucs for lambdas

type PatternWebhookRequest struct {
	MsgType string        `json:"msg_type" bson:"msg_type"`
	Data    []PatternData `json:"data" bson:"data"`
}

type PatternData struct {
	PatternType   string      `json:"patterntype" bson:"patterntype"`
	PatternName   string      `json:"patternname" bson:"patternname"`
	ProfitOne     float64     `json:"profit1" bson:"profit1"`
	DisplaySymbol string      `json:"displaysymbol" bson:"displaysymbol"`
	Symbol        string      `json:"symbol" bson:"symbol"`
	StopLoss      float64     `json:"stoploss" bson:"stoploss"`
	PatternUrl    string      `json:"url" bson:"url"`
	TimeFrame     string      `json:"timeframe" bson:"timeframe"`
	Status        string      `json:"status" bson:"status"`
	Entry         interface{} `json:"entry" bson:"entry"`
	PatternClass  string      `json:"patternclass" bsom:"patternclass"`
}

type AlpacaTrade struct {
	ObjectID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	EntryOrder      Order              `json:"entryOrder" bson:"entryOrder"`
	ExitOrder       Order              `json:"exitOrder" bson:"exitOrder"`
	TradeCompleted  bool               `json:"tradeCompleted" bson:"tradeCompleted"`
	TradeProfit     float64            `json:"tradeProfit" bson:"tradeProfit"`
	PatternData     PatternData        `json:"patternData" bson:"patternData"`
	RecordUpdatedAt *time.Time         `json:"recordUpdatedAt" bson:"recordUpdatedAt"`
}

type Order struct {
	ID             string             `json:"id"`
	ClientOrderID  string             `json:"client_order_id"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	SubmittedAt    time.Time          `json:"submitted_at"`
	FilledAt       *time.Time         `json:"filled_at"`
	ExpiredAt      *time.Time         `json:"expired_at"`
	CanceledAt     *time.Time         `json:"canceled_at"`
	FailedAt       *time.Time         `json:"failed_at"`
	ReplacedAt     *time.Time         `json:"replaced_at"`
	ReplacedBy     *string            `json:"replaced_by"`
	Replaces       *string            `json:"replaces"`
	AssetID        string             `json:"asset_id"`
	Symbol         string             `json:"symbol"`
	AssetClass     alpaca.AssetClass  `json:"asset_class"`
	OrderClass     alpaca.OrderClass  `json:"order_class"`
	Type           alpaca.OrderType   `json:"type"`
	Side           alpaca.Side        `json:"side"`
	TimeInForce    alpaca.TimeInForce `json:"time_in_force"`
	Status         string             `json:"status"`
	Qty            *decimal.Decimal   `json:"qty"`
	FilledQty      decimal.Decimal    `json:"filled_qty"`
	FilledAvgPrice *decimal.Decimal   `json:"filled_avg_price"`
	LimitPrice     *decimal.Decimal   `json:"limit_price"`
	StopPrice      *decimal.Decimal   `json:"stop_price"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
