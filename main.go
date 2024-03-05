package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

type PoolData struct {
	ID           uint64 `json:"ID"`
	USDLiquidity string `json:"USDLiquidity"`
	Type         string `json:"Type"`
}

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	http.HandleFunc("/pool-liquidity", func(w http.ResponseWriter, r *http.Request) {
		val, err := rdb.Get(ctx, "osmosis-pool-liquidities").Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Parse the JSON string into a slice of PoolData
		var arrPoolData []PoolData
		err = json.Unmarshal([]byte(val), &arrPoolData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(arrPoolData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
