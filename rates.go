package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
)

type CurrencyRateJson struct {
	Rate float32 `json:"rate"`
}

func storeCurrencyDateRate(w http.ResponseWriter, r *http.Request) {

	currency := chi.URLParam(r, "currency")
	if len(currency) == 0 {
		logger.Printf("storeCurrencyDateRate empty currency")
		w.WriteHeader(400)
		return
	}

	date, err := parseDate(chi.URLParam(r, "atDate"))
	if err != nil {
		logger.Printf("storeCurrencyDateRate invalid date %s", err)
		w.WriteHeader(400)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Printf("storeCurrencyDateRate cannot read body %s", err)
		w.WriteHeader(400)
		return
	}

	currencyRate := CurrencyRateJson{}
	if err := json.Unmarshal(body, &currencyRate); err != nil {
		w.WriteHeader(400)
		return
	}
	rateHistory.storeCurrencyDateRate(currency, date, currencyRate.Rate)
	w.WriteHeader(201)
}

func getLastRate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	date, rate, ok := rateHistory.getLastRate(currency)
	if !ok {
		w.WriteHeader(404)
		return
	}

	writeRate(w, date, rate)
}

func getDateRates(w http.ResponseWriter, r *http.Request) {
	strDate := chi.URLParam(r, "atDate")

	date, err := parseDate(strDate)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	rates := rateHistory.getDateRates(date)
	writeDateRates(w, rates)
}

func getRates(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	strDateFrom := chi.URLParam(r, "from")
	strDateTo := chi.URLParam(r, "to")

	dateFrom := Date(0)
	dateTo := Date(0)

	if len(strDateFrom) > 0 {
		dateFrom, _ = parseDate(strDateFrom)
	}

	if len(strDateTo) > 0 {
		dateTo, _ = parseDate(strDateTo)
	}

	rates := rateHistory.getRatesBetween(currency, dateFrom, dateTo)
	writeRates(w, rates)
}

func getRateDate(w http.ResponseWriter, r *http.Request) {
	currency := chi.URLParam(r, "currency")
	atDate, err := parseDate(chi.URLParam(r, "atDate"))
	if err != nil {
		w.WriteHeader(400)
		return
	}

	rate, found := rateHistory.getRate(currency, atDate)
	if !found {
		w.WriteHeader(404)
		return
	}

	writeRate(w, atDate, rate)
}

func writeRates(w http.ResponseWriter, rates DateRateMap) {
	daterates := make(map[string]float32)
	for date, rate := range rates {
		daterates[dateToStr(date)] = float32(rate)
	}
	outdata, err := json.Marshal(daterates)
	if err != nil {
		fmt.Printf("marshal json failed: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(200)
	w.Write(outdata)
}

func writeRate(w http.ResponseWriter, date Date, rate Rate) {
	rates := make(DateRateMap)
	rates[date] = rate
	writeRates(w, rates)
}

func writeDateRates(w http.ResponseWriter, rates DateCurrentyRate) {
	daterates := make(map[string]float32)
	for currency, rate := range rates {
		daterates[currency] = float32(rate)
	}
	outdata, err := json.Marshal(daterates)
	if err != nil {
		fmt.Printf("marshal json failed: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(200)
	w.Write(outdata)
}
