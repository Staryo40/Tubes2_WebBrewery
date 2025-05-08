package main

import (
    "fmt"
    "log"
    // "net/http"
    // "backend/api" 
    "backend/graph"
)

func main() {
    // target := "Earthquake"
    // paths := graph.BFSSearch(target)

    // if len(paths) == 0 {
    //     fmt.Printf("No path found to '%s'\n", target)
    //     return
    // }

    // fmt.Printf("Path to '%s':\n", target)
    // for pIndex, path := range paths {
    //     fmt.Printf("Path %d:\n", pIndex+1)
    //     for i, name := range path {
    //         el, exists := graph.Elements[name]
    //         if !exists {
    //             log.Printf("Element not found: %s\n", name)
    //             continue
    //         }
    //         fmt.Printf("  %d. %s (Category: %s, Image: %s)\n", i+1, el.Name, el.Category, el.Image)
    //     }
    //     fmt.Println()
    // }

    elementTier := make(map[string]int)
    elements := make(map[string]graph.Element)
    filePath := "../frontend/public/elements.json"
    graph.LoadTierMap(filePath, elementTier)
    graph.LoadElements(filePath, elements);

    // path := graph.BFSSearch(filePath, "Pond", elements)
    // target := elements["Pond"]
    // recipeElement := []string{}

    // for _, recipe := range target.Recipes{
    //     for i, el := range recipe{
    //         if (i == 0 && el == path[len(path)-2]){
    //             recipeElement = append(recipeElement, recipe[1])
    //         } else if (i == 1 && el == path[len(path)-2]){
    //             recipeElement = append(recipeElement, recipe[0])
    //         }
    //     }
    // }

    // fmt.Println("First path to target element")
    // for i, el := range path{
    //     fmt.Printf("%d. %s\n", i+1, el)
    // }

    // fmt.Println("Other element with first path element")
    // for i, el := range recipeElement{
    //     fmt.Printf("%d. %s\n", i+1, el)
    // }

    steps := graph.GetFullPathToTarget(filePath, "Katana", elements, elementTier)

    if len(steps) == 0 {
        fmt.Println("No crafting path found to 'Airplane'.")
    } else {
        for _, step := range steps {
            fmt.Printf("%s + %s → %s\n", step[0], step[1], step[2])
        }
    }

    err := graph.WriteGraphvizDOT("katana.dot", steps)
    if err != nil {
        log.Fatal(err)
    }
}
