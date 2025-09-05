package funcs

import (
	"net/http"
	"strconv"
	"strings"
)

func ArtistSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		RenderTemplate(w, "error.html", MethodNotAllowed)
		return
	}

	// Fetch all artists and locations data
	ArtistsAPIData, err := FetchAPI(w, r)
	if err != nil {
		return
	}

	data := PageData{
		Artist:   ArtistsAPIData.Artist,
		Location: ArtistsAPIData.Location,
	}

	// Get search query from URL
	query := r.URL.Query().Get("searchQuary")
	// Call search function
	data.Artist = SearchArtists(query, data)

	// When nothing match user quary search
	if len(data.Artist) == 0 {
		w.WriteHeader(http.StatusNotFound) // Look at later
	}

	// Apply pagentation to search result
	data, _ = Pagentation(r, data.Artist, 16)

	// Aplay autocomplete logic
	data.Suggestions = AutoComplete(ArtistsAPIData.Artist, ArtistsAPIData.Location.Index)

	RenderTemplate(w, "index.html", data)
}

func HandleMainSearchQuerie(w http.ResponseWriter, r *http.Request) {
	// Check for searchQuery URL parameter
	query, ok := r.URL.Query()["searchQuary"]
	if ok && len(query) > 0 && strings.TrimSpace(query[0]) != "" {
		ArtistSearchHandler(w, r) // Valid query
		return
	} else if ok && len(query) > 0 && strings.TrimSpace(query[0]) == "" {
		w.WriteHeader(http.StatusBadRequest) // searchQuary is present but empty
		Handler(w, r)
		return
	}

	Handler(w, r) // No searchQuary param at all: (load normal home page)
}

func HandlerFilltersQuerie(w http.ResponseWriter, r *http.Request) {
	APIData, err := FetchAPI(w, r)
	if err != nil {
		return
	}

	ArtistsAPIData := APIData.Artist

	erro := r.ParseForm()
	if erro != nil {
		w.WriteHeader(http.StatusBadRequest)
		RenderTemplate(w, "Error parsing form", BadRequest)
		return
	}

	// Creation date filter
	var creationDateStart int
	var creationDateEnd int
	// First album filter
	var firstAlbumStart string
	var firstAlbumEnd string
	// Checkbox number of members
	var checkbox []int
	// Location filter
	var locationFilter string

	for key := range r.URL.Query() {
		if key == "creationDateStart" && r.Form.Get(key) != "" {
			creationDateStart, err = strconv.Atoi(r.Form.Get("creationDateStart"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderTemplate(w, "error.html", BadRequest)
				return
			}
		} else if key == "creationDateEnd" && r.Form.Get(key) != "" {
			creationDateEnd, err = strconv.Atoi(r.Form.Get("creationDateEnd"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				RenderTemplate(w, "error.html", BadRequest)
				return
			}
		} else if key == "firstAlbumStart" && r.Form.Get(key) != "" {
			firstAlbumStart = r.Form.Get(key)
			if !isValidDate(firstAlbumStart) {
				w.WriteHeader(http.StatusBadRequest)
				RenderTemplate(w, "error.html", BadRequest)
				return
			}
		} else if key == "firstAlbumEnd" && r.Form.Get(key) != "" {
			firstAlbumEnd = r.Form.Get(key)
			if !isValidDate(firstAlbumEnd) {
				w.WriteHeader(http.StatusBadRequest)
				RenderTemplate(w, "error.html", BadRequest)
				return
			}
		} else if key == "location" && r.Form.Get(key) != "" {
			locationFilter = r.Form.Get(key)
		}

		number := filterCheckboxURL(key)
		if number != -1 {
			checkbox = append(checkbox, number)
		}
	}

	// Apply filters
	if creationDateStart != 0 || creationDateEnd != 0 {
		ArtistsAPIData = filterByCreationDate(ArtistsAPIData, creationDateStart, creationDateEnd)
	}

	if firstAlbumStart != "" || firstAlbumEnd != "" {
		ArtistsAPIData = filterByFirstAlbum(ArtistsAPIData, firstAlbumStart, firstAlbumEnd)
	}

	if len(checkbox) > 0 {
		ArtistsAPIData = filterByCheckbox(ArtistsAPIData, checkbox)
	}

	if locationFilter != "" {
		ArtistsAPIData = filterByLocation(ArtistsAPIData, APIData.Location, locationFilter)
	}

	if len(ArtistsAPIData) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	// Render the ArtistsAPIData results
	homeData := PageData{
		Artist:   ArtistsAPIData,
		Location: APIData.Location,
	}

	// Apply pagentation to search result
	homeData, _ = Pagentation(r, homeData.Artist, 16)

	// Applay Utocomplete Logic
	homeData.Suggestions = AutoComplete(APIData.Artist, homeData.Location.Index)

	RenderTemplate(w, "index.html", homeData)
}

func Pagentation(r *http.Request, items []Artists, itemsPerPage int) (PageData, bool) {
	pageStr := r.URL.Query().Get("page")
	page := 1

	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p < 1 {
			return PageData{}, false // invalid query
		}
		page = p
	}

	totalItems := len(items)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	if page > totalPages {
		return PageData{}, false // page not found
	}

	start := (page - 1) * itemsPerPage // Calculate the first item (index) in the page
	end := start + itemsPerPage
	if end > totalItems { // When cant fill the hole page
		end = totalItems
	}

	data := PageData{}

	data.Artist = items[start:end]
	data.Page = page
	data.TotalPages = totalPages
	data.HasNext = page < totalPages
	data.HasPrev = page > 1
	data.NextPage = page + 1
	data.PrevPage = page - 1

	// preserve search query if present
	if s := r.URL.Query().Get("searchQuary"); s != "" {
		data.SearchQuery = s
	}

	// build Prev/Next URLs manually
	baseQuery := r.URL.Query()
	var qParts []string
	for key, values := range baseQuery {
		if key == "page" {
			continue // skip, we will set it manually
		}
		for _, v := range values {
			qParts = append(qParts, key+"="+v)
		}
	}

	if data.HasPrev {
		qPrev := append([]string{}, qParts...)
		qPrev = append(qPrev, "page="+strconv.Itoa(data.PrevPage))
		data.PrevPageURL = "/?" + strings.Join(qPrev, "&")
	}
	if data.HasNext {
		qNext := append([]string{}, qParts...)
		qNext = append(qNext, "page="+strconv.Itoa(data.NextPage))
		data.NextPageURL = "/?" + strings.Join(qNext, "&")
	}

	return data, true
}

func AutoComplete(filtered []Artists, locationIndex []LocationIndex) []Suggestion {
	var suggestions []Suggestion

	// Add artist-related suggestions
	for _, art := range filtered {
		suggestions = append(suggestions, Suggestion{Value: art.Name, Label: art.Name + " - Brand Name"})
		suggestions = append(suggestions, Suggestion{Value: strconv.Itoa(art.CreationDate), Label: strconv.Itoa(art.CreationDate) + " - Creation Date"})
		suggestions = append(suggestions, Suggestion{Value: art.FirstAlbum, Label: art.FirstAlbum + " - First Album"})

		for _, m := range art.Members {
			suggestions = append(suggestions, Suggestion{Value: m, Label: m + " - Member"})
		}
	}

	// Add location suggestions
	for _, indexItem := range locationIndex {
		for _, location := range indexItem.Locations {
			suggestions = append(suggestions, Suggestion{Value: location, Label: location + " - Location"})
		}
	}

	return suggestions
}
