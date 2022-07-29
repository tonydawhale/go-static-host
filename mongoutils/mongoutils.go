package mongoutils

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MongoClient *mongo.Client

func Init() {
	c, err := mongo.Connect(
		context.TODO(), 
		options.Client().ApplyURI(os.Getenv("MONGO_URI")),
	)

	if err != nil {
		log.Fatal(err)
	}

	MongoClient = c

	// Ping the primary
	if err := MongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Successfully Connected To MongoDB")
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}

func CreateItemMetaData(uuid string, shortId string, contentType string) (res *mongo.InsertOneResult, err error) {
	database := MongoClient.Database(os.Getenv("MONGO_DATABASE"))
	collection := database.Collection(os.Getenv("MONGO_COLLECTION"))

	result, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "uuid", Value: uuid},
		{Key: "shortId", Value: shortId},
		{Key: "content-type", Value: contentType},
	})

	return result, err
}

type metadata struct {
	_id 		primitive.ObjectID
	ContentType string
	ShortId 	string
	Uuid		string
}

func FetchItemMetaData(shortId string) (metadata, error){
	database := MongoClient.Database(os.Getenv("MONGO_DATABASE"))
	collection := database.Collection(os.Getenv("MONGO_COLLECTION"))
	var result metadata
	err := collection.FindOne(
		context.TODO(), 
		bson.M{"shortId": shortId},
	).Decode(&result)

	return result, err
}