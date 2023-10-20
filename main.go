package main

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
)

type Hiking struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id"`
	Name        string             `json:"name" bson:"name"`
	Location    string             `json:"location" bson:"location"`
	Date        string             `json:"date" bson:"date"`
	Parking     bool               `json:"parking" bson:"parking"`
	Length      string             `json:"length" bson:"length"`
	Difficulty  int                `json:"difficulty" bson:"difficulty"`
	Description string             `json:"description" bson:"description"`
}

type Observation struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	HikingId string             `json:"hiking_id" bson:"hiking_id"`
	Name     string             `json:"name" bson:"name"`
	Comment  string             `json:"comment" bson:"comment"`
	Time     string             `json:"time" bson:"time"`
}

var hikingCollection *mongo.Collection
var observationCollection *mongo.Collection

func main() {
	//load go dotenv
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	url := os.Getenv("URL")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		fmt.Println(err)
		return
	}

	hikingCollection = client.Database("mobile-app").Collection("hiking")
	observationCollection = client.Database("mobile-app").Collection("observation")

	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	r.Use(gin.Logger())
	r.Use(cors.Default())
	r.GET("/hiking", getHikings())
	r.GET("/hiking/:id", getAHiking())
	r.POST("/hiking", createHiking())
	r.PATCH("/hiking/:id", updateHiking())
	r.DELETE("/hiking/:id", deleteHiking())

	r.GET("/observation", getObservations())
	r.GET("/observation/:id", getAnObservation())
	r.POST("/observation", createObservation())
	r.PATCH("/observation/:id", updateObservation())
	r.DELETE("/observation/:id", deleteObservation())

	r.Run(":8080")
}

func createHiking() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var hiking Hiking
		if err := c.BindJSON(&hiking); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hiking.ID = primitive.NewObjectID()
		_, err := hikingCollection.InsertOne(ctx, hiking)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, hiking)
	}
}

func getHikings() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var hikings []Hiking
		cursor, err := hikingCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err = cursor.All(ctx, &hikings); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, hikings)
	}
}

func getAHiking() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var hiking Hiking
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err = hikingCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&hiking); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, hiking)
	}
}

func updateHiking() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var hiking Hiking
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := c.BindJSON(&hiking); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hiking.ID = id

		if _, err := hikingCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": hiking}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, hiking)
	}
}

func deleteHiking() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if _, err := hikingCollection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Hiking deleted"})
	}
}

func createObservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var observation Observation
		if err := c.BindJSON(&observation); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		observation.ID = primitive.NewObjectID()
		_, err := observationCollection.InsertOne(ctx, observation)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, observation)
	}
}

func getObservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var observations []Observation
		cursor, err := observationCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		if err = cursor.All(ctx, &observations); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, observations)
	}
}

func getAnObservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var observation Observation
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err = observationCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&observation); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, observation)
	}
}

func updateObservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var observation Observation
		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if err := c.BindJSON(&observation); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		observation.ID = id

		if _, err := observationCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": observation}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, observation)
	}
}

func deleteObservation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		id, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if _, err := observationCollection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Observation deleted"})
	}
}
