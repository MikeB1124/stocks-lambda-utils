package stockslambdautils

import (
	"context"
	"fmt"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	*mongo.Client
}

func NewMongoClient(username string, password string) (MongoClient, error) {
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", username, password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return MongoClient{}, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return MongoClient{}, err
	}

	return MongoClient{client}, nil
}

func (client MongoClient) InsertEntryOrder(order AlpacaTrade) error {
	collection := client.Database("Stocks").Collection("orders")
	_, err := collection.InsertOne(context.TODO(), order)
	return err
}

func (client MongoClient) UpateOrder(order Order) (*mongo.UpdateResult, error) {
	collection := client.Database("Stocks").Collection("orders")
	filter := bson.M{"order.id": order.ID}
	update := bson.M{
		"$set": bson.M{
			"order":           order,
			"recordUpdatedAt": time.Now().UTC(),
		},
	}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func (client MongoClient) UpdateAllExpiredOrders() (*mongo.UpdateResult, error) {
	// Update all expired orders
	collection := client.Database("Stocks").Collection("orders")
	filter := bson.M{
		"$or": []bson.M{
			{"order.status": "expired"},
			{"order.status": "canceled"},
		},
		"tradeCompleted": false,
	}
	update := bson.M{
		"$set": bson.M{
			"tradeCompleted":  true,
			"recordUpdatedAt": time.Now().UTC(),
			"tradeProfit":     0,
		},
	}
	result, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func FormatAlpacaOrderForDB(alpacaOrder *alpaca.Order) *Order {
	var order Order
	order = formatOrderForDB(alpacaOrder)
	// for _, leg := range alpacaOrder.Legs {
	// 	order.Legs = append(order.Legs, formatOrderForDB(&leg))

	// }
	return &order
}

func formatOrderForDB(alpacaOrder *alpaca.Order) Order {
	var order Order
	order.ID = alpacaOrder.ID
	order.ClientOrderID = alpacaOrder.ClientOrderID
	order.CreatedAt = alpacaOrder.CreatedAt
	order.UpdatedAt = alpacaOrder.UpdatedAt
	order.SubmittedAt = alpacaOrder.SubmittedAt
	order.FilledAt = alpacaOrder.FilledAt
	order.ExpiredAt = alpacaOrder.ExpiredAt
	order.CanceledAt = alpacaOrder.CanceledAt
	order.FailedAt = alpacaOrder.FailedAt
	order.ReplacedAt = alpacaOrder.ReplacedAt
	order.ReplacedBy = alpacaOrder.ReplacedBy
	order.Replaces = alpacaOrder.Replaces
	order.AssetID = alpacaOrder.AssetID
	order.Symbol = alpacaOrder.Symbol
	order.AssetClass = alpacaOrder.AssetClass
	order.OrderClass = alpacaOrder.OrderClass
	order.Type = alpacaOrder.Type
	order.Side = alpacaOrder.Side
	order.TimeInForce = alpacaOrder.TimeInForce
	order.Status = alpacaOrder.Status
	order.Qty = alpacaOrder.Qty.String()
	order.FilledQty = alpacaOrder.FilledQty.String()
	if alpacaOrder.FilledAvgPrice != nil {
		order.FilledAvgPrice = alpacaOrder.FilledAvgPrice.String()
	}
	order.LimitPrice = alpacaOrder.LimitPrice.String()
	if alpacaOrder.StopPrice != nil {
		order.StopPrice = alpacaOrder.StopPrice.String()
	}
	return order
}
