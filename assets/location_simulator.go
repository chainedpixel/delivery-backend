package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type LocationUpdate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

const (
	baseLat  = 13.6929
	baseLng  = -89.2182
	maxDelta = 0.01
)

func main() {

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjaWQiOiJhMGVlYmM5OS05YzBiLTRlZjgtYmI2ZC02YmI5YmQzODBhMTEiLCJleHAiOjE3NDg3MTM4OTMsImlhdCI6MTc0NzQxNzg5Mywicm9sZSI6IkFETUlOIiwic3ViIjoiYjJjM2Q0ZTUtZjZhNy04YjljLTBkMWUtMmYzYTRiNWM2ZDdlIn0.quHe6wqTDDoW6k7MsAt6zB2G7Tvcf-9RKCdZ2go0i20"
	baseURL := "http://localhost:7319/api/v1"
	orderID := "0c97a632-4e53-4cd5-a8e0-9109bf55deee"

	latitude := baseLat + (rand.Float64() * 0.1) - 0.05
	longitude := baseLng + (rand.Float64() * 0.1) - 0.05

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	endTime := time.Now().Add(5 * time.Minute)
	for time.Now().Before(endTime) {

		latitude += (rand.Float64() * maxDelta) - (maxDelta / 2)
		longitude += (rand.Float64() * maxDelta) - (maxDelta / 2)

		update := LocationUpdate{
			Latitude:  latitude,
			Longitude: longitude,
		}

		updateJSON, err := json.Marshal(update)
		if err != nil {
			log.Printf("Error encoding location: %v", err)
			continue
		}

		url := baseURL + "/tracking/location/" + orderID
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(updateJSON))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending location update: %v", err)
		} else {
			log.Printf("Sent location update: Lat %.6f, Lng %.6f, Status: %d",
				latitude, longitude, resp.StatusCode)
			resp.Body.Close()
		}

		time.Sleep(5 * time.Second)
	}
}
