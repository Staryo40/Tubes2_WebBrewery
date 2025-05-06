package main

import (
	"recipe-scraper/analyzer"
	// "recipe-scraper/scrape"
)

func main() {
	// scrape.Scrape()
	// analyzer.AnalyzeRecipes("recipes.json")
	// analyzer.MinMaxRecipeCounter(0, "recipes.json");
	analyzer.FindRecipePath("Air", "recipes.json")
}
