package utils

import (
	"backend/models"
	"fmt"
	"strings"
	"sync"
)

func ConcurrentPathFinder(pathNumber, maxIteration int, finderFunc func(seed int) []models.Node, elementTier map[string]int) []models.CountedPath {
	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		result    [][]models.Node
		resultSet = make(map[string]bool) 
		resChan   = make(chan []models.Node, maxIteration)
	)

	// Spawn concurrent goroutines
	for i := 0; i < maxIteration; i++ {
		wg.Add(1)
		go func(seed int) {
			defer wg.Done()
			path := finderFunc(seed)
			if path != nil {
				resChan <- path
			}
		}(i)
	}

	// Close the result channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resChan)
	}()

	// Collect unique results until we hit the desired count
	for path := range resChan {
		fingerprint := PathFingerprint(path) 
		mu.Lock()
		if !resultSet[fingerprint] && len(result) < pathNumber {
			resultSet[fingerprint] = true
			result = append(result, path)
		}
		mu.Unlock()

		if len(result) >= pathNumber {
			break
		}
	}

	// Convert to CountedPath
	final := make([]models.CountedPath, len(result))
	for i, path := range result {
		final[i] = models.CountedPath{
			NodeCount: NodeCounter(path, elementTier),
			Path:      path,
		}
	}
	return final
}

func PathFingerprint(path []models.Node) string {
	var builder strings.Builder
	for _, node := range path {
		ing1, ing2 := node.Ingredient1, node.Ingredient2
		if ing1 > ing2 {
			ing1, ing2 = ing2, ing1
		}
		builder.WriteString(fmt.Sprintf("%s(%s+%s)->", node.Name, ing1, ing2))
	}
	return builder.String()
}