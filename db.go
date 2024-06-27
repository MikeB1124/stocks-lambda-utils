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

func (client MongoClient) UpateOrder(order alpaca.Order, orderType string) (*mongo.UpdateResult, error) {
	collection := client.Database("Stocks").Collection("orders")
	orderKey := fmt.Sprintf("%s.id", orderType)
	filter := bson.M{orderKey: order.ID}
	update := bson.M{
		"$set": bson.M{
			orderType:         order,
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
			{"entryOrder.status": "expired"},
			{"entryOrder.status": "canceled"},
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
