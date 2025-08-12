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
  Index []LocationIndex `json:"index"`
}

type LocationIndex struct {
    ID        int      `json:"id"`
    Locations []string `json:"locations"`
    Dates     string   `json:"dates"`
}

type Suggestion struct {
	Value string
	Label string
}
type PageData struct {
	Artist      []Artists
	Location    Locations
	Suggestions []Suggestion
	Page        int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
	NextPage    int
	PrevPage    int
	SearchQuery string
}

// Allowed query parameters
var validQueries = map[string]bool{
	"searchQuary": true,
	"page":        true,
	"place":       true,
}

type ErrorData struct {
	ErrorName string
	ErrorCode int
}

type LocationCodrds struct {
	Licence      string `json:"licence"` // Copy Right :)
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	Name         string `json:"name"`
	Display_Name string `json:"display_name"`
}
