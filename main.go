package main

import (
	"encoding/json" // Tratamento de dados em JSON / Data processing in JSON
	"fmt"           // Formatação e impressão de dados / Data formatting and processing
	"log"           // Registro de mensagens no terminal ou em arquivos, com timestamps / Messages logging to terminal or files, with timestamps
	"net/http"      // Cliente e Servidor HTTP (GET, POST, PUT, DELETE) / HTTP client and server
	"os"            // Interação com funcionalidades do sistema / Interaction with system features
	"strconv"       // Conversão de string / string conversion

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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
	log.Printf("Bot %s iniciado! Você já pode começar a usar", bot.Self.UserName)

	// Configura recebimento de mensagens
	// Configure receiving messages
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("erro ao iniciar canal de updates:", err)
	}

	moedas := map[string]struct {
		pair string
		key  string
		nome string
	}{
		"/dolar": {"USD-BRL", "USDBRL", "dólar"},
		"/euro":  {"EUR-BRL", "EURBRL", "euro"},
		"/libra": {"GBP-BRL", "GBPBRL", "libra esterlina"},
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if info, ok := moedas[update.Message.Text]; ok {
			valor, err := getExchangeRate(info.pair, info.key)
			if err != nil {
				valor = "Não foi possível buscar a cotação 😢"
			}
			exchangeRateFloat, err := strconv.ParseFloat(valor, 64)
			var resposta string
			if err != nil {
				resposta = fmt.Sprintf("💵 Cotação atual do %s: Não foi possível formatar a cotação 😢", info.nome)
			} else {
				resposta = fmt.Sprintf("💵 Cotação atual do %s: R$ %.2f", info.nome, exchangeRateFloat)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resposta)
			bot.Send(msg)
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Olá, veja as cotações através dos comandos /dolar, /euro, /libra.")
		bot.Send(msg)
	}
}

// Função de API que mostra a cotação atual de uma moeda
// API's function that shows the current rate of a currency
func getExchangeRate(pair, key string) (string, error) {
	url := fmt.Sprintf("https://economia.awesomeapi.com.br/json/last/%s", pair)
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

	valor, ok := resultado[key]["bid"].(string)
	if !ok {
		return "", fmt.Errorf("não consegui ler o valor")
	}

	return valor, nil
}
