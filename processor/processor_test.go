package processor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessor(t *testing.T) {
	t.Run("it creates correct graph with start and end", func(t *testing.T) {
		graph := NewGraph()
		expectedStartNeighbours := 0

		assert.Len(t, graph.Start.Neigbours, expectedStartNeighbours)
	})

	t.Run("it adds one word message to graph", func(t *testing.T) {
		graph := NewGraph()

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
		graph := NewGraph()

		msg := "hello, are you"

		node1 := NewNode("hello,", 1)
		start := NewNode("", 0, node1)
		node2 := NewNode("are", 1)
		node3 := NewNode("you", 1, graph.End)
		node1.Neigbours = append(node1.Neigbours, node2)
		node2.Neigbours = append(node2.Neigbours, node3)

		expctedMap := map[string]*Node{
			"hello,": node1,
			"are":    node2,
			"you":    node3,
		}
		expectedGraph := &Graph{
			Start: start,
			Nodes: expctedMap,
		}

		graph.AddMsg(msg)

		assert.Equal(t, expectedGraph.Start, graph.Start)
		assert.Equal(t, expectedGraph.Nodes, graph.Nodes)
	})
}
