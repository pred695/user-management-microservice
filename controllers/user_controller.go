package controllers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v3"
	"github.com/pred695/user-management-microservice/database"
	"github.com/pred695/user-management-microservice/models"
)
var redisContext = context.Background()
// GetUserById fetches a user by their ID with Redis caching
func GetUserById(ctx fiber.Ctx) error {
	contextMap := fiber.Map{
		"message":    "Get User by ID",
		"statusText": "Ok",
	}
	db := database.DbConn
	redisClient := database.RedisClient

	id := ctx.Params("user_id")
	cacheKey := "user_" + id

	// Check if user is cached in Redis
	cachedUser, err := redisClient.Get(redisContext, cacheKey).Result()
	if err == redis.Nil {
		// Not found in cache, query from database
		var user models.User
		if err := db.First(&user, id).Error; err != nil {
			contextMap["statusText"] = "Not Found"
			contextMap["message"] = "User not found"
			return ctx.Status(fiber.StatusNotFound).JSON(contextMap)
		}

		// Cache the user in Redis
		userJSON, _ := json.Marshal(user)
		redisClient.Set(redisContext, cacheKey, userJSON, 10*time.Minute) // Cache for 10 minutes

		contextMap["user"] = user
		return ctx.JSON(contextMap)
	} else if err != nil {
		// Redis error
		contextMap["statusText"] = "Internal Server Error"
		contextMap["message"] = "Redis error"
		return ctx.Status(fiber.StatusInternalServerError).JSON(contextMap)
	}

	// User found in Redis cache
	var cachedUserObj models.User
	json.Unmarshal([]byte(cachedUser), &cachedUserObj)

	contextMap["user"] = cachedUserObj
	return ctx.JSON(contextMap)
}


// GetUsers fetches all users with Redis caching
func GetUsers(ctx fiber.Ctx) error {
	contextMap := fiber.Map{
		"message":    "Get Users",
		"statusText": "Ok",
	}
	db := database.DbConn
	redisClient := database.RedisClient

	// Check if users are cached in Redis
	cacheKey := "users_list"
	cachedUsers, err := redisClient.Get(redisContext, cacheKey).Result()
	if err == redis.Nil {
		// Not found in cache, query from database
		var users []models.User
		result := db.Find(&users)
		if result.Error != nil {
			contextMap["statusText"] = "Internal Server Error"
			contextMap["message"] = "Error Fetching Users"
			return ctx.Status(fiber.StatusInternalServerError).JSON(contextMap)
		}

		// Cache the users list in Redis
		usersJSON, _ := json.Marshal(users)
		redisClient.Set(redisContext, cacheKey, usersJSON, 10*time.Minute) // Cache for 10 minutes

		contextMap["users"] = users
		return ctx.JSON(contextMap)
	} else if err != nil {
		// Redis error
		contextMap["statusText"] = "Internal Server Error"
		contextMap["message"] = "Redis error"
		return ctx.Status(fiber.StatusInternalServerError).JSON(contextMap)
	}

	// Users found in Redis cache
	var cachedUsersObj []models.User
	json.Unmarshal([]byte(cachedUsers), &cachedUsersObj)

	contextMap["users"] = cachedUsersObj
	return ctx.JSON(contextMap)
}
