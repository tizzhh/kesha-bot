package processor

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummyLogger = slog.New(slog.DiscardHandler)
)

func TestProcessor(t *testing.T) {
	t.Run("it creates correct graph with start and end", func(t *testing.T) {
		graph := NewGraph(dummyLogger)
		expectedStartNeighbours := 0

		assert.Len(t, graph.Start.Neigbours, expectedStartNeighbours)
	})

	t.Run("it adds one word message to graph", func(t *testing.T) {
		graph := NewGraph(dummyLogger)

		expectedStartNeighbours := 1
		expectedWord := "aboba"
		expectedNode := NewNode(expectedWord, 1, graph.End)
		expectedMap := map[string]*Node{
			expectedWord: expectedNode,
		}

		graph.AddMsg(expectedWord)

		assert.Len(t, graph.Start.Neigbours, expectedStartNeighbours)
		assert.Equal(t, expectedNode, graph.Start.Neigbours[0])
		assert.Equal(t, expectedMap, graph.Nodes)
	})

	t.Run("it creates correct graph with a provided message", func(t *testing.T) {
		graph := NewGraph(dummyLogger)

		msg := "hello, are you"
		graph.AddMsg(msg)

		node1 := NewNode("hello,", 1)
		start := NewNode("", 0, node1)
		node2 := NewNode("are", 1)
		node3 := NewNode("you", 1, graph.End)
		node1.AddNeighbour(node2)
		node2.AddNeighbour(node3)

		expctedMap := map[string]*Node{
			"hello,": node1,
			"are":    node2,
			"you":    node3,
		}
		expectedGraph := &Graph{
			Start: start,
			Nodes: expctedMap,
		}

		assert.Equal(t, expectedGraph.Start, graph.Start)
		assert.Equal(t, expectedGraph.Nodes, graph.Nodes)
	})

	t.Run("it creates correct graph with 2 colliding messages", func(t *testing.T) {
		graph := NewGraph(dummyLogger)

		msg1 := "hello, are you"
		graph.AddMsg(msg1)
		msg2 := "privet, are you"
		graph.AddMsg(msg2)

		node1 := NewNode("hello,", 1)
		start := NewNode("", 0, node1)
		node2 := NewNode("are", 2)
		node3 := NewNode("you", 2, graph.End)
		node1.AddNeighbour(node2)
		node2.AddNeighbour(node3)

		node4 := NewNode("privet,", 1, node2)
		start.AddNeighbour(node4)

		expctedMap := map[string]*Node{
			"hello,":  node1,
			"are":     node2,
			"you":     node3,
			"privet,": node4,
		}
		expectedGraph := &Graph{
			Start: start,
			Nodes: expctedMap,
		}

		assert.Equal(t, expectedGraph.Start, graph.Start)
		assert.Equal(t, expectedGraph.Nodes, graph.Nodes)
	})
}
