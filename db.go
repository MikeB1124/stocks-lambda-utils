package stockslambdautils

import (
	"context"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoClient struct {
	*mongo.Client
}

func (client MongoClient) InsertEntryOrder(order AlpacaEntryOrder) error {
	collection := client.Database("Stocks").Collection("orders")
	_, err := collection.InsertOne(context.TODO(), order)
	return err
}

func (client MongoClient) UpateOrder(order alpaca.Order) (*mongo.UpdateResult, error) {
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
	filter := bson.M{"order.status": "expired", "tradeCompleted": false}
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

func (client MongoClient) UpdateAllCancelledOrders() (*mongo.UpdateResult, error) {
	// Update all cancelled orders
	collection := client.Database("Stocks").Collection("orders")
	filter := bson.M{"order.status": "canceled", "tradeCompleted": false}
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
