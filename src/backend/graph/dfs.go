package graph

import (
	"backend/models"
	"backend/utils"
	"sort"
	"strings"
	// "fmt"
)

func ReverseDFS(target string, elements map[string]models.Element, elementTier map[string]int, seed int, strict bool) []models.Node {
	var stack []models.Node

	if elementTier[target] == 0{
		node := models.Node{
			Name: target,
			Ingredient1: "",
			Ingredient2: "",
		}
		stack = append(stack, node)
		return stack
	}

	var ok bool
	if strict {
		ok = ReverseDFSHelper(target, seed, elements, elementTier, &stack)
	} else {
		ok = LooseReverseDFSHelper(target, seed, elements, elementTier, &stack)
	}

	if ok {
		return stack
	}

	return nil
}

func BidirectionalDFS(target string, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {
	if elementTier[target] == 0 {
		return []models.Node{{Name: target}}
	}

	elementNames := make([]string, 0, len(elements))
    for name := range elements {
        elementNames = append(elementNames, name)
    }
    sort.Strings(elementNames)
	
	// Initiate forward stack
	forwardStack := [][]models.Node{}
    for _, elName := range elementNames {
        el := elements[elName]

        recipes := append([][]string(nil), el.Recipes...)
        sort.Slice(recipes, func(i, j int) bool {
            return strings.Join(recipes[i], "+") < strings.Join(recipes[j], "+")
        })

        for _, recipe := range recipes {
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
                    forwardStack = append(forwardStack, entry)
                }
            }
        }
    }

	// Initiate reverse stack
	reverseStack := [][]models.Node{}
	trecipes := append([][]string(nil), elements[target].Recipes...)
    sort.Slice(trecipes, func(i, j int) bool {
        return strings.Join(trecipes[i], "+") < strings.Join(trecipes[j], "+")
    })
    
    for _, recipe := range trecipes {
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
                reverseStack = append(reverseStack, entry)
            }
        }
    }

	// SEARCH
	results := [][]models.Node{}
	meetingLimit := 10
	for len(forwardStack) > 0 && len(reverseStack) > 0{
		result := BidirectionalDFSHelper(forwardStack, reverseStack, target, elements, elementTier, seed, results, meetingLimit)
		if result != nil {
			return result
		}

		// === Forward DFS Expansion ===
		currentForward := forwardStack[len(forwardStack)-1]
		forwardStack = forwardStack[:len(forwardStack)-1] 

		last := currentForward[len(currentForward)-1]
		for _, elName := range elementNames {
			el := elements[elName]
		
			recipes := append([][]string(nil), el.Recipes...)
            sort.Slice(recipes, func(i, j int) bool {
                a1, a2 := recipes[i][0], recipes[i][1]
                if a1 > a2 { a1, a2 = a2, a1 }
                b1, b2 := recipes[j][0], recipes[j][1]
                if b1 > b2 { b1, b2 = b2, b1 }
                return a1+"+"+a2 < b1+"+"+b2
            })

			for _, recipe := range recipes {
				if len(recipe) != 2 || elementTier[el.Name] >= elementTier[target] {
					continue
				}
				if (recipe[0] == last.Name || recipe[1] == last.Name) && elementTier[recipe[0]] < elementTier[el.Name] && elementTier[recipe[1]] < elementTier[el.Name] {
					newNode := models.Node{
						Name:        el.Name,
						Ingredient1: recipe[0],
						Ingredient2: recipe[1],
					}
					newPath := append([]models.Node{}, currentForward...)
					newPath = append(newPath, newNode)
					forwardStack = append(forwardStack, newPath)
				}
			}
		}

		// === Reverse DFS Expansion ===
		currentReverse := reverseStack[len(reverseStack)-1]
		reverseStack = reverseStack[:len(reverseStack)-1] // pop from end

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
			if elementTier[node.Ingredient2] != 0 && !nameSet[node.Ingredient2]  && elementTier[node.Ingredient2] < elementTier[node.Name] {
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
				reverseStack = append(reverseStack, newPath)
			}
		}
	}

	return nil
}

func LooseReverseDFSHelper(current string, seed int, elements map[string]models.Element, elementTier map[string]int, stack *[]models.Node) bool {
	if elementTier[current] == 0 {
		return true
	}

	el, exists := elements[current]
	if !exists || len(el.Recipes) == 0 {
		return false
	}

	validRecipe := [][]string{}
	for _, recipe := range el.Recipes {
		if len(recipe) == 2 {
			validRecipe = append(validRecipe, recipe)
		}
	}

	if len(validRecipe) == 0 {
		return false
	}

	recipeIndex := seed % len(validRecipe)
	chosen := validRecipe[recipeIndex]

	node := models.Node{
		Name:        current,
		Ingredient1: chosen[0],
		Ingredient2: chosen[1],
	}

	// Still use strict logic for sub-ingredients
	ok1 := true
	if elementTier[chosen[0]] != 0 {
		ok1 = ReverseDFSHelper(chosen[0], seed+1, elements, elementTier, stack)
	}

	ok2 := true
	if elementTier[chosen[1]] != 0 {
		ok2 = ReverseDFSHelper(chosen[1], seed+2, elements, elementTier, stack)
	}

	if ok1 && ok2 {
		*stack = append(*stack, node)
		return true
	}

	return false
}

func ReverseDFSHelper(current string, seed int, elements map[string]models.Element, elementTier map[string]int, stack *[]models.Node) bool {
	if elementTier[current] == 0 {
		return true
	}

	el, exists := elements[current]
	if !exists || len(el.Recipes) == 0 {
		return false
	}

	validRecipe := [][]string{}
	for _, recipe := range el.Recipes {
		if len(recipe) == 2 && elementTier[recipe[0]] < elementTier[current] && elementTier[recipe[1]] < elementTier[current] {
			validRecipe = append(validRecipe, recipe)
		}
	}

	if len(validRecipe) == 0 {
		return false
	}

	recipeIndex := seed % len(validRecipe)
	chosen := validRecipe[recipeIndex]

	node := models.Node{
		Name:        current,
		Ingredient1: chosen[0],
		Ingredient2: chosen[1],
	}

	// Recurse only if ingredients are not base
	ok1 := true
	if elementTier[chosen[0]] != 0 && !IsInStack(chosen[0], stack) {
		ok1 = ReverseDFSHelper(chosen[0], seed+1, elements, elementTier, stack)
	}

	ok2 := true
	if elementTier[chosen[1]] != 0 && !IsInStack(chosen[1], stack) {
		ok2 = ReverseDFSHelper(chosen[1], seed+1, elements, elementTier, stack)
	}

	if ok1 && ok2 {
		*stack = append(*stack, node)
		return true
	}

	return false
}

func BidirectionalDFSHelper(forwardStack [][]models.Node, reverseStack [][]models.Node, target string, elements map[string]models.Element, elementTier map[string]int,
	seed int, results [][]models.Node, meetingLimit int) []models.Node {
	for _, fpath := range forwardStack {
		forwardLast := fpath[len(fpath)-1]

		for _, rpath := range reverseStack {
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
						expanded = DFSExpandMissing(target, expanded, elements, elementTier, resolvedSeed)
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

	// No path selected yet
	if len(results) > 0 {
        if len(results[seed % len(results)]) == 0 {
            return nil
        }
        return results[seed%len(results)]
    }

	return nil
}

func MergePathsSmart(fpath []models.Node, rpath []models.Node, target string, elementTier map[string]int, elements map[string]models.Element) []models.Node {
	merged := []models.Node{}
	seen := make(map[string]bool)

	if len(fpath) > 0 {
		merged = append(merged, fpath[0])
		seen[fpath[0].Name] = true
	}

	if len(rpath) > 0 {
		last := rpath[len(rpath)-1]
		if !seen[last.Name] {
			merged = append(merged, last)
			seen[last.Name] = true
		}

		if last.Ingredient1 == fpath[0].Name || last.Ingredient2 == fpath[0].Name {
			return merged
		}
	}

	i, j := 1, len(rpath)-2
	for i < len(fpath) || j >= 0 {
		if i < len(fpath) {
			node := fpath[i]
			if seen[node.Name] {
				return merged
			} else {
				if inIngredient(node, merged){
					insertIndex := max((len(rpath)-j)-1, 1)
					merged = utils.InsertAt(merged, insertIndex, node)
					return merged
				} else {
					insertIndex := max((len(rpath)-j)-1, 1)
					merged = utils.InsertAt(merged, insertIndex, node)
					seen[node.Name] = true
				}
			}
			i++
		}
		if j >= 0 {
			node := rpath[j]
			if seen[node.Name] {
				return merged
			} else {
				if ingredientAvailable(node, merged){
					insertIndex := min(i, len(merged))
					merged = utils.InsertAt(merged, insertIndex, node)
					return merged
				} else {
					insertIndex := min(i, len(merged))
					merged = utils.InsertAt(merged, insertIndex, node)
					seen[node.Name] = true
				}
			}
			j--
		}
	}

	return merged
}

func inIngredient(node models.Node, path []models.Node) bool {
	for _, p := range path {
		if p.Ingredient1 == node.Name || p.Ingredient2 == node.Name {
			return true
		}
	}
	return false
}

func ingredientAvailable(node models.Node, path []models.Node) bool {
	ing1 := node.Ingredient1
	ing2 := node.Ingredient2
	for _, el := range path{
		if el.Name == ing1 || el.Name == ing2 {
			return true
		}
	}
	return false
}

func DFSExpandMissing(target string, path []models.Node, elements map[string]models.Element, elementTier map[string]int, seed int) []models.Node {
	// for _, entry := range path{ 
	// 	fmt.Printf("%s", entry)
	// }
	// fmt.Println()

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

	// fmt.Println(missing)

	if missing == "" {
		return path
	}

	var expansion []models.Node
	ok := ReverseDFSHelper(missing, seed, elements, elementTier, &expansion)
	if !ok {
		return path 
	}

	return append(expansion, path...)
}
// HELPER
func IsInStack(name string, stack *[]models.Node) bool {
	for _, existing := range *stack {
		if existing.Name == name {
			return true
		}
	}
	return false
}

