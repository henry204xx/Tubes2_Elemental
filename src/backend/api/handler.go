// backend/api/handler.go
package api

import (
	"backend/bfs"
	"backend/dfs"
	"encoding/json"
	"net/http"
	"time"
)

type QueryRequest struct {
	Target       string `json:"target"`
	MaxSolutions int    `json:"maxSolutions"`
	Method       string `json:"method"` // "bfs" or "dfs"
}

type QueryResponse struct {
	Trees       []*dfs.TreeNode `json:"trees"`
	NodeCount   int             `json:"nodeCount"`
	ElapsedTime string          `json:"elapsedTime"`
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	start := time.Now()

	var trees []*dfs.TreeNode
	switch req.Method {
	case "dfs":
		trees, _ = dfs.DFS(req.Target, req.MaxSolutions) // returns trees and visited count
	case "bfs":
		trees, _ = bfs.BFS(req.Target, req.MaxSolutions)
	default:
		http.Error(w, "Unknown method", http.StatusBadRequest)
		return
	}

	res := QueryResponse{
		Trees:       trees,
		NodeCount:   dfs.GlobalVisitedCount,
		ElapsedTime: time.Since(start).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
