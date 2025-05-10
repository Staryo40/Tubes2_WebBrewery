package api

import (
	"strings"
	"time"
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
	

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(models.ResponsePayload{Paths: result})
}