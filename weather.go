package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Weather struct {
	City        string  `json:"city"`
	Temperature float32 `json:"temperature"`
	Description string  `json:"description"`
}

func getWeather(city string) (*Weather, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	url := fmt.Sprintf("<https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric>", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	weather := &Weather{
		City:        city,
		Temperature: data["main"].(map[string]interface{})["temp"].(float32),
		Description: data["weather"].([]interface{})[0].(map[string]interface{})["description"].(string),
	}

	return weather, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city := vars["city"]

	weather, err := getWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/weather/{city}", weatherHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}
