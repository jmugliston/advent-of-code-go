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

	"github.com/atheius/aoc/parsing"
	"github.com/atheius/aoc/utils"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "input", "input.txt")

	input, err := os.ReadFile(path)

	if err != nil {
		panic("Could not find the input file")
	}

	if *partFlag == "1" {
		fmt.Println(Part1(string(input)))
	}
}

type Graph struct {
	Nodes []*Node
	Edges []*Edge
}

type Node struct {
	Name string
	Data []string
}

type Edge struct {
	Source string
	Target string
	Data   []string
}

func (g *Graph) Clone() Graph {
	clone := Graph{
		Nodes: make([]*Node, len(g.Nodes)),
		Edges: make([]*Edge, len(g.Edges)),
	}

	for i, node := range g.Nodes {
		clone.Nodes[i] = &Node{
			Name: node.Name,
			Data: make([]string, len(node.Data)),
		}
		copy(clone.Nodes[i].Data, node.Data)
	}

	for i, edge := range g.Edges {
		clone.Edges[i] = &Edge{
			Source: edge.Source,
			Target: edge.Target,
			Data:   edge.Data,
		}
	}

	return clone
}

func (g *Graph) GetNode(n string) (*Node, error) {
	for _, node := range g.Nodes {
		if node.Name == n {
			return node, nil
		}
	}
	return &Node{}, fmt.Errorf("Node not found")
}

func (g *Graph) AddNode(n string) error {
	existing, _ := g.GetNode(n)

	if existing.Name != "" {
		return nil
	}

	g.Nodes = append(g.Nodes, &Node{
		Name: n,
	})

	return nil
}

func (g *Graph) RemoveNode(n string) error {
	g.Nodes = utils.Filter(g.Nodes, func(node *Node) bool {
		return node.Name != n
	})

	return nil
}

func (g *Graph) AddEdge(name, source, target string) error {
	_, err := g.GetNode(source)

	if err != nil {
		return err
	}

	_, err = g.GetNode(target)

	if err != nil {
		return err
	}

	g.Edges = append(g.Edges, &Edge{
		Source: source,
		Target: target,
		// Store the original node soure/target in data (used for final output)
		Data: []string{source, target},
	})

	return nil
}

func (g *Graph) RemoveEdge(source, target string) {
	g.Edges = utils.Filter(g.Edges, func(edge *Edge) bool {
		return !(edge.Source == source && edge.Target == target)
	})
}

func parseInput(input string) Graph {
	lines := parsing.ReadLines(input)

	g := Graph{
		Nodes: make([]*Node, 0),
	}

	for _, line := range lines {
		split := strings.Split(line, ":")

		source := split[0]
		destinations := strings.Fields(split[1])

		g.AddNode(source)

		for _, dest := range destinations {
			g.AddNode(dest)
			// Set the edge name as the original source/target node names
			g.AddEdge(source+"-"+dest, source, dest)
		}
	}

	return g
}

// Karger's algorithm for finding the minimum cut in a graph
// https://en.wikipedia.org/wiki/Karger%27s_algorithm
// Returns the edges that were cut, and the two partitions
func minCut(g Graph) ([]*Edge, [][]string) {

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

	cuts := make([]*Edge, 0)
	for _, edge := range g.Edges {
		cuts = append(cuts, &Edge{
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
