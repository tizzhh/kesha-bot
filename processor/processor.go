package processor

import (
	"fmt"
	"log/slog"
	"os"
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

func NewGraph() *Graph {
	start := &Node{}
	end := &Node{}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

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

	for i := range len(neighbours) {
		newNode.Neigbours[i] = neighbours[i]
	}

	return newNode
}

func (g *Graph) AddMsg(msg string) {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		g.logger.Warn(fmt.Sprintf("[graph] empty msg %q in Add", msg))
	}

	prevNode := g.Start

	for _, token := range tokens {
		_, exists := g.Nodes[token]
		if !exists {
			newNode := NewNode(token, 0)
			g.Nodes[token] = newNode
		}

		node := g.Nodes[token]

		prevNode.AddNeighbour(node)
		prevNode = node

		g.Nodes[token].Weight++
	}

	lastNode := g.Nodes[tokens[len(tokens)-1]]
	lastNode.AddNeighbour(g.End)
}

func (n *Node) AddNeighbour(node *Node) {
	if !slices.Contains(n.Neigbours, node) {
		n.Neigbours = append(n.Neigbours, node)
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d %v", n.Word, n.Weight, n.Neigbours)
}
