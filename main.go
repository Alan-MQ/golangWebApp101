package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/golangWebApp101/golangWebApp101/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoClient *mongo.Client
	redisClient *redis.Client
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()
	mongoConn := options.Client().ApplyURI(config.DBUri)
	mongoClient, err := mongo.Connect(ctx, mongoConn)

	if err != nil {
		panic(err)
	}

	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	redisClient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})
	if _, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisClient.Set(ctx, "test", "welcome to golang with redis and mongodb", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")
	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config in main", err)
	}

	defer mongoClient.Disconnect(ctx)

	value, err := redisClient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("Key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": value})
	})
	log.Fatal(server.Run(":" + config.Port))
}
