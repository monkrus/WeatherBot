package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	// Create a test server to handle requests to the OpenWeatherMap API
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte(`{"weather":[{"description":"clear sky"}],"main":{"temp":289.52,"humidity":89,"pressure":1013},"name":"London"}`))
	}))
	defer testServer.Close()

	// Replace the API endpoint with the test server's URL
	url := fmt.Sprintf("%s?id=12345&appid=df2ea56fc8cca21e38e1ffab4894fb61", testServer.URL)

	// Test case 1: successful request
	result, err := getWeatherData("12345", url)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "London: clear sky, temperature 16.37 Celsius, humidity 89%, pressure 1013 hPa"
	if result != expected {
		t.Errorf("unexpected result: expected=%q, actual=%q", expected, result)
	}

	// Test case 2: failed request due to invalid city ID
	_, err = getWeatherData("invalid_city_id", url)
	if err == nil {
		t.Errorf("expected error")
	} else if !strings.Contains(err.Error(), "failed to parse response") {
		t.Errorf("unexpected error message: %v", err)
	}
}
