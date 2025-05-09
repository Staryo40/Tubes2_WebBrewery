package utils

import (
	"backend/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func LoadElements(filePath string, elements map[string]models.Element) {
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("Failed to open %s: %v", filePath, err)
    }
    defer file.Close()

    var all []models.Element
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

    var elements []models.Element
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
            tier = 0 // will adjust later based on maxTier
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

func WriteGraphvizImage(nodes []models.Node, dotPath, pngPath string) error {
    file, err := os.Create(dotPath)
    if err != nil {
        return fmt.Errorf("failed to create dot file: %v", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)

    // Write DOT format
    writer.WriteString("digraph G {\n")
    writer.WriteString("    rankdir=TB;\n") // top-to-bottom layout
    writer.WriteString("    node [shape=box, style=filled, fillcolor=lightblue];\n\n")

    for _, node := range nodes {
        if node.Ingredient1 == "" || node.Ingredient2 == "" || node.Name == "" {
            continue
        }

        writer.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", node.Ingredient1, node.Name))
        writer.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\";\n", node.Ingredient2, node.Name))
    }

    writer.WriteString("}\n")
    writer.Flush()
	file.Close()

    fmt.Println("✅ DOT file written to:", dotPath)
	
    // Run Graphviz to produce PNG
    cmd := exec.Command("dot", "-Tpng", dotPath, "-o", pngPath)
	cmd.Env = os.Environ() 

    out, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to generate PNG using Graphviz: %v\n%s", err, string(out))
    }	

    fmt.Println("✅ PNG image generated at:", pngPath)
    return nil
}