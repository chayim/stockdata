package stockdata_test

import (
	"testing"
	"time"

	"github.com/chayim/stockdata"
)

// test fetching future exdivs
func TestExDividendCAData(t *testing.T) {
	data := stockdata.FetchTSXExDividends()
	if len(data) <= 1 {
		t.Error("Failed to parse exdividend data.")
	}
	for _, item := range data {
		if item.Symbol == "" {
			t.Error("Failed to parse symbol for exdivs")
		}

		if item.Yield <= 0.0 {
			t.Errorf("Invalid dividend for symbol %s.", item.Symbol)
		}

		if item.ExDividendDate.Before(time.Now()) {
			t.Errorf("Invalid date %s provided in exdivdiend history data for %s.", item.ExDividendDate, item.Symbol)
		}
	}
}
