package main

import (
	"backend/bfs"
	"backend/dfs"
	"fmt"
	"log"
	"time"
	// "backend/scraper"
)

func main() {
	// scraper.Run() // Uncomment if needed

	fmt.Println("=== Recipe Query Tool ===")
	fmt.Print("Choose query type [1: Single, 2: Multiple]: ")
	var choice int
	if _, err := fmt.Scan(&choice); err != nil {
		log.Fatal("Invalid input for choice:", err)
	}

	fmt.Print("Enter the target recipe: ")
	var target string
	if _, err := fmt.Scan(&target); err != nil {
		log.Fatal("Invalid input for recipe:", err)
	}

	var maxSolution int
	if choice == 2 {
		fmt.Print("Enter max number of solutions: ")
		if _, err := fmt.Scan(&maxSolution); err != nil {
			log.Fatal("Invalid input for max solutions:", err)
		}
	} else {
		maxSolution = 1
	}

	fmt.Print("Choose search method [1: BFS, 2: DFS]: ")
	var searchChoice int
	if _, err := fmt.Scan(&searchChoice); err != nil {
		log.Fatal("Invalid input for search method:", err)
	}

	start := time.Now()

	switch searchChoice {
	case 1:
		fmt.Println("Running BFS...")
		bfs.BFS(target, maxSolution)
	case 2:
		fmt.Println("Running DFS...")
		dfs.DFS(target, maxSolution)
	default:
		fmt.Println("Invalid search choice. Must be 1 (BFS) or 2 (DFS).")
		return
	}

	fmt.Printf("Execution time: %s\n", time.Since(start))
}
