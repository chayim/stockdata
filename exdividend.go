package stockdata

import (
	"os"
	"strings"
	"strconv"
	"time"
	"github.com/anaskhan96/soup"
)

func FetchTSXExDividends() []ExDividend {

	resp, err := soup.Get("https://tsx.exdividend.ca/")
	if err != nil {
		os.Exit(1)
	}

	m := []ExDividend{}

	doc := soup.HTMLParse(resp)
	divdata := doc.Find("table", "id", "stock-table").Find("tbody")
	rows := divdata.FindAll("tr")
	for _, row := range rows {
		cols := row.FindAll("td")

		// the dividend values in the table are strings of DIV currency (i.e 0.64 USD)
		if cols[7].Text() == "N/A" {
			continue
		}
		_divdata := strings.Fields(cols[7].Text())

		yield, err := strconv.ParseFloat(cols[5].Text(), 64)
		if err != nil {
			continue
		}
		dividend, _ := strconv.ParseFloat(_divdata[0], 64)

		var currency string
		if len(_divdata) == 2 {
			currency = _divdata[1]
		} else {
			currency = "CAD"
		}

		exdivdate, _ := time.Parse("Jan 2, 2006", cols[3].Text())

		// currently not in the tsx exdiv data, but hopefully one day (this is their bug)
		paydate, _ := time.Parse("Jan 2, 2006", cols[4].Text())
		symbol := strings.TrimSpace(cols[1].Find("a").Text())

		exd := ExDividend{
			Symbol:      symbol,
			PayableDate: paydate,
			ExDividendDate:   exdivdate,
			Yield:       	  yield,
			Currency:         currency,
			Dividend:    	  dividend,
		}

		m = append(m, exd)

	}

	return m
}
