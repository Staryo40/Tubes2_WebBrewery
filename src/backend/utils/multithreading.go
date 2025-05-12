package utils

import (
	"backend/models"
	"fmt"
	"runtime"
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
        sem     = make(chan struct{}, runtime.NumCPU())      
        resChan = make(chan seedPath, maxIteration)         
    )

    for seed := 0; seed < maxIteration; seed++ {
        wg.Add(1)
        go func(s int) {
            defer wg.Done()

            sem <- struct{}{}
            defer func() { <-sem }()

            defer func() {
                if r := recover(); r != nil {
                    fmt.Printf("Panic in seed=%d: %v\n", s, r)
                }
            }()

            if p := finderFunc(s); p != nil {
                resChan <- seedPath{seed: s, path: p}
            }
        }(seed)
    }

    go func() {
        wg.Wait()
        close(resChan)
    }()

    all := make([]seedPath, 0, maxIteration)
    for sp := range resChan {
        all = append(all, sp)
    }

    sort.Slice(all, func(i, j int) bool {
        return all[i].seed < all[j].seed
    })

    seen := make(map[string]bool)
    final := make([]models.CountedPath, 0, pathNumber)
    for _, sp := range all {
        fp := PathFingerprint(sp.path)
        if !seen[fp] {
            seen[fp] = true
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