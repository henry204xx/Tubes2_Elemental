package bfs

import (
	"fmt"
	"os"
	"strings"
	"backend/scraper"
	"backend/dfs"
)


type Queue struct {
	Elements  []string
	Path      []dfs.Step
	Visited   map[string]bool
	Ancestry  []string
}

func copyQueue(q Queue) Queue {
	newElements := make([]string, len(q.Elements))
	copy(newElements, q.Elements)

	newPath := make([]dfs.Step, len(q.Path))
	copy(newPath, q.Path)

	newVisited := make(map[string]bool)
	for k, v := range q.Visited {
		newVisited[k] = v
	}

	return Queue{
		Elements:  newElements,
		Path:      newPath,
		Visited:   newVisited,
		Ancestry:  append([]string{}, q.Ancestry...),
	}
}

func BFS(root string, maxSolution int) {
	var allPaths []*dfs.TreeNode
	count := 0

	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
		return
	}

	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
		return
	}

	queue := Queue{
		Visited: make(map[string]bool),
	}
	queue.Elements = append(queue.Elements, root)
	queue.Path = append(queue.Path, dfs.Step{Result: root, Components: []string{}})

	type Node struct {
		queue Queue
	}

	var queues []Node
	queues = append(queues, Node{queue: queue})

	for len(queues) > 0 {
		currentNode := queues[len(queues)-1]
		queues = queues[:len(queues)-1]
		current := currentNode.queue

		if len(current.Elements) == 0 {
			if dfs.ContainsElementInPath(current.Path, "Time") {
				afterTimeCount := dfs.UniqueElementsAfterTime(current.Path)
				if afterTimeCount < 100 {
					continue // invalid path: not enough elements created after Time
				}
			}

			fmt.Println(count+1, "Found complete path:")
			PrintQueue(current)

			tree := dfs.BuildTree(current.Path, root)
			allPaths = append(allPaths, tree)

			count++
			if count >= maxSolution {
				break
			}
			continue
		}

		elem := current.Elements[0]
		current.Elements = current.Elements[1:]

		if current.Visited[elem] {
			queues = append(queues, Node{queue: current})
			continue
		}

		dfs.GlobalVisitedCount++

		if dfs.IsBasicElement(elem) || elem == "Time" {
			queues = append(queues, Node{queue: current})
			continue
		}

		recipes, exists := scraper.ReverseMapping[elem]
		if !exists {
			continue
		}

		for _, recipe := range recipes {
			parts := strings.Split(recipe, " + ")
			if len(parts) != 2 {
				continue
			}
			a, b := parts[0], parts[1]

			tierA, tierB, tierElem := scraper.ElemTier[a], scraper.ElemTier[b], scraper.ElemTier[elem]
			if tierA >= tierElem || tierB >= tierElem {
				continue // Skip processing if a or b has a higher tier than elem
			}
			newQueue := copyQueue(current)
			newQueue.Ancestry = append(append([]string{}, current.Ancestry...), elem)
			newQueue.Elements = append(newQueue.Elements, a, b)

			newQueue.Path = append(newQueue.Path, dfs.Step{
				Result:     elem,
				Components: []string{a, b},
			})

			newQueue.Visited[elem] = true

			queues = append(queues, Node{queue: newQueue})
		}
	}

	fmt.Printf("Total paths found: %d\n", count)
	fmt.Printf("Total nodes visited: %d\n", dfs.GlobalVisitedCount)

	f, err := os.Create("paths.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	for i, tree := range allPaths {
		fmt.Fprintf(f, "Path %d:\n", i+1)
		fmt.Fprintf(f, "%s", root + "\n")
		dfs.PrintTree(tree, "", f)
		fmt.Fprintln(f)
	}
}

func PrintQueue(q Queue) {
    fmt.Println("Elements:", q.Elements)
    fmt.Println("Paths:")
    for _, p := range q.Path {
        fmt.Println(p.Result, "->", p.Components)
    }
}
