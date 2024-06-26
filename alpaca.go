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

func (client AlpacaClient) GetAlpacaAccount() *alpaca.Account {
	acct, err := client.GetAccount()
	if err != nil {
		panic(err)
	}
	return acct
}

func (client AlpacaClient) CreateAlpacaOrder(
	symbol string,
	entryPrice float64,
	stopPrice float64,
	takeProfit float64,
	qty int,
	orderSide string,
	orderType string,
	orderClass string,
	timeInForce string,
) (*alpaca.Order, error) {
	qtyDecimal := decimal.NewFromInt(int64(qty))

	orderRequest := alpaca.PlaceOrderRequest{
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        alpaca.Side(orderSide),
		Type:        alpaca.OrderType(orderType),
		TimeInForce: alpaca.TimeInForce(timeInForce),
	}

	if orderType == "limit" || orderType == "stop_limit" {
		entryPriceDecimal := decimal.NewFromFloat(entryPrice)
		orderRequest.LimitPrice = &entryPriceDecimal
	}

	if orderClass == "bracket" {
		orderRequest.OrderClass = alpaca.OrderClass(orderClass)

		stopPriceDecimal := decimal.NewFromFloat(stopPrice)
		orderRequest.StopLoss = &alpaca.StopLoss{
			StopPrice: &stopPriceDecimal,
		}

		takeProfitDecimal := decimal.NewFromFloat(takeProfit)
		orderRequest.TakeProfit = &alpaca.TakeProfit{
			LimitPrice: &takeProfitDecimal,
		}
	}

	order, err := client.PlaceOrder(orderRequest)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (client AlpacaClient) GetAlpacaOrders(status string, symbols []string, nested bool) ([]alpaca.Order, error) {
	orderRequest := alpaca.GetOrdersRequest{
		Status: status,
		Nested: nested,
	}

	if len(symbols) != 0 {
		orderRequest.Symbols = symbols
	}

	orders, err := client.GetOrders(orderRequest)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (client AlpacaClient) GetAlpacaOrderByID(orderID string) (*alpaca.Order, error) {
	order, err := client.GetOrder(orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (client AlpacaClient) CancelAlpacaOrder(orderID string) error {
	err := client.CancelOrder(orderID)
	if err != nil {
		return err
	}
	return nil
}
