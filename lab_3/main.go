package main

import (
	"decision-theory/binrels"
	"fmt"
)

func main() {
	rels := binrels.Zero(4)

	//0011
	rels[0][2] = true
	rels[0][3] = true

	//1001
	rels[1][0] = true
	rels[1][3] = true

	//0010
	rels[2][2] = true

	//1100
	rels[3][0] = true
	rels[3][1] = true

	fmt.Println("\nInitial Relation:")
	binrels.Print(rels)

	mrRels := binrels.MutualReachability(rels)

	fmt.Println("\nMutual Reachability Relation:")
	binrels.Print(mrRels)
}
