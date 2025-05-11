package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var mapping map[string]string // data map elemen to result

// Fungsi untuk membaca data hasil scrapping
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

var reverseMapping map[string][]string // data result to kumpulan resep

// Fungsi untuk membaca data hasil scarpping
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

type ElementInfo struct { // struct untuk menyimpan data informasi setiap elemen
	Name    string    `json:"name"`
	Recipes *[]string `json:"recipes"`
	Tier    string    `json:"tier"`
}

var elementsInfo []ElementInfo          // data untuk menyimpan info elemen (nama, resep, tier)
var elementsTier = make(map[string]int) // data untuk menyimpan informasi tier setiap elemen

// Fungsi untuk membaca data hasil scrapping
func LoadElementInfo() error {
	file, err := os.Open("elements.json")
	if err != nil {
		return fmt.Errorf("error opening elements.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&elementsInfo); err != nil {
		return fmt.Errorf("error decoding elements.json: %v", err)
	}
	return nil
}

type StateNode struct { // struct untuk menyimpan state
	elements []string // Elemen yang ada saat ini
	path     []string // resep menuju elemen target
}

// Fungsi untuk melihat apakah path saat ini mengandung resep target
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

// Cek slice string apakah mengandung suatu elemen tertentu
func ContainElements(sl []string, val string) bool {
	for _, s := range sl {
		if s == val {
			return true
		}
	}
	return false
}

// Untuk memeriksa path, apakah hasil combine
// digunakan pada resep lain untuk mencapai elemen target
func PathValid(resultPath []string) bool {
	// resultPath: slice of string, ["A + B = C", "B + D = F"]
	if len(resultPath) == 0 { // Jika path kosong, path tidak valid -> bukan solusi
		return false
	}

	lastPath := resultPath[len(resultPath)-1] // ambil path terakhir

	partsPath := strings.Split(lastPath, " = ") // split path menjadi dua bagian dengan pemisah "="
	if len(partsPath) != 2 {                    // Jika potongan path setelah di split tidak terbagi menjadi 2, path tidak valid -> bukan solusi
		return false
	}
	targetPath := partsPath[1]                    // ambil elemen hasil combine
	usedElem := map[string]bool{targetPath: true} // tandai bahwa elemen tersebut sudah pernah digunakan

	// loop untuk path yang ada di belakangnya
	for i := len(resultPath) - 1; i >= 0; i-- {
		partsPath := strings.Split(resultPath[i], " = ")
		if len(partsPath) != 2 {
			continue
		}
		recipe := strings.Split(partsPath[0], " + ") // pisahkan elemen menjadi dua bagian dengan pemisah "+"
		res := partsPath[1]

		if usedElem[res] { // Jika elemen hasil combine sudah pernah digunakan, tandai resepnya menjadi true
			usedElem[recipe[0]] = true
			usedElem[recipe[1]] = true
		}
	}

	// Cek seluruh komponen path
	for _, step := range resultPath {
		partsPath := strings.Split(step, " = ") // pisahkan path berdasarkan tanda "="
		if len(partsPath) != 2 {
			continue
		}
		res := partsPath[1] // Ambil hasil combine-nya
		if !usedElem[res] {
			return false // Jika hasil combine tersebut tidak digunakan pada komponen path lain, path tidak valid
		}
	}
	return true // semua path valid (hasil combine digunakan pada komponen path yang lain)
}

// Fungsi yang memanfaatkan go routine untuk mengecek, apakah state hasil kombinasi
// sudah pernah diunjungi atau belum. Jika belum, return State hasil kombinasi tersebut.
func TryCombinePar(currState StateNode, elem1, elem2 string, visited *sync.Map) *StateNode {
	keys := []string{ // Keys untuk hasil kombinasi
		fmt.Sprintf("%s + %s", elem1, elem2),
		fmt.Sprintf("%s + %s", elem2, elem1),
	}

	for _, key := range keys {
		res, ok := mapping[key]

		// periksa  apakah kombinasi sudah ada di state path
		if !ok || ContainElements(currState.path, fmt.Sprintf("%s = %s", key, res)) || ContainElements(currState.elements, res) {
			continue
		}

		// periksa apakah tier elemen bahan melebihi tier elemen hasil combine, jika ya -> continue
		if elementsTier[elem1] >= elementsTier[res] || elementsTier[elem2] >= elementsTier[res] {
			continue
		}

		// Tambahkan elemen ke current state node
		newElements := append([]string{}, currState.elements...)
		newElements = append(newElements, res)

		newPath := append([]string{}, currState.path...)
		newPath = append(newPath, fmt.Sprintf("%s = %s", key, res))

		keyElems := append([]string{}, newElements...)
		sort.Strings(keyElems)

		keyPaths := append([]string{}, newPath...)
		sort.Strings(keyPaths)

		hash := strings.Join(keyElems, ",") + "|" + strings.Join(keyPaths, ",")

		// Jika node kombinasi element state belum pernah dikunjungi, tambahkan pada queue
		_, exists := visited.Load(hash)
		if !exists {
			visited.Store(hash, true)
			return &StateNode{elements: newElements, path: newPath}
		}
	}
	return nil
}

// Fungsi untuk mengecek, apakah state hasil kombinasi sudah pernah diunjungi atau belum.
// Jika belum, return State hasil kombinasi tersebut.
func TryCombine(currState StateNode, elem1, elem2 string, visited map[string]bool) *StateNode {
	keys := []string{ // Keys untuk hasil kombinasi
		fmt.Sprintf("%s + %s", elem1, elem2),
		fmt.Sprintf("%s + %s", elem2, elem1),
	}

	for _, key := range keys { // periksa masing-masing key
		res, ok := mapping[key]

		// periksa apakah kombinasi valid atau sudah ada di state path
		if !ok || ContainElements(currState.path, fmt.Sprintf("%s = %s", key, res)) || ContainElements(currState.elements, res) {
			continue
		}

		// periksa apakah tier elemen bahan melebihi tier elemen hasil combine, jika ya -> continue
		if elementsTier[elem1] >= elementsTier[res] || elementsTier[elem2] >= elementsTier[res] {
			continue
		}

		// Tambahkan elemen ke current state node
		newElements := append([]string{}, currState.elements...)
		newElements = append(newElements, res)

		newPath := append([]string{}, currState.path...)
		newPath = append(newPath, fmt.Sprintf("%s = %s", key, res))

		keyElems := append([]string{}, newElements...)
		sort.Strings(keyElems)

		keyPaths := append([]string{}, newPath...)
		sort.Strings(keyPaths)

		hash := strings.Join(keyPaths, ",") + "|" + strings.Join(keyPaths, ",")

		// Jika node kombinasi element state belum pernah dikunjungi, tambahkan pada queue
		if !visited[hash] {
			visited[hash] = true
			return &StateNode{elements: newElements, path: newPath}
		}
	}
	return nil
}

// Fungsi rekursif untuk menghitung banyaknya kombinasi/tree unik yang dapat dibentuk
func CountWays(target string, memo map[string]int) int {
	// memo untuk menyimpan setiap elemen sudah muncul berapa kali (value count)

	baseSet := map[string]bool{ // tandai seluruh elemen dasar sebagai true
		"Air": true, "Water": true, "Earth": true, "Fire": true,
	}

	if baseSet[target] { // apabila sudah mencapai elemen dasar, tambah jumlah kombinasi dengan 1
		return 1
	}

	if val, ok := memo[target]; ok { // Apabila target sudah pernah dikunjungi, return jumlah kemunculan elemen tersebut
		return val
	}

	total := 0                                      // total pembentukan elemen dengan resep yang unik
	for _, recipe := range reverseMapping[target] { // ambil semua resep dari suatu elemen target
		parts := strings.Split(recipe, " + ") // pisahkan kedua resep
		if len(parts) != 2 {
			continue
		}
		a, b := parts[0], parts[1]

		// periksa, apakah tier elemen a dan/atau b lebih dari tier elemen target.
		// Jika ya, lewati resep tersebut.
		if elementsTier[a] >= elementsTier[target] || elementsTier[b] >= elementsTier[target] {
			continue
		}

		var waysA, waysB int

		if a != b { // jika kedua elemen tidak sama, hitung resep unik keduanya
			waysA = CountWays(a, memo)
			waysB = CountWays(b, memo)
		} else { // jika kedua elemen sama, hitung salah satunya saja
			waysA = CountWays(a, memo)
			waysB = 1
		}

		total += waysA * waysB // jumlahkan banyak resep unik dari elemen target
	}

	memo[target] = total // simpan jumlah resep unik elemen target
	return total
}

type RecipeNode struct { // struct tree untuk path final
	Result string      // result elemen hasil combine
	Recipe string      // resep elemen
	Left   *RecipeNode // sub-pohon kiri
	Right  *RecipeNode // sub-pohon kanan
}

// Fungsi untuk membuat pohon path final
func BuildTree(path []string) *RecipeNode {
	stepMap := make(map[string][2]string) // "result": ["A", "B"]
	for _, step := range path {           // periksa seluruh path
		parts := strings.Split(step, " = ") // pisahkan hasil combine dengan resep
		if len(parts) != 2 {
			continue
		}
		result := parts[1]
		ingredients := strings.Split(parts[0], " + ")
		if len(ingredients) != 2 {
			continue
		}
		stepMap[result] = [2]string{ingredients[0], ingredients[1]} // simpan resep elemen result dalam map
	}

	visited := make(map[string]bool) // visited agar tidak mencari elemen yang sama dua kali

	// Fungsi rekursif untuk membuat pohon dari elemen target
	var build func(res string) *RecipeNode
	build = func(res string) *RecipeNode {
		if visited[res] { // jika elemen sudah dikunjungi, return nil
			return nil
		}
		visited[res] = true               // tandai sebagai sudah dikunjungi
		node := &RecipeNode{Result: res}  // buat sub-pohon elemen res
		if pair, ok := stepMap[res]; ok { // jika elemen bukan elemen dasar bangkitkan sub-pohon kiri dan kanan
			node.Recipe = fmt.Sprintf("%s + %s", pair[0], pair[1])
			node.Left = build(pair[0])
			node.Right = build(pair[1])
		}
		return node
	}

	// Bangun pohon dari hasil terakhir (elemen target)
	if len(path) == 0 {
		return nil
	}

	last := path[len(path)-1] // ambil elemen terakhir
	parts := strings.Split(last, " = ")
	if len(parts) != 2 {
		return nil
	}
	target := parts[1]   // ambil hasil combine/elemen
	return build(target) // bangun pohon
}

const workerCount = 12 // Jumlah worker paralel

type StateMessage struct { // queue untuk bfs paralel
	State StateNode
}

var nodeCount int64 // jumlah state node yang dikunjungi

func bfsManyPar(startElements []string, elementTarget string, resultTarget []string, maxRecipe int) [][]string {
	stateCh := make(chan StateMessage, int(math.Pow(10, 7))) // queue
	resultCh := make(chan []string, maxRecipe)               // path hasil temuan
	visited := &sync.Map{}                                   // map visited
	resultSet := make(map[string]bool)                       // map solusi agar tidak duplikat
	var resultMu sync.Mutex                                  // pengaman untuk mengubah variabel resultSet dan resultCh, agar tidak konflik antar worker

	startState := StateNode{elements: startElements} // state awal
	stateCh <- StateMessage{State: startState}       // tambahkan state awal ke queue

	var wg sync.WaitGroup // penanda untuk menunggu semua worker selesai

	for i := 0; i < workerCount; i++ { // masing-masing worker bekerja
		wg.Add(1)

		// Fungsi untuk mencari kombinasi elemen dari stateCh menggunakan masing-masing worker
		go func() {
			defer wg.Done()
			for msg := range stateCh {
				curr := msg.State

				// Periksa, apakah path sudah ada di result, sudah valid, dan hasil akhirnya adalah elemen target.
				if IsContainTarget(curr.path, resultTarget) &&
					PathValid(curr.path) &&
					strings.HasSuffix(curr.path[len(curr.path)-1], "= "+elementTarget) {

					hash := strings.Join(curr.path, "|")

					resultMu.Lock()
					if !resultSet[hash] { // Jika path belum pernah disimpan
						resultSet[hash] = true
						resultCh <- curr.path // kirim ke resultCh
					}
					resultMu.Unlock()
				}

				// Mencari seluruh kombinasi elemen yang ada di state sekarang
				for i := 0; i < len(curr.elements); i++ {
					for j := i; j < len(curr.elements); j++ {
						newNode := TryCombinePar(curr, curr.elements[i], curr.elements[j], visited)
						if newNode != nil {
							atomic.AddInt64(&nodeCount, 1)           // tambah node state
							stateCh <- StateMessage{State: *newNode} // kirim node baru ke queue
						}
					}
				}
			}
		}()
	}

	var results [][]string
	done := make(chan struct{}) // sinyal bahwa len(results) sudah mencapai maxRecipe

	// Fungsi untuk mengirim hasil ke results
	go func() {
		for res := range resultCh {
			results = append(results, res)
			if len(results) >= maxRecipe {
				close(done)
				return
			}
		}
	}()

	// Fungsi untuk menutup channel resultCh (semua hasil sudah dipindahkan ke results) setelah semua worker selesai
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Menunggu salah satu dari dua case terpenuhi, yaitu sudah ada resultCh atau semua worker selesai
	select {
	case <-done: // keluar program apabila sudah menemukan semua results
	case <-func() chan struct{} {
		done := make(chan struct{})
		// Fungsi untuk menutup channel queue (stateCh) apabila semua worker sudah selesai
		go func() {
			wg.Wait()
			close(stateCh)
			close(done)
		}()
		return done
	}():
	}

	return results
}

func bfsshort(startElements []string, resultTarget []string) []string {
	// BFS: tinjau dari awal elemen
	queue := []StateNode{{elements: startElements}} // queue = [StateNode {elements: ["Air", "Earth", "Fire", "Water"], path: []}]

	var resultAllPath []string
	visited := make(map[string]bool)
	for len(queue) > 0 {
		currentState := queue[0] // simpan current state, StateNode{elements: ["Air", "Earth", "Fire", "Water"], path : []}
		queue = queue[1:]        // pop elemen awal queue

		// Cek, apakah path sudah mengandung target
		if IsContainTarget(currentState.path, resultTarget) {
			resultAllPath = currentState.path
			break
		}

		// Mencari kombinasi seluruh elemen pada currentState.elements
		for i := 0; i < len(currentState.elements); i++ {
			for j := i; j < len(currentState.elements); j++ {
				elem1 := currentState.elements[i]
				elem2 := currentState.elements[j]

				newCurrentState := TryCombine(currentState, elem1, elem2, visited)
				if newCurrentState != nil {
					nodeCount += 1
					queue = append(queue, *newCurrentState)
				}
			}
		}
	}
	return resultAllPath
}

func main() {
	// // Membaca file data JSON
	if err := LoadMapping(); err != nil {
		fmt.Println("Error loading mapping: ", err)
		return
	}

	if err := LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping: ", err)
		return
	}

	if err := LoadElementInfo(); err != nil {
		fmt.Println("Error loading elements info: ", err)
		return
	}

	// Mengambil tier setiap elemen
	for _, info := range elementsInfo {
		s := strings.Split(info.Tier, " ")
		if len(s) > 1 {
			tier, err := strconv.Atoi(s[1])
			if err != nil {
				elementsTier[info.Name] = 0
			} else {
				elementsTier[info.Name] = tier
			}
		} else {
			elementsTier[info.Name] = 0
		}
	}

	fmt.Println("Mappings loaded successfully")

	elementTarget := "Stone" // element target
	method := "many"         // fast: jalur tercepat ke elemen, many: banyak resep yang ditampilkan
	maxRecipe := 3
	recipeTarget := reverseMapping[elementTarget] // resep-resep untuk membentuk elemen target

	//loop membuang resep yang tier bahannya melebihi tier elemen target
	var recipeTargetFiltered []string
	if len(recipeTarget) > 0 {
		for _, t := range recipeTarget {
			s := strings.Split(t, " + ")
			if !((elementsTier[s[0]] >= elementsTier[elementTarget]) || (elementsTier[s[1]] >= elementsTier[elementTarget])) {
				recipeTargetFiltered = append(recipeTargetFiltered, t)
			}
		}
	}

	// join elementtarget
	var resultTarget []string
	if len(recipeTargetFiltered) > 0 {
		for _, recipe := range recipeTargetFiltered {
			resultTarget = append(resultTarget, strings.Join([]string{recipe, elementTarget}, " = "))
		}
	}

	// menghitung jumlah kombinasi yang unik pada elemen target
	memo := make(map[string]int)
	count := CountWays(elementTarget, memo)
	if maxRecipe > count {
		maxRecipe = count
	} else if maxRecipe <= 0 {
		fmt.Println("Banyak resep maksimum tidak valid. Masukan angka >= 1")
	}

	startElements := []string{"Air", "Earth", "Fire", "Water"}

	var resultAllPath []string // untuk menyimpan path jalur resep terpendek
	var resultAll [][]string   // untuk menyimpan banyak resep
	var trees []*RecipeNode
	var duration string

	// Jika elemen target tidak memiliki resep atau elemen target adalah elemen dasar -> elemen berhasil ditemukan
	if recipeTargetFiltered == nil || ContainElements(startElements, elementTarget) {
		t := time.Now()
		resultAllPath = recipeTargetFiltered[:0]
		duration = time.Since(t).String()

		fmt.Printf("Resep %v ditemukan:\n%v\n", elementTarget, resultAllPath)
		trees = append(trees, BuildTree(resultAllPath))
		fmt.Printf("Total Node State: %v\n", nodeCount)
		fmt.Printf("Durasi          : %s", duration)
		return
	}

	// BFS: Jalur tercepat
	if method == "fast" {
		t := time.Now()
		resultAllPath = bfsshort(startElements, resultTarget)
		duration = time.Since(t).String()

		if resultAllPath != nil {
			fmt.Printf("Resep %v ditemukan:\n%v\n", elementTarget, resultAllPath)
		} else {
			fmt.Printf("Resep %v tidak ditemukan\n", elementTarget)
		}

		trees = append(trees, BuildTree(resultAllPath))
		fmt.Printf("Total Node State: %v\n", nodeCount)
		fmt.Printf("Durasi          : %s", duration)

		// BFS: Semua resep
	} else if method == "many" {
		t := time.Now()
		resultAll = bfsManyPar(startElements, elementTarget, resultTarget, maxRecipe)
		duration = time.Since(t).String()

		if resultAll != nil {
			fmt.Printf("Resep %v ditemukan:\n%v\n", elementTarget, resultAll)
		} else {
			fmt.Printf("Resep %v tidak ditemukan", elementTarget)
		}

		for _, path := range resultAll {
			trees = append(trees, BuildTree(path))
		}
		fmt.Printf("Total resep unik: %v\n", maxRecipe)
		fmt.Printf("Total Node State: %v\n", nodeCount)
		fmt.Printf("Durasi          : %s", duration)
	}
}
