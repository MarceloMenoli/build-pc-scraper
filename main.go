package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

// Produto representa as informações de um produto
type Produto struct {
	Nome  string `json:"nome"`
	Preco string `json:"preco"`
}

func main() {
	// URL da primeira página
	baseURL := "https://www.kabum.com.br/hardware/processadores/processador-amd?page_number=1&page_size=20&facet_filters=&sort=most_searched"

	// Cria o contexto do chromedp
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Define um timeout global para evitar travamentos
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	// Navega para a primeira página
	if err := chromedp.Run(ctx, chromedp.Navigate(baseURL)); err != nil {
		log.Fatal("Erro ao acessar a página inicial:", err)
	}

	pageNumber := 1
	var produtos []Produto

	for {
		// Aguarda que os produtos estejam visíveis na página
		if err := chromedp.Run(ctx,
			chromedp.WaitVisible(`a.productLink`, chromedp.ByQuery),
		); err != nil {
			log.Println("Nenhum produto encontrado na página ou timeout:", err)
			break
		}

		// Extrai os nomes e preços dos produtos
		var pageNames, pagePrices []string
		err := chromedp.Run(ctx,
			chromedp.Evaluate(
				`Array.from(document.querySelectorAll("a.productLink span.nameCard")).map(el => el.innerText.trim())`,
				&pageNames,
			),
			chromedp.Evaluate(
				`Array.from(document.querySelectorAll("a.productLink span.priceCard")).map(el => el.innerText.trim())`,
				&pagePrices,
			),
		)
		if err != nil {
			log.Println("Erro ao extrair dados da página:", err)
			break
		}

		// Se não houver produtos na página, encerra o loop
		if len(pageNames) == 0 {
			log.Printf("Nenhum produto encontrado na página %d. Encerrando.\n", pageNumber)
			break
		}

		// Organiza os dados extraídos em objetos Produto
		for i, nome := range pageNames {
			preco := ""
			if i < len(pagePrices) {
				preco = pagePrices[i]
			}
			produtos = append(produtos, Produto{
				Nome:  nome,
				Preco: preco,
			})
		}

		// Verifica se o botão "nextLink" existe e está habilitado (aria-disabled="false")
		var hasNext bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(
				`document.querySelector('a.nextLink[aria-disabled="false"]') !== null`, &hasNext,
			),
		)
		if err != nil {
			log.Println("Erro ao verificar botão 'nextLink':", err)
			break
		}

		if !hasNext {
			log.Println("Botão 'nextLink' não encontrado ou está desabilitado. Encerrando a iteração.")
			break
		}

		// Clica no botão "nextLink" para ir para a próxima página e aguarda o carregamento
		err = chromedp.Run(ctx,
			chromedp.Click(`a.nextLink`, chromedp.ByQuery),
			chromedp.Sleep(2*time.Second),
		)
		if err != nil {
			log.Println("Erro ao clicar no botão 'nextLink':", err)
			break
		}

		pageNumber++
	}

	// Converte os dados coletados para JSON
	resultJSON, err := json.MarshalIndent(produtos, "", "  ")
	if err != nil {
		log.Fatal("Erro ao converter os dados para JSON:", err)
	}

	// Exibe o JSON no console
	fmt.Println(string(resultJSON))
	log.Println("Processo de scraping finalizado.")
}
