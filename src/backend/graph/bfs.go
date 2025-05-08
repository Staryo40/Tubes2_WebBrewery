package graph

import (
    "encoding/json"
    "log"
    "os"
    "regexp"
    "strconv"
    "fmt"
    "bufio"
)

type Element struct {
    Name     string     `json:"name"`
    Category string     `json:"category"`
    Recipes  [][]string `json:"recipes"` 
    Image    string     `json:"image"` 
}

var Elements map[string]Element
var ElementTier map[string]int

func BFSSearch(filePath string, target string, elements map[string]Element) []string {
    var queue [][]string

    // Initialize queue with starting elements
    for name, el := range elements {
        if el.Category == "Starting elements" {
            queue = append(queue, []string{name})
        }
    }

    visited := make(map[string]bool)

    for len(queue) > 0 {
        path := queue[0]
        queue = queue[1:]

        current := path[len(path)-1]

        if visited[current] {
            continue
        }
        visited[current] = true

        if current == target {
            return path 
        }

        for _, el := range elements {
            for _, recipe := range el.Recipes {
                for _, ing := range recipe {
                    if ing == current {
                        newPath := append([]string{}, path...)
                        newPath = append(newPath, el.Name)
                        queue = append(queue, newPath)
                        break 
                    }
                }
            }
        }
    }

    return nil 
}

func GetRecipe(filePath string, target string, elements map[string]Element, elementTier map[string]int) []string{
    path := BFSSearch(filePath, target, elements)
    if path == nil || len(path) < 2 {
        return nil
    }

    targetEl, ok := elements[target]
    if !ok {
        return nil
    }

    secondTier := 100
    secondRecipe := ""
    for _, recipe := range targetEl.Recipes {
        if len(recipe) != 2 {
            continue
        }

        if recipe[0] == path[len(path)-2] {
            if elementTier[recipe[1]] < secondTier {
                secondTier = elementTier[recipe[1]]
                secondRecipe = recipe[1]
            }
        } else if recipe[1] == path[len(path)-2] {
            if elementTier[recipe[0]] < secondTier {
                secondTier = elementTier[recipe[0]]
                secondRecipe = recipe[0]
            }
        }
    }

    if secondRecipe == "" {
        return nil
    }

    return []string{secondRecipe, path[len(path)-2]}
}

func FullPathHelper(filePath string, target string, elements map[string]Element, elementTier map[string]int, result *[][]string, visited map[string]bool) {
    if elementTier[target] == 0 || visited[target] {
        return
    }

    visited[target] = true

    recipe := GetRecipe(filePath, target, elements, elementTier)
    if len(recipe) != 2 {
        return
    }

    step := []string{recipe[0], recipe[1], target}
    *result = append(*result, step)

    // fmt.Printf("Added step: %s + %s → %s\n", recipe[0], recipe[1], target)

    if elementTier[recipe[0]] > 0 {
        FullPathHelper(filePath, recipe[0], elements, elementTier, result, visited)
    }
    if elementTier[recipe[1]] > 0 {
        FullPathHelper(filePath, recipe[1], elements, elementTier, result, visited)
    }
}

func GetFullPathToTarget(filePath string, target string, elements map[string]Element, elementTier map[string]int) [][]string {
    var result [][]string
    sited := make(map[string]bool)
    FullPathHelper(filePath, target, elements, elementTier, &result, sited)
    return result
}

func LoadElements(filePath string, elements map[string]Element) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open %s: %v", filePath, err)
    }
    defer file.Close()

    var all []Element
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&all); err != nil {
        log.Fatalf("Failed to parse elements JSON: %v", err)
    }

    for _, el := range all {
        elements[el.Name] = el
    }
}

func LoadTierMap(filePath string, elementTier map[string]int){
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open JSON: %v", err)
    }
    defer file.Close()

    var elements []Element
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&elements); err != nil {
        log.Fatalf("Failed to decode JSON: %v", err)
    }

    maxTier := 0
    tierRegex := regexp.MustCompile(`Tier (\d+) elements`)

    for _, el := range elements {
        var tier int
        switch {
        case el.Category == "Starting elements":
            tier = 0
        case el.Category == "Special element":
            tier = 16 // will adjust later based on maxTier
        case tierRegex.MatchString(el.Category):
            numStr := tierRegex.FindStringSubmatch(el.Category)[1]
            parsed, err := strconv.Atoi(numStr)
            if err != nil {
                log.Printf("Failed to parse tier from %q", el.Category)
                continue
            }
            tier = parsed
            if tier > maxTier {
                maxTier = tier
            }
        default:
            log.Printf("Unknown category: %s", el.Category)
            continue
        }

        elementTier [el.Name] = tier
    }

    // Update special elements if maxTier found is different
    for name, tier := range elementTier  {
        if tier == 16 {
            elementTier [name] = maxTier + 1
        }
    }
}

func WriteGraphvizDOT(filename string, steps [][]string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    writer.WriteString("digraph G {\n")
    writer.WriteString("    rankdir=TB;\n") // top to bottom layout
    writer.WriteString("    node [shape=box, style=filled, fillcolor=lightblue];\n\n")

    // Write all edges
    for _, step := range steps {
        if len(step) != 3 {
            continue
        }
        from1 := step[0]
        from2 := step[1]
        to := step[2]

        // Draw edges with labels
        writer.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", from1, to))
        writer.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", from2, to))
    }

    writer.WriteString("\n}")
    writer.Flush()
    fmt.Println("DOT graph written to", filename)
    return nil
}