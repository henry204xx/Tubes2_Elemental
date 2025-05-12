package bfs

import (
	"backend/dfs"
	"backend/scraper"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	MaxWorkers    = 8 
	QueueCapacity = 1000
)

type Queue struct {
	Elements []string
	Path     []dfs.Step
	Visited  map[string]bool
}

type WorkerPool struct {
	jobs          chan Node
	wg            sync.WaitGroup
	results       chan *dfs.TreeNode
	Count         int32 // found solutions
	maxSolutions  int
	done          chan struct{}
	activeJobs    int32 
	jobsSubmitted int32 
	depthTracker  *DepthTracker // depth levels for BFS ordering
	mu            sync.Mutex   
}

type Node struct {
	queue Queue
	depth int 
}

func NewWorkerPool(maxSolutions int, depthTracker *DepthTracker) *WorkerPool {
	return &WorkerPool{
		jobs:         make(chan Node, QueueCapacity),
		results:      make(chan *dfs.TreeNode, maxSolutions),
		maxSolutions: maxSolutions,
		done:         make(chan struct{}),
		depthTracker: depthTracker,
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
				return
			}

			atomic.AddInt32(&p.activeJobs, 1)

			p.depthTracker.StartProcessingAtDepth(job.depth)
			for job.depth > p.depthTracker.GetCurrentDepth() {
			select {
			case <-p.done:
				p.depthTracker.FinishProcessingAtDepth(job.depth)
				return
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}

			// Process the job
			p.processNode(job)

			atomic.AddInt32(&p.activeJobs, -1)

			// Check if reached max solutions or all work is done
			if atomic.LoadInt32(&p.Count) >= int32(p.maxSolutions) || 
				(atomic.LoadInt32(&p.activeJobs) == 0 && p.isQueueEmpty() && atomic.LoadInt32(&p.jobsSubmitted) > 0) {
				p.Stop()
				return
			}
		default:
			if atomic.LoadInt32(&p.Count) >= int32(p.maxSolutions) || 
				(atomic.LoadInt32(&p.activeJobs) == 0 && p.isQueueEmpty() && atomic.LoadInt32(&p.jobsSubmitted) > 0) {
				p.Stop()
				return
			}
			time.Sleep(1 * time.Millisecond) // Consider reducing or removing
		}
	}
}

func (p *WorkerPool) processNode(node Node) {
	defer p.depthTracker.FinishProcessingAtDepth(node.depth)
	select {
	case <-p.done:
		return
	default:
	}

	current := node.queue

	

	
	// Check for path completion
	if len(current.Elements) == 0 {
		if dfs.ContainsElementInPath(current.Path, "Time") {
			afterTimeCount := dfs.UniqueElementsAfterTime(current.Path)
			if afterTimeCount < 100 {
				return
			}
		}
	
		newCount := atomic.AddInt32(&p.Count, 1)
		fmt.Println(newCount, "Found complete path:")
		PrintQueue(current)

		tree := dfs.BuildTree(current.Path, current.Path[0].Result)
		select {
		case p.results <- tree:
		case <-p.done:
			return
		default:
		}

		return
	}

	elem := current.Elements[0]
	current.Elements = current.Elements[1:]

	if current.Visited[elem] {
		p.Submit(Node{queue: current, depth: node.depth})
		return
	}

	atomic.AddInt64(&dfs.GlobalVisitedCount, 1)

	if dfs.IsBasicElement(elem) || elem == "Time" {
		p.Submit(Node{queue: current, depth: node.depth})
		return
	}

	recipes, exists := scraper.ReverseMapping[elem]
	if !exists {
		return
	}

	branchesSubmitted := false
	for _, recipe := range recipes {
		parts := strings.Split(recipe, " + ")
		if len(parts) != 2 {
			continue
		}
		a, b := parts[0], parts[1]

		tierA, tierB, tierElem := scraper.ElemTier[a], scraper.ElemTier[b], scraper.ElemTier[elem]
		if tierA >= tierElem || tierB >= tierElem {
			continue
		}

		newQueue := copyQueue(current)
		newQueue.Elements = append(newQueue.Elements, a, b)
		newQueue.Path = append(newQueue.Path, dfs.Step{Result: elem, Components: []string{a, b}})
		newQueue.Visited[elem] = true

		nextDepth := node.depth + 1
		p.depthTracker.AddNodeAtDepth(nextDepth)
		p.Submit(Node{queue: newQueue, depth: nextDepth})
		branchesSubmitted = true
	}

	if !branchesSubmitted && len(current.Elements) == 0 {
		return
	}
}

func (p *WorkerPool) Submit(node Node) {
	select {
	case <-p.done:
		return
	default:
	}

	atomic.AddInt32(&p.jobsSubmitted, 1)

	select {
	case <-p.done:
		return
	case p.jobs <- node:
	default:
		go func() {
			select {
			case <-p.done:
				return
			case p.jobs <- node:
			}
		}()
	}
}

func (p *WorkerPool) Wait() {
	p.wg.Wait()
	close(p.results)
}

func (p *WorkerPool) Stop() {
	select {
	case <-p.done:
	default:
		close(p.done)
	}
}

func (p *WorkerPool) GetResults() []*dfs.TreeNode {
	var results []*dfs.TreeNode
	for result := range p.results {
		results = append(results, result)
	}
	return results
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
		Elements: newElements,
		Path:     newPath,
		Visited:  newVisited,
	}
}

type DepthTracker struct {
	currentDepth    int
	nodesAtDepth    map[int]int32
	processingDepth map[int]int32
	lock            sync.Mutex
}

// GetCurrentDepth returns the current BFS depth being processed.
func (dt *DepthTracker) GetCurrentDepth() int {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	return dt.currentDepth
}

func NewDepthTracker() *DepthTracker {
	return &DepthTracker{
		currentDepth:    0,
		nodesAtDepth:    make(map[int]int32),
		processingDepth: make(map[int]int32),
	}
}

func (dt *DepthTracker) AddNodeAtDepth(depth int) {
	dt.lock.Lock() 
	defer dt.lock.Unlock()
	dt.nodesAtDepth[depth]++
}

func (dt *DepthTracker) StartProcessingAtDepth(depth int) {
	dt.lock.Lock()
	defer dt.lock.Unlock()
	dt.processingDepth[depth]++
}

func (dt *DepthTracker) FinishProcessingAtDepth(depth int) {
	dt.lock.Lock() 
	defer dt.lock.Unlock()
	dt.processingDepth[depth]--
	if depth == dt.currentDepth && dt.processingDepth[depth] <= 0 {
		for dt.processingDepth[dt.currentDepth] <= 0 && dt.nodesAtDepth[dt.currentDepth+1] > 0 {
			dt.currentDepth++
			fmt.Printf("Go to depth %d\n", dt.currentDepth)
		}
	}
}

func (p *WorkerPool) isQueueEmpty() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.jobs) == 0
}

func BFS(root string, maxSolution int) {

	dfs.GlobalVisitedCount = 0
	
	// Load necessary data
	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
		return
	}
	
	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
		return
	}
	
	
	depthTracker := NewDepthTracker()
	
	// Create initial queue
	queue := Queue{
		Visited: make(map[string]bool),
	}
	depthTracker.AddNodeAtDepth(0)
	queue.Elements = append(queue.Elements, root)
	queue.Path = append(queue.Path, dfs.Step{Result: root, Components: []string{}})
	
	// Create and start the worker pool
	pool := NewWorkerPool(maxSolution, depthTracker)
	pool.Start()
	
	
	depthTracker.AddNodeAtDepth(0)
	
	
	initialNode := Node{
		queue: queue,
		depth: 0,
	}
	
	
	pool.Submit(initialNode)
	
	// Use a separate goroutine to periodically check if all work is done
	go func() {
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				// Check if we've reached max solutions
				if atomic.LoadInt32(&pool.Count) >= int32(maxSolution) {
					// fmt.Println("Reached max solutions:", maxSolution)
					pool.Stop()
					return
				}
				
				// Check if we've explored all possibilities
				if atomic.LoadInt32(&pool.activeJobs) == 0 && pool.isQueueEmpty() && 
				   atomic.LoadInt32(&pool.jobsSubmitted) > 0 {
					pool.Stop()
					return
				}
			case <-pool.done:
				return
			}
		}
	}()
	
	// Wait for all workers to finish
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
	fmt.Printf("Total nodes visited: %d\n", dfs.GlobalVisitedCount)
	
	// Write results to file
	f, err := os.Create("paths.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()
	
	for i, tree := range results {
		fmt.Fprintf(f, "Path %d:\n", i+1)
		fmt.Fprintf(f, "%s\n", root)
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

func SaveResultsToJSON(results []*dfs.TreeNode, filename string) error {
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