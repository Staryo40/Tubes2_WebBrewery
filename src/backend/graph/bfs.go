package graph

import (
	"backend/models"
	"backend/utils"
	"sort"
	"strings"
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
func HeuristicForwardBFS(target string, elements map[string]models.Element, elementTier map[string]int, num int) []models.Node {
    if elementTier[target] == 0 {
        return []models.Node{{Name: target}}
    }

    availableRecipes := []models.Node{}
    seen := make(map[string]bool) 

    // Initialization: base combinations
    for _, el := range elements {
        if (elementTier[el.Name] == 0){
            seen[el.Name] = true
        }
    }

    for _, el := range elements {
        for _, recipe := range el.Recipes {
            if len(recipe) == 2 &&
                elementTier[recipe[0]] == 0 && elementTier[recipe[1]] == 0 &&
                elementTier[el.Name] <= elementTier[target] {
                
                if !seen[el.Name] {
                    newNode := models.Node{
                        Name:        el.Name,
                        Ingredient1: recipe[0],
                        Ingredient2: recipe[1],
                    }
                    availableRecipes = append(availableRecipes, newNode)
                    seen[el.Name] = true
                }
            }
        }
    }

    addition := true
    for addition {
        addition = false

        if seen[target] {
            result := HeuristicForwardBFSHelper(target, availableRecipes, elements, elementTier)
            if IsFullyExpanded(result, elementTier) {
                return result
            }
        }
    
        current := make([]models.Node, len(availableRecipes))
        copy(current, availableRecipes)

        for _, el := range elements {
            if seen[el.Name] || elementTier[el.Name] > elementTier[target] {
                continue
            }

            for _, recipe := range el.Recipes {
                if len(recipe) != 2 || seen[el.Name] {
                    continue
                }

                if seen[recipe[0]] && seen[recipe[1]] {
                    newNode := models.Node{
                        Name:        el.Name,
                        Ingredient1: recipe[0],
                        Ingredient2: recipe[1],
                    }
                    
                    availableRecipes = append(availableRecipes, newNode)
                    seen[el.Name] = true
                    addition = true
                }
            }
        }
    }

    return nil // No valid recipe found
}

// BFS FROM TARGET WITH HEURISTICS (FAST)
func HeuristicReverseBFS(target string, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {   
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

    count := 0
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
    
        // Expand current path using the helper
        expansions := HeuristicReverseBFSHelper(target, current, elements, elementTier)

        for _, expanded := range expansions {
            if IsFullyExpanded(expanded, elementTier) {
                if count == seed {
                    return expanded
                }
                count++
            } else {
                queue = append(queue, expanded)
            }
        }
    }

    return nil
}

// BIDIRECTIONAL
func HeuristicBidirectionalBFS(target string, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {
	if elementTier[target] == 0 {
		return []models.Node{{Name: target}}
	}

    elementNames := make([]string, 0, len(elements))
    for name := range elements {
        elementNames = append(elementNames, name)
    }
    sort.Strings(elementNames)

    // Initialize Forward Queue
    forwardQueue := [][]models.Node{}
    for _, elName := range elementNames {
        el := elements[elName]

        sort.Slice(el.Recipes, func(i, j int) bool {
            return strings.Join(el.Recipes[i], "+") < strings.Join(el.Recipes[j], "+")
        })

        for _, recipe := range el.Recipes {
            if len(recipe) == 2 {
                ing1, ing2 := recipe[0], recipe[1]

                if ing1 > ing2 {
                    ing1, ing2 = ing2, ing1
                }

                if elementTier[ing1] == 0 && elementTier[ing2] == 0 {
                    newNode := models.Node{
                        Name:        el.Name,
                        Ingredient1: ing1,
                        Ingredient2: ing2,
                    }
                    entry := []models.Node{newNode}
                    forwardQueue = append(forwardQueue, entry)
                }
            }
        }
    }

    // Initialize Reverse Queue
    reverseQueue := [][]models.Node{}
    sort.Slice(elements[target].Recipes, func(i, j int) bool {
        return strings.Join(elements[target].Recipes[i], "+") < strings.Join(elements[target].Recipes[j], "+")
    })
    
    for _, recipe := range elements[target].Recipes {
        if len(recipe) == 2 {
            ing1, ing2 := recipe[0], recipe[1]
    
            if ing1 > ing2 {
                ing1, ing2 = ing2, ing1
            }
    
            if elementTier[ing1] < elementTier[target] && elementTier[ing2] < elementTier[target] {
                newNode := models.Node{
                    Name:        target,
                    Ingredient1: ing1,
                    Ingredient2: ing2,
                }
                entry := []models.Node{newNode}
                reverseQueue = append(reverseQueue, entry)
            }
        }
    }

	results := [][]models.Node{}
    meetingLimit := 10
	for len(forwardQueue) > 0 && len(reverseQueue) > 0 {
        // CHECK FOR MERGE
        result := BidirectionalBFSHelper(forwardQueue, reverseQueue, target, elements, elementTier, seed, results, meetingLimit)
		if result != nil {
			return result
		}

		// FORWARD STEP
        currentForward := forwardQueue[0]
		forwardQueue = forwardQueue[1:] 

		last := currentForward[len(currentForward)-1]
		for _, el := range elements {
			for _, recipe := range el.Recipes {
				if len(recipe) != 2 || elementTier[el.Name] >= elementTier[target] {
					continue
				}
				if (recipe[0] == last.Name || recipe[1] == last.Name) && 
                elementTier[recipe[0]] < elementTier[target] && elementTier[recipe[1]] < elementTier[target] &&
                elementTier[recipe[0]] < elementTier[el.Name] && elementTier[recipe[1]] < elementTier[el.Name] {
					newNode := models.Node{
						Name:        el.Name,
						Ingredient1: recipe[0],
						Ingredient2: recipe[1],
					}
					newPath := append([]models.Node{}, currentForward...)
					newPath = append(newPath, newNode)
					forwardQueue = append(forwardQueue, newPath) 
				}
			}
		}

		// REVERSE STEP
		currentReverse := reverseQueue[0]
		reverseQueue = reverseQueue[1:] // pop from end

		nameSet := make(map[string]bool)
		for _, node := range currentReverse {
			nameSet[node.Name] = true
		}

		var missing string
		for _, node := range currentReverse {
			if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] && elementTier[node.Ingredient1] < elementTier[node.Name] {
				missing = node.Ingredient1
				break
			}
			if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] && elementTier[node.Ingredient2] < elementTier[node.Name]  {
				missing = node.Ingredient2
				break
			}
		}

		if missing != "" {
			for _, recipe := range elements[missing].Recipes {
				if len(recipe) != 2 || elementTier[recipe[0]] >= elementTier[missing] || elementTier[recipe[1]] >= elementTier[missing] {
					continue
				}

				newNode := models.Node{
					Name:        missing,
					Ingredient1: recipe[0],
					Ingredient2: recipe[1],
				}
				newPath := append([]models.Node{newNode}, currentReverse...) 
				reverseQueue = append(reverseQueue, newPath)
			}
		}
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
    producedBy := make(map[string]models.Node)
    for _, node := range path {
        nameSet[node.Name] = true
        producedBy[node.Name] = node
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

            if (elementTier[recipe[0]] > elementTier[missing] || elementTier[recipe[1]] > elementTier[missing]){
                continue
            }

            if existing, exists := producedBy[newNode.Name]; exists {
                // If already produced with same ingredients, skip
                if (existing.Ingredient1 == newNode.Ingredient1 && existing.Ingredient2 == newNode.Ingredient2) ||
                   (existing.Ingredient1 == newNode.Ingredient2 && existing.Ingredient2 == newNode.Ingredient1) {
                    continue
                }
            }

            extended := append([]models.Node{newNode}, path...)
            expandedPaths = append(expandedPaths, extended)
        }
    }

    return expandedPaths
}

func HeuristicForwardBFSHelper(target string, availableRecipe []models.Node, elements map[string]models.Element, elementTier map[string]int) []models.Node {
	recipeMap := make(map[string]models.Node)
	for _, node := range availableRecipe {
		recipeMap[node.Name] = node
	}

	var presentRecipe models.Node
	found := false

	for _, recipe := range elements[target].Recipes {
		if len(recipe) != 2 {
			continue
		}
		if (IngredientInRecipeList(recipe[0], availableRecipe) || elementTier[recipe[0]] == 0) && (IngredientInRecipeList(recipe[1], availableRecipe) || elementTier[recipe[1]] == 0) {
			presentRecipe = models.Node{
				Name:        target,
				Ingredient1: recipe[0],
				Ingredient2: recipe[1],
			}
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	visited := make(map[string]bool)
	result := reconstructPath(presentRecipe, recipeMap, visited)
	return result
}

func reconstructPath(node models.Node, recipeMap map[string]models.Node, visited map[string]bool) []models.Node {
	if visited[node.Name] {
		return nil 
	}
	visited[node.Name] = true

	var result []models.Node

	if subNode, ok := recipeMap[node.Ingredient1]; ok {
		result = append(result, reconstructPath(subNode, recipeMap, visited)...)
	}

	if subNode, ok := recipeMap[node.Ingredient2]; ok {
		result = append(result, reconstructPath(subNode, recipeMap, visited)...)
	}

	result = append(result, node)
	return result
}

func HeuristicReverseBFSHelper(target string, path []models.Node, elements map[string]models.Element, elementTier map[string]int) [][]models.Node {
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

func BidirectionalBFSHelper(forwardQueue [][]models.Node, reverseQueue [][]models.Node, target string, elements map[string]models.Element, elementTier map[string]int, 
    seed int, results [][]models.Node, meetingLimit int) []models.Node {
    for _, fpath := range forwardQueue {
        forwardLast := fpath[len(fpath)-1]

        for _, rpath := range reverseQueue {
            for _, rnode := range rpath {
                if forwardLast.Name == rnode.Name || forwardLast.Name == rnode.Ingredient1 || forwardLast.Name == rnode.Ingredient2 {
                    merged := MergePathsSmart(fpath, rpath, target, elementTier, elements)

                    // Ensure target is last
                    if merged[len(merged)-1].Name != target {
                        for _, recipe := range elements[target].Recipes {
                            if len(recipe) == 2 && elementTier[recipe[0]] < elementTier[target] && elementTier[recipe[1]] < elementTier[target] {

                                nameSet := make(map[string]bool)
                                for _, n := range merged {
                                    nameSet[n.Name] = true
                                }

                                if nameSet[recipe[0]] && nameSet[recipe[1]] {
                                    merged = append(merged, models.Node{
                                        Name:        target,
                                        Ingredient1: recipe[0],
                                        Ingredient2: recipe[1],
                                    })
                                    break
                                }
                            }
                        }
                    }

                    // Expand until fully resolved
                    expanded := merged
                    resolvedSeed := seed
                    for !IsFullyExpanded(expanded, elementTier) {
                        expanded = BFSExpandMissing(target, expanded, elements, elementTier, resolvedSeed)
                        resolvedSeed++
                    }

                    // Deduplication
                    comboSeen := make(map[string]bool)
                    deduped := []models.Node{}
                    for _, node := range expanded {
                        key := ComboKey(node.Name, node.Ingredient1, node.Ingredient2)
                        if !comboSeen[key] {
                            deduped = append(deduped, node)
                            comboSeen[key] = true
                        }
                    }

                    alreadyExists := false
                    for _, r := range results {
                        if utils.PathsEqual(r, deduped) {
                            alreadyExists = true
                            break
                        }
                    }

                    if !alreadyExists {
                        results = append(results, deduped)
                        if len(results) >= meetingLimit {
                            index := seed % len(results)
                            return results[index]
                        }
                    }
                }
            }
        }
    }

    if len(results) > 0 {
        return results[seed%len(results)]
    }

    return nil
}

func BFSExpandMissing(target string, path []models.Node, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {
    nameSet := make(map[string]bool)
    for _, node := range path {
        nameSet[node.Name] = true
    }

    var missing string
    for _, node := range path {
        if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] && elementTier[node.Ingredient1] < elementTier[node.Name] {
            missing = node.Ingredient1
            break
        }
        if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] && elementTier[node.Ingredient2] < elementTier[node.Name] {
            missing = node.Ingredient2
            break
        }
    }

    if missing == "" {
        return path
    }

    expansion := BFSGetMissingPath(missing, elements, elementTier, seed)

    if expansion == nil {
        return path 
    }

    return append(expansion, path...)
}

func BFSGetMissingPath(current string, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {
    // fmt.Println("Called")
	queue := [][]models.Node{}

	// Initialize with all valid recipes for the target
	el, exists := elements[current]
	if !exists || len(el.Recipes) == 0 {
		return nil
	}

    sort.Slice(el.Recipes, func(i, j int) bool {
		return strings.Join(el.Recipes[i], "+") < strings.Join(el.Recipes[j], "+")
	})

    for _, recipe := range el.Recipes {
        if len(recipe) != 2 {
            continue
        }
        if elementTier[recipe[0]] >= elementTier[current] || elementTier[recipe[1]] >= elementTier[current] {
            continue
        }
        node := models.Node{
            Name:        current,
            Ingredient1: recipe[0],
            Ingredient2: recipe[1],
        }
        queue = append(queue, []models.Node{node})
    }

    // RESOLVING CURRENT
    for len(queue) > 0 {
        path := queue[0]
        queue = queue[1:]

        if IsFullyExpanded(path, elementTier) {
            return path
        }

        // Find first missing ingredient in path
        nameSet := make(map[string]bool)
        for _, node := range path {
            nameSet[node.Name] = true
        }

        var missing string
        for _, node := range path {
            if elementTier[node.Ingredient1] != 0 && !nameSet[node.Ingredient1] && elementTier[node.Ingredient1] < elementTier[node.Name] {
                missing = node.Ingredient1
                break
            }
            if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2] && elementTier[node.Ingredient2] < elementTier[node.Name] {
                missing = node.Ingredient2
                break
            }
        }

        if missing != "" {
            missingEl := elements[missing]
            sort.Slice(missingEl.Recipes, func(i, j int) bool {
                return strings.Join(missingEl.Recipes[i], "+") < strings.Join(missingEl.Recipes[j], "+")
            }) 

            for _, recipe := range missingEl.Recipes {
                if len(recipe) != 2 || elementTier[recipe[0]] >= elementTier[missing] || elementTier[recipe[1]] >= elementTier[missing] {
                    continue
                }
                newNode := models.Node{
                    Name:        missing,
                    Ingredient1: recipe[0],
                    Ingredient2: recipe[1],
                }

                skip := false 
                for _, entry := range path {
                    if utils.NodeEqual(entry, newNode) {
                        skip = true
                        break
                    }
                }
                if skip {
                    continue
                }

                newPath := append([]models.Node{newNode}, path...)
                queue = append(queue, newPath)
            }
        }
    }

	return nil // no valid path found
}

func IsFullyExpanded(path []models.Node, elementTier map[string]int) bool {
    if len(path) == 0 || path == nil {
        return false
    }

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

        if tier1 == 0 && tier2 == 0 {
            continue
        }

        for _, ing := range []string{ing1, ing2} {
            if elementTier[ing] == 0 || used[ing] {
                continue
            }

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

func ComboKey(product, ing1, ing2 string) string {
	if ing1 > ing2 {
		ing1, ing2 = ing2, ing1
	}
	return product + "|" + ing1 + "+" + ing2
}

func IngredientInRecipeList(ing string, recipeList []models.Node) bool {
    for _, node := range recipeList{
        if ing == node.Name {
            return true
        }
    }
    return false
}
