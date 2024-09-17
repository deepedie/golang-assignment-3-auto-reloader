package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type WeatherData struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

func main() {
	ticker := time.NewTicker(15 * time.Second)
	for range ticker.C {
		water := rand.Intn(100) + 1
		wind := rand.Intn(100) + 1

		weather := WeatherData{
			Water: water,
			Wind:  wind,
		}

		data, err := json.Marshal(weather)
		if err != nil {
			logrus.Error("Failed to marshal weather data: ", err)
			continue
		}

		resp, err := http.Post("http://localhost:8081/weather", "application/json", bytes.NewBuffer(data))
		if err != nil {
			logrus.Error("Failed to send data: ", err)
			continue
		}
		defer resp.Body.Close()

		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			logrus.Error("Failed to decode response: ", err)
			continue
		}

		// Log status
		logrus.Infof("Sent data: %+v", weather)
		logrus.Infof("Received response: %+v", result)

		status := result["status"].(map[string]interface{})
		fmt.Printf(`
		{
		   "water": %d,
		   "wind": %d
		}
		
		status water: %s
		status wind: %s
		`, weather.Water, weather.Wind, status["water_status"], status["wind_status"])
	}
}
