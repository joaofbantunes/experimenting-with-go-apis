package internal

import (
	"errors"
	"log/slog"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/auth"
	"github.com/joaofbantunes/experimenting-with-go-apis/fiber/internal/features/shared/problems"
)

func NewApp(root *CompositionRoot) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler(root.CreateLogger),
	})

	app.Use(otelfiber.Middleware(
		otelfiber.WithTracerProvider(root.O11y.TracerProvider),
		otelfiber.WithMeterProvider(root.O11y.MeterProvider)))

	app.Use(auth.Authenticate(root.CreateLogger))

	addRoutes(app, root)
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

func errorHandler(loggerProvider func(name string) *slog.Logger) func(c *fiber.Ctx, err error) error {
	logger := loggerProvider("fiber_error_handler")
	return func(c *fiber.Ctx, err error) error {
		var pe *problems.ProblemError
		if errors.As(err, &pe) {
			return c.Status(pe.Problem.Status).JSON(pe.Problem)
		}
		var vpe *problems.ValidationProblemError
		if errors.As(err, &vpe) {
			return c.Status(vpe.ValidationProblem.Status).JSON(vpe.ValidationProblem)
		}
		// Fallback to default 500 Internal Server Error
		logger.ErrorContext(c.UserContext(), "internal server error", slog.Any("error", err))
		return c.Status(fiber.StatusInternalServerError).SendString("internal server error")
	}
}
