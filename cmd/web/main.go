package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/pred695/user-management-microservice/database"
	"github.com/pred695/user-management-microservice/routes"
	"gorm.io/gorm"
)

var DbConn *gorm.DB
var Config database.Config
const maxRetries = 5
const retryDelay = 5 * time.Second

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Config = database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// Retry connection logic
	for i := 0; i < maxRetries; i++ {
		DbConn, err = database.Connect(&Config)
		if err == nil {
			log.Printf("Successfully connected to database on attempt %d\n", i+1)
			break
		}

		log.Printf("Failed to connect to database on attempt %d: %v\n", i+1, err)
		time.Sleep(retryDelay) // Wait before retrying
	}

	if err != nil {
		log.Fatalf("Could not establish a database connection after %d attempts: %v", maxRetries, err)
	} else {
		log.Printf("DB Config: %+v\n", Config)
	}
	database.InitRedis()
}

func main() {
	app := fiber.New()
	db, err := DbConn.DB()
	if err != nil {
		log.Fatal("Error connecting to the database", err)
	}
	defer db.Close()
	routes.SetUpRoutes(app)
	app.Listen(":3001")

}
