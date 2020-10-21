package sushi

// Sushi defines the properties of a sushi to be listed
type Sushi struct {
	ID          string 		`json:"id"`
	ImageNumber string 		`json:"imageNumber,omitempty"`
	Name        string 		`json:"name,omitempty"`
	Ingredients []string 	`json:"ingredients,omitempty"`
}

// New creates a sushi
func New(ID, ImageNumber, Name string, Ingredients []string) *Sushi {
	return &Sushi{
		ID: 			ID,
		ImageNumber: 	ImageNumber,
		Name: 			Name,
		Ingredients: 	Ingredients,
	}
}

// Repository provides access to the sushi storage
type Repository interface {
	CreateSushi(s *Sushi) error
	GetSushis() ([]Sushi, error)
	DeleteSushi(ID string) error
	UpdateSushi(ID string, s *Sushi) error
	GetSushiByID(ID string) (*Sushi, error)
}