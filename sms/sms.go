package sms

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// SendSMS envia uma mensagem SMS utilizando a API do Twilio.
func SendSMS(body, to string) error {
	accountSID := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_PHONE_NUMBER")

	if accountSID == "" || authToken == "" || from == "" {
		return fmt.Errorf("credenciais do Twilio não configuradas")
	}

	// Monta os parâmetros da requisição
	msgData := url.Values{}
	msgData.Set("To", to)
	msgData.Set("From", from)
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Cria a requisição para o endpoint do Twilio
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json"
	req, err := http.NewRequest("POST", urlStr, &msgDataReader)
	if err != nil {
		return err
	}
	req.SetBasicAuth(accountSID, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Executa a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("SMS enviado com sucesso!")
		return nil
	}
	return fmt.Errorf("falha ao enviar SMS, status: %d", resp.StatusCode)
}
