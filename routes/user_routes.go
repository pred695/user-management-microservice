package routes

import "github.com/gofiber/fiber/v3"
import "github.com/pred695/user-management-microservice/controllers"

func SetUpRoutes(app *fiber.App) {
	private := app.Group("/private")
	private.Get("/user/:user_id", controllers.GetUserById);
	private.Get("/users", controllers.GetUsers);
	// private.Put("/user/:user_id", controllers.UpdateUser);
	// private.Delete("/user/:user_id", controllers.DeleteUser);
}
