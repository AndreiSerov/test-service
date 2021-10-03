package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CurrencyData struct {
	Price24h       float64 `json:"price_24h"`
	Volume24h      float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trade_price"`
}

type Currency struct {
	Symbol string `json:"symbol"`
	CurrencyData
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var c []Currency
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Can not unmarshal JSON")
		return
	}
	json.Unmarshal(body, &c)

	// save in REPO
	var m = make(map[string]CurrencyData)
	for _, s := range c {
		m[s.Symbol] = CurrencyData{
			s.Price24h,
			s.Volume24h,
			s.LastTradePrice,
		}
	}

	if resp, err := json.Marshal(m); err != nil {
		log.Println(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func main() {
	http.HandleFunc("/", PostHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
