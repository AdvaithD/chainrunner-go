package main

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/graph/simple"
)

func main() {
	fmt.Println("Starting graph")
	graph := simple.NewWeightedDirectedGraph(0, math.Inf(1))
}
