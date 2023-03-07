package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/tidwall/gjson"
)

func main() {
	// Replace YOUR_TELEGRAM_BOT_TOKEN with your actual bot token
	bot, err := tgbotapi.NewBotAPI("6298432449:AAFefsa7SWWBeOuenk6w-wll2ioW6G42Djc")
	if err != nil {
		log.Panic(err)
	}

	// Replace YOUR_OPENWEATHERMAP_API_KEY with your actual API key
	apiKey := "df2ea56fc8cca21e38e1ffab4894fb61"

	// Set up an HTTP client to make API requests to OpenWeatherMap
	client := &http.Client{}

	// Set up a message handler for the /weather command
	updates, err := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if err != nil {
		log.Panic(err)
	}

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
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, an error occurred. Please try again later."))
				continue
			}
			defer resp.Body.Close()

			// Parse the response JSON and get the temperature and description
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Sorry, an error occurred. Please try again later."))
				continue
			}
			temp := gjson.GetBytes(body, "main.temp").Float()
			desc := gjson.GetBytes(body, "weather.0.description").String()

			// Convert the temperature from Kelvin to Celsius and format it
			tempC := temp - 273.15
			tempStr := strconv.FormatFloat(tempC, 'f', 1, 64)

			// Send the weather forecast to the user
			message := "The temperature in " + location + " is " + tempStr + "°C and the weather is " + desc + "."
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
		}
	}

	// Start the bot
	log.Printf("Starting bot %s", bot.Self.UserName)
}
