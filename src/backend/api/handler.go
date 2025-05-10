package api

import (
	"strings"
	"time"
	"fmt"
    "encoding/json"
    "net/http"
    "backend/graph"
    "backend/models"
    // "backend/utils"
)

var (
    Elements     map[string]models.Element
    ElementTiers map[string]int
)

func InitData(elements map[string]models.Element, tiers map[string]int) {
    Elements = elements
    ElementTiers = tiers
}

func RecipeHandler(w http.ResponseWriter, r *http.Request) {
	// GENERAL REQUEST CHECKS
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    var req models.RequestPayload
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Target == "" {
        http.Error(w, "Invalid JSON or missing target", http.StatusBadRequest)
        return
    }
	// REQUEST SPECIFIC INPUT CHECKS
	if _, exists := Elements[req.Target]; !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "Target element not found",
			Request: req,
		})
		return
	}
	method := strings.ToLower(req.Method)
	if method != "bfs" && method != "dfs" {
		http.Error(w, "Method must be 'BFS' or 'DFS'", http.StatusBadRequest)
		return
	}
	req.Method = strings.ToUpper(method) 
	if req.PathNumber <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "PathNumber must be greater than 0",
			Request: req,
		})
		return
	}
	if req.PathNumber > 100 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Error:   "PathNumber is too large (max 100)",
			Request: req,
		})
		return
	}
	// GET OUTPUT
	start := time.Now()
    var result [][]models.Node
    switch strings.ToUpper(req.Method) {
    case "BFS":
		if req.Bidirectional{
			// Fill later
		} else {
			result = graph.FindPathsBFS(req.Target, Elements, ElementTiers, req.PathNumber, req.Bidirectional)
		}
        

    case "DFS":
		if req.Bidirectional{
			// Fill later
		} else {
			result = graph.FindPathsDFS(req.Target, Elements, ElementTiers, req.PathNumber, req.Bidirectional)
		}

    default:
        http.Error(w, "Invalid method: must be 'BFS' or 'DFS'", http.StatusBadRequest)
        return
    }
	elapsed := time.Since(start)

    if len(result) == 0 {
        http.Error(w, "No paths found for target", http.StatusNotFound)
        return
    }

    response := models.ResponsePayload{
		Count:     len(result),
		ElapsedTime: elapsed.Microseconds(),
		Paths:     result,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Println("tes7")
	fmt.Printf("Response Count: %d\n", response.Count)
	fmt.Printf("Response ElapsedTime (microseconds): %d\n", response.ElapsedTime)
}

func WithCORS(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}