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

func (g *Graph) Add(msg string) {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		g.logger.Warn(fmt.Sprintf("[graph] empty msg %q in Add", msg))
	}
	g.logger.Debug(fmt.Sprintf("[graph] tokens: %v", tokens))

	lastNode := g.Start

	for i, token := range tokens {
		g.logger.Debug(fmt.Sprintf("[graph] token: %q", token))
		_, exists := g.Nodes[token]
		if !exists {
			newNode := NewNode(token, 0)
			g.Nodes[token] = newNode
		}

		node := g.Nodes[token]

		lastNode.AddNeightbour(node)
		lastNode = node

		if i == len(tokens)-1 {
			g.AddEndMsg(node)
		}

		g.Nodes[token].Weight++
	}
}

func (g *Graph) AddEndMsg(node *Node) {
	if !slices.Contains(node.Neigbours, g.End) {
		node.Neigbours = append(node.Neigbours, g.End)
	}
}

func (n *Node) AddNeightbour(node *Node) {
	if !slices.Contains(n.Neigbours, node) {
		n.Neigbours = append(n.Neigbours, node)
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d %v", n.Word, n.Weight, n.Neigbours)
}
