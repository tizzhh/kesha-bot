package processor

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type Node struct {
	Word      string
	Weight    int
	Neigbours []*Node
}

type Graph struct {
	Nodes map[string]*Node
	Start *Node
	End   *Node

	logger *slog.Logger
}

func NewGraph(logger *slog.Logger) *Graph {
	start := &Node{}
	end := &Node{}

	return &Graph{
		Nodes: map[string]*Node{},
		Start: start,
		End:   end,

		logger: logger,
	}
}

func NewNode(word string, weight int, neighbours ...*Node) *Node {
	newNode := &Node{
		Word:      word,
		Neigbours: make([]*Node, len(neighbours)),
		Weight:    weight,
	}

	copy(newNode.Neigbours, neighbours)

	return newNode
}

func (g *Graph) AddMsg(msg string) {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		g.logger.Warn(fmt.Sprintf("[graph] empty msg %q in Add", msg))
	}

	prevNode := g.Start

	for _, token := range tokens {
		node := g.GetOrCreateNode(token)

		prevNode.AddNeighbour(node)
		prevNode = node

		g.Nodes[token].Weight++
	}

	lastNode := g.Nodes[tokens[len(tokens)-1]]
	lastNode.AddNeighbour(g.End)
}

func (g *Graph) GetOrCreateNode(token string) *Node {
	_, exists := g.Nodes[token]
	if !exists {
		newNode := NewNode(token, 0)
		g.Nodes[token] = newNode
	}

	return g.Nodes[token]
}

func (n *Node) AddNeighbour(node *Node) {
	if !slices.Contains(n.Neigbours, node) {
		n.Neigbours = append(n.Neigbours, node)
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d %v", n.Word, n.Weight, n.Neigbours)
}
