package funcs

import (
	"encoding/json"
	"net/http"
)

// Fetches artists data from the API
func FetchAPI(w http.ResponseWriter, r *http.Request) (FilltersData, error) {
	// Check HTTP method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "error.html", MethodNotAllowed)
		return FilltersData{}, nil
	}

	// Call the API to get the data
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		InternalServerError.Message = "Failed to get API data"
		RenderTemplate(w, "error.html", InternalServerError)
		return FilltersData{}, err
	}
	defer resp.Body.Close()

	// Decode the JSON into struct
	var apiData []Artists
	if err := json.NewDecoder(resp.Body).Decode(&apiData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		InternalServerError.Message = "Failed to parse API response"
		RenderTemplate(w, "error.html", InternalServerError)
		return FilltersData{}, err
	}

	// Get the location data
	location, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil || location.StatusCode != http.StatusOK {
		w.WriteHeader(http.StatusNotFound)
		RenderTemplate(w, "error.html", NotFound)
		return FilltersData{}, err
	}
	defer location.Body.Close()

	var locations Locations
	if err := json.NewDecoder(location.Body).Decode(&locations); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		InternalServerError.Message = "Failed to get artist relation data"
		RenderTemplate(w, "error.html", InternalServerError)
		return FilltersData{}, err
	}

	// Apply initial filtering
	filtered := FilterArtists(apiData, ExcludeIDs)

	resultData := FilltersData{
		Artist:   filtered,
		Location: locations,
	}

	return resultData, nil
}
