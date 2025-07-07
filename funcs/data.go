package funcs

var ExcludeIDs = []int{11, 12, 21, 22, 49}

const Items_Per_Page = 20

type Artists struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type RelationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistDetailPage struct {
	Artist   Artists
	Relation RelationsData
}

type Locations struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type PageData struct {
	Artist   []Artists
	Location Locations
	// Suggestions []string
	Page       int
	TotalPages int
	HasNext    bool
	HasPrev    bool
	NextPage   int
	PrevPage   int
}

// Allowed query parameters
var validQueries = map[string]bool{
	"searchQuary": true,
	"page":        true,
}
