package main

import (
	// "recipe-scraper/analyzer"
	"recipe-scraper/scrape"
)

func main() {
	// scrape.ScrapeRecipe()
	// analyzer.AnalyzeRecipes("recipes.json")
	// analyzer.MinMaxRecipeCounter(0, "recipes.json");
	// analyzer.FindRecipePath("Air", "recipes.json")
	scrape.ScrapeElementsAndImages();
}
