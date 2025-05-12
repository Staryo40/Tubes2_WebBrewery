package utils

import (
	"backend/models"
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ConcurrentPathFinder(pathNumber, maxIteration int, finderFunc func(seed int) []models.Node, elementTier map[string]int) []models.CountedPath {
    type seedPath struct {
        seed int
        path []models.Node
    }

    var (
        wg      sync.WaitGroup
        resChan = make(chan seedPath, maxIteration)
    )

    // 1) Launch all finderFunc(seed) in parallel
    for seed := 0; seed < maxIteration; seed++ {
        wg.Add(1)
        go func(seed int) {
            defer wg.Done()
            if p := finderFunc(seed); p != nil {
                resChan <- seedPath{seed, p}
            }
        }(seed)
    }

    // 2) Close channel once all goroutines finish
    go func() {
        wg.Wait()
        close(resChan)
    }()

    // 3) Collect all results
    all := make([]seedPath, 0, maxIteration)
    for sp := range resChan {
        all = append(all, sp)
    }

    // 4) Sort by seed ascending
    sort.Slice(all, func(i, j int) bool {
        return all[i].seed < all[j].seed
    })

    // 5) Dedupe & pick first pathNumber unique paths
    resultSet := make(map[string]bool)
    final := make([]models.CountedPath, 0, pathNumber)

    for _, sp := range all {
        fp := PathFingerprint(sp.path)
        if !resultSet[fp] {
            resultSet[fp] = true
            final = append(final, models.CountedPath{
                NodeCount: NodeCounter(sp.path, elementTier),
                Path:      sp.path,
            })
            if len(final) >= pathNumber {
                break
            }
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