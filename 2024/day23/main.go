package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jmugliston/aoc/graph"
	"github.com/jmugliston/aoc/parsing"
	"golang.org/x/exp/slices"
	"gonum.org/v1/gonum/stat/combin"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")
var exampleFlag = flag.Bool("example", false, "Use the example instead of the puzzle input")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)

	inputFile := "input.txt"
	if *exampleFlag {
		inputFile = "example.txt"
	}
	path := filepath.Join(filepath.Dir(filename), "input", inputFile)

	input, err := os.ReadFile(path)

	if err != nil {
		panic("Could not find the input file")
	}

	if *partFlag == "1" {
		fmt.Println(Part1(string(input)))
	} else {
		fmt.Println(Part2(string(input)))
	}
}

func parseInput(input string, directed bool) graph.Graph {
	lines := parsing.ReadLines(input)

	g := graph.Graph{
		Nodes: make([]*graph.Node, 0),
	}

	for _, line := range lines {
		split := strings.Split(line, "-")

		source := split[0]
		destination := split[1]

		g.AddNode(source)
		g.AddNode(destination)
		g.AddEdge(source+"-"+destination, source, destination, []string{})
		if !directed {
			g.AddEdge(destination+"-"+source, destination, source, []string{})
		}
	}

	return g
}

func FindCycles(g *graph.Graph, source *graph.Node, maxPathLen int) [][]string {

	type QueueItem struct {
		Node *graph.Node
		Path []string
	}

	queue := []QueueItem{}

	queue = append(queue, QueueItem{Node: source, Path: []string{}})

	possiblePaths := [][]string{}

	firstCheck := true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if len(current.Path) > maxPathLen {
			continue
		}

		if current.Node.Name == source.Name && !firstCheck {
			possiblePaths = append(possiblePaths, current.Path)
			continue
		}

		for _, edge := range g.Edges {
			if edge.Source == current.Node.Name {
				targetNode, _ := g.GetNode(edge.Target)

				alreadyVisited := false
				for _, visitedNode := range current.Path {
					if visitedNode == targetNode.Name {
						alreadyVisited = true
						break
					}
				}

				if alreadyVisited {
					continue
				}

				nextPath := make([]string, len(current.Path))
				copy(nextPath, current.Path)
				nextPath = append(nextPath, edge.Target)

				queue = append(queue, QueueItem{Node: targetNode, Path: nextPath})
			}
		}

		firstCheck = false
	}

	return possiblePaths
}

func ReachableNodes(g *graph.Graph, source *graph.Node, maxHops int) []*graph.Node {
	type QueueItem struct {
		Node *graph.Node
		Hops int
	}

	visited := map[string]bool{}
	queue := []QueueItem{}

	queue = append(queue, QueueItem{Node: source, Hops: 0})

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current.Node.Name] {
			continue
		}

		visited[current.Node.Name] = true

		if current.Hops < maxHops {
			for _, edge := range g.Edges {
				if edge.Source == current.Node.Name {
					targetNode, _ := g.GetNode(edge.Target)
					queue = append(queue, QueueItem{Node: targetNode, Hops: current.Hops + 1})
				}
			}
		}
	}

	reachable := []*graph.Node{}
	for _, node := range g.Nodes {
		if visited[node.Name] {
			reachable = append(reachable, node)
		}
	}

	return reachable
}

func findIntersection(slice1, slice2 []string) []string {
	set := make(map[string]bool)
	for _, v := range slice1 {
		set[v] = true
	}

	intersection := []string{}
	for _, v := range slice2 {
		if set[v] {
			intersection = append(intersection, v)
		}
	}

	return intersection
}

func Part1(input string) int {
	lanGraph := parseInput(input, false)

	uniquePaths := make(map[string]bool)

	for _, node := range lanGraph.Nodes {
		// Very slow for the full input - could definitely be optimised
		nextPaths := FindCycles(&lanGraph, node, 3)

		for _, path := range nextPaths {
			if len(path) == 3 {
				containsT := false
				for _, p := range path {
					if strings.HasPrefix(p, "t") {
						containsT = true
						break
					}
				}
				if containsT {
					slices.Sort(path)
					key := strings.Join(path, ":")
					uniquePaths[key] = true
				}
			}
		}
	}

	return len(uniquePaths)
}

func Part2(input string) string {
	lanGraph := parseInput(input, false)

	reachableNodesMap := make(map[string][]string)
	for _, node := range lanGraph.Nodes {
		reachableNodes := ReachableNodes(&lanGraph, node, 1)
		reachableNodeNames := make([]string, len(reachableNodes))
		for i, n := range reachableNodes {
			reachableNodeNames[i] = n.Name
		}
		reachableNodesMap[node.Name] = reachableNodeNames
	}

	intersectionMap := make(map[string][]string)
	for k, v := range reachableNodesMap {
		for k2, v2 := range reachableNodesMap {
			if k == k2 {
				continue
			}
			intersection := findIntersection(v, v2)
			if len(intersection) > 0 {
				intersectionMap[k+":"+k2] = intersection
			}
		}
	}

	maxIntersection := []string{}
	for _, v := range intersectionMap {
		if len(v) > 2 {
			allEqual := true

			c := combin.Combinations(len(v), 2)

			for _, pair := range c {
				key := v[pair[1]] + ":" + v[pair[0]]
				if !slices.Equal(v, intersectionMap[key]) {
					allEqual = false
					break
				}
			}

			if allEqual {
				if len(v) > len(maxIntersection) {
					maxIntersection = v
				}
			}
		}
	}

	slices.Sort(maxIntersection)

	return strings.Join(maxIntersection, ",")
}
