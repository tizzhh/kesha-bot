package processor

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

var (
	ErrEmptyMessage   = errors.New("empty message")
	ErrNoNeighbours   = errors.New("node has no neighbours")
	ErrRandomWalkFail = errors.New("failed to get a random neighbour")
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
}

func NewProcessor() *Processor {
	start := NewNode("", 1)
	end := NewNode("", 1)

	return &Processor{
		Nodes: map[string]*Node{},
		Start: start,
		End:   end,
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

func (p *Processor) AddMsg(msg string) error {
	tokens := strings.Fields(msg)
	if len(tokens) == 0 {
		return fmt.Errorf("[processor] empty msg %q in Add: %w", msg, ErrEmptyMessage)
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

	return nil
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

func (p *Processor) Generate() (string, error) {
	newMsg := strings.Builder{}
	cur := p.Start

	var err error
	for cur != p.End {
		cur, err = cur.chooseRandomNeighbour()
		if err != nil {
			return "", fmt.Errorf("[processor] failed to choose random neighbour: %w", err)
		}
		newMsg.WriteString(cur.Word)
		newMsg.WriteString(" ")
	}

	res := newMsg.String()
	return res[:len(res)-2], nil
}

func (n *Node) chooseRandomNeighbour() (*Node, error) {
	if len(n.Neigbours) == 0 {
		return nil, fmt.Errorf("[node %v]: 0 neighbours: %w", n, ErrNoNeighbours)
	}

	sumOfWeights := 0
	for _, neighbour := range n.Neigbours {
		sumOfWeights += neighbour.Weight
	}

	rnd := rand.Intn(sumOfWeights) //nolint:gosec
	for _, neighbour := range n.Neigbours {
		if rnd < neighbour.Weight {
			return neighbour, nil
		}
		rnd -= neighbour.Weight
	}

	return nil, fmt.Errorf("[node %v]: failed to choose a neighbour: %w", n, ErrRandomWalkFail)
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %q %d", n.Word, n.Weight)
}
