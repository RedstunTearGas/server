package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

type (
	StripeToken struct {
		Id    string `json:"id"`
		Email string `json:"email"`
	}

	ChargeRequest struct {
		Token StripeToken `json:"token" validate:"required"`
	}
	ChargeResponse struct {
		Message string `json:"message"`
	}
)

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.POST("/charge", func(c echo.Context) error {
		cr := new(ChargeRequest)

		if err := c.Bind(cr); err != nil {
			return c.JSON(http.StatusForbidden, &ChargeResponse{
				Message: "wrong data",
			})
		}

		stripe.Key = os.Getenv("STRIPE_KEY")

		customerParams := &stripe.CustomerParams{
			Email: stripe.String(cr.Token.Email),
		}
		customerParams.SetSource(cr.Token.Id)

		newCustomer, err := customer.New(customerParams)

		if err != nil {
			c.JSON(http.StatusInternalServerError, &ChargeResponse{
				Message: err.Error(),
			})
			return nil
		}

		chargeParams := &stripe.ChargeParams{
			Amount:      stripe.Int64(108 * 100),
			Currency:    stripe.String(string(stripe.CurrencyEUR)),
			Description: stripe.String("Bombe Lacry 450ml 500g"),
			Customer:    stripe.String(newCustomer.ID),
		}

		if _, err := charge.New(chargeParams); err != nil {
			c.JSON(http.StatusInternalServerError, &ChargeResponse{
				Message: err.Error(),
			})
			return nil
		}

		return c.JSON(http.StatusOK, &ChargeResponse{
			Message: cr.Token.Id,
		})
	})

	e.Logger.Fatal(e.Start(":1323"))
}
