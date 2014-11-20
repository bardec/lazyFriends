package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

type githubWebhook struct {
	Action string `json:"action"`
}

func main() {
	// http.HandleFunc("/", hello)
	// http.HandleFunc("/weather/", func(w http.ResponseWriter, r *http.Request) {
	// 	city := strings.SplitN(r.URL.Path, "/", 3)[2]
	//
	// 	data, err := query(city)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	//
	// 	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// 	json.NewEncoder(w).Encode(data)
	// })
	http.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var github githubWebhook

		if err := json.NewDecoder(r.Body).Decode(&github); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Printf("%s", github.Action)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(github)
	})
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func query(city string) (weatherData, error) {
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city)
	if err != nil {
		return weatherData{}, err
	}

	defer resp.Body.Close()

	var d weatherData

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, err
	}
	fahrenheit := (d.Main.Kelvin-273.15)*1.8 + 32
	fmt.Printf("%f in city %s\n", fahrenheit, city)

	// kelvin := &d.Main.Kelvin
	// *kelvin = fahrenheit
	*(&d.Main.Kelvin) = fahrenheit
	return d, nil
}
