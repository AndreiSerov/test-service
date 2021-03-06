package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// DB CONNECT
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongodb:27017"))
    if err != nil {
        log.Fatal(err)
    }
    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)

	// DB interaction
    testDatabase := client.Database("test")
    currencyCollection := testDatabase.Collection("curencies")
	for _, target := range c {
		targetResult, err := currencyCollection.InsertOne(ctx, target)
		if err != nil { 
			fmt.Println("InsertOne ERROR:", err)
			return
		} 
		fmt.Printf("targetResult: %v\n", targetResult)
		newID := targetResult.InsertedID 
		fmt.Println("InsertOne() newID:", newID) 
		fmt.Println("InsertOne() newID type:", reflect.TypeOf(newID)) 
	}

	// Map response
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
