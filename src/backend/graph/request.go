package graph

import (
	"backend/models"
	"backend/utils"
)

func FindPathsBFS(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) [][]models.Node {
	if bi {

	} else {
		resEntry := HeuristicReverseBFS(target, elements, elementTier)
		if resEntry == nil {
			return nil
		}
		return [][]models.Node{resEntry}
	}

	return nil
}

func FindPathsDFS(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) [][]models.Node {
	if bi {

	} else {
		restriction := true
		if pathNumber > 1 {
			restriction = false
		}
		startRecipe := 0
		startN := 0
		i := 0
		maxIteration := 150
		result := [][]models.Node{}
		for len(result) < pathNumber && i < maxIteration {
			resEntry := ReverseDFS(target, elements, elementTier, startRecipe, startN, restriction)
			if resEntry == nil {
				i++
				continue
			}
		
			duplicate := false
			for _, entry := range result {
				if utils.PathsEqual(entry, resEntry) {
					duplicate = true
					break
				}
			}
		
			if !duplicate {
				result = append(result, resEntry)
			}
		
			i++
		}
		
		if len(result) == 0{
			return nil
		}
		return result
	}

	return nil
}