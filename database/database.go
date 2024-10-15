package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/pred695/user-management-microservice/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

var (
	DbConn      *gorm.DB         // a public variable to hold the Postgres connection
	RedisClient *redis.Client     // a public variable for Redis client
	ctx         = context.Background() // context for Redis operations
)

// Connect initializes the PostgreSQL connection
func Connect(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to the PostgreSQL database:", err)
		return nil, err
	}
	fmt.Println("Successfully connected to PostgreSQL")

	// Automatically migrate the User model (ensure your User model is defined)
	db.AutoMigrate(new(models.User))

	// Assign the global DbConn variable
	DbConn = db
	return db, nil
}

// InitRedis initializes the Redis client
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"), // Fetch Redis host and port from environment variables
		Password: os.Getenv("REDIS_PASSWORD"), // Fetch Redis password from environment variables
		DB:       0,                           // Default Redis DB
	})

	// Test the Redis connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis successfully")
}