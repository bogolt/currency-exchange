package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var rateHistory RateHistory
var logger *log.Logger

func main() {

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	rateHistory = newRateHistory()
	rateHistory.load()

	r := chi.NewRouter()
	r.Route("/rates", func(r chi.Router) {
		r.Post("/cur/{currency}/{atDate}", storeCurrencyDateRate)

		r.Get("/cur/{currency}/from/{from}/to/{to}", getRates)
		r.Get("/cur/{currency}/{atDate}", getRateDate)
		r.Get("/cur/{currency}", getLastRate)
		r.Get("/all/{atDate}", getDateRates)
	})

	logger.Printf("currency-history listen on port 8080")
	http.ListenAndServe(":8080", r)
}
