package graph

import (
	"fmt"

	"github.com/jmugliston/aoc/utils"
)

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

func (g *Graph) AddEdge(name, source, target string, data []string) error {
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
		Data:   data,
	})

	return nil
}

func (g *Graph) RemoveEdge(source, target string) {
	g.Edges = utils.Filter(g.Edges, func(edge *Edge) bool {
		return !(edge.Source == source && edge.Target == target)
	})
}
