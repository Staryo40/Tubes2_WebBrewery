package analyzer;

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type RecipeEntry struct {
	Product  string     `json:"product"`
	Category string     `json:"category"`
	Recipes  [][]string `json:"recipes"`
}

func AnalyzeRecipes(recipeFile string) {
	// Load JSON file
	data, err := os.ReadFile(recipeFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var recipes []RecipeEntry
	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Analyze
	totalProducts := len(recipes)
	productNames := make([]string, 0, totalProducts)
	categorySet := make(map[string]struct{})

	for _, entry := range recipes {
		productNames = append(productNames, entry.Product)
		categorySet[entry.Category] = struct{}{}
	}

	// Convert category set to slice
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}

	// Output
	fmt.Printf("Total products: %d\n", totalProducts)
	fmt.Printf("Unique categories: %d\n", len(categories))
	fmt.Println("All products:")
	for _, name := range productNames {
		fmt.Println("-", name)
	}
}
