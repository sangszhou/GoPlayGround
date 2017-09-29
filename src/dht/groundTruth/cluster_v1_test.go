package groundTruth

import (
	"testing"
	"dht/v1"
	"log"
	"math/big"
)

func TestFingerTableStart(t *testing.T) {
	id1 := big.NewInt(0)
	log.Print(FingerTableStart(id1, 0).String())
	log.Print(FingerTableStart(id1, 1).String())
	log.Print(FingerTableStart(id1, 2).String())
}

func TestSortNodes(t *testing.T) {
	node1 := v1.CreateNode(1)
	node2 := v1.CreateNode(3)
	node3 := v1.CreateNode(5)


	var slice []*v1.Node

	slice = append(slice, node2, node3, node1)

	for _, node := range slice {
		log.Printf("[" + node.Address + "]")
	}

	SortNodes(slice)

	for _, node := range slice {
		log.Printf("[" + node.Address + "]")
	}
}

func TestCreateCluster(t *testing.T)  {
	cluster := &Cluster{}

	node1 := v1.CreateNode(3)
	node0 := v1.CreateNode(0)


	cluster.AddNode(node1)
	cluster.AddNode(node0)
	cluster.Describe()
}

/**
	如果不按照 debug 的方式来走，可能会有打印不出 log 的问题
 */
func TestFindPredecessor(t *testing.T)  {
	cluster := &Cluster{}

	node1 := v1.CreateNode(1)
	node0 := v1.CreateNode(0)

	cluster.AddNode(node1)
	cluster.AddNode(node0)

	SortNodes(cluster.Nodes)

	predecessor := cluster.Predecessor(node0.Id)
	// should be 0
	log.Printf("predecessor id: " + predecessor.Address)

	predecessor = cluster.Predecessor(node1.Id)
	// should be 1
	log.Printf("predecessor id: " + predecessor.Address)

}


func TestFindSuccessor(t *testing.T) {
	cluster := &Cluster{}

	node1 := v1.CreateNode(1)
	node0 := v1.CreateNode(0)

	cluster.AddNode(node1)
	cluster.AddNode(node0)

	SortNodes(cluster.Nodes)

	successor := cluster.Successor(big.NewInt(0).Add(node0.Id, big.NewInt(1)))
	// should be 0
	log.Print("successor id: " + successor.Address)

	successor2 := cluster.Successor(big.NewInt(0).Add(node1.Id, big.NewInt(1)))
	// should be 1
	log.Print("successor id: " + successor2.Address)

}

func TestStableCluster(t *testing.T)  {
	cluster := &Cluster{}

	node1 := v1.CreateNode(1)
	node0 := v1.CreateNode(0)


	cluster.AddNode(node1)
	cluster.AddNode(node0)

	cluster.WaitStable()

	cluster.Describe()
}

func TestStableCluster2(t *testing.T)  {
	cluster := &Cluster{}

	node1 := v1.CreateNode(1)
	node0 := v1.CreateNode(0)
	node3 := v1.CreateNode(3)


	cluster.AddNode(node1)
	cluster.AddNode(node0)
	cluster.AddNode(node3)

	cluster.WaitStable()

	cluster.Describe()
}

func TestStableCluster3(t *testing.T)  {
	cluster := &Cluster{}

	node1 := v1.CreateNode(1)
	node0 := v1.CreateNode(0)
	node3 := v1.CreateNode(3)
	node6 := v1.CreateNode(6)


	cluster.AddNode(node1)
	cluster.AddNode(node0)
	cluster.AddNode(node3)
	cluster.AddNode(node6)

	cluster.WaitStable()

	cluster.Describe()
}

func TestStableCluster4(t *testing.T)  {
	cluster := &Cluster{}

	node0 := v1.CreateNode(0)
	node3 := v1.CreateNode(3)
	node6 := v1.CreateNode(6)


	cluster.AddNode(node0)
	cluster.AddNode(node3)
	cluster.AddNode(node6)

	cluster.WaitStable()

	cluster.Describe()
}