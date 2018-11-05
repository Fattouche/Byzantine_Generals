package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// The flag to specify the number of generals
var G = flag.Int("G", 4, "The number of generals. First is commander, rest are lieutenants")

// The flag to specify how many faulty generals there are
var M = flag.Int("M", 1, "The number of faulty generals")

// The flag to specify the commany that will be given by the commander(general 0)
var O = flag.String("O", "ATTACK", "The order by the general(ATTACK or RETREAT)")

type Node struct {
	inputValue  string
	outputValue string
	process_ids map[int]int
	children    []*Node
}

type General struct {
	id   int
	node *Node
}

func genFaultyIndexes(numFaultyGenerals, numGenerals int) map[int]int {
	faultyIndexes := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(numGenerals)
	for _, r := range p[:numFaultyGenerals] {
		faultyIndexes[r] = 1
	}
	return faultyIndexes
}

func initializeGenerals(numGenerals int, order string) []General {
	generals := make([]General, numGenerals)
	for i := 0; i < numGenerals; i++ {
		generals[i] = General{i, nil}
	}
}

func oppositeOrder(order string) string {
	if order == "RETREAT" {
		return "ATTACK"
	} else {
		return "RETREAT"
	}
}

func (general General) broadcastOrder(generals []General, faultyGenerals map[int]int) {
	faultyGeneral := false
	if _, ok := faultyGenerals[0]; ok {
		faultyGeneral = true
	}
	for i := 0; i < len(generals); i++ {
		if generals[i].node == nil {
			generals[i].node = new(Node)
		}
		if generals[i].node.process_ids[general.id] == 1 {
			continue
		}
		if i%2 == 0 && faultyGeneral {
			generals[i].node.inputValue = oppositeOrder(general.node.inputValue)
		} else {
			generals[i].node.inputValue = general.node.inputValue
		}
		generals[i].node.process_ids[general.id] = 1
		general.node.children = append(general.node.children, generals[i].node)
	}
}

func main() {
	flag.Parse()
	var numGenerals int = *G
	var numFaultyGenerals int = *M
	var order string = *O
	if 3*numFaultyGenerals > numGenerals {
		fmt.Println("Too many faulty generals, can only have a third of the total generals being faulty")
		return
	}
	//Randomly choose which generals are faulty
	faultyIndexes := genFaultyIndexes(numFaultyGenerals, numGenerals)
	//Initialize the generals
	generals := initializeGenerals(numGenerals, order)
	generals[0].node = &Node{order, "", map[int]int{0: 1}, nil}
	//Using the commander, initialize the broadcast to all other generals
	generals[0].broadcastOrder(generals, faultyIndexes)

	for i := 1; i <= numFaultyGenerals; i++ {
		for j := 1; j < numGenerals; j++ {
			node := generals[j].node.children
			generals[j].broadcastOrder(generals)
		}
	}
}
