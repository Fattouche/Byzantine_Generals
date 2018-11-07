package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

const (
	ATTACK  = "ATTACK"
	RETREAT = "RETREAT"
)

// The flag to specify the number of generals
var G = flag.Int("G", 7, "The number of generals. First is commander, rest are lieutenants")

// The flag to specify how many faulty generals there are
var M = flag.Int("M", 2, "The number of faulty generals")

// The flag to specify the commany that will be given by the commander(general 0)
var O = flag.String("O", "ATTACK", "The order by the general(ATTACK or RETREAT)")

// Values:
//	inputValue(string): The order the node has been given
//	outputValue(string): The order the node decides on
//	processIDs(map[int]int): Keeps track of which nodes have seen this message before
//	children([]*Node): Keeps track of the children of the node
//	id(int): The id of the corresponding general
type Node struct {
	inputValue  string
	outputValue string
	processIDs  map[int]int
	children    []*Node
	id          int
}

// Parameters:
//		numFaultyGenerals(int): The number of faulty generals
//		numGenerals(int): The total number of generals

// This function will randomly select numFaultyGenerals generals to be the faulty ones.
// The commander can also be chosen.
func genFaultyIndexes(numFaultyGenerals, numGenerals int) map[int]int {
	faultyIndexes := make(map[int]int)
	rand.Seed(time.Now().UnixNano())
	p := rand.Perm(numGenerals)
	for _, r := range p[:numFaultyGenerals] {
		faultyIndexes[r] = 1
	}
	return faultyIndexes
}

// Paramters:
//		order(string): The order that will be give

// This function returns the opposite of the given order.
func oppositeOrder(order string) string {
	if order == RETREAT {
		return ATTACK
	} else {
		return RETREAT
	}
}

// Parameters:
//		id(int): The id of the node you want to message
//		faultyGenerals(map[int]int): A map of all the generals that are known to be faulty.

// This function will simulate sending a message by updating the information of the recieving node.
// It will then add the recieivng node as one of its children which will be used later in the decision.
func (node Node) sendMessage(id int, faultyGenerals map[int]int) *Node {
	var order string
	if _, ok := node.processIDs[id]; ok {
		return nil
	}
	faultyGeneral := false
	if _, ok := faultyGenerals[node.id]; ok {
		faultyGeneral = true
	}
	if id%2 == 0 && faultyGeneral {
		order = oppositeOrder(node.inputValue)
	} else {
		order = node.inputValue
	}
	processIDs := make(map[int]int)
	for k, v := range node.processIDs {
		processIDs[k] = v
	}
	processIDs[id] = 1
	return &Node{order, "", processIDs, nil, id}
}

// This is responsible for deciding the decision of each of the nodes within our created tree.
func (node *Node) decide() string {
	// Base case: If this is a leaf node, then the output value is the input value
	if len(node.children) == 0 {
		node.outputValue = node.inputValue
		return node.outputValue
	}
	// If it isn't a leaf node it means we have children that we need to use to decide.
	decisions := make(map[string]int)
	for _, child := range node.children {
		decision := child.decide()
		if _, ok := decisions[decision]; ok {
			decisions[decision]++
		} else {
			decisions[decision] = 1
		}
	}
	// Get the majority vote from this nodes children and assign it to the output value.
	node.outputValue = majorityDecision(decisions)
	return node.outputValue
}

// Parameters:
//		decisions(map[string]int): A map of the count of each ATTACK or RETREAT decisions

// This helper function is used to return the majority decision among a nodes children.
func majorityDecision(decisions map[string]int) string {
	if decisions[ATTACK] > decisions[RETREAT] {
		return ATTACK
	} else {
		return RETREAT
	}
}

// Parameters:
//		numGenerals(int): The number of generals(including commander)
//		numFaultyGenerals(int): The number of faulty generals
//		order(string): The order that the commander will give to his generals

// This is the main orchestrator for the byzantine generals algorithm
func byzantineGenerals(numGenerals, numFaultyGenerals int, order string) {
	//Randomly choose which generals are faulty(can include the commander)
	faultyIndexes := genFaultyIndexes(numFaultyGenerals, numGenerals)

	// Initialize the commander with the original order and itself as a visited commander in the map.
	// This will ensure the commander is never messaged again.
	commander := &Node{order, "", map[int]int{0: 1}, nil, 0}
	if faultyIndexes[0] == 1 {
		fmt.Println("Commander is faulty")
	}
	queue := []*Node{commander}

	currDepth := 0
	elemDepth := 1
	nextElemDepth := 0

	// This is a depth-limited breadth first search so that we sendMessages and create nodes for all generals,
	// only up to the depth defined as the number of faulty generals.
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		nextElemDepth += numGenerals
		elemDepth--
		if elemDepth == 0 {
			currDepth++
			if currDepth > numFaultyGenerals {
				break
			}
			elemDepth = nextElemDepth
			nextElemDepth = 0
		}
		// Iterate through each general and send the message, add to end of BFS queue.
		for i := 1; i < numGenerals; i++ {
			childNode := node.sendMessage(i, faultyIndexes)
			if childNode == nil {
				continue
			}
			node.children = append(node.children, childNode)
			queue = append(queue, childNode)
		}
	}
	// We have built the tree, now we can decide what to the consensus decision is.
	commander.outputValue = commander.decide()

	for i, general := range commander.children {
		if _, ok := faultyIndexes[i+1]; ok {
			fmt.Print("Faulty ")
		}
		fmt.Printf("general %d decides on %s\n", general.id, general.outputValue)
	}
	fmt.Printf("Consensus decision of %s", commander.outputValue)
}

func main() {
	flag.Parse()
	// Parse flags in to better variable names
	var numGenerals int = *G
	var numFaultyGenerals int = *M
	var order string = *O
	// Make sure that the number of faulty generals is less than a third of the number of generals
	if 3*numFaultyGenerals > numGenerals {
		fmt.Println("Too many faulty generals, can only have a third of the total generals being faulty")
		return
	}
	// Call the orchestrator
	byzantineGenerals(numGenerals, numFaultyGenerals, order)
}
