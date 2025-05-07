package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

var mapping map[string]string

func LoadMapping() error {
	file, err := os.Open("ElemToRes.json")
	if err != nil {
		return fmt.Errorf("error opening ElemToRes.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&mapping); err != nil {
		return fmt.Errorf("error decoding ElemToRes.json: %v", err)
	}
	return nil
}

var reverseMapping map[string][]string

func LoadReverseMapping() error {
	file, err := os.Open("ResToElem.json")
	if err != nil {
		return fmt.Errorf("error opening ResToElem.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&reverseMapping); err != nil {
		return fmt.Errorf("error decoding ResToElem.json: %v", err)
	}
	return nil
}

type StateNode struct {
	elements []string // Elemen yang ada saat ini
	path     []string // resep menuju elemen target
}

func IsInMap(s string) bool {
	for mapElem := range mapping {
		if s == mapElem {
			return true
		}
	}
	return false
}

func IsContainTarget(sPath []string, target []string) bool {
	for _, s := range sPath {
		for _, t := range target {
			if s == t {
				return true
			}
		}
	}
	return false
}

// Cek slice string apakah mengandung elemen val
func ContainElements(sl []string, val string) bool {
	for _, s := range sl {
		if s == val {
			return true
		}
	}
	return false
}

func IsContainAllTarget(pathResult [][]string, target []string) bool {
	count := 0
	for _, pathRes := range pathResult {
		for _, t := range target {
			if pathRes[len(pathRes)-1] == t {
				count += 1
			}
		}
	}
	return count == len(target)
}

func IsContainTargetElements(pathResult [][]string, currPath []string) bool {
	for _, pathRes := range pathResult {
		path := fmt.Sprintf("%v", pathRes)
		current := fmt.Sprintf("%v", currPath)
		if path == current {
			return true
		}
	}
	return false
}

func PathValid(result []string) bool {
	//count := 0
	for i := 0; i < len(result)-1; i++ {
		res1 := strings.Split(result[i], " ")
		for j := i + 1; j < len(result); j++ {
			res2 := strings.Split(result[j], " ")
			if !(res2[0] == res1[len(res2)-1] || res2[2] == res1[len(res2)-1]) {
				return false //count += 1
			}
		}
	}
	return true //(count <= len(result)-1 && count >= 1)
}

func main() {
	// Read JSON data
	if err := LoadMapping(); err != nil {
		fmt.Println("Error loading mapping: ", err)
		return
	}

	if err := LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping: ", err)
		return
	}

	fmt.Println("Mappings loaded successfully")

	elementTarget := "Earth"
	recipeTarget := reverseMapping[elementTarget]
	var resultTarget []string

	for _, recipe := range recipeTarget {
		resultTarget = append(resultTarget, strings.Join([]string{recipe, elementTarget}, " = "))
	}
	fmt.Println(resultTarget)

	method := "shortest"
	maxRecipe := 5
	if maxRecipe > len(resultTarget) {
		maxRecipe = len(resultTarget)
	} else if maxRecipe <= 0 {
		fmt.Println("Banyak resep maksimum tidak valid. Masukan angka >= 1")
	}

	startElements := []string{"Air", "Earth", "Fire", "Water"}

	// BFS: tinjau dari awal elemen
	queue := []StateNode{{elements: startElements}} // queue = [StateNode {elements: ["Air", "Earth", "Fire", "Water"], path: []}]

	var currentState, newCurrentState StateNode
	var combine1, combine2, res, elem1, elem2, resultPath string
	var resultAllPath []string
	visited := make(map[string]bool)

	if recipeTarget == nil || ContainElements(startElements, elementTarget) {
		resultAllPath = recipeTarget[:0]
		fmt.Println(resultAllPath)
		return
	}

	// BFS: Jalur terpendek
	if method == "shortest" {
		for len(queue) > 0 {
			currentState = queue[0] // StateNode{elements: ["Air", "Earth", "Fire", "Water"], path : []}
			queue = queue[1:]       // pop elemen awal
			fmt.Println(len(queue), currentState)

			// Cek, apakah path sudah mengandung target
			if IsContainTarget(currentState.path, resultTarget) && PathValid(currentState.path) {
				resultAllPath = currentState.path
				break
			}

			// Mencari kombinasi seluruh elemen pada currentState.elements
			for i := 0; i < len(currentState.elements); i++ {
				for j := i; j < len(currentState.elements); j++ {
					elem1 = currentState.elements[i]
					elem2 = currentState.elements[j]

					if elem1 == elem2 {
						combine1 = strings.Join([]string{elem1, elem2}, " + ")

						if IsInMap(combine1) {
							// tambahkan elemen ke current state node
							res = mapping[combine1]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine1, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else {
							continue
						}
					} else {
						combine1 = strings.Join([]string{elem1, elem2}, " + ")
						combine2 = strings.Join([]string{elem2, elem1}, " + ")

						if IsInMap(combine1) {
							// tambahkan elemen ke current state node
							res = mapping[combine1]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine1, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else if IsInMap(combine2) {
							res = mapping[combine2]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine2, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else {
							continue
						}
					}
				}
			}
		}
		fmt.Println(resultAllPath)

		// BFS: Semua resep
	} else if method == "many" {
		var resultAll [][]string

		for len(queue) > 0 {
			currentState = queue[0] // StateNode{elements: ["Air", "Earth", "Fire", "Water"], path : []}
			queue = queue[1:]       // pop elemen awal
			fmt.Println(len(queue), currentState)

			// Cek, apakah path mengandung target
			if IsContainTarget(currentState.path, resultTarget) && PathValid(currentState.path) {
				resultAllPath = append([]string{}, currentState.path...)

				// Cek, apakah path sudah ditambahkan ke sclice hasil
				if !IsContainTargetElements(resultAll, resultAllPath) {
					resultAll = append(resultAll, resultAllPath)
				}

				// Jika result all sudah mengandung semua resulttarget atau sudah
				// ditemukan sebanyak maxRecipe elemen, break, solusi ditemukan
				if IsContainAllTarget(resultAll, resultTarget) || len(resultAll) == maxRecipe {
					break
				}
			}
			fmt.Println(resultAll)

			// Mencari kombinasi seluruh elemen pada currentState.elements
			for i := 0; i < len(currentState.elements); i++ {
				for j := i; j < len(currentState.elements); j++ {
					elem1 = currentState.elements[i]
					elem2 = currentState.elements[j]

					if elem1 == elem2 {
						combine1 = strings.Join([]string{elem1, elem2}, " + ")

						if IsInMap(combine1) {
							// tambahkan elemen ke current state node
							res = mapping[combine1]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine1, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else {
							continue
						}
					} else {
						combine1 = strings.Join([]string{elem1, elem2}, " + ")
						combine2 = strings.Join([]string{elem2, elem1}, " + ")

						if IsInMap(combine1) {
							// tambahkan elemen ke current state node
							res = mapping[combine1]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine1, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else if IsInMap(combine2) {
							res = mapping[combine2]
							if ContainElements(startElements, res) {
								continue
							}
							resultPath = strings.Join([]string{combine2, res}, " = ")

							if ContainElements(currentState.elements, res) {
								continue
							}

							newElements := append([]string{}, currentState.elements...)
							newElements = append(newElements, res)
							newPath := append([]string{}, currentState.path...)
							newPath = append(newPath, resultPath)

							newElementsCopy := append([]string{}, newElements...)
							sort.Strings(newElementsCopy)

							key := strings.Join(newElementsCopy, ",")

							if !visited[key] {
								visited[key] = true
								newCurrentState = StateNode{elements: newElements, path: newPath}
								queue = append(queue, newCurrentState)
							}

						} else {
							continue
						}
					}
				}
			}
			fmt.Println()
		}
		fmt.Println(resultAll)
	}
}
