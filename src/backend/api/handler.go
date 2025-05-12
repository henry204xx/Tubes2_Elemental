// backend/api/handler.go

package api

import (
	"backend/bfs"
	"backend/dfs"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type QueryRequest struct {
	Target       string `json:"target"`
	MaxSolutions int    `json:"maxSolutions"`
	Method       string `json:"method"` // "bfs" or "dfs"
}

type QueryResponse struct {
	Trees        []*dfs.TreeNode `json:"trees"`
	NumSolutions int             `json:"numSolutions"`
	NodeCount    int             `json:"nodeCount"`
	ElapsedTime  string          `json:"elapsedTime"`
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req QueryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("masuk eror decoder:", err.Error())
		fmt.Println(r.Body)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	start := time.Now()

	var trees []*dfs.TreeNode
	var numSolutions, visitedNodes int

	switch req.Method {
	case "dfs":
		trees, numSolutions, visitedNodes = dfs.DFS(req.Target, req.MaxSolutions)
	case "bfs":
		trees, numSolutions, visitedNodes = bfs.BFS(req.Target, req.MaxSolutions)
	default:
		http.Error(w, "Unknown method", http.StatusBadRequest)
		return
	}
	fmt.Println("MASUK NI ABIS FUNGSINYA")
	res := QueryResponse{
		Trees:        trees,
		NumSolutions: numSolutions, // Number of solutions found
		NodeCount:    visitedNodes, // Number of visited nodes
		ElapsedTime:  time.Since(start).String(),
	}
	fmt.Println("Berhasil res:=")
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Berhasil w.header")
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		fmt.Println("GAGAL ENCODE JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Header written? %v\n", w.Header())
	data, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Println("GAGAL MARSHAL:", err)
	} else {
		fmt.Println("HASIL JSON:\n", string(data))
	}
}
