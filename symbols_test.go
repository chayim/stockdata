package stockdata_test

import (
	"testing"

	"github.com/chayim/stockdata"
)

func validateSymbols(t *testing.T, symbols []stockdata.Symbol) {
	length := len(symbols)
	if length < 200 {
		t.Errorf("Failed to return a reasonable minimum number of stock symbols (%d returned).", length)
	}
	for _, item := range symbols {
		if len(item.Name) == 0 {
			t.Error("Fetched an invalid stock from the NYSE.")
		}

		if len(item.Symbol) == 0 || len(item.Symbol) >= 10 {
			t.Errorf("Fetched an invalid stock (%s) from the NYSE.", item.Symbol)
		}
	}
}

func TestTSXETFs(t *testing.T) {
	symbols := stockdata.FetchTSXETFs()
	if len(symbols) < 50 {
		t.Errorf("Only %d symbols retrieved", len(symbols))
	}
	for _, item := range symbols {
		if len(item) <= 1 && len(item) >= 5 {
			t.Errorf("%s is an invalid symbol", item)
		}
	}
}

func TestFetchingTSX(t *testing.T) {
	symbols := stockdata.FetchTSXSymbols()
	validateSymbols(t, symbols)
}

func TestFetchingNYSE(t *testing.T) {
	symbols := stockdata.FetchNYSESymbols()
	validateSymbols(t, symbols)
}

func TestFetchingAMEX(t *testing.T) {
	symbols := stockdata.FetchAMEXSymbols()
	validateSymbols(t, symbols)
}

func TestFetchingNasdaq(t *testing.T) {
	symbols := stockdata.FetchNasdaqSymbols()
	validateSymbols(t, symbols)
}
