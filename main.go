package main

import (
	"log"
	"net/http"
	"os"

	"build-pc-scraper/handlers"
	"build-pc-scraper/scraper"

	"github.com/robfig/cron/v3"
)

func main() {
	// Configura e inicia o cron job para atualizar os produtos a cada 30 minutos.
	c := cron.New()
	_, err := c.AddFunc("*/30 * * * *", func() {
		if err := scraper.UpdateProducts(); err != nil {
			log.Println("Erro ao atualizar produtos:", err)
		}
	})
	if err != nil {
		log.Fatal("Erro ao agendar o cron job:", err)
	}
	c.Start()
	defer c.Stop()

	// Executa uma atualização inicial antes de iniciar o servidor.
	if err := scraper.UpdateProducts(); err != nil {
		log.Fatal("Erro na atualização inicial:", err)
	}

	// Configura o endpoint HTTP.
	http.HandleFunc("/produtos", handlers.ProductsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Println("Servidor rodando na porta", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}
}
