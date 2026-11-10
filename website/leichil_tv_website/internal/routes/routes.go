// package routes

// import (
// 	"leichil_tv_website/internal/handlers"
// 	"leichil_tv_website/views/pages"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// )

// func Register(app *fiber.App) {
// 	app.Get("/", handlers.Home)
// 	app.Get("/privacy-policy", handlers.Privacy)
// 	app.Get("/refund-policy", handlers.Refund)
// 	app.Get("/terms-and-conditions", handlers.Terms)

// 	app.Get("/checkout/premium", func(c *fiber.Ctx) error {
// 		var html strings.Builder

// 		err := pages.PremiumCheckout(
// 			"movie_001",
// 			"Example Manipuri Movie",
// 			"/static/images/movies/example.jpg",
// 			"99",
// 			"Jamson",
// 			"9876543210",
// 		).Render(c.UserContext(), &html)

// 		if err != nil {
// 			return err
// 		}

// 		c.Set("Content-Type", "text/html; charset=utf-8")
// 		return c.SendString(html.String())
// 	})

// 	app.Get("/checkout/subscription", func(c *fiber.Ctx) error {
// 		var html strings.Builder

// 		err := pages.SubscriptionCheckout(
// 			"basic_monthly",
// 			"Basic Monthly",
// 			"149",
// 			"1 Month",
// 		).Render(c.UserContext(), &html)

// 		if err != nil {
// 			return err
// 		}

// 		c.Set("Content-Type", "text/html; charset=utf-8")
// 		return c.SendString(html.String())
// 	})
// }

package routes

import (
	"strings"

	"leichil_tv_website/internal/handlers"
	"leichil_tv_website/views/pages"

	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	app.Get("/", handlers.Home)
	app.Get("/privacy-policy", handlers.Privacy)
	app.Get("/refund-policy", handlers.Refund)
	app.Get("/terms-and-conditions", handlers.Terms)

	// Premium Content Checkout
	app.Get("/checkout/premium", func(c *fiber.Ctx) error {
		var html strings.Builder

		err := pages.PremiumCheckout(
			"movie_001",
			"Example Manipuri Movie",
			"/static/images/movies/example.jpg",
			"99",
			"Jamson",
			"9876543210",
		).Render(c.UserContext(), &html)

		if err != nil {
			return err
		}

		c.Type("html", "utf-8")
		return c.SendString(html.String())
	})

	// Subscription Checkout
	app.Get("/checkout/subscription", func(c *fiber.Ctx) error {
		var html strings.Builder

		err := pages.SubscriptionCheckout(
			"basic_monthly",
			"Basic Monthly",
			"149",
			"1 Month",
		).Render(c.UserContext(), &html)

		if err != nil {
			return err
		}

		c.Type("html", "utf-8")
		return c.SendString(html.String())
	})
}
