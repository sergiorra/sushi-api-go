package sushi

import (
	"context"
	"time"
)

// Sushi defines the properties of a sushi to be listed
type Sushi struct {
	ID          string 		`json:"id"`
	ImageNumber string 		`json:"imageNumber,omitempty"`
	Name        string 		`json:"name,omitempty"`
	Ingredients []string 	`json:"ingredients,omitempty"`
	CreatedAt 	*time.Time `json:"-"`
	UpdatedAt 	*time.Time `json:"-"`
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
	CreateSushi(ctx context.Context, s *Sushi) error
	GetSushis(ctx context.Context) ([]Sushi, error)
	DeleteSushi(ctx context.Context, ID string) error
	UpdateSushi(ctx context.Context, ID string, s *Sushi) error
	GetSushiByID(ctx context.Context, ID string) (*Sushi, error)
}
