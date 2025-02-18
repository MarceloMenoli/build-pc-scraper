package handlers

import (
	"encoding/json"
	"net/http"

	"build-pc-scraper/scraper"
)

// ProductsHandler expõe os produtos em formato JSON.
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	products := scraper.GetProducts()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
