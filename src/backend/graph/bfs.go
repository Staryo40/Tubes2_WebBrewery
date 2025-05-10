package graph

import (
    "backend/models"
    "backend/utils"
    // "fmt"
)

// BFS FROM BASIC ELEMENTS (VERY SLOW)
func ForwardBFS(target string, elements map[string]models.Element, elementTier map[string]int) []models.Node {
    var queue [][]models.Node
    
    if elementTier[target] == 0 {
        // Return something
        node := models.Node{
            Name:        target,
            Ingredient1: "",
            Ingredient2: "",
        }
        res := []models.Node{node}
        return res
    }

    // Initialize queue with elements that are made directly with base elements
    for _, el := range elements{
        for _, recipe := range el.Recipes{
            if (len(recipe) == 2){
                if (elementTier[recipe[0]] == 0 || elementTier[recipe[1]] == 0){
                    node := models.Node{
                        Name:        el.Name,
                        Ingredient1: recipe[0],
                        Ingredient2: recipe[1],
                    }
                    nodeList := []models.Node{node}
                    queue = append(queue, nodeList)
                }
            }
        }
    }

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        // fmt.Printf("Dequeued path (last: %s):\n", current[len(current)-1].Name)
        // for _, n := range current {
        //     fmt.Printf("  %s ← %s + %s\n", n.Name, n.Ingredient1, n.Ingredient2)
        // }

        // special expand if the target has already been found
        last := current[len(current)-1]
        if last.Name == target {
            expanded := BFSHandleTarget(target, current, elements, elementTier)
            if len(expanded) == 1 && utils.PathsEqual(expanded[0], current) {
                return current
            } else {
                queue = append(queue, expanded...)
            }
        }

        // Expand path
        for _, el := range elements {
            for _, recipe := range el.Recipes {
                if elementTier[last.Name] >= elementTier[el.Name] {
                    continue
                }

                if len(recipe) == 2 {
                    if recipe[0] == last.Name || recipe[1] == last.Name {
                        node := models.Node{
                            Name:        el.Name,
                            Ingredient1: recipe[0],
                            Ingredient2: recipe[1],
                        }
                        newPath := append([]models.Node{}, current...)
                        newPath = append(newPath, node)
                        queue = append(queue, newPath)
                    }
                }
            }
        }
    }

    // No valid recipe found
    return nil
}

// BFS FROM TARGET WITHOUT HEURISICS (QUITE SLOW)
func ReverseBFS(target string, elements map[string]models.Element, elementTier map[string]int) []models.Node {
    var queue [][]models.Node
    
    if elementTier[target] == 0 {
        node := models.Node{
            Name:        target,
            Ingredient1: "",
            Ingredient2: "",
        }
        res := []models.Node{node}
        return res
    }

    // Initialize queue direct recipes of the target
    for _, recipe := range elements[target].Recipes{
        if len(recipe) == 2{
            node := models.Node{
                Name:        target,
                Ingredient1: recipe[0],
                Ingredient2: recipe[1],
            }
            entry := []models.Node{node}
            queue = append(queue, entry)
        }
    }

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        expanded := BFSHandleTarget(target, current, elements, elementTier)
        if len(expanded) == 1 && utils.PathsEqual(expanded[0], current) {
            return current
        } else {
            queue = append(queue, expanded...)
        }
    }

    // did not find anything
    return nil
}

// BFS FROM BASIC ELEMENTS WITH HEURISTICS (#)
func HeuristicForwardBFS(target string, elements map[string]models.Element, elementTier map[string]int) []models.Node {
    var queue [][]models.Node
    
    if elementTier[target] == 0 {
        // Return something
        node := models.Node{
            Name:        target,
            Ingredient1: "",
            Ingredient2: "",
        }
        res := []models.Node{node}
        return res
    }

    // Initialize queue with elements that are made directly with base elements
    for _, el := range elements{
        for _, recipe := range el.Recipes{
            if (len(recipe) == 2){
                if (elementTier[recipe[0]] == 0 || elementTier[recipe[1]] == 0 && elementTier[el.Name] < elementTier[target]){
                    node := models.Node{
                        Name:        el.Name,
                        Ingredient1: recipe[0],
                        Ingredient2: recipe[1],
                    }
                    nodeList := []models.Node{node}
                    queue = append(queue, nodeList)
                }
            }
        }
    }

    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        // fmt.Printf("Dequeued path (last: %s):\n", current[len(current)-1].Name)
        // for _, n := range current {
        //     fmt.Printf("  %s ← %s + %s\n", n.Name, n.Ingredient1, n.Ingredient2)
        // }

        // special expand if the target has already been found
        last := current[len(current)-1]
        if last.Name == target {
            expanded := BFSHandleTarget(target, current, elements, elementTier)
            if len(expanded) == 1 && utils.PathsEqual(expanded[0], current) {
                return current
            } else {
                queue = append(queue, expanded...)
            }
        }

        // Expand path
        for _, el := range elements {
            for _, recipe := range el.Recipes {
                if elementTier[last.Name] >= elementTier[el.Name] {
                    continue
                }

                if len(recipe) == 2 {
                    if recipe[0] == last.Name || recipe[1] == last.Name {
                        node := models.Node{
                            Name:        el.Name,
                            Ingredient1: recipe[0],
                            Ingredient2: recipe[1],
                        }
                        newPath := append([]models.Node{}, current...)
                        newPath = append(newPath, node)
                        queue = append(queue, newPath)
                    }
                }
            }
        }
    }

    // No valid recipe found
    return nil
}

// BFS FROM TARGET WITH HEURISTICS (FAST)
func HeuristicReverseBFS(target string, elements map[string]models.Element, elementTier map[string]int) []models.Node {   
    queue := [][]models.Node{}
    if elementTier[target] == 0 {
        // Return something
        node := models.Node{
            Name:        target,
            Ingredient1: "",
            Ingredient2: "",
        }
        res := []models.Node{node}
        return res
    }

    visited := make(map[string]bool)
    visited[target] = true

    entries := FindExpansionNodes(target, elements, elementTier)
    for _, entry := range entries{
        newEntry := []models.Node{entry}
        queue = append(queue, newEntry)
    }
  

        for len(queue) > 0 {
            current := queue[0]
            queue = queue[1:]
        
            // Expand current path using the helper
            expansions := HeuristicBFSHelper(target, current, elements, elementTier)

            if len(expansions) == 1 && IsFullyExpanded(expansions[0], elementTier) {
                return expansions[0]
            }
        
            queue = append(queue, expansions...)
        }

    return nil
}

// HELPER FUNCTIONS
func BFSHandleTarget(target string, path []models.Node, elements map[string]models.Element, elementTier map[string]int) [][]models.Node {
    // fmt.Println("🔍 Attempting to expand full path:")
    // for i := len(path) - 1; i >= 0; i-- {
    //     node := path[i]
    //     fmt.Printf("  %d. %s ← %s + %s\n", len(path)-i, node.Name, node.Ingredient1, node.Ingredient2)
    // }
    // fmt.Println()

    nameSet := make(map[string]bool)
    for _, node := range path {
        nameSet[node.Name] = true
    }

    // Find recipe with missing path
    var missing string
    for _, node := range path {
        if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] {
            missing = node.Ingredient1
            break
        }
        if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] {
            missing = node.Ingredient2
            break
        }
    }

    if missing == "" {
        return [][]models.Node{path}
    }

    // Find recipe for missing
    var expandedPaths [][]models.Node
    for _, recipe := range elements[missing].Recipes {
        if len(recipe) == 2 {
            newNode := models.Node{
                Name:        missing,
                Ingredient1: recipe[0],
                Ingredient2: recipe[1],
            }

            if nameSet[newNode.Name] {
                continue 
            }

            if (elementTier[recipe[0]] > elementTier[missing] || elementTier[recipe[1]] > elementTier[missing]){
                continue
            }

            // fmt.Printf("🔁 Expanding missing: %s → %s + %s\n", newNode.Name, newNode.Ingredient1, newNode.Ingredient2)
            extended := append([]models.Node{newNode}, path...)
            expandedPaths = append(expandedPaths, extended)
        }
    }

    return expandedPaths
}

func HeuristicBFSHelper(target string, path []models.Node, elements map[string]models.Element, elementTier map[string]int) [][]models.Node {
    // for i := len(path) - 1; i >= 0; i-- {
    //     node := path[i]
    //     fmt.Printf("  %d. %s ← %s + %s\n", len(path)-i, node.Name, node.Ingredient1, node.Ingredient2)
    // }
    // fmt.Println()

    nameSet := make(map[string]bool)
    for _, node := range path {
        nameSet[node.Name] = true
    }

    // Find recipe with missing path
    var missing string
    for _, node := range path {
        if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] {
            missing = node.Ingredient1
            break
        }
        if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] {
            missing = node.Ingredient2
            break
        }
    }

    if missing == "" {
        return [][]models.Node{path}
    }

    // Find recipe for missing
    var expandedPaths [][]models.Node
    for _, newNode := range FindExpansionNodes(missing, elements, elementTier) {
        if nameSet[newNode.Name] {
            continue
        }

        extended := append([]models.Node{newNode}, path...)
        expandedPaths = append(expandedPaths, extended)
    }

    return expandedPaths
}

func IsFullyExpanded(path []models.Node, elementTier map[string]int) bool {
    nameSet := make(map[string]bool)
    for _, node := range path {
        nameSet[node.Name] = true
    }

    for _, node := range path {
        if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] {
            return false
        }
        if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] {
            return false
        }
    }
    return true
}

func FindExpansionNodes(target string, elements map[string]models.Element, elementTier map[string]int) []models.Node {
    used := make(map[string]bool)
    var result []models.Node

    recipes := elements[target].Recipes
    targetTier := elementTier[target]

    // Look for tier 0 + tier 0, if found just return
    for _, r := range recipes {
        if len(r) != 2 {
            continue
        }
        a, b := r[0], r[1]
        if elementTier[a] == 0 && elementTier[b] == 0 {
            return []models.Node{{
                Name:        target,
                Ingredient1: a,
                Ingredient2: b,
            }}
        }
    }

    // Otherwise, apply best-recipe-per-ingredient logic
    for _, r := range recipes {
        if len(r) != 2 {
            continue
        }

        ing1, ing2 := r[0], r[1]
        tier1, tier2 := elementTier[ing1], elementTier[ing2]

        // Skip both-tier-0 combos (already handled above)
        if tier1 == 0 && tier2 == 0 {
            continue
        }

        for _, ing := range []string{ing1, ing2} {
            if elementTier[ing] == 0 || used[ing] {
                continue
            }

            // fmt.Println(ing1 + " + " + ing2 + " -> " + target)

            // Heuristic: ingredient must be strictly lower than product
            if elementTier[ing] >= targetTier {
                continue
            }

            bestRecipe := []string{}
            bestSum := 100

            for _, alt := range recipes {
                if len(alt) != 2 {
                    continue
                }
                a, b := alt[0], alt[1]
                ta, tb := elementTier[a], elementTier[b]

                if a == ing && ta >= tb && ta < targetTier {
                    sum := ta + tb
                    if sum < bestSum {
                        bestRecipe = []string{a, b}
                        bestSum = sum
                    }
                } else if b == ing && tb >= ta && tb < targetTier {
                    sum := ta + tb
                    if sum < bestSum {
                        bestRecipe = []string{a, b}
                        bestSum = sum
                    }
                }
            }

            if len(bestRecipe) == 2 {
                result = append(result, models.Node{
                    Name:        target,
                    Ingredient1: bestRecipe[0],
                    Ingredient2: bestRecipe[1],
                })
                used[ing] = true
            }
        }
    }

    return result
}