package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/tidwall/gjson"
)

func main() {
	path_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	envFile := filepath.Join(path_dir, ".env")
	log.Println("Loading", envFile, "file...")

	err = godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	botToken := os.Getenv("BOT_TOKEN")
	apiKey := os.Getenv("API_KEY")

	// Replace BOT_TOKEN with your actual bot token
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	// Replace YOUR_OPENWEATHERMAP_API_KEY with your actual API key
	//apiKey := "API_KEY"

	// Set up an HTTP client to make API requests to OpenWeatherMap
	client := &http.Client{}

	// Set up a message handler for the /weather command
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		log.Panic(err)
	}

	// Defer the closing of the HTTP client at the end of the function
	defer client.CloseIdleConnections()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if strings.HasPrefix(update.Message.Text, "/weather") {
			// Get the location from the message text
			location := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/weather"))

			// Build the API request URL
			url := "https://api.openweathermap.org/data/2.5/weather?q=" + location + "&appid=" + apiKey

			// Make the API request and get the response
			resp, err := client.Get(url)
			if err != nil {
				log.Println(err)
				_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, an error occurred. Please try again later."))
				if err != nil {
					log.Println(err)
				}
				continue
			}

			// Parse the response JSON and get the temperature and description
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
				_, err = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, an error occurred. Please try again later."))
				if err != nil {
					log.Println(err)
				}
				//bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, an error occurred. Please try again later."))
				continue
			}
			temp := gjson.GetBytes(body, "main.temp").Float()
			desc := gjson.GetBytes(body, "weather.0.description").String()

			// Convert the temperature from Kelvin to Celsius and format it
			tempC := temp - 273.15
			tempStr := strconv.FormatFloat(tempC, 'f', 1, 64)

			// Send the weather forecast to the user
			message := "The temperature in " + location + " is " + tempStr + "°C and the weather is " + desc + "."
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)

			}
		}
	}
}
