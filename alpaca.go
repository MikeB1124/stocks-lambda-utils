package stockslambdautils

import (
	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/shopspring/decimal"
)

type AlpacaClient struct {
	*alpaca.Client
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

func (client AlpacaClient) CreateBracketOrder(symbol string, entryPrice float64, qty int, stopLoss float64, takeProfit float64) (*alpaca.Order, error) {
	entryPriceDecimal := decimal.NewFromFloat(entryPrice)
	stopLossDecimal := decimal.NewFromFloat(stopLoss)
	takeProfitDecimal := decimal.NewFromFloat(takeProfit)
	qtyDecimal := decimal.NewFromInt(int64(qty))

	order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
		OrderClass:  alpaca.Bracket,
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        alpaca.Buy,
		Type:        alpaca.Limit,
		TimeInForce: alpaca.Day,
		LimitPrice:  &entryPriceDecimal,
		TakeProfit: &alpaca.TakeProfit{
			LimitPrice: &takeProfitDecimal,
		},
		StopLoss: &alpaca.StopLoss{
			StopPrice: &stopLossDecimal,
		},
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
