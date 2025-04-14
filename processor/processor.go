package processor

import (
	"fmt"
	"log"
	"math/rand"
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

	logger *log.Logger
}

func NewProcessor(logger *log.Logger) *Processor {
	start := NewNode("", 1)
	end := NewNode("", 1)

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

func (p *Processor) AddMsg(msg string) {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		p.logger.Printf("[processor] empty msg %q in Add", msg)
	}

	prevNode := p.Start

	for _, token := range tokens {
		node := p.GetOrCreateNode(token)

		prevNode.AddNeighbour(node)
		prevNode = node

		p.Nodes[token].Weight++
	}

	lastNode := p.Nodes[tokens[len(tokens)-1]]
	lastNode.AddNeighbour(p.End)
}

func (p *Processor) GetOrCreateNode(token string) *Node {
	_, exists := p.Nodes[token]
	if !exists {
		newNode := NewNode(token, 0)
		p.Nodes[token] = newNode
	}

	return p.Nodes[token]
}

func (n *Node) AddNeighbour(node *Node) {
	if _, exists := n.Neigbours[node.Word]; !exists {
		n.Neigbours[node.Word] = node
	}
}

func (p *Processor) Generate() string {
	newMsg := strings.Builder{}
	cur := p.Start

	for cur != p.End {
		cur = cur.chooseRandomNeighbour(p.logger)
		newMsg.WriteString(cur.Word)
		newMsg.WriteString(" ")
	}

	res := newMsg.String()
	return res[:len(res)-2]
}

func (n *Node) chooseRandomNeighbour(logger *log.Logger) *Node {
	if len(n.Neigbours) == 0 {
		log.Fatalf("[node %v]: chooseRandomNeighbour 0 neighbours\n", n)
	}

	sumOfWeights := 0
	for _, neighbour := range n.Neigbours {
		sumOfWeights += neighbour.Weight
	}

	rnd := rand.Intn(sumOfWeights) //nolint:gosec
	for _, neighbour := range n.Neigbours {
		if rnd < neighbour.Weight {
			return neighbour
		}
		rnd -= neighbour.Weight
	}

	logger.Fatalf("[node %v]: chooseRandomNeighbour failed\n", n)
	return nil
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d %v", n.Word, n.Weight, n.Neigbours)
}
