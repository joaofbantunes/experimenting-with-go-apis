package internal

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/orders/register_order"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/auth"
)

func addRoutes(app *fiber.App, root *CompositionRoot) {

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/hello/:stuff", func(c *fiber.Ctx) error {
		return c.SendString(fmt.Sprintf("Hello, %s!", c.Params("stuff")))
	})

	v1 := app.Group("/api/v1")

	// can be applied to api or v1 if all endpoints require auth (not app, because of health checks)
	requireAuthMw := auth.RequireAuthentication()

	v1Orders := v1.Group("/orders", requireAuthMw)
	v1Orders.Post("",
		auth.RequirePermission("orders.register", root.CreateLogger),
		register_order.NewRegisterOrderEndpoint(
			root.CreateLogger,
			root.TimeProvider))
}
