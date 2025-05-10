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
	// API
    elements := make(map[string]models.Element)
	tiers := make(map[string]int)
    jsonPath := "../frontend/public/elements.json"

	utils.LoadElements(jsonPath, elements)
	utils.LoadTierMap(jsonPath, tiers)	

	api.InitData(elements, tiers)

	// Register route
	http.HandleFunc("/api/recipe", api.WithCORS(api.RecipeHandler))

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// target := "Aquarium"
	// start := time.Now()
	// elapsed := time.Since(start)

	// result := graph.HeuristicForwardBFS(target, elements, tiers)
	// for i, node := range result {
	// 	fmt.Printf("%d. %s + %s → %s\n", i+1, node.Ingredient1, node.Ingredient2, node.Name)
	// }
	// fmt.Printf("⏱ Execution time: %s\n", elapsed)
}
