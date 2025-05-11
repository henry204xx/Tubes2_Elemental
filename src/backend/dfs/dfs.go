package dfs

import (
	// "encoding/json"
	"backend/scraper"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var GlobalVisitedCount int

type Step struct {
	Result     string
	Components []string
}

type Stack struct {
	Elements  []string
	Path      []Step
	Visited   map[string]bool
	Ancestry  []string
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

func (s *Stack) push(element string, step Step) {
	s.Elements = append(s.Elements, element)
	s.Path = append(s.Path, step)
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
		Ancestry:  append([]string{}, s.Ancestry...),
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

type TreeNode struct {
	Result     string     `json:"result"`
	Components []*TreeNode `json:"components,omitempty"`
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

        // If we've already expanded this element, return just the node
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


func DFS(root string, maxSolution int) {
	var allPaths []*TreeNode
	Count := 0

	// Load data as needed
	if err := scraper.LoadReverseMapping(); err != nil {
		fmt.Println("Error loading reverse mapping:", err)
		return
	}

	if err := scraper.LoadTierElem(); err != nil {
		fmt.Println("Error loading tiers:", err)
		return
	}

	stack := Stack{
		Visited: make(map[string]bool),
	}
	stack.push(root, Step{Result: root, Components: []string{}})

	type Node struct {
		stack Stack
	}

	var stacks []Node
	stacks = append(stacks, Node{stack: stack})

	for len(stacks) > 0 {
		currentNode := stacks[len(stacks)-1]
		stacks = stacks[:len(stacks)-1]
		current := currentNode.stack

		if len(current.Elements) == 0 {

			if ContainsElementInPath(current.Path, "Time") {
				afterTimeCount := UniqueElementsAfterTime(current.Path)
				if afterTimeCount < 100 {
					continue // invalid path: not enough elements created after Time
				}
			}

			fmt.Println(Count+1, "Found complete path:")
			printStack(current)

			tree := BuildTree(current.Path, root)
			allPaths = append(allPaths, tree)

			Count++
			if Count >= maxSolution {
				break
			}
			continue
		}

		elem := current.Elements[len(current.Elements)-1]
		current.Elements = current.Elements[:len(current.Elements)-1]

		if current.Visited[elem] {
			stacks = append(stacks, Node{stack: current})
			continue
		}

		GlobalVisitedCount++

		if IsBasicElement(elem) || elem == "Time" {
			stacks = append(stacks, Node{stack: current})
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
			newStack := copyStack(current)
			newStack.Ancestry = append(append([]string{}, current.Ancestry...), elem)
			newStack.Elements = append(newStack.Elements, b, a)

			newStack.Path = append(newStack.Path, Step{
				Result:     elem,
				Components: []string{a, b},
			})

			newStack.Visited[elem] = true

			stacks = append(stacks, Node{stack: newStack})
		}
	}

	// Write merged paths to files
	f, err := os.Create("paths.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer f.Close()

	for i, tree := range allPaths {
		fmt.Fprintf(f, "Path %d:\n", i+1)
		fmt.Fprintf(f, "%s", root + "\n")
		PrintTree(tree, "", f)
		fmt.Fprintln(f)
	}

	jsonFile, err := os.Create("paths.json")
    if err != nil {
        fmt.Println("Error creating JSON file:", err)
        return
    }
    defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(allPaths); err != nil {
		fmt.Println("Error encoding to JSON:", err)
		return
	}
}

func printStack(s Stack) {
    fmt.Println("Elements:", s.Elements)
    fmt.Println("Paths:")
    for _, p := range s.Path {
        fmt.Println(p.Result, "->", p.Components)
    }
}
