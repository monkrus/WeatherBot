package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	
)

const (
	webhookURL        = "/webhook"
	openWeatherMapURL = "https://api.openweathermap.org/data/2.5/weather"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create a new bot instance
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Set up a webhook to receive updates
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(os.Getenv("NGROK_URL") + webhookURL + "/" + bot.Token))
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP server to listen for updates
	http.HandleFunc(webhookURL+"/"+bot.Token, func(w http.ResponseWriter, r *http.Request) {
		update := tgbotapi.Update{}
		err := json.NewDecoder(r.Body).Decode(&update)
		if err != nil {
			log.Println(err)
			return
		}

		if update.Message == nil {
			return
		}

		// Check if the message is a location
		if update.Message.Location != nil {
			// Get the weather for the location
			weather, err := getWeather(update.Message.Location.Latitude, update.Message.Location.Longitude, os.Getenv("OPENWEATHERMAP_API_KEY"))
			if err != nil {
				log.Println(err)
				return
			}

			// Send the weather as a message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, weather)
			bot.Send(msg)
		}
	})

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getWeather(lat, long float64, apiKey string) (string, error) {
	// Build the URL for the OpenWeatherMap API request
	url := fmt.Sprintf("%s?lat=%f&lon=%f&appid=%s", openWeatherMapURL, lat, long, apiKey)

	// Send the request and parse the response
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openweathermap API returned non-200 status code: %d", resp.StatusCode)
	}

	var weatherResp struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return "", err
	}

	// Convert temperature from Kelvin to Celsius
	temperature := weatherResp.Main.Temp - 273.15

	return fmt.Sprintf("The weather for your location is %.1f°C", temperature), nil
}
