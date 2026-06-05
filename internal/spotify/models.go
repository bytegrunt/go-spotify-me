package spotify

// Artist represents an artist's details
type Artist struct {
	Name       string
	Genres     string
	Popularity int
}

// Song represents a song's details
type Song struct {
	Name       string
	Artist     string
	Album      string
	Popularity int
}

// Me represents the user information from the /me endpoint
type Me struct {
	Country     string
	DisplayName string
	Email       string
	Product     string
	ProfileURL  string
}
