package main

import (
    // "fmt"
    // "time"
    "log"
    "net/http"
    "backend/api" 
    // "backend/graph"
    "backend/models"
    "backend/utils"
)

func main() {
    elements := make(map[string]models.Element)
	tiers := make(map[string]int)
    jsonPath := "../frontend/public/elements.json"

	utils.LoadElements(jsonPath, elements)
	utils.LoadTierMap(jsonPath, tiers)	

	api.InitData(elements, tiers)

	// Register route
	http.HandleFunc("/api/recipe", api.RecipeHandler)

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
