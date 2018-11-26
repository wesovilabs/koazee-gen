package main

import (
	"fmt"
	"github.com/wesovilabs/koazee-gen/operation"
)

const basePath = "/Users/ivan/Workspace/Wesovilabs/koazee-workspace/koazee/operation"

var reduceOutputPath = fmt.Sprintf("%s/%s", basePath, "reduce")
var mapOutputPath = fmt.Sprintf("%s/%s", basePath, "maps")
var filterOutputPath = fmt.Sprintf("%s/%s", basePath, "filter")
var addOutputPath = fmt.Sprintf("%s/%s", basePath, "add")
var containsOutputPath = fmt.Sprintf("%s/%s", basePath, "contains")
var dropOutputPath = fmt.Sprintf("%s/%s", basePath, "drop")

//go:generate go run koazee-gen.go

func main() {
	fmt.Println("Generating code")

	operation.GenerateAddDispatcher(addOutputPath)
	/**
	operation.GenerateContainsDispatcher(containsOutputPath)
	operation.GenerateMapDispatcher(mapOutputPath)
	operation.GenerateFilterDispatcher(filterOutputPath)
	operation.GenerateReduceDispatcher(reduceOutputPath)
	**/

}
