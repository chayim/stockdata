package stockdata

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tealeg/xlsx"
)

// Given the way we fetch symbols and parse a list, we need to fetch symbols
// via a URL. This function make sit possible to do that.
func fetchTSXSymbolUrl(tsxURL string) []Symbol {

	type tmxjson struct {
		last_updated int
		length       int
		Results      []struct {
			Symbol      string
			Name        string
			Instruments []struct {
				Symbol string
				Name   string
			}
		}
	}

	var m []Symbol

	resp, err := http.Get(tsxURL)
	if err != nil {
		return nil
	}

	if resp.StatusCode > 204 && resp.StatusCode < 200 {
		return nil
	}
	defer resp.Body.Close()

	tmxobj := tmxjson{}

	json.NewDecoder(resp.Body).Decode(&tmxobj)
	for i := 0; i < len(tmxobj.Results); i++ {
		name, err := url.QueryUnescape(tmxobj.Results[i].Name)
		if err != nil {
			continue
		}
		stock := Symbol{
			Symbol: tmxobj.Results[i].Symbol,
			Name:   name,
		}
		m = append(m, stock)

		for inst := 0; inst < len(tmxobj.Results[i].Instruments); inst++ {
			name, err := url.QueryUnescape(tmxobj.Results[i].Instruments[inst].Name)

			if err != nil {
				continue
			}
			stock := Symbol{
				Symbol: tmxobj.Results[i].Instruments[inst].Symbol,
				Name:   name,
			}

			m = append(m, stock)
		}
	}

	return m
}

func FetchTSXETFs() []string {

	url := "https://www.tsx.com/resource/en/1168"

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	xlsFile, err := xlsx.OpenBinary(body)
	if err != nil {
		return nil
	}

	var etfs []string
	for _, sheet := range xlsFile.Sheets {

		// if Issuers is in the name - we're good to go
		if strings.Index(sheet.Name, "TSX") == -1 {
			continue
		}
		for _, row := range sheet.Rows {
			cells := row.Cells
			etfs = append(etfs, cells[2].String())
		}
	}

	return etfs

}

// Fetch symbols for the Toronto Stock Exchange
func FetchTSXSymbols() []Symbol {
	keys := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	baseURL := "https://www.tsx.com/json/company-directory/search/tsx/"
	var url string
	var m []Symbol

	for i := 0; i < len(keys); i++ {
		url = fmt.Sprintf("%s%s*", baseURL, keys[i:i+1])
		m = fetchTSXSymbolUrl(url)
	}
	return m
}

func FetchAMEXSymbols() []Symbol {
	return fetchUSSymbols("amex")
}

func FetchNYSESymbols() []Symbol {
	return fetchUSSymbols("nyse")
}

func FetchNasdaqSymbols() []Symbol {
	return fetchUSSymbols("nasdaq")
}

func fetchUSSymbols(exchange string) []Symbol {
	var m []Symbol

	url := fmt.Sprintf("https://old.nasdaq.com/screening/companies-by-name.aspx?letter=0&exchange=%s&render=download", strings.ToLower(exchange))

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	if resp.StatusCode > 204 && resp.StatusCode < 200 {
		return nil
	}

	defer resp.Body.Close()
	reader := csv.NewReader(resp.Body)
	body, err := reader.ReadAll()
	for i := 0; i < len(body); i++ {
		etf := false
		symbol := strings.TrimSpace(body[i][0])
		name := strings.TrimSpace(body[i][1])
		sector := body[i][5]
		if sector == "n/a" {
			sector = ""
		}
		industry := body[i][6]
		if industry == "n/a" {
			industry = ""
		}
		if strings.Index(name, "ETF") != -1 {
			etf = true
		}

		stock := Symbol{
			Sector:   sector,
			Industry: industry,
			Symbol:   symbol,
			Name:     name,
			ETF:      etf,
		}

		m = append(m, stock)
	}

	return m

}
