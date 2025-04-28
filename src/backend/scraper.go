package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	// "time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

type Element struct {
	Name    string   `json:"name"`
	Recipes []string `json:"recipes"`
	Tier    string   `json:"tier"`
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
	elemFile, err := os.Create("elements.json")
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

	// Create mapping: A + B -> Result and Result -> [A + B]
	validNames := make(map[string]bool)
	for _, el := range elements {
		validNames[el.Name] = true
	}

	mapping := make(map[string]string)
	reverseMapping := make(map[string][]string)

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
			mapping[key] = el.Name
			reverseMapping[el.Name] = append(reverseMapping[el.Name], key)
		}
	}


	// Save mapping (A + B -> Result) to ElemToRes.json
	mapFile, err := os.Create("ElemToRes.json")
	if err != nil {
		fmt.Println("Error creating ElemToRes.json:", err)
		return
	}
	defer mapFile.Close()

	mapEncoder := json.NewEncoder(mapFile)
	mapEncoder.SetIndent("", "  ")
	if err := mapEncoder.Encode(mapping); err != nil {
		fmt.Println("Error writing ElemToRes.json:", err)
		return
	}

	// Save reverse mapping (Result -> [A + B]) to ResToElem.json
	reverseMapFile, err := os.Create("ResToElem.json")
	if err != nil {
		fmt.Println("Error creating ResToElem.json:", err)
		return
	}
	defer reverseMapFile.Close()

	reverseEncoder := json.NewEncoder(reverseMapFile)
	reverseEncoder.SetIndent("", "  ")
	if err := reverseEncoder.Encode(reverseMapping); err != nil {
		fmt.Println("Error writing ResToElem.json:", err)
		return
	}

	fmt.Println("Mappings saved to ElemToRes.json and ResToElem.json")
}


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

// func main() {
// 	Run()
// 	if err := LoadMapping(); err != nil {
// 		fmt.Println("Error loading mapping:", err)
// 		return
// 	}
// 	if err := LoadReverseMapping(); err != nil {
// 		fmt.Println("Error loading reverse mapping:", err)
// 		return
// 	}
// 	fmt.Println("Mappings loaded successfully")
// }
