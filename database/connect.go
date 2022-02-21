package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectUsers() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb+srv://huy:123@cluster0.zfyia.mongodb.net/golang?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connect
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("golang").Collection("users")
	return collection
}

func ConnectPosts() *mongo.Collection {
	clientOptions := options.Client().ApplyURI("mongodb+srv://huy:123@cluster0.zfyia.mongodb.net/golang?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("golang").Collection("posts")
	return collection
}
