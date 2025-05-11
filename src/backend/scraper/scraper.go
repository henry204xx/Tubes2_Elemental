package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	// "time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parseParts(recipe string) (string, string) {
	parts := strings.Split(recipe, "+")
	a := strings.TrimSpace(parts[0])
	b := strings.TrimSpace(parts[1])
	return a, b
}

type Element struct {
	Name    string   `json:"name"`
	Recipes []string `json:"recipes"`
	Tier    string   `json:"tier"`
}

// getScraperDir returns the absolute path to the directory of this file.
func getScraperDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Cannot get current file info")
	}
	return filepath.Dir(filename)
}

// helper to open a file in the scraper directory
func openInScraperDir(filename string) (*os.File, error) {
	return os.Open(filepath.Join(getScraperDir(), filename))
}

// helper to create a file in the scraper directory
func createInScraperDir(filename string) (*os.File, error) {
	return os.Create(filepath.Join(getScraperDir(), filename))
}

func Run() {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

	var elements []Element

	c.OnHTML("div.mw-parser-output", func(e *colly.HTMLElement) {
		currentTier := ""
		e.DOM.Children().Each(func(_ int, s *goquery.Selection) {
			if goquery.NodeName(s) == "h3" {
				span := s.Find(".mw-headline")
				text := span.Text()
				if strings.HasPrefix(text, "Tier") {
					currentTier = text
				}
			}

			if goquery.NodeName(s) == "table" && s.HasClass("list-table") {
				s.Find("tbody tr").Each(func(_ int, tr *goquery.Selection) {
					name := tr.Find("td:nth-child(1) a").Text()
					if name == "" {
						return
					}

					var recipes []string
					tr.Find("td:nth-child(2) li").Each(func(_ int, li *goquery.Selection) {
						recipes = append(recipes, strings.TrimSpace(li.Text()))
					})

					elem := Element{
						Name:    name,
						Recipes: recipes,
						Tier:    currentTier,
					}
					elements = append(elements, elem)
				})
			}
		})
	})

	if err := c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"); err != nil {
		fmt.Println("Error visiting the page:", err)
		return
	}

	fmt.Printf("Scraped %d elements\n", len(elements))

	// Save elements.json
	elemFile, err := createInScraperDir("elements.json")
	if err != nil {
		fmt.Println("Error creating elements.json:", err)
		return
	}
	defer elemFile.Close()

	encoder := json.NewEncoder(elemFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(elements); err != nil {
		fmt.Println("Error writing elements.json:", err)
		return
	}

	// Create Mapping: A + B -> Result and Result -> [A + B]
	validNames := make(map[string]bool)
	for _, el := range elements {
		validNames[el.Name] = true
	}

	Mapping := make(map[string]string)
	ReverseMapping := make(map[string][]string)

	for _, el := range elements {
		for _, recipe := range el.Recipes {
			trimmed := strings.TrimSpace(recipe)
			parts := strings.Split(trimmed, "+")
			if len(parts) != 2 {
				continue
			}
			a := strings.TrimSpace(parts[0])
			b := strings.TrimSpace(parts[1])

			// Check if both parts are valid
			if !validNames[a] || !validNames[b] {
				continue // skip if invalid component
			}

			key := a + " + " + b
			Mapping[key] = el.Name
			ReverseMapping[el.Name] = append(ReverseMapping[el.Name], key)
		}
	}


	// Save Mapping (A + B -> Result) to ElemToRes.json
	mapFile, err := createInScraperDir("ElemToRes.json")
	if err != nil {
		fmt.Println("Error creating ElemToRes.json:", err)
		return
	}
	defer mapFile.Close()

	mapEncoder := json.NewEncoder(mapFile)
	mapEncoder.SetIndent("", "  ")
	if err := mapEncoder.Encode(Mapping); err != nil {
		fmt.Println("Error writing ElemToRes.json:", err)
		return
	}

	// Save reverse Mapping (Result -> [A + B]) to ResToElem.json
	reverseMapFile, err := createInScraperDir("ResToElem.json")
	if err != nil {
		fmt.Println("Error creating ResToElem.json:", err)
		return
	}
	defer reverseMapFile.Close()

	reverseEncoder := json.NewEncoder(reverseMapFile)
	reverseEncoder.SetIndent("", "  ")
	if err := reverseEncoder.Encode(ReverseMapping); err != nil {
		fmt.Println("Error writing ResToElem.json:", err)
		return
	}

	fmt.Println("Mappings saved to ElemToRes.json and ResToElem.json")

	// Map textual tiers to numeric values
	TierMap := make(map[string]int)
	for _, el := range elements {
		t := el.Tier
		if t == "" {
			TierMap[el.Name] = 0
			if (el.Name == "Time") {
				TierMap[el.Name] = 16
			}
		} else {
			var tierNum int
			fmt.Sscanf(t, "Tier %d", &tierNum)
			TierMap[el.Name] = tierNum
		}
	}

	// Create sorted Mapping
	SortedMapping := make(map[string][]string)
	for result, recipes := range ReverseMapping {
		// Sort recipes by max(tier(A), tier(B))
		sorted := make([]string, len(recipes))
		copy(sorted, recipes)

		// Custom sort
		sort.Slice(sorted, func(i, j int) bool {
			a1, a2 := parseParts(sorted[i])
			b1, b2 := parseParts(sorted[j])
		
			tA1, tA2 := TierMap[a1], TierMap[a2]
			tB1, tB2 := TierMap[b1], TierMap[b2]
		
			maxA := max(tA1, tA2)
			minA := min(tA1, tA2)
		
			maxB := max(tB1, tB2)
			minB := min(tB1, tB2)
		
			if maxA != maxB {
				return maxA > maxB
			}
			if minA != minB {
				return minA > minB
			}
			return sorted[i] > sorted[j] // alphabetical fallback
		})

		SortedMapping[result] = sorted
	}

	// Save to sort.json
	sortFile, err := createInScraperDir("sort.json")
	if err != nil {
		fmt.Println("Error creating sort.json:", err)
		return
	}
	defer sortFile.Close()

	sortEncoder := json.NewEncoder(sortFile)
	sortEncoder.SetIndent("", "  ")
	if err := sortEncoder.Encode(SortedMapping); err != nil {
		fmt.Println("Error writing sort.json:", err)
		return
	}

	fmt.Println("Sorted recipe Mapping saved to sort.json")
	
}

var Mapping map[string]string

func LoadMapping() error {
	file, err := openInScraperDir("ElemToRes.json")
	if err != nil {
		return fmt.Errorf("error opening ElemToRes.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Mapping); err != nil {
		return fmt.Errorf("error decoding ElemToRes.json: %v", err)
	}
	return nil
}

var ReverseMapping map[string][]string
func LoadReverseMapping() error {
	file, err := openInScraperDir("ResToElem.json")
	if err != nil {
		return fmt.Errorf("error opening ResToElem.json: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&ReverseMapping); err != nil {
		return fmt.Errorf("error decoding ResToElem.json: %v", err)
	}
	return nil
}

var ElemTier map[string]int

func LoadTierElem() error {
	// Open elements.json to read data
	file, err := openInScraperDir("elements.json")
	if err != nil {
		return fmt.Errorf("error opening elements.json: %v", err)
	}
	defer file.Close()

	// Read and decode the JSON into a slice of Element
	var elements []Element
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&elements); err != nil {
		return fmt.Errorf("error decoding elements.json: %v", err)
	}

	// Initialize the ElemTier map
	ElemTier = make(map[string]int)

	// Map each element to its tier (as an integer)
	for _, el := range elements {
		// If the tier is not empty, extract the numeric tier value
		if el.Tier != "" {
			var tierNum int
			// Parse tier value (e.g., "Tier 8 elements" -> 8)
			fmt.Sscanf(el.Tier, "Tier %d", &tierNum)
			ElemTier[el.Name] = tierNum
		}
	}

	return nil
}