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
	NextPageURL string
	PrevPageURL string
}

type FilltersData struct {
	Artist   []Artists
	Location Locations
}

// Allowed query parameters
var validQueries = map[string]bool{
	"searchQuary":       true,
	"page":              true,
	"place":             true,
	"filterQuerys":      true,
	"creationDateStart": true,
	"creationDateEnd":   true,
	"firstAlbumStart":   true,
	"firstAlbumEnd":     true,
	"cb1":               true,
	"cb2":               true,
	"cb3":               true,
	"cb4":               true,
	"cb5":               true,
	"cb6":               true,
	"cb7":               true,
	"cb8":               true,
	"location":          true,
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

var NotFound = struct {
	Code    string
	Title   string
	Message string
}{
	Code:    "404",
	Title:   "Page not found",
	Message: "Sorry, this page was not found",
}

var BadRequest = struct {
	Code    string
	Title   string
	Message string
}{
	Code:    "400",
	Title:   "Bad Request",
	Message: "Sorry, Bad request.",
}

var MethodNotAllowed = struct {
	Code    string
	Title   string
	Message string
}{
	Code:    "405",
	Title:   "Method Not Allowed",
	Message: "Sorry, This method is not allowed.",
}

var InternalServerError = struct {
	Code    string
	Title   string
	Message string
}{
	Code:    "500",
	Title:   "Internal Error",
	Message: "Zzz, Internal Server Error.",
}
