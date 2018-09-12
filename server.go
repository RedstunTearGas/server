package main

import (
	"fmt"
	"net/http"
	"os"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

func main() {
	stripe.Key = os.Getenv("SECRET_KEY")

	http.HandleFunc("/charge", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		customerParams := &stripe.CustomerParams{
			Email: stripe.String(r.Form.Get("stripeEmail")),
		}
		customerParams.SetSource(r.Form.Get("stripeToken"))

		newCustomer, err := customer.New(customerParams)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		chargeParams := &stripe.ChargeParams{
			Amount:      stripe.Int64(500),
			Currency:    stripe.String(string(stripe.CurrencyUSD)),
			Description: stripe.String("Sample Charge"),
			Customer:    stripe.String(newCustomer.ID),
		}

		if _, err := charge.New(chargeParams); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Charge completed successfully!")
	})

	http.ListenAndServe(":4567", nil)
}
