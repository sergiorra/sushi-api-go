package sample

import (
	sushi "github.com/sergiorra/sushi-api-go/pkg"
)

var Sushis = map[string]sushi.Sushi{
	"01D3XZ38KDR": sushi.Sushi{
		ID:    "01D3XZ38KDR",
		ImageNumber:  "1",
		Name:   "California Roll",
		Ingredients: []string {"Crab", "Avocado", "Cucumber", "Sesame seeds"},
	},
	"01D3XZ38TRE": sushi.Sushi{
		ID:    "01D3XZ38TRE",
		ImageNumber:  "2",
		Name:  "Tiger Roll",
		Ingredients: []string {"Avocado", "Cucumber", "Tobiko", "Shrimp tempura"},
	},
	"01D3XZ38KLE": sushi.Sushi{
		ID:    "01D3XZ38KLE",
		ImageNumber:  "3",
		Name:   "Crunch Roll",
		Ingredients: []string {"Spicy tuna", "Crispy seaweed", "Tempura"},
	},
}
