package main

import (
    "fmt"
    "time"
    // "log"
    // "net/http"
    // "backend/api" 
    "backend/graph"
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

	// api.InitData(elements, tiers)

	// http.HandleFunc("/api/recipe", api.WithCORS(api.RecipeHandler))

	// log.Println("Server is running at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))

	target := "Smoke"
	start := time.Now()
	elapsed := time.Since(start)

	result := graph.HeuristicBidirectionalBFS(target, elements, tiers, 0)
	// result := graph.ReverseDFS(target, elements, tiers, 2, 2, true)
	if result == nil {
		fmt.Println("Kok kosong")
	}

	for i, node := range result {
		fmt.Printf("%d. %s (%d) + %s (%d) → %s (%d)\n", i+1, node.Ingredient1, tiers[node.Ingredient1], node.Ingredient2, tiers[node.Ingredient2], node.Name, tiers[node.Name])
	}
	fmt.Printf("⏱ Execution time: %s\n", elapsed)

	dotPath := "test/grilled.dot"
	pngPath := "test/grilled.png"
	err := utils.WriteGraphvizImage(result, dotPath, pngPath)
	if err != nil {
		fmt.Printf("Warning: could not write graphviz image: %v\n", err)
	}
}
