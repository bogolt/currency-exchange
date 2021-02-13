package main

import (
	"log"
	"os"
	"strings"
	"testing"
)

func testDate(t *testing.T, strDate string) {
	date, err := parseDate(strDate)
	if err != nil {
		t.Errorf("Failed to parse date %s %s", strDate, err)
		return
	}

	dateStr := dateToStr(date)
	if dateStr != strDate {
		t.Errorf("date convert failed %s != %s", dateStr, strDate)
	}
}

func TestParseDate(t *testing.T) {

	dates := []string{"2018-08-01", "2016-01-29"}
	for _, date := range dates {
		testDate(t, date)
	}
}

func TestParser(t *testing.T) {
	in := `DATE,SEKUSD
2016-01-29,8.5709
2016-02-01,8.5385
2016-02-02,8.5749
2016-02-03,8.4745
2016-02-04,8.3983
`

	curName, rates, err := parseRatesFromReader(strings.NewReader(in))
	if err != nil {
		t.Errorf("parse failed: %s", err)
		return
	}
	if curName != "SEK" {
		t.Errorf("currency name invalid, expected 'SEK' got '%s'", curName)
	}
	if len(rates) != 5 {
		t.Errorf("Invalid number of dates parsed")
		return
	}

	//2016-02-02,8.5749
	dt, _ := parseDate("2016-02-02")
	rate := rates[dt]
	if rate != 8.5749 {
		t.Errorf("expect %f got %f", 8.5749, rate)
	}

	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	rh := newRateHistory()
	for date, rate := range rates {
		rh.storeCurrencyDateRate("SEK", date, float32(rate))
	}

	date, _ := parseDate("2016-02-02")
	if rate, ok := rh.getRate("SEK", date); !ok || rate != 8.5749 {
		t.Errorf("expect 8.5749 for SEK-2016-02-02")
	}

	date, _ = parseDate("2016-02-12")
	if _, ok := rh.getRate("SEK", date); ok {
		t.Errorf("expect not found")
	}

	date, _ = parseDate("2016-02-02")
	if _, ok := rh.getRate("USD", date); ok {
		t.Errorf("expect not found")
	}

	expectDate, _ := parseDate("2016-02-04")
	if date, rate, ok := rh.getLastRate("SEK"); !ok || date != expectDate || rate != 8.3983 {
		t.Errorf("Expect 2016-02-04 8.3983 being last exchange rate")
	}
}
