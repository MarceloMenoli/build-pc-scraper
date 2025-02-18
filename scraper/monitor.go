package scraper

import (
	"log"
	"strconv"
	"strings"

	"build-pc-scraper/sms"
)

// converterPreco converte o preço (string) para float64.
// Exemplo: "R$ 599,00" -> 599.00
func converterPreco(precoStr string) (float64, error) {
	limpo := strings.ReplaceAll(precoStr, "R$", "")
	limpo = strings.TrimSpace(limpo)
	limpo = strings.ReplaceAll(limpo, ".", "")  // remove separador de milhar, se houver
	limpo = strings.ReplaceAll(limpo, ",", ".") // substitui a vírgula pelo ponto decimal
	return strconv.ParseFloat(limpo, 64)
}

// VerificaPrecos monitora os produtos e envia um SMS se o preço de um processador que contenha
// "Ryzen 5 7600" estiver abaixo de um valor definido.
func VerificaPrecos() {
	processadorAlvo := "Ryzen 5 7600"
	valorAlvo := 1700.00 // valor alvo em reais
	prods := GetProducts()
	for _, prod := range prods {
		if strings.Contains(prod.Nome, processadorAlvo) {
			preco, err := converterPreco(prod.Preco)
			log.Println("Preço: ", preco)

			if err != nil {
				log.Println("Erro ao converter o preço:", err)
				continue
			}
			if preco < valorAlvo {
				log.Println("Preço: ", preco, "valor alvo: ", valorAlvo)

				msg := "Alerta: O preço do " + prod.Nome + " está abaixo de R$" +
					strconv.FormatFloat(valorAlvo, 'f', 2, 64) + "! Preço atual: " + prod.Preco
				// Substitua pelo número de destino desejado
				if err := sms.SendSMS(msg, "+5511945725712"); err != nil {
					log.Println("Erro ao enviar SMS:", err)
				}
			}
		}
	}
}
