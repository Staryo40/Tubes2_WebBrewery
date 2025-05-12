package graph

import (
	"backend/models"
	"backend/utils"
)

func FindPathsBFS(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) []models.CountedPath {
	if bi {
		i := 0
		maxIteration := 150
		result := [][]models.Node{}
		for len(result) < pathNumber && i < maxIteration {
			resEntry := BidirectionalBFS(target, elements, elementTier, i)
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

		if len(result) == 0 {
			return nil
		}

		res := []models.CountedPath{}
		for _, path := range result{
			newEntry := models.CountedPath{
				NodeCount: utils.NodeCounter(path, elementTier),
				Path: path,
			}
			res = append(res, newEntry)
		}

		return res
	} else {
		i := 0
		maxIteration := 150
		result := [][]models.Node{}
		for len(result) < pathNumber && i < maxIteration {
			resEntry := HeuristicReverseBFS(target, elements, elementTier, i)
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

		if len(result) == 0 {
			return nil
		}

		res := []models.CountedPath{}
		for _, path := range result{
			newEntry := models.CountedPath{
				NodeCount: utils.NodeCounter(path, elementTier),
				Path: path,
			}
			res = append(res, newEntry)
		}

		return res
	}
}

func FindPathsDFS(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) []models.CountedPath {
	if bi {
		i := 0
		maxIteration := 150
		result := [][]models.Node{}
		for len(result) < pathNumber && i < maxIteration {
			resEntry := BidirectionalDFS(target, elements, elementTier, i)
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

		if len(result) == 0 {
			return nil
		}

		res := []models.CountedPath{}
		for _, path := range result{
			newEntry := models.CountedPath{
				NodeCount: utils.NodeCounter(path, elementTier),
				Path: path,
			}
			res = append(res, newEntry)
		}

		return res
	} else {
		i := 0
		maxIteration := 150
		result := [][]models.Node{}
		for len(result) < pathNumber && i < maxIteration {
			resEntry := ReverseDFS(target, elements, elementTier, i, true)
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

		if len(result) == 0 {
			return nil
		}

		res := []models.CountedPath{}
		for _, path := range result{
			newEntry := models.CountedPath{
				NodeCount: utils.NodeCounter(path, elementTier),
				Path: path,
			}
			res = append(res, newEntry)
		}

		return res
	}
}

func FindPathsBFSConcurrent(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) []models.CountedPath {
	if bi {
		return utils.ConcurrentPathFinder(pathNumber, 150,
			func(seed int) []models.Node {
				return BidirectionalBFS(target, elements, elementTier, seed)
			}, elementTier)
	} else {
		return utils.ConcurrentPathFinder(pathNumber, 150,
			func(seed int) []models.Node {
				return HeuristicReverseBFS(target, elements, elementTier, seed)
			}, elementTier)
	}
}

func FindPathsDFSConcurrent(target string, elements map[string]models.Element, elementTier map[string]int, pathNumber int, bi bool) []models.CountedPath {
	if bi {
		return utils.ConcurrentPathFinder(pathNumber, 150,
			func(seed int) []models.Node {
				return BidirectionalDFS(target, elements, elementTier, seed)
			}, elementTier)
	} else {
		return utils.ConcurrentPathFinder(pathNumber, 150,
			func(seed int) []models.Node {
				return ReverseDFS(target, elements, elementTier, seed, true)
			}, elementTier)
	}
}