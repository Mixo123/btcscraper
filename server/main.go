package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/rs/cors"
)

type Price struct {
	BTC string `json:"BTC"`
	LTC string `json:"LTC"`
	ADA string `json:"ADA"`
}

type CoinInfo struct {
	Prices Price `json:"prices"`
}

var (
	coinInfo CoinInfo
	mu       sync.Mutex
)

func scrapeCoinData(symbol string) {
	site := "https://coinmarketcap.com/currencies/%s"
	fullUrl := fmt.Sprintf(site, symbol)
	c := colly.NewCollector(colly.AllowedDomains("www.coinmarketcap.com", "coinmarketcap.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML("span.clvjgF", func(h *colly.HTMLElement) {
		switch symbol {
			case "bitcoin":
				mu.Lock()
				coinInfo.Prices.BTC = h.Text
				mu.Unlock()
			case "litecoin":
				mu.Lock()
				coinInfo.Prices.LTC = h.Text
				mu.Unlock()
			case "cardano":
				mu.Lock()
				coinInfo.Prices.ADA = h.Text
				mu.Unlock()
		}
	})

	c.OnScraped(func(r *colly.Response) {
		mu.Lock()
		fmt.Printf("\rBTC Price: %s, LTC Price: %s, ADA Price: %s", coinInfo.Prices.BTC, coinInfo.Prices.LTC, coinInfo.Prices.ADA)
		mu.Unlock()
	})

	c.Visit(fullUrl)
}

func coinHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coinInfo)
}

func main() {
	go func() {
		for {
			scrapeCoinData("bitcoin")
			scrapeCoinData("litecoin")
			scrapeCoinData("cardano")
			time.Sleep(10 * time.Second)
		}
	}()

	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}).Handler(http.DefaultServeMux)

	http.HandleFunc("/coins", coinHandler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", handler)
}
