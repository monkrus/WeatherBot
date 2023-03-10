package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func getWeatherData(cityID string, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response data
	var data struct {
		Name    string `json:"name"`
		Weather []struct {
			Description string `json:"description"`
		} `json:"weather"`
		Main struct {
			Temperature float64 `json:"temp"`
			Humidity    int     `json:"humidity"`
			Pressure    int     `json:"pressure"`
		} `json:"main"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Format the weather data string
	return fmt.Sprintf("%s: %s, temperature %.2f Kelvin, humidity %d%%, pressure %d hPa",
		data.Name, data.Weather[0].Description, data.Main.Temperature, data.Main.Humidity, data.Main.Pressure), nil
}
