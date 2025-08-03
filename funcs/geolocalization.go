package funcs

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleLocation(w http.ResponseWriter, r *http.Request) {
	place := r.URL.Query().Get("place")
	if place == "" {
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "400 Bad Request", ErrorCode: 400})
		return
	}

	coordinates, err := GetCoordinates(place)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "500 Internal Server Error", ErrorCode: 500})
		return
	} else if coordinates == (LocationCodrds{}) { // handle empty result
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", ErrorData{ErrorName: "404 Not Found", ErrorCode: 404})
		return
	}

	// RenderTemplate(w, "test.html", coordinates)

	// Good way to view the JSON without using RenderTemplete(e, "index.html", data)
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(coordinates)

	// http.Redirect(w, r, url, http.StatusSeeOther) // 303 standard for redirects after a POST or internal logic

	// Redirect to Google Maps
	googleMapsURL := fmt.Sprintf("https://www.google.com/maps?q=%s,%s", coordinates.Lat, coordinates.Lon)
	http.Redirect(w, r, googleMapsURL, http.StatusSeeOther)
}

func GetCoordinates(place string) (LocationCodrds, error) {
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json", place)

	resp, err := http.Get(url)
	if err != nil {
		return LocationCodrds{}, err
	}
	defer resp.Body.Close()

	var results []LocationCodrds
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return LocationCodrds{}, err
	}

	if len(results) == 0 {
		return LocationCodrds{}, nil // not found
	}

	return results[0], nil
}
