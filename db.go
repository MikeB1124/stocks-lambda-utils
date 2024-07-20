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
	Env string
}

func NewMongoClient(username string, password string, env string) (MongoClient, error) {
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb+srv://%s:%s@cluster0.du0vf.mongodb.net", username, password))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return MongoClient{}, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return MongoClient{}, err
	}

	return MongoClient{client, env}, nil
}

func (client MongoClient) InsertEntryOrder(order AlpacaTrade) error {
	collection := client.Database("Stocks").Collection("orders" + client.Env)
	_, err := collection.InsertOne(context.TODO(), order)
	return err
}

func (client MongoClient) UpdateAllExpiredOrders() (*mongo.UpdateResult, error) {
	// Update all expired orders
	collection := client.Database("Stocks").Collection("orders" + client.Env)
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
			"tradeCanceled":   true,
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

func (client MongoClient) GetFilledTradesFromDB() ([]AlpacaTrade, error) {
	collection := client.Database("Stocks").Collection("orders" + client.Env)
	filter := bson.M{
		"$or": []bson.M{
			{"order.legs.0.status": "filled"},
			{"order.legs.0.status": "partially_filled"},
			{"order.legs.1.status": "filled"},
			{"order.legs.1.status": "partially_filled"},
		},
		"tradeCompleted": false,
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var orders []AlpacaTrade
	err = cursor.All(context.TODO(), &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (client MongoClient) GetOpenTrades() ([]AlpacaTrade, error) {
	collection := client.Database("Stocks").Collection("orders" + client.Env)
	filter := bson.M{
		"tradeCompleted": false,
	}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	var orders []AlpacaTrade
	err = cursor.All(context.TODO(), &orders)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (client MongoClient) BulkUpdateTradeProfits(trades []AlpacaTrade) (*mongo.BulkWriteResult, error) {
	collection := client.Database("Stocks").Collection("orders" + client.Env)
	var updates []mongo.WriteModel
	for _, trade := range trades {
		filter := bson.M{"_id": trade.ObjectID}
		update := bson.M{
			"$set": bson.M{
				"tradeProfit":     trade.TradeProfit,
				"tradeCompleted":  true,
				"recordUpdatedAt": time.Now().UTC(),
			},
		}
		updates = append(updates, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	bulkWriteResult, err := collection.BulkWrite(context.TODO(), updates)
	if err != nil {
		return nil, err
	}
	return bulkWriteResult, nil
}

func (client MongoClient) BulkUpdateTrades(trades []AlpacaTrade) (*mongo.BulkWriteResult, error) {
	collection := client.Database("Stocks").Collection("orders" + client.Env)
	var updates []mongo.WriteModel
	for _, trade := range trades {
		filter := bson.M{"_id": trade.ObjectID}
		update := bson.M{
			"$set": bson.M{
				"order":           trade.Order,
				"recordUpdatedAt": time.Now().UTC(),
			},
		}
		updates = append(updates, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update))
	}
	bulkWriteResult, err := collection.BulkWrite(context.TODO(), updates)
	if err != nil {
		return nil, err
	}
	return bulkWriteResult, nil
}

func (client MongoClient) GetSettings() (Settings, error) {
	var settings Settings
	coll := client.Database("Stocks").Collection("settings")
	err := coll.FindOne(context.TODO(), bson.M{"env": client.Env}).Decode(&settings)
	if err != nil {
		return Settings{}, err
	}
	return settings, nil
}

func (client MongoClient) GetAllSettings() ([]Settings, error) {
	var settings []Settings
	coll := client.Database("Stocks").Collection("settings")
	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		return []Settings{}, err
	}
	err = cursor.All(context.TODO(), &settings)
	if err != nil {
		return []Settings{}, err
	}
	return settings, nil
}

func FormatAlpacaOrderForDB(alpacaOrder *alpaca.Order) *Order {
	var order Order
	order = formatOrderForDB(alpacaOrder)
	for _, leg := range alpacaOrder.Legs {
		order.Legs = append(order.Legs, formatOrderForDB(&leg))

	}
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
	if alpacaOrder.LimitPrice != nil {
		order.LimitPrice = alpacaOrder.LimitPrice.String()
	}
	if alpacaOrder.StopPrice != nil {
		order.StopPrice = alpacaOrder.StopPrice.String()
	}
	return order
}
