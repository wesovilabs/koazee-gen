package main

import (
	"fmt"
	"github.com/wesovilabs/koazee-gen/internal/reduce"
)

const outputPath = "/Users/ivan/Workspace/Wesovilabs/koazee-workspace/koazee/internal/reduce"

//go:generate go run main.go

func main() {
	fmt.Println("Generating code")
	reduce.GenerateDispatcher(outputPath)
	//reduce.GenerateDispatcher(outputPath)
}
