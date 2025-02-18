package scraper

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

// Product representa um produto com nome e preço.
type Product struct {
	Nome  string `json:"nome"`
	Preco string `json:"preco"`
}

var (
	products []Product
	mu       sync.RWMutex
)

// UpdateProducts realiza o scraping e atualiza a lista de produtos.
func UpdateProducts() error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ExecPath("/usr/bin/chromium"), // ou use a variável de ambiente CHROME_PATH se preferir
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Cria um contexto para o chromedp com timeout.
	ctx, cancel := chromedp.NewContext(allocCtx)
	ctx, cancel = context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	var newProducts []Product

	// URL da primeira página.
	baseURL := "https://www.kabum.com.br/hardware/processadores/processador-amd?page_number=1&page_size=20&facet_filters=&sort=most_searched"

	// Navega para a página inicial.
	if err := chromedp.Run(ctx, chromedp.Navigate(baseURL)); err != nil {
		log.Println("Erro ao acessar a página inicial:", err)
		return err
	}

	pageNumber := 1
	for {
		// Aguarda que os produtos estejam visíveis na página.
		if err := chromedp.Run(ctx, chromedp.WaitVisible(`a.productLink`, chromedp.ByQuery)); err != nil {
			log.Printf("Nenhum produto encontrado ou timeout na página %d: %v\n", pageNumber, err)
			break
		}

		// Extrai os nomes e preços dos produtos.
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
			log.Printf("Erro ao extrair dados da página %d: %v\n", pageNumber, err)
			break
		}

		// Se não houver produtos, encerra o loop.
		if len(pageNames) == 0 {
			log.Printf("Nenhum produto encontrado na página %d. Encerrando.\n", pageNumber)
			break
		}

		// Agrupa os dados extraídos em objetos Product.
		for i, name := range pageNames {
			price := ""
			if i < len(pagePrices) {
				price = pagePrices[i]
			}
			newProducts = append(newProducts, Product{Nome: name, Preco: price})
		}

		// Verifica se existe o botão "nextLink" habilitado.
		var hasNext bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(
				`document.querySelector('a.nextLink[aria-disabled="false"]') !== null`,
				&hasNext,
			),
		)
		if err != nil || !hasNext {
			log.Println("Não há próxima página. Encerrando atualização.")
			break
		}

		// Clica no botão "nextLink" para ir para a próxima página e aguarda o carregamento.
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

	// Atualiza a variável global com os novos dados.
	mu.Lock()
	products = newProducts
	mu.Unlock()

	log.Printf("Atualização concluída. Produtos atualizados: %d\n", len(newProducts))
	return nil
}

// GetProducts retorna uma cópia da lista de produtos.
func GetProducts() []Product {
	mu.RLock()
	defer mu.RUnlock()
	// Retorna uma cópia para evitar data races.
	result := make([]Product, len(products))
	copy(result, products)
	return result
}
