package internal

import (
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func NewApp(root *CompositionRoot) *fiber.App {
	app := fiber.New()

	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(root.O11y.TracerProvider),
		otelfiber.WithMeterProvider(root.O11y.MeterProvider)))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//app.Get("/:cena", func(c *fiber.Ctx) error {
	//	return c.SendString(fmt.Sprintf("Hello, %s!", c.Params("cena")))
	//})

	addSwagger(app)

	return app
}

func addSwagger(app *fiber.App) {
	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "api/v1.yaml",
		Path:     "swagger",
		Title:    "Swagger API Docs",
	}

	app.Use(swagger.New(cfg))
}
