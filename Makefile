BOT_TOKEN=your_telegram_bot_token_here
WEATHER_API_KEY=your_openweathermap_api_key_here

build:
	go build -o bot main.go

run:
	go run main.go

clean:
	rm -f bot

.PHONY: build run clean