package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

//Date year{4 bits} ,month{4bits}, day{5 bits}
type Date uint32

//Rate exchange rate, 32bit.32bit
type Rate float32

//DateRateMap specific currency date-exchagne values
type DateRateMap map[Date]Rate

//CurrencyHistory store history for all currencies
type CurrencyHistory map[string]DateRateMap

//DateCurrencyRate shows rates for all currencies for the same date
type DateCurrentyRate map[string]Rate

//RateHistory is our main storage
type RateHistory struct {
	history CurrencyHistory
	mutex   sync.RWMutex

	writer *os.File
}

func (rh *RateHistory) OpenWriter() bool {
	f, err := os.OpenFile("_history.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		logger.Printf("cannot open history file: %s", err)
		return false
	}

	finfo, err := f.Stat()
	if err != nil {
		logger.Printf("cannot stat history file: %s", err)
		return false
	}
	if finfo.Size() == 0 {
		f.WriteString("curreny,date,rate\n")
	}
	logger.Printf("opened _history with total %d bytes", finfo.Size())
	rh.writer = f
	return true
}

func newRateHistory() RateHistory {
	return RateHistory{
		history: make(CurrencyHistory),
	}
}

func (rh *RateHistory) getDateRates(date Date) DateCurrentyRate {
	rates := make(DateCurrentyRate)

	rh.mutex.RLock()
	defer rh.mutex.RUnlock()

	for cur, h := range rh.history {
		if rate, ok := h[date]; ok {
			rates[cur] = rate
		}
	}
	return rates
}

func (rh *RateHistory) getLastRate(cur string) (Date, Rate, bool) {
	rh.mutex.RLock()
	defer rh.mutex.RUnlock()

	h, ok := rh.history[cur]
	if !ok {
		return 0, 0, false
	}

	maxDate := Date(0)
	found := false
	var maxDateRate Rate
	for date, rate := range h {
		if date > maxDate {
			maxDate = date
			maxDateRate = rate
			found = true
		}
	}
	return maxDate, maxDateRate, found
}

func (rh *RateHistory) getRate(cur string, date Date) (Rate, bool) {
	rh.mutex.RLock()
	defer rh.mutex.RUnlock()

	h, ok := rh.history[cur]
	if !ok {
		return 0, false
	}

	rate, found := h[date]
	return rate, found
}

func (rh *RateHistory) getRatesFrom(cur string, date Date) DateRateMap {
	rh.mutex.RLock()
	defer rh.mutex.RUnlock()

	h, ok := rh.history[cur]
	if !ok {
		return nil
	}

	outData := make(DateRateMap)
	for dt, rate := range h {
		if dt < date {
			continue
		}
		outData[dt] = rate
	}
	return outData
}

func (rh *RateHistory) getRatesBetween(cur string, fromDate Date, toDate Date) DateRateMap {
	rh.mutex.RLock()
	defer rh.mutex.RUnlock()

	h, ok := rh.history[cur]
	if !ok {
		return nil
	}

	outData := make(DateRateMap)
	for dt, rate := range h {
		if dt < fromDate {
			continue
		}
		if dt >= toDate {
			continue
		}

		outData[dt] = rate
	}
	return outData
}

func (rh *RateHistory) loadFrom(filename string) error {
	curname, data, err := parseCsv(filename)
	if err != nil {
		return err
	}
	rh.mutex.Lock()
	defer rh.mutex.Unlock()
	if existingCurrencyHistory, ok := rh.history[curname]; ok {
		for date, rate := range data {
			existingCurrencyHistory[date] = rate
		}
	} else {
		rh.history[curname] = data
		logger.Printf("loaded %d items for currency %s \n", len(data), curname)
	}
	return nil
}

func (rh *RateHistory) load() {
	files, _ := filepath.Glob("data/*.csv")
	for _, filename := range files {
		err := rh.loadFrom(filename)
		if err != nil {
			logger.Printf("load from %s failed %s", filename, err)
		}
	}
	rh.loadExtraHistory()
}

func (rh *RateHistory) loadExtraHistory() error {
	history, err := parseExtraRates("_history.csv")
	if err != nil {
		return err
	}
	rh.mutex.Lock()
	defer rh.mutex.Unlock()

	for cur, h := range history {
		if exHistory, ok := rh.history[cur]; ok {
			for date, rate := range h {
				exHistory[date] = rate
			}
		} else {
			rh.history[cur] = h
		}

	}
	return nil
}

func (rh *RateHistory) storeCurrencyDateRate(currency string, date Date, rate float32) {
	logger.Printf("store %s at %s equal %f", currency, dateToStr(date), rate)
	rh.mutex.Lock()
	defer rh.mutex.Unlock()

	h, ok := rh.history[currency]
	if !ok {
		h = make(DateRateMap)
		rh.history[currency] = h
	}

	h[date] = Rate(rate)

	rh.appendCurrencyValue(currency, date, rate)
}

func (rh *RateHistory) appendCurrencyValue(currency string, date Date, rate float32) {
	if rh.writer == nil && !rh.OpenWriter() {
		// writer is not working, just skip it
		logger.Printf("Cannot store to file")
		return
	}

	s := fmt.Sprintf("%s,%s,%f\n", currency, dateToStr(date), rate)
	logger.Printf("storing to file %s", s)
	rh.writer.WriteString(s)
	rh.writer.Sync()
}
