package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

type BTCInfo struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

var (
	btcinfo BTCInfo
	mu      sync.Mutex
)

func scrapeBitcoinData() {
	c := colly.NewCollector(colly.AllowedDomains("www.coinmarketcap.com", "coinmarketcap.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scraping: %s\n", e.Error())
	})

	c.OnHTML("span.lsTl", func(h *colly.HTMLElement) {
		mu.Lock()
		btcinfo.Title = h.Text
		mu.Unlock()
	})

	c.OnHTML("span.clvjgF", func(h *colly.HTMLElement) {
		mu.Lock()
		btcinfo.Price = h.Text
		mu.Unlock()
	})

	c.OnScraped(func(r *colly.Response) {
		mu.Lock()
		fmt.Printf("\rTitle: %s, Price: %s", btcinfo.Title, btcinfo.Price)
		mu.Unlock()
	})

	c.Visit("https://coinmarketcap.com/currencies/bitcoin/")
}

func btcHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(btcinfo)
}

func main() {
	go func() {
		for {
			scrapeBitcoinData()
			time.Sleep(5 * time.Second)
		}
	}()

	http.HandleFunc("/btc", btcHandler)
	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}