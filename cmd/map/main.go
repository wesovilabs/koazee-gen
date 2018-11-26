package main

import (
	"fmt"
	"reflect"
)
//go:generate go run main.go
func main(){
	acc := false
	accPtr := &acc
	fmt.Println(reflect.TypeOf(acc))
	fmt.Println(reflect.TypeOf(accPtr).String())
}
