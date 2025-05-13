package dfs

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"backend/scraper"
	"encoding/json"
)

var GlobalVisitedCount int64

// Thread pool configuration
const (
	MaxWorkers    = 4
	QueueCapacity = 100 
)

type Step struct {
	Result     string
	Components []string
}

type Stack struct {
	Elements  []string
	Path      []Step
	Visited   map[string]bool
}


type WorkerPool struct {
	jobs          chan Node
	wg            sync.WaitGroup
	results       chan *TreeNode
	Count         int32 // found solutions
	maxSolutions  int
	mu            sync.Mutex 
	done          chan struct{}
	activeJobs    int32 
	jobsSubmitted int32 
}

type Node struct {
	stack Stack
}

type TreeNode struct {
	Result     string      `json:"result"`
	Components []*TreeNode `json:"components,omitempty"`
}

func NewWorkerPool(maxSolutions int) *WorkerPool {
	return &WorkerPool{
		jobs:    make(chan Node, QueueCapacity),
		results: make(chan *TreeNode, maxSolutions),
		maxSolutions: maxSolutions,
		done: make(chan struct{}),
	}
}

func (p *WorkerPool) Start() {
	// Start worker goroutines
	for i := 0; i < MaxWorkers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *WorkerPool) worker() {
	defer p.wg.Done()
	
	for {
		select {
		case <-p.done:
			return
		case job, ok := <-p.jobs:
			if !ok {
				// Channel closed, exit worker
				return 
			}
			
			
			atomic.AddInt32(&p.activeJobs, 1)
			
			
			p.processNode(job)
			
			
			atomic.AddInt32(&p.activeJobs, -1)
			
			// Check if reached the max solutions
			if atomic.LoadInt32(&p.Count) >= int32(p.maxSolutions) {
				p.Stop()
				return
			}
			
		
			if atomic.LoadInt32(&p.activeJobs) == 0 && p.isQueueEmpty() && atomic.LoadInt32(&p.Count) < int32(p.maxSolutions) {
				// If no active jobs, no queued jobs, and have not found max solutions,
				// explored all possibilities
				p.mu.Lock()
				fmt.Println("All paths explored. Found", p.Count, "solutions.")
				p.mu.Unlock()
				p.Stop() 
				return
			}
		default:
			
			if atomic.LoadInt32(&p.Count) >= int32(p.maxSolutions) {
				p.Stop()
				return
			}
			
			time.Sleep(5 * time.Millisecond)
		}
	}
}


func (p *WorkerPool) isQueueEmpty() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.jobs) == 0
}

func (p *WorkerPool) processNode(node Node) {
	//  check if we should stop processing
	select {
	case <-p.done:
	
		return
	default:
		
	}

	current := node.stack
	
	
	if len(current.Elements) == 0 {
		// Validate the path
		if ContainsElementInPath(current.Path, "Time") {
			afterTimeCount := UniqueElementsAfterTime(current.Path)
			if afterTimeCount < 100 {
				return // invalid path: not enough elements created after Time
			}
		}
		
		newCount := atomic.AddInt32(&p.Count, 1)
		
		p.mu.Lock()
		fmt.Println(newCount, "Found complete path:")
		printStack(current)
		p.mu.Unlock()
		
		// Create tree and send to results channel
		tree := BuildTree(current.Path, current.Path[0].Result)
		select {
		case p.results <- tree:
			

		case <- p.done:
			
			return
		default:
			
		}
		
		return
	}
	
	elem := current.Elements[len(current.Elements)-1]
	current.Elements = current.Elements[:len(current.Elements)-1]
	
	if current.Visited[elem] {
	
	select {
	case <-p.done:
	
		return
	default:
		
		p.Submit(Node{stack: current})
		return
	}
	}
	
	atomic.AddInt64(&GlobalVisitedCount, 1)

	
	if IsBasicElement(elem) || elem == "Time" {
		p.Submit(Node{stack: current})
		return
	}
	
	recipes, exists := scraper.ReverseMapping[elem]
	if !exists {
		return
	}
	
	branchesSubmitted := false
	
	for _, recipe := range recipes {
		select {
		case <-p.done:
			
			return
		default:
			
		}
		
		parts := strings.Split(recipe, " + ")
		if len(parts) != 2 {
			continue
		}
		a, b := parts[0], parts[1]
		
		tierA, tierB, tierElem := scraper.ElemTier[a], scraper.ElemTier[b], scraper.ElemTier[elem]
		if tierA >= tierElem || tierB >= tierElem {
			continue // Skip processing if a or b has a higher tier than elem
		}
		
		newStack := copyStack(current)
		newStack.Elements = append(newStack.Elements, b, a)
		
		newStack.Path = append(newStack.Path, Step{
			Result:     elem,
			Components: []string{a, b},
		})
		
		newStack.Visited[elem] = true
		
	
		p.Submit(Node{stack: newStack})
		branchesSubmitted = true
	}
	
	// check if done
	if !branchesSubmitted && len(current.Elements) == 0 {
		// This path is exhausted, and no new branches were created
		if atomic.LoadInt32(&p.activeJobs) == 1 && p.isQueueEmpty() {
			// If this is the last active job and queue is empty, done
			p.mu.Lock()
			fmt.Println("All paths explored. Found", p.Count, "solutions.")
			p.mu.Unlock()
			p.Stop()
		}
	}
}

func (p *WorkerPool) Submit(node Node) {
	// Only submit if haven't reached max solutions
	if atomic.LoadInt32(&p.Count) < int32(p.maxSolutions) {
	
		select {
		case <-p.done:
		
			return
		default:
		}
		

		atomic.AddInt32(&p.jobsSubmitted, 1)	
		
		select {
		case p.jobs <- node:
		
		case <-p.done:
			return
		default:
		}
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.results) 
}

func (p *WorkerPool) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	select {
	case <-p.done:
	
	default:
		close(p.done)
		close(p.jobs)
	}
}

func (p *WorkerPool) GetResults() []*TreeNode {
	var results []*TreeNode
	for result := range p.results {
		results = append(results, result)
	}
	return results
}

func PrintTree(node *TreeNode, prefix string, f *os.File) {
	if node == nil {
		return
	}
	
	// Print children
	for i, child := range node.Components {
		isLast := i == len(node.Components)-1
		
		// Determine connectors
		var connector, newPrefix string
		if prefix == "" {
			// First level children
			connector = "├── "
			newPrefix = "|   "
		} else {
			connector = "├── "
			newPrefix = prefix + "|   "
			if isLast {
				connector = "└── "
				newPrefix = prefix + "    "
			}
		}
		
		// Print the connector and the child node's result
		fmt.Fprintf(f, "%s%s%s\n", prefix, connector, child.Result)
		
		// Print vertical line if not last child
		if !isLast {
			fmt.Fprintf(f, "%s|\n", newPrefix)
		}
		
		// Recurse for children
		PrintTree(child, newPrefix, f)
	}
}

func UniqueElementsAfterTime(path []Step) int {
	foundTime := false
	seen := make(map[string]bool)
	
	for _, step := range path {
		if step.Result == "Time" || Contains(step.Components, "Time") {
			foundTime = true
			continue
		}
		if foundTime {
			seen[step.Result] = true
			for _, comp := range step.Components {
				seen[comp] = true
			}
		}
	}
	return len(seen)
}

func Contains(arr []string, target string) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

func ContainsElementInPath(path []Step, elem string) bool {
	for _, step := range path {
		if step.Result == elem || Contains(step.Components, elem) {
			return true
		}
	}
	return false
}

func copyStack(s Stack) Stack {
	newElements := make([]string, len(s.Elements))
	copy(newElements, s.Elements)
	
	newPath := make([]Step, len(s.Path))
	copy(newPath, s.Path)
	
	newVisited := make(map[string]bool)
	for k, v := range s.Visited {
		newVisited[k] = v
	}
	
	return Stack{
		Elements:  newElements,
		Path:      newPath,
		Visited:   newVisited,
	}
}

func IsBasicElement(elem string) bool {
	switch elem {
	case "Water", "Earth", "Fire", "Air":
		return true
	default:
		return false
	}
}

func BuildTree(path []Step, target string) *TreeNode {
	stepMap := make(map[string]Step)
	for _, s := range path {
		stepMap[s.Result] = s
	}
	
	// Track which elements we've already expanded (not just processed)
	expanded := make(map[string]bool)
	
	var build func(string) *TreeNode
	build = func(elem string) *TreeNode {
		// Retrieve the step for the current node
		step, exists := stepMap[elem]
		if !exists || len(step.Components) == 0 {
			return &TreeNode{Result: elem} // Leaf node (no components)
		}
		
		// Create the node first
		node := &TreeNode{Result: elem}
		
		// If already expanded this element, return just the node
		if expanded[elem] {
			return node
		}
		expanded[elem] = true
		
		// Create children nodes
		children := []*TreeNode{}
		for _, c := range step.Components {
			childNode := build(c) // Recursive call to build child nodes
			if childNode != nil {
				children = append(children, childNode)
			}
		}
		
		node.Components = children
		return node
	}
	
	return build(target) // Start from the target node
}

func DFS(root string, maxSolution int) ([]*TreeNode, int, int) {

	GlobalVisitedCount = 0
	
	// Load necessary data
	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
		return nil, 0, 0
	}
	
	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
		return nil, 0, 0
	}
	
	// Create initial stack
	stack := Stack{
		Visited: make(map[string]bool),
	}
	stack.push(root, Step{Result: root, Components: []string{}})
	
	// Create and start the worker pool
	pool := NewWorkerPool(maxSolution)
	pool.Start()
	
	// Submit the initial node
	pool.Submit(Node{stack: stack})
	
	
	// Use a separate goroutine to periodically check if all work is done
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// Check if we've reached max solutions
				if atomic.LoadInt32(&pool.Count) >= int32(maxSolution) {
					pool.Stop()
					return
				}
				
				// Check if there are no active jobs and queue is empty
				if atomic.LoadInt32(&pool.activeJobs) == 0 && pool.isQueueEmpty() && atomic.LoadInt32(&pool.jobsSubmitted) > 0 {
					fmt.Println("All paths explored. Found", pool.Count, "solutions.")
					pool.Stop()
					fmt.Println("Stopping worker pool...")
					return
				}
			case <-pool.done:
				// Pool has been stopped elsewhere
				return
			}
		}
	}()
	

	pool.Wait()
	
	// Get results and save to JSON
	
	results := pool.GetResults()
	err := SaveResultsToJSON(results, "paths.json")
	if err != nil {
		fmt.Println("Error saving results to JSON:", err)
	} else {
		fmt.Println("Results successfully saved to paths.json")
	}
	
	fmt.Printf("Total paths found: %d\n", pool.Count)
	fmt.Printf("Total nodes visited: %d\n", GlobalVisitedCount)

	pool.Stop()

	
	
	// Write results to file
	f, err := os.Create("paths.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, 0, 0
	}
	defer f.Close()
	
	for i, tree := range results {
		fmt.Fprintf(f, "Path %d:\n", i+1)
		fmt.Fprintf(f, "%s\n", root)
		PrintTree(tree, "", f)
		fmt.Fprintln(f)
	}

	if (results == nil) {
		return results, -1, int(GlobalVisitedCount)
	}
	return results, int(pool.Count), int(GlobalVisitedCount)
}


func (s *Stack) push(element string, step Step) {
	s.Elements = append(s.Elements, element)
	s.Path = append(s.Path, step)
}

func printStack(s Stack) {
	fmt.Println("Elements:", s.Elements)
	fmt.Println("Paths:")
	for _, p := range s.Path {
		fmt.Println(p.Result, "->", p.Components)
	}
}

func SaveResultsToJSON(results []*TreeNode, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(results)
	return err
}