package v1

import (
	"math/big"
	"log"
	"net/rpc"
)

const m  = 128

type Node struct {
	// "host:port"
	Address string
	Id *big.Int
	Successor string
	Predecessor string
	finger [m]string
	// 作为数据时自身保存的数据，可以扩展之
	Data map[string]string
	nextFingerIndex2Fix int
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
	log.Print("node %s create cluster", n.Address)
	// 把 fingertable 的所有 entry 都指向自己
	n.Predecessor = n.addr()
	n.Successor = n.addr()
	for i := 0; i < len(n.finger); i ++ {
		n.finger[i] = n.addr()
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
	n.finger[1] = findPredecessor(bootstrap, FingerEntry(n.addr(), 1))
	n.Successor = n.finger[1]
	n.Predecessor = findPredecessor(n.Successor, nil)
	n.stabilize()

	// 似乎下面的更新都不需要, 因为 fix table 会更新的
	for i := 2; i < len(n.finger); i ++ {
		start := FingerEntry(n.addr(), i)
		// quick path, no need to update
		if Hash(n.finger[i-1]).Cmp(start) > 0 {
			n.finger[i] = n.finger[i-1]
		} else {
			n.finger[i] = findSuccessor(bootstrap, start)
		}
	}
}

/**
	id 可以为空, 为空表示查找自己的位置
	node 是连接点，借助它找到目标节点
	找到 id 的 successor
 */
func findSuccessor(node string, id *big.Int) string {
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
	err = client.Call("Node.FindSuccessor", id, &response)
	if err != nil {
		log.Print("connected to %s, but failed to find successor", node)
	}
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
	random index  > 1 into finger table
	finger[i].node = findSuccessor(finger[i].start)
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
		n.finger[n.nextFingerIndex2Fix] = response
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
func (n *Node) FindPredecessor(id *big.Int, response *string) error {
	// quick path
	if InclusiveBetween(Hash(n.Predecessor), id, n.Id) {
		response = &n.Predecessor
	} else {
		// 需要遍历 fingerTable 来查看最近的点在哪，然后递归查找
		// how to reverse fingerTable
		closestPreceding := n.closestPrecedingNode(id)
		client , err := rpc.Dial("tcp", closestPreceding)
		if err != nil {
			log.Print("failed to connect to remote node %v", closestPreceding)
			return err
		}

		err = client.Call("Node.FindPredecessor", id, response)
		if err != nil {
			log.Print("Failed to find predecessor from node %v", closestPreceding)
		}
	}
	return nil
}

/**
	这个需要考虑 currentNode 和 id 之间的关系
	Inclusive 已经考虑了过 0 的问题了
 */
func (n *Node) closestPrecedingNode(id *big.Int) string {
	for i:= len(n.finger)-1; i >= 0; i -- {
		if ExclusiveBetween(id, Hash(n.finger[i]), n.Id) {
			return n.finger[i]
		}
	}

	return ""
}

// 新节点加入时，请求自己帮忙找它的 successor
func (n *Node) FindSuccessor(id *big.Int, response *string) error {
	// 这个判断会不会使得更新暂停呢 ?
	if ExclusiveBetween(n.Id, id, Hash(n.Successor)) {
		*response = n.Successor
	} else {
		// find successor 依靠 predecessor 来完成
		var predecessor string
		err := n.FindPredecessor(id, &predecessor)

		if err != nil {
			log.Print("Failed to find predecessor from node %v", n.addr())
		}

		client, err2 := rpc.Dial("tcp", predecessor)

		if err2 != nil {
			log.Print("Failed to connec to node %v", predecessor)
		}

		err2 = client.Call("Node.FindSuccessor", id, response)
		if err2 != nil {
			log.Print("Failed to call find successor on node %v", predecessor)
		}
	}

	return nil
}

