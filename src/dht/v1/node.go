package v1

import (
	"math/big"
	"log"
	"net/rpc"
)

const m  = 3

type Node struct {
	// "host:port"
	Address     string
	Id          *big.Int
	Successor   string
	Predecessor string
	Finger      [m]string
	// 作为数据时自身保存的数据，可以扩展之
	Data map[string]string
	// 用来持续的 fix data table
	nextFingerIndex2Fix int
}

/**
	deep copy part of original node structure
 */
func copyNode(origin *Node) *Node {
	return &Node{
		Address:origin.Address,
		Id:origin.Id,
		Successor:origin.Successor,
		Predecessor:origin.Predecessor,
	}
}


type PutArg struct {
	Key, Value string
}

func (n *Node) addr() string {
	return n.Address
}

// 创建 ring
// 注意不要出现递归查找的情况
// @todo make node identifier a method, instead type
func (n *Node) create() {
	log.Printf("node %s create cluster", n.Address)
	// 把 fingertable 的所有 entry 都指向自己
	n.Predecessor = n.addr()
	n.Successor = n.addr()
	for i := 0; i < len(n.Finger); i ++ {
		n.Finger[i] = n.addr()
	}
}

// 加入一个已有的环中
func (n *Node) join(bootstrap string)  {
	// 初始化，只能把它设置成  "", 而不能是 nil
	//n.initFingerTable(bootstrap)
	n.Predecessor = ""
	n.Successor = findSuccessor(bootstrap, n.Id)
}

/**
	this init is no logger needed in multinode leave and join serinario
	@deprecated
 */
func (n *Node) initFingerTable(bootstrap string)  {
	n.Finger[1] = findPredecessor(bootstrap, FingerEntry(n.addr(), 1))
	n.Successor = n.Finger[1]
	n.Predecessor = findPredecessor(n.Successor, nil)
	n.stabilize()

	// 似乎下面的更新都不需要, 因为 fix table 会更新的
	for i := 2; i < len(n.Finger); i ++ {
		start := FingerEntry(n.addr(), i)
		// quick path, no need to update
		if Hash(n.Finger[i-1]).Cmp(start) > 0 {
			n.Finger[i] = n.Finger[i-1]
		} else {
			n.Finger[i] = findSuccessor(bootstrap, start)
		}
	}
}

/**
	id 可以为空, 为空表示查找自己的位置
	node 是连接点，借助它找到目标节点
	找到 id 的 successor
 */
func findSuccessor(helper string, id *big.Int) string {
	if id == nil {
		id = Hash(helper)
	}

	client, err := rpc.Dial("tcp", helper)
	if err != nil {
		log.Print("Failed to connect to bootstrap %s", helper)
	}

	defer client.Close()
	var response string
	// 在 predecessor 上查找 successor 应该一下就能找到
	err = client.Call("Node.FindSuccessor", id, &response)
	if err != nil {
		log.Print("connected to %s, but failed to find successor", helper)
	}

	// 找到 predecessor.successor



	return response
}

func findPredecessor(node string, id *big.Int) string {
	if id == nil {
		id = Hash(node)
	}
	client, err := rpc.Dial("tcp", node)
	if err != nil {
		log.Print("Failed to connect to bootstrap %s", node)
	}

	defer client.Close()
	var response string
	// 在 predecessor 上查找 successor 应该一下就能找到
	err = client.Call("Node.FindPredecessor", id, &response)
	if err != nil {
		log.Print("connected to %s, but failed to find successor", node)
	}
	return response
}


// tell successor that i am your need predecessor
// called periodically
// x = successor.predecessor
// if ( x belong to (currentNode, successor)
//     successor = x
// successor.notify(n)
func (n *Node) stabilize()  {
	successorCandidate := findPredecessor(n.Successor, nil)

	// 如果自己是仅有的节点，那么不需要 stabilize
	if(successorCandidate == n.Address) {
		return
	} else if successorCandidate != "" {
		if ExclusiveBetween(n.Id, Hash(successorCandidate), Hash(n.Successor)) {
			log.Print("successor changed from %v to %v for node %v",
				n.Predecessor, successorCandidate, n.addr())
				n.Successor = successorCandidate
				// notify new predecessor that its successor has changed
				client, err := rpc.Dial("tcp", n.Successor)
				if err != nil {
					log.Print("Failed to connect node %v", n.Predecessor)
				}


				//@todo 常用的做法是什么
				err = client.Call("Node.Notify", n.addr(), nil)
		}
	}
}

/*
	random index  > 1 into Finger table
	Finger[i].node = findSuccessor(Finger[i].start)
 */
func (n *Node)fixFingerTable()  {
	// 最小是 1
	if n.nextFingerIndex2Fix >= m {
		n.nextFingerIndex2Fix = 1
	}

	start := FingerEntry(n.Address, n.nextFingerIndex2Fix)

	// 是不是要进行 force 一下呢
	response := findSuccessor(n.Successor, start)
	if response != "" {
		n.Finger[n.nextFingerIndex2Fix] = response
	}

}

// public methods, message receiver
// predecessor 告诉当前节点的 successor 应该更新了
func (n*Node) Notify(predecessorCandidate string, response *interface{}) error {

	// 再次确认 successor 离自己更近
	if n.Predecessor == "" || ExclusiveBetween(Hash(n.Predecessor), Hash(predecessorCandidate), n.Id) {
		log.Print("node %v update predecessor from %v to %v", n.addr(), n.Predecessor, predecessorCandidate)
		n.Predecessor = predecessorCandidate
	} else {
		log.Print("node %v, don't need to update predecessor", n.addr())
	}
	*response = nil
	return nil
}

// 需要递归的遍历，直到找到某个节点的 id' 和其 successor' 满足 (id' < id < successor')
// 只有当前节点的 successor 和 predecessor 可信，以此来找到节点
// 实际上， predecessor 和 successor 两个方法，只要有一个可用就行了，另外一个可以依靠这个方法实现
// FindSuccessor = FindPredecessor.successor
func (n *Node) FindPredecessor(id *big.Int, predecessor *Node) error {
	// quick path
	// 此时，节点 N 还没有把 id 放到自己的 Finger table 中，所以自己的 successor 不会是 id 对应的节点
	if InclusiveBetween(n.Id, id, Hash(n.Successor)) {
		//deep copy struct
		predecessor = copyNode(n)
	} else {
		// 需要遍历 fingerTable 来查看最近的点在哪，然后递归查找
		// how to reverse fingerTable
		closestPreceding := n.closestPrecedingNode(id)
		//*** 判断此节点是不是就是自身
		if closestPreceding == n.Address {
			log.Print("it seems there is only one node in cluster, set it to you predecessor")
			predecessor = copyNode(n)
			return nil
		} else {
			client , err := rpc.Dial("tcp", closestPreceding)
			if err != nil {
				log.Print("failed to connect to remote node %v", closestPreceding)
				return err
			}

			err = client.Call("Node.FindPredecessor", id, predecessor)
			if err != nil {
				log.Print("Failed to find predecessor from node %v", closestPreceding)
			}
		}
	}
	return nil
}

/**
	这个需要考虑 currentNode 和 id 之间的关系
	Inclusive 已经考虑了过 0 的问题了
	这个节点找的是 predecessor
 */
func (n *Node) closestPrecedingNode(id *big.Int) string {
	for i:= len(n.Finger)-1; i >= 0; i -- {
		if ExclusiveBetween(id, Hash(n.Finger[i]), n.Id) {
			return n.Finger[i]
		}
	}

	return n.Address
}

// 新节点加入时，请求自己帮忙找它的 successor
func (n *Node) FindSuccessor(id *big.Int, successor *Node) error {

	// find successor 依靠 predecessor 来完成
	var predecessor *Node
	err := n.FindPredecessor(id, predecessor)

	if err != nil {
		log.Print("Failed to find predecessor from node %v", n.addr())
	}

	successor = n.identify(predecessor.Successor)
	return nil
}

/**
	返回自身的地址
	this code is not cool
 */
func (n *Node) identify(addr string) *Node {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		log.Print("node: %v failed to connect to addr: %v ", n.Address, addr)
		return nil
	}
	var targetNode *Node
	err = client.Call("Node.Identify", addr, targetNode)

	if err != nil {
		log.Print("failed to identify %v", addr)
	}

	return targetNode
}

func (n *Node) Identify(addr string, node *Node) error {
	if addr == n.Address {
		node = copyNode(n)
		return nil
	} else {
		log.Print("the addr is not me")

	}

	return nil
}


