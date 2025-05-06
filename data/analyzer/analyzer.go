package analyzer;

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type RecipeEntry struct {
	Product  string     `json:"product"`
	Category string     `json:"category"`
	Recipes  [][]string `json:"recipes"`
}

type ProductCount struct {
	Product string
	Count   int
}

func AnalyzeRecipes(recipeFile string) {
	// Load JSON file
	data, err := os.ReadFile(recipeFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var recipes []RecipeEntry
	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Analyze
	totalProducts := len(recipes)
	productNames := make([]string, 0, totalProducts)
	categorySet := make(map[string]struct{})

	for _, entry := range recipes {
		productNames = append(productNames, entry.Product)
		categorySet[entry.Category] = struct{}{}
	}

	// Convert category set to slice
	categories := make([]string, 0, len(categorySet))
	for cat := range categorySet {
		categories = append(categories, cat)
	}

	// Output
	fmt.Printf("Total products: %d\n", totalProducts)
	fmt.Printf("Unique categories: %d\n", len(categories))
	fmt.Println("All products:")
	for _, name := range productNames {
		fmt.Println("-", name)
	}
}

func MinMaxRecipeCounter(mode int, recipeFile string) {
	data, err := os.ReadFile(recipeFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var recipes []RecipeEntry
	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Step 1: Build a slice of product-count pairs
	var counts []ProductCount
	for _, recipe := range recipes {
		counts = append(counts, ProductCount{
			Product: recipe.Product,
			Count:   len(recipe.Recipes),
		})
	}

	// Step 2: Sort descending by count
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Count > counts[j].Count
	})

	// Step 3: Print top 5
	limit := 5
	if len(counts) < 5 {
		limit = len(counts)
	}

	fmt.Printf("Top %d products with the most recipes:\n", limit)
	for i := 0; i < limit; i++ {
		fmt.Printf("%d. %s — %d recipe entries\n", i+1, counts[i].Product, counts[i].Count)
	}
}

func FindRecipePath(target string, recipeFile string) {
	data, err := os.ReadFile(recipeFile)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	var recipes []RecipeEntry
	if err := json.Unmarshal(data, &recipes); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	// Step 1: Build product name → RecipeEntry map for quick lookup
	productMap := make(map[string]RecipeEntry)
	for _, r := range recipes {
		productMap[r.Product] = r
	}

	// Step 2: Build reverse dependency graph
	graph := make(map[string][]string)
	for _, r := range recipes {
		for _, ingredient := range r.Recipes {
			if len(ingredient) > 0 {
				inputProduct := ingredient[0]
				graph[inputProduct] = append(graph[inputProduct], r.Product)
			}
		}
	}

	// Step 3: Find all starting elements
	var queue []string
	visited := make(map[string]bool)
	parent := make(map[string]string) // to reconstruct path

	for _, r := range recipes {
		if r.Category == "Starting Elements" {
			queue = append(queue, r.Product)
			visited[r.Product] = true
		}
	}

	// Step 4: BFS
	found := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == target {
			found = true
			break
		}

		for _, next := range graph[current] {
			if !visited[next] {
				visited[next] = true
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}

	// Step 5: Reconstruct path if found
	if found {
		path := []string{}
		for cur := target; cur != ""; cur = parent[cur] {
			path = append([]string{cur}, path...)
		}
		fmt.Println("Path to", target, ":", path)
	} else {
		fmt.Println("No path to product:", target)
	}
}