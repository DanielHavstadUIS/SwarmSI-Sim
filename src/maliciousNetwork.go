package main

import (
	"fmt"
	"math/rand"
)

// freezing params
//uint256 public penaltyMultiplierDisagreement = 1; is already 1
//   uint256 public penaltyMultiplierNonRevealed = 2;

// The length of a round in blocks.
//uint256 public roundLength = 152; abstract to be 1

//call to freeze
//Stakes.freezeDeposit(
//	currentReveals[revIndex].overlay,                 //number of nodes in neighbourhood // only this matters
//	penaltyMultiplierDisagreement * roundLength * uint256(2 ** truthRevealedDepth)
//);

type FrozenStatus struct {
	isFrozen bool
	duration int
}

type FixedIdealSwarmNetworkMalicious struct {
	// Amount of nodes needs to be divisible by four
	networkNodeCount  int
	neighbourhoods    []neighbourhood
	nodeAddressMap    map[uint64]*node
	nodes             []*node
	stakeDistribution StakeCreator

	// gonna go with an initial scheme where 1 is a correct proof, and any other number is a unique malicious reveal
	revealMap map[*node]int
	// need to later do some logic so a node can be unfrozen
	// might need slashable stake
	frozenMap map[*node]*FrozenStatus
}

func (sn *FixedIdealSwarmNetworkMalicious) CreateSwarmNetwork() {
	stakes := make([]int, 0, NODECOUNT)
	rand.Seed(int64(STAKESEED))
	for i := 0; i < sn.networkNodeCount; i++ {
		stakes = append(stakes, sn.stakeDistribution.GetStake(i))
	}

	// To generate the same network
	rand.Seed(SETUPSEED)
	sn.nodeAddressMap = make(map[uint64]*node)

	buckets := make([]*neighbourhood, sn.networkNodeCount/4)

	// newhood := neighbourhood{}
	for i := 0; i < sn.networkNodeCount; i++ {
		nn := node{
			Id:       uint64(i),
			Earnings: 0,
			stake:    stakes[i],
		}

		h := rand.Intn(len(buckets))
		if buckets[h] == nil {
			buckets[h] = &neighbourhood{nodes: make([]*node, 0, 4)}
		}
		buckets[h].nodes = append(buckets[h].nodes, &nn)

		sn.nodeAddressMap[nn.Id] = &nn
		sn.nodes = append(sn.nodes, &nn)

		if len(buckets[h].nodes) == 4 {
			buckets[h].nodeCount = len(buckets[h].nodes)
			sn.neighbourhoods = append(sn.neighbourhoods, *buckets[h])
			buckets = append(buckets[:h], buckets[h+1:]...)
		}
		// Initially set malicious attributes to default
		fmt.Println(&nn)
		sn.revealMap[&nn] = 0
		sn.frozenMap[&nn] = &FrozenStatus{false, 0}
	}
}

func (sn *FixedIdealSwarmNetworkMalicious) SelectNeighbourhood() *neighbourhood {
	// Anchor is selected at random. Here it is assumed the chance is 1/len(neighbourhoods)
	ind := rand.Intn(len(sn.neighbourhoods))
	return &sn.neighbourhoods[ind]

}

func (sn *FixedIdealSwarmNetworkMalicious) selectUnfrozenNodes(nbhood *neighbourhood) []*node {
	// find eligible unfrozen nodes
	unfrozenNodes := make([]*node, 0)
	for i := 0; i < nbhood.nodeCount; i++ {
		if sn.frozenMap[nbhood.nodes[i]].isFrozen == false {
			unfrozenNodes = append(unfrozenNodes, nbhood.nodes[i])
		}
	}
	return unfrozenNodes
}

func (sn *FixedIdealSwarmNetworkMalicious) unfreezeNodes() {
	for i := 0; i < len(sn.nodes); i++ {
		currentStatus := sn.frozenMap[sn.nodes[i]]

		if currentStatus.isFrozen == true {
			currentStatus.duration = currentStatus.duration - 1
			if currentStatus.duration == 0 {
				currentStatus.isFrozen = false
			}
		}
	}
}

func (sn *FixedIdealSwarmNetworkMalicious) freezeDeposit(node *node, duration int) {
	currentStatus := sn.frozenMap[node]
	currentStatus.isFrozen = true
	currentStatus.duration = duration

}

func (sn *FixedIdealSwarmNetworkMalicious) SelectWinner() *node {
	nbhood := sn.SelectNeighbourhood()
	fmt.Println(nbhood.nodeCount)
	sn.unfreezeNodes()
	//init reveals
	//select one neighbourhood node to be malicous
	for i := 0; i < nbhood.nodeCount; i++ {
		//if not already initialized
		if sn.revealMap[nbhood.nodes[i]] == 0 {
			if i == 0 {
				sn.revealMap[nbhood.nodes[0]] = 2
			} else {
				sn.revealMap[nbhood.nodes[i]] = 1

			}
		}
	}

	// find eligible unfrozen nodes
	unfrozenNodes := sn.selectUnfrozenNodes(nbhood)

	// select truth
	//fmt.Println(len(unfrozenNodes))
	truthCursor := rand.Intn(len(unfrozenNodes))
	//fmt.Println(truthCursor)
	truth := sn.revealMap[unfrozenNodes[truthCursor]]

	//freeze those not following truth

	for i := 0; i < len(unfrozenNodes); i++ {
		if sn.revealMap[unfrozenNodes[i]] != truth {
			sn.freezeDeposit(unfrozenNodes[i], nbhood.nodeCount)
		}
	}

	// find eligible unfrozen nodes again
	unfrozenNodes = sn.selectUnfrozenNodes(nbhood)

	fmt.Println(len(unfrozenNodes))

	// It's weigthed by the stake of the nodes.
	weigthSum := 0
	for i := 0; i < len(unfrozenNodes); i++ {
		weigthSum += unfrozenNodes[i].stake
	}

	//fmt.Println(weigthSum)
	num := rand.Intn(weigthSum)

	// Should always return a winner.
	// Since num should be less than total
	// weighted sum.
	for i := 0; i < len(unfrozenNodes); i++ {
		num -= unfrozenNodes[i].stake
		if num <= 0 {
			return unfrozenNodes[i]
		}
	}

	// If it gets here, something is wrong
	panic("Found no winning node")
}
func (sn *FixedIdealSwarmNetworkMalicious) UpdateNetwork() {
	// Fixed network, no change
}

// Creates an array of nodes at their current state.
// Used for storing nodes data at each round.
func (sn *FixedIdealSwarmNetworkMalicious) GetNodeArray() *[]node {
	nodes := make([]node, 0, len(sn.nodes))
	for _, v := range sn.nodes {
		nodes = append(nodes, *v)
	}
	return &nodes
}
func (sn *FixedIdealSwarmNetworkMalicious) GetNodeAdressMap() map[uint64]*node {
	return sn.nodeAddressMap
}

type KademSwarmTreeStorageDepthMalicious struct {
	addressLength     int
	nodeCount         int
	stakeDistribution StakeCreator
	kademTree         bintree
	fullySaturate     bool
	storageDepth      int

	// addressBook is mapping nodeID to the node
	addressBook map[uint64]*node // TODO: is it needed for this imp=?
	// Need extra book for kademlia address
	kademAddress map[string]*node
	nodes        []*node

	// gonna go with an initial scheme where 1 is a correct proof, and any other number is a unique malicious reveal
	revealMap map[*node]int
	// need to later do some logic so a node can be unfrozen
	// might need slashable stake
	frozenMap map[*node]bool
}

func (kdst *KademSwarmTreeStorageDepthMalicious) CreateSwarmNetwork() {
	stakes := make([]int, 0, NODECOUNT)
	rand.Seed(int64(STAKESEED))
	for i := 0; i < NODECOUNT; i++ {
		stakes = append(stakes, kdst.stakeDistribution.GetStake(i))
	}

	// To generate the same network
	rand.Seed(SETUPSEED)
	for i := 0; i < kdst.nodeCount; i++ {
		// Create node
		n := &node{Id: uint64(i), stake: stakes[i]}

		//Create Kademlia address.
		nAdd := ""
		if kdst.fullySaturate {
			str := "%0"
			str += fmt.Sprintf("%db", kdst.addressLength)

			nAdd = fmt.Sprintf(str, i)
		} else {
			nAdd = randomBitString(kdst.addressLength)
		}
		for j := 0; j < 100; j++ {
			//for { // Avoid infinite loop
			if _, ok := kdst.kademAddress[nAdd]; !ok {
				break
			}
			nAdd = randomBitString(kdst.addressLength)
		}
		n.address = nAdd

		// Add node to data structures
		kdst.nodes = append(kdst.nodes, n)
		kdst.addressBook[uint64(i)] = n // TODO: might not be used, del.

		kdst.kademAddress[nAdd] = n
		kdst.kademTree.InsertNode(n, n.address)
		// Initially set malicious attributes to default
		kdst.revealMap[n] = 0
		kdst.frozenMap[n] = false
	}
}
func (kdst KademSwarmTreeStorageDepthMalicious) UpdateNetwork() {
	return
}

func (kdst *KademSwarmTreeStorageDepthMalicious) SelectNeighbourhood() *neighbourhood {
	anch := randomBitString(kdst.addressLength)
	nodes := kdst.kademTree.navigateWithStop(anch, kdst.storageDepth).allNodeBelowArr
	nei := neighbourhood{nodes: nodes, nodeCount: len(nodes)}
	return &nei
}

func (kdst KademSwarmTreeStorageDepthMalicious) SelectWinner() *node {
	nbhood := kdst.SelectNeighbourhood()
	fmt.Println(nbhood.nodeCount)
	//init reveals
	//select one neighbourhood node to be malicous
	for i := 0; i < nbhood.nodeCount; i++ {
		//if not already initialized
		if kdst.revealMap[nbhood.nodes[i]] == 0 {
			if i == 0 {
				kdst.revealMap[nbhood.nodes[0]] = 2
			} else {
				kdst.revealMap[nbhood.nodes[i]] = 1

			}
		}
	}

	// find eligible unfrozen nodes
	unfrozenNodes := make([]*node, 0)

	for i := 0; i < nbhood.nodeCount; i++ {
		if kdst.frozenMap[nbhood.nodes[i]] == false {
			unfrozenNodes = append(unfrozenNodes, nbhood.nodes[i])
		}
	}

	// select truth
	//fmt.Println(len(unfrozenNodes))
	truthCursor := rand.Intn(len(unfrozenNodes))
	//fmt.Println(truthCursor)
	truth := kdst.revealMap[unfrozenNodes[truthCursor]]

	//freeze those not following truth

	for i := 0; i < len(unfrozenNodes); i++ {
		if kdst.revealMap[unfrozenNodes[i]] != truth {
			kdst.frozenMap[unfrozenNodes[i]] = true
		}
	}

	//comment original code in case i need it
	// It's weigthed by the stake of the nodes.
	//weigthSum := 0
	//for i := 0; i < nbhood.nodeCount; i++ {
	//	weigthSum += nbhood.nodes[i].stake
	//	}
	// Should always return a winner.
	// Since num should be less than total
	// weighted sum.
	// for i := 0; i < nbhood.nodeCount; i++ {
	// 	num -= nbhood.nodes[i].stake
	// 	if num <= 0 {
	// 		return nbhood.nodes[i]
	// 	}
	// }

	// find eligible unfrozen nodes again
	unfrozenNodes = make([]*node, 0)

	for i := 0; i < nbhood.nodeCount; i++ {
		if kdst.frozenMap[nbhood.nodes[i]] == false {
			unfrozenNodes = append(unfrozenNodes, nbhood.nodes[i])
		}
	}

	fmt.Println(len(unfrozenNodes))

	// It's weigthed by the stake of the nodes.
	weigthSum := 0
	for i := 0; i < len(unfrozenNodes); i++ {
		weigthSum += unfrozenNodes[i].stake
	}

	//fmt.Println(weigthSum)
	num := rand.Intn(weigthSum)

	// Should always return a winner.
	// Since num should be less than total
	// weighted sum.
	for i := 0; i < len(unfrozenNodes); i++ {
		num -= unfrozenNodes[i].stake
		if num <= 0 {
			return unfrozenNodes[i]
		}
	}

	// If it gets here, something is wrong
	panic("Found no winning node")
}

func (kdst *KademSwarmTreeStorageDepthMalicious) GetNodeAdressMap() map[uint64]*node {
	return kdst.addressBook
}
func (kdst *KademSwarmTreeStorageDepthMalicious) GetNodeArray() *[]node {
	nodes := make([]node, 0, len(kdst.nodes))
	for _, v := range kdst.nodes {
		nodes = append(nodes, *v)
	}
	return &nodes
}
