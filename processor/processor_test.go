package processor

import (
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	dummyLogger = log.New(io.Discard, "", 0)
)

//nolint:gocognit
func TestProcessor(t *testing.T) {
	t.Run("it creates correct processor with start and end", func(t *testing.T) {
		processor := NewProcessor(dummyLogger)
		expectedStartNeighbours := 0

		assert.Len(t, processor.Start.Neigbours, expectedStartNeighbours)
	})

	t.Run("it adds one word message to processor", func(t *testing.T) {
		processor := NewProcessor(dummyLogger)

		expectedStartNeighbours := 1
		expectedWord := "aboba"
		expectedNode := NewNode(expectedWord, 1, processor.End)
		expectedMap := map[string]*Node{
			expectedWord: expectedNode,
		}

		processor.AddMsg(expectedWord)

		assert.Len(t, processor.Start.Neigbours, expectedStartNeighbours)
		assert.Equal(t, expectedNode, processor.Start.Neigbours[expectedWord])
		assert.Equal(t, expectedMap, processor.Nodes)
	})

	t.Run("it creates correct processor with a provided message", func(t *testing.T) {
		processor := NewProcessor(dummyLogger)

		msg := "hello, are you"
		processor.AddMsg(msg)

		node1 := NewNode("hello,", 1)
		start := NewNode("", 1, node1)
		node2 := NewNode("are", 1)
		node3 := NewNode("you", 1, processor.End)
		node1.AddNeighbour(node2)
		node2.AddNeighbour(node3)

		expctedMap := map[string]*Node{
			"hello,": node1,
			"are":    node2,
			"you":    node3,
		}
		expectedprocessor := &Processor{
			Start: start,
			Nodes: expctedMap,
		}

		assert.Equal(t, expectedprocessor.Start, processor.Start)
		assert.Equal(t, expectedprocessor.Nodes, processor.Nodes)
	})

	t.Run("it creates correct processor with 2 colliding messages", func(t *testing.T) {
		processor := NewProcessor(dummyLogger)

		msg1 := "hello, are you"
		processor.AddMsg(msg1)
		msg2 := "privet, are you"
		processor.AddMsg(msg2)

		node1 := NewNode("hello,", 1)
		start := NewNode("", 1, node1)
		node2 := NewNode("are", 2)
		node3 := NewNode("you", 2, processor.End)
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
		expectedprocessor := &Processor{
			Start: start,
			Nodes: expctedMap,
		}

		assert.Equal(t, expectedprocessor.Start, processor.Start)
		assert.Equal(t, expectedprocessor.Nodes, processor.Nodes)
	})

	t.Run("it generates new messages correctly", func(t *testing.T) {
		const (
			helloHowAreYou = "hello, how are you?"
			hiHowAreYou    = "hi, how are you?"

			numberOfGenerations = 10_000
			allowedDelta        = 2
		)

		tCases := []struct {
			Name                   string
			Msgs                   []string
			ExpectedMsgPercentages map[string]float64
		}{
			{
				Name: "one message",
				Msgs: []string{helloHowAreYou},
				ExpectedMsgPercentages: map[string]float64{
					helloHowAreYou: 100,
				},
			},
			{
				Name: "two messages",
				Msgs: []string{helloHowAreYou, hiHowAreYou},
				ExpectedMsgPercentages: map[string]float64{
					helloHowAreYou: 50,
					hiHowAreYou:    50,
				},
			},
			{
				Name: "three messages",
				Msgs: []string{helloHowAreYou, helloHowAreYou, hiHowAreYou},
				ExpectedMsgPercentages: map[string]float64{
					helloHowAreYou: 66,
					hiHowAreYou:    33,
				},
			},
		}

		for _, tCase := range tCases {
			t.Run(tCase.Name, func(t *testing.T) {
				processor := NewProcessor(dummyLogger)

				for _, msg := range tCase.Msgs {
					processor.AddMsg(msg)
				}

				gotMsgCalls := map[string]int{}
				for range numberOfGenerations {
					newMsg := processor.Generate()
					gotMsgCalls[newMsg]++
				}

				totalNumberOfCalls := 0
				for _, calls := range gotMsgCalls {
					totalNumberOfCalls += calls
				}

				gotMsgPercentages := map[string]float64{}
				for msg, calls := range gotMsgCalls {
					gotMsgPercentages[msg] = (float64(calls) / float64(totalNumberOfCalls)) * 100
				}

				for msg, percentage := range tCase.ExpectedMsgPercentages {
					assert.InDelta(t, percentage, gotMsgPercentages[msg], allowedDelta)
				}
			})
		}
	})
}
