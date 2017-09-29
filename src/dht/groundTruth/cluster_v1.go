package groundTruth

import (
	"dht/v1"
	"math/big"
	"sort"
	"fmt"
)

type Cluster struct {
	Nodes []*v1.Node
}

const keySize = 3
var hashMod = new(big.Int).Exp(big.NewInt(2), big.NewInt(keySize), nil)

func SortNodes(nodes []*v1.Node) {
	// and the node is sorted
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id.Cmp(nodes[j].Id) < 0
	})
}


func ReverseSortNodes(nodes []*v1.Node)  {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Id.Cmp(nodes[j].Id) > 0
	})
}

func (cluster *Cluster)AddNode(newNode *v1.Node)  {
	// append slice using append
	cluster.Nodes = append(cluster.Nodes, newNode)
	SortNodes(cluster.Nodes)
}

func (cluster *Cluster) WaitStable()  {


	for _, node := range cluster.Nodes {
		// iterate node finger table
		for i := 0; i < len(node.Finger); i ++ {
			start := FingerTableStart(node.Id, i)
			successor := cluster.Successor(start)
			if successor == nil {
				// only one node in cluster
				node.Finger[i] = node.Address
			} else {
				node.Finger[i] = successor.Address
			}
		}

		//successor and predecessor
		node.Successor = cluster.Successor(big.NewInt(0).Add(node.Id, big.NewInt(1))).Address
		node.Predecessor = cluster.Predecessor(node.Id).Address
	}
}



func FingerTableStart(id *big.Int, entryIdx int) *big.Int {
	two := big.NewInt(2)
	exponent := big.NewInt(int64(entryIdx))
	two.Exp(two, exponent, hashMod)
	newId := big.NewInt(0).Add(id, two)
	return newId.Mod(newId, hashMod)
}


/**
	ensure sorted
 */
func (cluster *Cluster)Successor(id *big.Int) *v1.Node {

	for _, node := range cluster.Nodes {
		if node.Id.Cmp(id) >= 0 {
			return node
		}
	}

	return cluster.Nodes[0]
}

/**
	ensure sorted
 */
func (cluster * Cluster)Predecessor(id *big.Int) *v1.Node {
	for i := len(cluster.Nodes)-1; i >= 0; i -- {
		if cluster.Nodes[i].Id.Cmp(id) < 0 {
			return cluster.Nodes[i]
		}
	}

	return cluster.Nodes[len(cluster.Nodes)-1]
}

func (cluster *Cluster) Describe()  {
	for _, node := range cluster.Nodes {
		DescribeNode(node)
		fmt.Print("\n\n")
	}
}

func DescribeNode(node *v1.Node)  {
	nodeName := fmt.Sprintf("Node addr: %v\n", node.Address)
	nodeId := fmt.Sprintf("Node id: %v\n", node.Id.String())
	fmt.Print(nodeName)
	fmt.Print(nodeId)
	fmt.Println("successor: " + node.Successor)
	fmt.Println("predecessor: " + node.Predecessor)
	for idx, tableEntry := range node.Finger {
		entry := fmt.Sprintf("\t finger[%v] = start(%v): %v\n",
			idx, FingerTableStart(node.Id, idx), tableEntry)
		fmt.Print(entry)
	}
}

