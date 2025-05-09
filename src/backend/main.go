package main

import (
    "fmt"
    "time"
    "log"
    // "net/http"
    // "backend/api" 
    "backend/graph"
    "backend/models"
    "backend/utils"
)

func main() {
    elementTier := make(map[string]int)
    elements := make(map[string]models.Element)
    filePath := "../frontend/public/elements.json"
    utils.LoadTierMap(filePath, elementTier)
    utils.LoadElements(filePath, elements);

    // err := graph.WriteGraphvizDOT("grilled.dot", steps)
    // if err != nil {
    //     log.Fatal(err)
    // }

    // start := time.Now()
    // result := graph.ReverseBFS("Grilled cheese", elements, elementTier)
    // elapsed := time.Since(start)
    
    // for i, node := range result {
    //     fmt.Printf("%d. %s + %s → %s\n", i+1, node.Ingredient1, node.Ingredient2, node.Name)
    // }
    // fmt.Printf("⏱ Execution time: %s\n", elapsed)

    // PER TARGET 
    start := time.Now()
    result := graph.HeuristicBFS("Algae", elements, elementTier)
    elapsed := time.Since(start)
    for i, node := range result {
        fmt.Printf("%d. %s + %s → %s\n", i+1, node.Ingredient1, node.Ingredient2, node.Name)
    }
    fmt.Printf("⏱ Execution time: %s\n", elapsed)

    err := utils.WriteGraphvizImage(result, "algae.dot", "algae.png")
    if err != nil {
        log.Fatal(err)
    }
}
