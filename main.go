package main

import (
	"encoding/json"
)

type Sushi struct {
	ID          string `json:"id"`
	ImageNumber string `json:"imageNumber"`
	Name        string `json:"name"`
	Ingredients string `json:"ingredients"`
}

var sushis []Sushi

func main() {

}
