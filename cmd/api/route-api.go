package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "C-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Post("/api/payment-intent", app.GetPaymentIntent)
	mux.Get("/api/widget/{id}", app.GetWidgetByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)
	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is-autheticated", app.CheckAuthentication)
	mux.Post("/api/forget-password", app.SendPasswordResetEmail)
	mux.Post("/api/reset-password", app.ResetPassword)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.Auth)

		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Loggin"))
		})
		mux.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)
		mux.Post("/all-sales", app.AllSales)
		mux.Post("/all-subscription", app.AllSucription)
		mux.Post("/get-sale/{id}", app.GetSale)
		mux.Post("/refund", app.RefundCharge)
		mux.Post("/cancel-subscription", app.CancelSubscription)
		mux.Post("/all-users", app.AllUsers)
		mux.Post("/all-users/{id}", app.DetailUser)
		mux.Post("/all-users/edit/{id}", app.EditUser)
		mux.Post("/all-users/delete/{id}", app.DeleteUser)
	})
	return mux
}
