package main

type neighbourhood struct {
	nodeCount int
	nodes     []*node
}

type node struct {
	Id       uint64
	Earnings int
	stake    int
	address  string
}
