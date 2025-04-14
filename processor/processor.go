package processor

import (
	"fmt"
	"log/slog"
	"strings"
)

type Node struct {
	Word      string
	Weight    int
	Neigbours map[string]*Node
}

type Processor struct {
	Nodes map[string]*Node
	Start *Node
	End   *Node

	logger *slog.Logger
}

func NewProcessor(logger *slog.Logger) *Processor {
	start := NewNode("", 0)
	end := NewNode("", 0)

	return &Processor{
		Nodes: map[string]*Node{},
		Start: start,
		End:   end,

		logger: logger,
	}
}

func NewNode(word string, weight int, neighbours ...*Node) *Node {
	newNode := &Node{
		Word:      word,
		Neigbours: map[string]*Node{},
		Weight:    weight,
	}

	for _, v := range neighbours {
		newNode.Neigbours[v.Word] = v
	}

	return newNode
}

func (g *Processor) AddMsg(msg string) {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		g.logger.Warn(fmt.Sprintf("[processor] empty msg %q in Add", msg))
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

func (g *Processor) GetOrCreateNode(token string) *Node {
	_, exists := g.Nodes[token]
	if !exists {
		newNode := NewNode(token, 0)
		g.Nodes[token] = newNode
	}

	return g.Nodes[token]
}

func (n *Node) AddNeighbour(node *Node) {
	if _, exists := n.Neigbours[node.Word]; !exists {
		n.Neigbours[node.Word] = node
	}
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d %v", n.Word, n.Weight, n.Neigbours)
}
