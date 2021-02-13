package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func parseDate(dateString string) (Date, error) {
	layout := "2006-01-02"
	t, err := time.Parse(layout, dateString)
	if err != nil {
		return 0, err
	}
	return Date(uint32(t.Year())<<9 | uint32(t.Month())<<5 | uint32(t.Day())), nil
}

func dateToStr(dt Date) string {
	year := dt >> 9
	month := (dt >> 5) & 0xf
	day := dt & 31
	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

func parseRate(s string) (Rate, error) {
	frate, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return Rate(frate), nil
}

func parseCsv(filename string) (string, DateRateMap, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()

	return parseRatesFromReader(f)
}

func parseRatesFromReader(f io.Reader) (string, DateRateMap, error) {
	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return "", nil, err
	}

	if len(records) == 0 {
		return "", nil, fmt.Errorf("No data")
	}

	currency := records[0][1][:3]

	dateRateMap := make(DateRateMap)
	for _, r := range records[1:] {
		if len(r) < 2 {
			continue
		}
		date, err := parseDate(r[0])
		if err != nil {
			return "", nil, fmt.Errorf("failed to parse date %s %s", r[0], err)
		}
		strDate := dateToStr(date)
		if strDate[5] == '4' {
			fmt.Printf("Parsed %s as %d as %s", r[0], date, strDate)
			continue
		}
		rate, err := parseRate(r[1])
		if err != nil {
			// fmt.Printf("Failed to parse rate %s (%s) for date %s \n", r[1], err, r[0])
			continue
		}

		dateRateMap[date] = rate
	}
	return currency, dateRateMap, nil
}

func parseExtraRates(filename string) (CurrencyHistory, error) {

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("No data")
	}

	history := make(CurrencyHistory)
	for _, r := range records[1:] {
		if len(r) < 3 {
			continue
		}
		currency := r[0]
		date, err := parseDate(r[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %s %s", r[0], err)
		}
		strDate := dateToStr(date)
		if strDate[5] == '4' {
			fmt.Printf("Parsed %s as %d as %s", r[0], date, strDate)
			continue
		}
		rate, err := parseRate(r[2])
		if err != nil {
			// fmt.Printf("Failed to parse rate %s (%s) for date %s \n", r[1], err, r[0])
			continue
		}

		h, ok := history[currency]
		if !ok {
			h = make(DateRateMap)
			history[currency] = h
		}
		h[date] = rate
	}
	return history, nil
}
