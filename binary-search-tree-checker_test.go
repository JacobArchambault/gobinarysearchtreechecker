package main

import (
	"fmt"
	"testing"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func isBST(node *Node) bool {
	if node == nil || (node.Left == nil && node.Right == nil) {
		return true
	}
	if node.Left != nil && node.Value < node.Left.Value {
		return false
	}
	if node.Right != nil && node.Value > node.Right.Value {
		return false
	}

	return isBST(node.Left) && isBST(node.Right)

}

func main() {
	result := isBST(nil) // Example usage
	fmt.Println("empty tree is BST:", result)
}

func TestEmptyNode_IsBST(t *testing.T) {
	if !isBST(nil) {
		t.Error("empty tree should be BST")
	}
}