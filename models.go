package stockslambdautils

import (
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
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
	Order           Order              `json:"order" bson:"order"`
	TradeCompleted  bool               `json:"tradeCompleted" bson:"tradeCompleted"`
	TradeProfit     float64            `json:"tradeProfit" bson:"tradeProfit"`
	PatternData     PatternData        `json:"patternData" bson:"patternData"`
	RecordUpdatedAt *time.Time         `json:"recordUpdatedAt" bson:"recordUpdatedAt"`
}

type Order struct {
	ID             string             `json:"id" bson:"id"`
	ClientOrderID  string             `json:"client_order_id" bson:"client_order_id"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	SubmittedAt    time.Time          `json:"submitted_at" bson:"submitted_at"`
	FilledAt       *time.Time         `json:"filled_at" bson:"filled_at"`
	ExpiredAt      *time.Time         `json:"expired_at" bson:"expired_at"`
	CanceledAt     *time.Time         `json:"canceled_at" bson:"canceled_at"`
	FailedAt       *time.Time         `json:"failed_at" bson:"failed_at"`
	ReplacedAt     *time.Time         `json:"replaced_at" bson:"replaced_at"`
	ReplacedBy     *string            `json:"replaced_by" bson:"replaced_by"`
	Replaces       *string            `json:"replaces" bson:"replaces"`
	AssetID        string             `json:"asset_id" bson:"asset_id"`
	Symbol         string             `json:"symbol" bson:"symbol"`
	AssetClass     alpaca.AssetClass  `json:"asset_class" bson:"asset_class"`
	OrderClass     alpaca.OrderClass  `json:"order_class" bson:"order_class"`
	Type           alpaca.OrderType   `json:"type" bson:"type"`
	Side           alpaca.Side        `json:"side" bson:"side"`
	TimeInForce    alpaca.TimeInForce `json:"time_in_force" bson:"time_in_force"`
	Status         string             `json:"status" bson:"status"`
	Qty            string             `json:"qty" bson:"qty"`
	FilledQty      string             `json:"filled_qty" bson:"filled_qty"`
	FilledAvgPrice string             `json:"filled_avg_price" bson:"filled_avg_price"`
	LimitPrice     string             `json:"limit_price" bson:"limit_price"`
	StopPrice      string             `json:"stop_price" bson:"stop_price"`
	Legs           []Order            `json:"legs" bson:"legs"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}
