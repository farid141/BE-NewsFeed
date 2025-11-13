package router

import (
	"github.com/farid141/go-rest-api/controller"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	userController *controller.UserController
	authController *controller.AuthController
	postController *controller.PostController
}

func NewRouter(
	userController *controller.UserController,
	authController *controller.AuthController,
	postController *controller.PostController,
) *Router {
	return &Router{
		userController: userController,
		authController: authController,
		postController: postController,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", r.authController.Login)
	api.Post("/register", r.authController.Register)
	api.Post("/refresh_token", r.authController.RefreshToken)
	api.Post("/logout", r.authController.Logout)
	api.Get("/users", r.userController.GetUsers)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("secret")},
		TokenLookup: "cookie:token", // 👈 look in cookie named "token"
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "access forbidden, invalid or expired token",
				})
			}
			return nil
		},
	}))

	api.Get("/me", r.authController.Me)

	api.Post("/follow/:id", r.userController.FollowUser(true))
	api.Delete("/follow/:id", r.userController.FollowUser(false))
	api.Post("/posts", r.postController.CreatePost)
	api.Get("/feed", r.postController.GetFeed)
}
