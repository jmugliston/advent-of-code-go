package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/jmugliston/aoc/graph"
	"github.com/jmugliston/aoc/parsing"
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
	}
}

func parseInput(input string) graph.Graph {
	lines := parsing.ReadLines(input)

	g := graph.Graph{
		Nodes: make([]*graph.Node, 0),
	}

	for _, line := range lines {
		split := strings.Split(line, ":")

		source := split[0]
		destinations := strings.Fields(split[1])

		g.AddNode(source)

		for _, dest := range destinations {
			g.AddNode(dest)
			// Set the edge name as the original source/target node names
			g.AddEdge(source+"-"+dest, source, dest, []string{source, dest})
		}
	}

	return g
}

// Karger's algorithm for finding the minimum cut in a graph
// https://en.wikipedia.org/wiki/Karger%27s_algorithm
// Returns the edges that were cut, and the two partitions
func minCut(g graph.Graph) ([]*graph.Edge, [][]string) {

	for {
		if len(g.Nodes) == 2 {
			break
		}

		edge := g.Edges[rand.IntN(len(g.Edges))]

		a := edge.Source
		b := edge.Target

		// Merge b into a
		for _, edge := range g.Edges {
			if edge.Source == b {
				edge.Source = a
			}
			if edge.Target == b {
				edge.Target = a
			}
		}

		// Remove self loops
		for _, edge := range g.Edges {
			if edge.Source == edge.Target {
				g.RemoveEdge(edge.Source, edge.Target)
			}
		}

		aNode, _ := g.GetNode(a)
		bNode, _ := g.GetNode(b)

		aNode.Data = append(aNode.Data, bNode.Data...)
		aNode.Data = append(aNode.Data, b)

		g.RemoveNode(b)
	}

	cuts := make([]*graph.Edge, 0)
	for _, edge := range g.Edges {
		cuts = append(cuts, &graph.Edge{
			// Get the original node names from the edge data
			Source: edge.Data[0],
			Target: edge.Data[1],
		})
	}

	partitions := make([][]string, 2)

	partitions[0] = g.Nodes[0].Data
	partitions[0] = append(partitions[0], g.Nodes[0].Name)

	partitions[1] = g.Nodes[1].Data
	partitions[1] = append(partitions[1], g.Nodes[1].Name)

	return cuts, partitions
}

func Part1(input string) int {

	g := parseInput(input)

	// Karger's algorithm is randomized so we run it until we find
	// the min cut of 3 (which we know is the right number of cuts)
	for {
		result := make(chan int)
		wait := make(chan struct{})

		var waitGroup sync.WaitGroup
		go func() {
			// Try 100 iterations in parallel
			for i := 0; i < 100; i++ {
				waitGroup.Add(1)
				go func() {
					defer waitGroup.Done()
					cuts, partitions := minCut(g.Clone())
					if len(cuts) == 3 {
						result <- len(partitions[0]) * len(partitions[1])
					}
				}()
			}
			waitGroup.Wait()
			close(wait)
		}()

		// Wait for threads to finish or a result to be found
		select {
		case answer := <-result:
			return answer
		case <-wait:
			// No result found, try again
		}
	}

}
