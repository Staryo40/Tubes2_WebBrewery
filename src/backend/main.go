package main

import (
	// "backend/scrape"
    // "fmt"
    // "time"
	// "backend/graph"
    "log"
    "net/http"
    "backend/api" 
    "backend/models"
    "backend/utils"
)

func main() { 
	// ------------------------ SCRAPING ------------------------
	// elementsJsonPath := "elements.json"
	// scrape.ScrapeElements(elementsJsonPath)

	// ------------------------ API ------------------------
    elements := make(map[string]models.Element)
	tiers := make(map[string]int)
    jsonPath := "../frontend/public/elements.json"

	utils.LoadElements(jsonPath, elements)
	utils.LoadTierMap(jsonPath, tiers)	

	api.InitData(elements, tiers)

	http.HandleFunc("/api/recipe", api.WithCORS(api.RecipeHandler))

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// -------------------- SINGULAR TESTS --------------------
	// target := "Chameleon"
	// start := time.Now()

	// // result := graph.HeuristicReverseBFS(target, elements, tiers, 0)
	// result := graph.BidirectionalBFS(target, elements, tiers, 1)
	// // result := graph.ReverseDFS(target, elements, tiers, 0, true)
	// // result := graph.BidirectionalDFS(target, elements, tiers, 1)
	// elapsed := time.Since(start)
	// if result == nil {
	// 	fmt.Println("Kok kosong")
	// }

	// for i, node := range result {
	// 	fmt.Printf("%d. %s (%d) + %s (%d) → %s (%d)\n", i+1, node.Ingredient1, tiers[node.Ingredient1], node.Ingredient2, tiers[node.Ingredient2], node.Name, tiers[node.Name])
	// }
	// fmt.Printf("⏱ Execution time: %s\n", elapsed)
	// fmt.Printf("Total nodes: %d\n", utils.NodeCounter(result, tiers))

	// -------------------- PRODUCE RESULT IMAGE --------------------
	// dotPath := "test/grilled.dot"
	// pngPath := "test/grilled.png"
	// err := utils.WriteGraphvizImage(result, dotPath, pngPath)
	// if err != nil {
	// 	fmt.Printf("Warning: could not write graphviz image: %v\n", err)
	// }

	// ------------------------ TIME TEST ------------------------
	// const runs = 1000
	// var total time.Duration
	// for i := 0; i < runs; i++ {
	// start := time.Now()
	// graph.ReverseDFS(target, elements, tiers, 0, true)
	// total += time.Since(start)
	// }
	// fmt.Printf("⏱ Avg: %s\n", total/time.Duration(runs))
}
