package main

import (
	"backend/scraper"
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	Name       string
	Components [][]*Node
}

var nodeCache = make(map[string]*Node)
var nodeMu sync.Mutex
var activeJobMu sync.Mutex
var waitGroupMu sync.Mutex
var componentsMu sync.Mutex
var CU sync.Mutex
var LU sync.Mutex
var wg sync.WaitGroup
var activeJob = 1

func incrementWg(wg *sync.WaitGroup) {
	waitGroupMu.Lock()
	wg.Add(1)
	waitGroupMu.Unlock()
}

func decreaseWg(wg *sync.WaitGroup) {
	waitGroupMu.Lock()
	wg.Done()
	waitGroupMu.Unlock()
}

func incrementActiveJob() {
	activeJobMu.Lock()
	activeJob++
	activeJobMu.Unlock()
}

func decrementActiveJob() {
	activeJobMu.Lock()
	activeJob--
	activeJobMu.Unlock()
}

func newNode(id int, name string) *Node {
	nodeMu.Lock()
	defer nodeMu.Unlock()

	if n, ok := nodeCache[name]; ok {
		fmt.Println("go ID:", id, "Says :", name, " already process")
		return n
	}
	fmt.Println("go ID:", id, "IS PROCESSING :", name)
	n := &Node{Name: name}
	nodeCache[name] = n
	return n
}

// jobsQueue := make(chan string, 100)

// read only no need mutex
func getElementComponent(parentName string, mapTier map[string]int,
	resToElemMap map[string][]string) []string {

	component := resToElemMap[parentName]
	parentTier := mapTier[parentName]
	var validCombos []string

	for _, elem := range component {
		parts := strings.Split(elem, " + ")
		if mapTier[parts[0]] < parentTier && mapTier[parts[1]] < parentTier {
			validCombos = append(validCombos, elem)
		}
	}

	return validCombos //["a+b","b+c"]<-returnnya

}

func worker(id int, jobsQue chan string, nodeCache map[string]*Node,
	mapTier map[string]int, resToElemMap map[string][]string, wg *sync.WaitGroup) {

	for elementName := range jobsQue {

		if _, ok := nodeCache[elementName]; !ok {
			//get parent component
			parentNode := newNode(id, elementName)
			component := getElementComponent(elementName, mapTier, resToElemMap)

			for _, element := range component {
				fmt.Println("go ID:", id, "INI COMPONENT PARENT :", element)

				parts := strings.Split(element, " + ")

				CU.Lock()
				if _, ok := nodeCache[parts[0]]; !ok {
					if parts[0] != "Earth" && parts[0] != "Air" && parts[0] != "Water" && parts[0] != "Fire" {
						jobsQue <- parts[0]
						fmt.Println("go ID:", id, "yang dimasukin :", parts[0])
						incrementActiveJob()
						fmt.Println("go ID:", id, "actuve job :", activeJob)
						incrementWg(wg)
					}

				}

				if _, ok := nodeCache[parts[1]]; !ok {
					if parts[1] != "Earth" && parts[1] != "Air" && parts[1] != "Water" && parts[1] != "Fire" {
						jobsQue <- parts[1]
						fmt.Println("go ID:", id, "yang dimasukin :", parts[1])
						incrementActiveJob()
						fmt.Println("go ID:", id, "actuve job :", activeJob)
						incrementWg(wg)
					}

				}
				CU.Unlock()

				// make his component
				child_1 := newNode(id, parts[0])
				child_2 := newNode(id, parts[1])

				// add his component
				componentsMu.Lock()
				parentNode.Components = append(parentNode.Components, []*Node{child_1, child_2})
				componentsMu.Unlock()
			}
		} else {

		}

		decrementActiveJob()
		decreaseWg(wg)
		fmt.Println("INI ACTIVE JOB", activeJob)

	}
}

func main() {
	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
	}

	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
	}

	mapTier := scraper.ElemTier
	resToElemMap := scraper.ReverseMapping
	done := make(chan struct{})
	jobsQue := make(chan string, 100)
	var wg sync.WaitGroup

	// Tambahkan worker terlebih dahulu
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go worker(i, jobsQue, nodeCache, mapTier, resToElemMap, &wg)
	}

	wg.Add(1)
	jobsQue <- "Sandwich"

	go func() {
		wg.Wait()
		close(jobsQue)
		close(done) // Kirim sinyal selesai
	}()

	// Tunggu sampai semua selesai
	<-done
	fmt.Println("Semua pekerjaan selesai")
}
