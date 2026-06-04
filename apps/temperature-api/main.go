package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"time"
    "context"
    "log"
    "os/signal"
    "syscall"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/temperature", temperatureHandler)

	port := os.Getenv("PORT")
	if port == "" {
        port = "8081"
    }
	
  server := &http.Server{
        Addr: ":" + port,
    }

    go func() {
        sigint := make(chan os.Signal, 1)
        signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
        <-sigint

        if err := server.Shutdown(context.Background()); err != nil {
            log.Printf("Server shutdown error: %v", err)
        }
    }()

    log.Printf("Server is running on port %s", port)
    if err := server.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatalf("HTTP server error: %v", err)
    }
    log.Println("Server stopped")

}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")

	if location == "" {
		switch sensorID {
		case "1":
			location = "Living Room"
		case "2":
			location = "Bedroom"
		case "3":
			location = "Kitchen"
		default:
			location = "Unknown"
		}
	}

	if sensorID == "" {
		switch location {
		case "Living Room":
			sensorID = "1"
		case "Bedroom":
			sensorID = "2"
		case "Kitchen":
			sensorID = "3"
		default:
			sensorID = "0"
		}
	}

	response := struct {
		Value       float64   `json:"value"`
		Unit        string    `json:"unit"`
		Timestamp   time.Time `json:"timestamp"`
		Location    string    `json:"location"`
		Status      string    `json:"status"`
		SensorID    string    `json:"sensor_id"`
		SensorType  string    `json:"sensor_type"`
		Description string    `json:"description"`
	}{
		Value:       20.0 + (rand.Float64()*10 - 5),
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Temperature sensor in " + location,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}