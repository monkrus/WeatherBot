package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Test cases for the main function

// Test case to check if the program can read the .env file
func TestLoadEnvFile(t *testing.T) {
	err := os.Setenv("BOT_TOKEN", "test")
	assert.NoError(t, err)
	err = os.Setenv("API_KEY", "test")
	assert.NoError(t, err)

	path_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	envFile := filepath.Join(path_dir, ".env")
	err = godotenv.Load(envFile)
	assert.NoError(t, err)
}

// Test case to check if the bot token and API key can be retrieved from the .env file
func TestGetEnvVars(t *testing.T) {
	err := os.Setenv("BOT_TOKEN", "test")
	assert.NoError(t, err)
	err = os.Setenv("API_KEY", "test")
	assert.NoError(t, err)

	botToken := os.Getenv("BOT_TOKEN")
	apiKey := os.Getenv("API_KEY")

	assert.Equal(t, "test", botToken)
	assert.Equal(t, "test", apiKey)
}

// Test case to check if the program can handle missing .env file
func TestLoadEnvFileMissing(t *testing.T) {
	path_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	envFile := filepath.Join(path_dir, "missing.env")
	err = godotenv.Load(envFile)
	assert.Error(t, err)
}
func TestGetWeatherData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"name": "London","weather":[{"description":"clear sky"}],"main":{"temp": 289.52,"humidity": 89,"pressure": 1013}}`)
	}))
	defer ts.Close()

	url := fmt.Sprintf("%s/weather/%s", ts.URL, "invalid_city_id")
	_, err := getWeatherData("invalid_city_id", url)

	if err == nil {
		t.Error("expected error")
	}
}
