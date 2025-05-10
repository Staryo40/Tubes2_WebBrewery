package graph

import (
    "backend/models"
    // "backend/utils"
    // "fmt"
)

func ReverseDFS(target string, elements map[string]models.Element, elementTier map[string]int, recipe int, num int, strict bool) []models.Node {
	var stack []models.Node

	var ok bool
	if strict {
		ok = ReverseDFSHelper(target, recipe, num, elements, elementTier, &stack)
	} else {
		ok = LooseReverseDFSHelper(target, recipe, num, elements, elementTier, &stack)
	}
	if ok {
		// Reverse the stack to get the path from base → target
		// for i, j := 0, len(stack)-1; i < j; i, j = i+1, j-1 {
		// 	stack[i], stack[j] = stack[j], stack[i]
		// }
		return stack
	}

	return nil
}

func LooseReverseDFSHelper(current string, r int, n int, elements map[string]models.Element, elementTier map[string]int, stack *[]models.Node) bool {
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

	recipeIndex := r % len(validRecipe)
	chosen := validRecipe[recipeIndex]

	node := models.Node{
		Name:        current,
		Ingredient1: chosen[0],
		Ingredient2: chosen[1],
	}

	// Still use strict logic for sub-ingredients
	ok1 := true
	if elementTier[chosen[0]] != 0 {
		ok1 = ReverseDFSHelper(chosen[0], r+1, n+1, elements, elementTier, stack)
	}

	ok2 := true
	if elementTier[chosen[1]] != 0 {
		ok2 = ReverseDFSHelper(chosen[1], r+2, n+2, elements, elementTier, stack)
	}

	if ok1 && ok2 {
		*stack = append(*stack, node)
		return true
	}

	return false
}

func ReverseDFSHelper(current string,r int,n int,elements map[string]models.Element,elementTier map[string]int, stack *[]models.Node) bool {
	if elementTier[current] == 0 {
		return true
	}

	el, exists := elements[current]
	if !exists || len(el.Recipes) == 0 {
		return false
	}

	validRecipe := [][]string{}
	for _, recipe := range el.Recipes {
		if len(recipe) == 2 &&
			(elementTier[recipe[0]] < elementTier[current] && elementTier[recipe[1]] < elementTier[current]) {
			validRecipe = append(validRecipe, recipe)
		}
	}

	if len(validRecipe) == 0 {
		return false
	}

	recipeIndex := (r + n) % len(validRecipe)
	chosen := validRecipe[recipeIndex]

	node := models.Node{
		Name:        current,
		Ingredient1: chosen[0],
		Ingredient2: chosen[1],
	}

	// Recurse only if ingredients are not base
	ok1 := true
	if elementTier[chosen[0]] != 0 && !IsInStack(chosen[0], stack) {
		ok1 = ReverseDFSHelper(chosen[0], r+1, n+1, elements, elementTier, stack)
	}

	ok2 := true
	if elementTier[chosen[1]] != 0 && !IsInStack(chosen[1], stack) {
		ok2 = ReverseDFSHelper(chosen[1], r+2, n+2, elements, elementTier, stack)
	}

	if ok1 && ok2 {
		*stack = append(*stack, node)
		return true
	}

	return false
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