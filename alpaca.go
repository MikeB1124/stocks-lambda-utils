package stockslambdautils

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaClient struct {
	*alpaca.Client
}

func NewAlpacaClient(apiKey string, apiSecret string, baseUrl string) AlpacaClient {
	client := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    apiKey,
		APISecret: apiSecret,
		BaseURL:   baseUrl,
	})
	return AlpacaClient{client}
}

func (client AlpacaClient) GetAllAlpacaOrders() ([]alpaca.Order, error) {
	orders, err := client.GetOrders(alpaca.GetOrdersRequest{
		Nested: true,
		Status: "all",
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (client AlpacaClient) GetAlpacaAccount() *alpaca.Account {
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	return acct
}

func (client AlpacaClient) CreateBracketOrder(symbol string, entryPrice float64, qty int, orderType string) (*alpaca.Order, error) {
	entryPriceDecimal := decimal.NewFromFloat(entryPrice)
	qtyDecimal := decimal.NewFromInt(int64(qty))

	order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        alpaca.Side(orderType),
		Type:        alpaca.Limit,
		TimeInForce: alpaca.Day,
		LimitPrice:  &entryPriceDecimal,
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (client AlpacaClient) GetAlpacaOrders(status string, symbols []string) ([]alpaca.Order, error) {
	orders, err := client.GetOrders(alpaca.GetOrdersRequest{
		Status:  status,
		Nested:  true,
		Symbols: symbols,
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (client AlpacaClient) CancelAlpacaOrder(orderID string) error {
	err := client.CancelOrder(orderID)
	if err != nil {
		return err
	}
	return nil
}
