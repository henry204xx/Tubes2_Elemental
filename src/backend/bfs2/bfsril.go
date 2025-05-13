package main

import (
	"backend/scraper"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type JSONNode struct {
	Result     string     `json:"result"`
	Components []JSONNode `json:"components,omitempty"`
}

type NodeElement struct {
	Name       string
	Components [][]string
}

var counter int
var counterMutex sync.Mutex

func incrementCounter() {
	counterMutex.Lock()
	counter++
	counterMutex.Unlock()
}

var nodeCounter int64
var isShuttingDown atomic.Bool
var queueNotEmpty = make(chan struct{}, 1)
var done = make(chan struct{})

var visitedElem sync.Map
var mapAllReceipe sync.Map
var jobsChannel = make(chan string, 100)
var waitingQueue []string
var waitingQueueMu sync.Mutex
var nodeMutex sync.Mutex
var wgSync sync.WaitGroup

var maxPaths int

func getNodeComponents(parentName string, elementTier map[string]int, res2ElemMap map[string][]string) [][]string {
	component := res2ElemMap[parentName]
	parentTier := elementTier[parentName]
	var validCombos [][]string

	for _, elem := range component {
		parts := strings.Split(elem, " + ")
		if elementTier[parts[0]] < parentTier && elementTier[parts[1]] < parentTier {
			validCombos = append(validCombos, parts)
		}
	}

	return validCombos
}

func makeNode(name string, elementTier map[string]int,
	res2ElemMap map[string][]string) *NodeElement {

	temp := getNodeComponents(name, elementTier, res2ElemMap)
	fmt.Println(temp)
	time.Sleep(50000)
	elem := &NodeElement{
		Name:       name,
		Components: temp,
	}
	visitedElem.Store(name, true)
	mapAllReceipe.Store(name, elem)
	return elem
}

func Add(elem string) {
	waitingQueueMu.Lock()
	defer waitingQueueMu.Unlock()

	for _, item := range waitingQueue {
		if item == elem {
			return
		}
	}

	waitingQueue = append(waitingQueue, elem)
	select {
	case queueNotEmpty <- struct{}{}:
	default: // avoid blocking if already signaled
	}
}

func dispatchJobs() {
	defer close(jobsChannel)

	for {
		fmt.Println("ini coo")
		select {
		case <-done:
			return
		case <-queueNotEmpty:
			for {
				waitingQueueMu.Lock()
				if len(waitingQueue) == 0 {
					waitingQueueMu.Unlock()
					break
				}
				elem := waitingQueue[0]
				waitingQueue = waitingQueue[1:]
				waitingQueueMu.Unlock()

				select {
				case jobsChannel <- elem:
				case <-done:
					return
				}
			}
		}
	}
}

func workerElement(workerId int, elementTier map[string]int, res2ElemMap map[string][]string) {
	defer wgSync.Done()
	fmt.Printf("Worker %d: ready to work\n", workerId)

	for elem := range jobsChannel {
		if _, loaded := visitedElem.LoadOrStore(elem, true); !loaded {
			fmt.Println("wokerID: ", workerId, " IS PROCESSING ", elem)
			nodeElement := makeNode(elem, elementTier, res2ElemMap)
			for _, component := range nodeElement.Components {
				Add(component[0])
				Add(component[1])
			}
		}
	}
	fmt.Printf("Worker %d: channel closed, exiting\n", workerId)
}

func RunBFS(start string) {
	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
	}

	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
	}

	elementTier := scraper.ElemTier
	res2ElemMap := scraper.ReverseMapping

	go dispatchJobs()
	for i := 0; i < 5; i++ {
		wgSync.Add(1)
		go workerElement(i, elementTier, res2ElemMap)
	}
	Add(start)

	go func() {
		for {
			waitingQueueMu.Lock()
			empty := len(waitingQueue) == 0
			waitingQueueMu.Unlock()

			if empty {
				// Double check
				time.Sleep(500 * time.Millisecond)
				waitingQueueMu.Lock()
				empty = len(waitingQueue) == 0
				waitingQueueMu.Unlock()

				if empty {
					close(done)
					return
				}
			} else {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	wgSync.Wait()
}

var buildCache sync.Map // map[string][]JSONNode

func buildTreeLimited(name string, maxResults int) []JSONNode {
	// Cek cache
	if val, ok := buildCache.Load(name); ok {
		return val.([]JSONNode)
	}

	nodeRaw, ok := mapAllReceipe.Load(name)
	if !ok {
		return []JSONNode{{Result: name}} // leaf node
	}
	node := nodeRaw.(*NodeElement)

	var results []JSONNode

	for _, combo := range node.Components {
		if len(combo) != 2 {
			continue
		}
		leftList := buildTreeLimited(combo[0], maxResults)
		rightList := buildTreeLimited(combo[1], maxResults)

		for _, left := range leftList {
			for _, right := range rightList {
				results = append(results, JSONNode{
					Result:     name,
					Components: []JSONNode{left, right},
				})
				if len(results) >= maxResults {
					buildCache.Store(name, results)
					return results
				}
			}
		}
	}

	if len(results) == 0 {
		results = append(results, JSONNode{Result: name})
	}

	buildCache.Store(name, results)
	return results
}

func stringfyATree(tree []JSONNode) (string, error) {
	data, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}
	return string(data), nil
}

func main() {
	RunBFS("Sandwich")
	val, ok := mapAllReceipe.Load("Sandwich")
	if ok {
		fmt.Println("Value found:", val)
	} else {
		fmt.Println("Key not found")
	}
	a := buildTreeLimited("Sandwich", 2)
	b, _ := stringfyATree(a)
	fmt.Println(b)
}
