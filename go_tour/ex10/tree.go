package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

func Walk(t *tree.Tree, ch chan int) {
	walkRecursive(t, ch)
	close(ch)
}

func walkRecursive(t *tree.Tree, ch chan int) {
	if t.Left != nil {
		walkRecursive(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		walkRecursive(t.Right, ch)
	}
}

func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go Walk(t1, ch1)
	go Walk(t2, ch2)
	for i := range ch1 {
		if i != <-ch2 {
			return false
		}
	}
	return true
}

func main() {
	ch := make(chan int)
	go Walk(tree.New(1), ch)
	for i := range ch {
		fmt.Println(i)
	}
	fmt.Println("Test completed if TRUE:", Same(tree.New(1), tree.New(1)))
	fmt.Println("Test completed if FALSE:", Same(tree.New(1), tree.New(2)))
}
