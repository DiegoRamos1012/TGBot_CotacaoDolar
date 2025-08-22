package main

import (
	"encoding/json" // Tratamento de dados em JSON / Data processing in JSON
	"fmt"           // Formatação e impressão de dados / Data formatting and processing
	"log"           // Registro de mensagens no terminal ou em arquivos, com timestamps / Messages logging to terminal or files, with timestamps
	"net/http"      // Cliente e Servidor HTTP (GET, POST, PUT, DELETE) / HTTP client and server
	"os"            // Interação com funcionalidades do sistema / Interaction with system features
	"strconv"       // Conversão de string / string conversion

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

/* Projeto de Bot no Telegram que mostra a cotação atual do dólar em real brasileiro através de uma API pública */
/* Telegram Bot project that shows the current exchange rate of the dollar in Brazilian real through a public API*/

func main() {
	// Carrega o .env do projeto e utiliza a variável que contém o token do bot
	// Loads the project's env and uses the variable what containing the bot's token
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o .env")
	}

	token := os.Getenv("TGBOT_TOKEN")
	if token == "" {
		log.Fatal("Token do bot não definido!")
	}

	// Conecta com o bot "Sábio do Dólar"
	// Starts connection with bot "Sábio do Dólar"
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Bot %s iniciado!", bot.Self.UserName)

	// Configura recebimento de mensagens
	// Configure receiving messages
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("erro ao iniciar canal de updates:", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case "/dolar":
			// Pega cotação via função "getExchangeRateDollarBRL"
			// Gets exchange Rate by function "getExchangeRateDollarBRL"
			exchangeRateDollarBRL, err := getExchangeRateDollarBRL()
			if err != nil {
				exchangeRateDollarBRL = "Não foi possível buscar a cotação 😢"
			}

			// Converte a string para float64
			// Convert string to float64
			exchangeRateFloat, err := strconv.ParseFloat(exchangeRateDollarBRL, 64)
			if err != nil {
				exchangeRateDollarBRL = "Não foi possível formatar a cotação 😢"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "💵 Cotação do dólar: " + exchangeRateDollarBRL)
				bot.Send(msg)
				continue
			}
			exchangeRateDolBRLFormatted := fmt.Sprintf("R$ %.2f", exchangeRateFloat)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "💵 Cotação atual do dólar: " + exchangeRateDolBRLFormatted)
			bot.Send(msg)

		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá, envie /dolar para ver a cotação atual do dólar")
			bot.Send(msg)
		}
	}
}

// Função de API
// API's function
func getExchangeRateDollarBRL() (string, error) {
	// Public API
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status: %d", resp.StatusCode)
	}

	var resultado map[string]map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&resultado); err != nil {
		return "", err
	}

	// O JSON retorna: {"USDBRL": {"bid": "5.48", ...}}
	// The JSON returns: {"USDBRL": {"bid": "5.48", ...}}
	valor, ok := resultado["USDBRL"]["bid"].(string)
	if !ok {
		return "", fmt.Errorf("não consegui ler o valor")
	}

	return valor, nil
}
