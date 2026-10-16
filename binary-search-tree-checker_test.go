package main

import (
	"testing"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func isBST(node *Node) bool {
	return isBSTInRange(node, nil, nil)
}

func isBSTInRange(node *Node, min, max *int) bool {
	if node == nil {
		return true
	}
	if min != nil && node.Value < *min {
		return false
	}
	if max != nil && node.Value > *max {
		return false
	}

	results := make(chan bool, 2)
	go func() {
		results <- isBSTInRange(node.Left, min, &node.Value)
	}()
	go func() {
		results <- isBSTInRange(node.Right, &node.Value, max)
	}()

	leftOK := <-results
	rightOK := <-results
	return leftOK && rightOK
}

func TestEmptyNode_IsBST(t *testing.T) {
	if !isBST(nil) {
		t.Error("empty tree should be BST")
	}
}

func TestSingleNode_IsBST(t *testing.T) {
	node := &Node{Value: 10}
	if !isBST(node) {
		t.Error("single node tree should be BST")
	}
}

func TestInvalidLeftNode_IsBST(t *testing.T) {
	node := &Node{Value: 10, Left: &Node{Value: 15}}
	if isBST(node) {
		t.Error("tree with left child greater than parent should not be BST")
	}
}

func TestInvalidRightNode_IsBST(t *testing.T) {
	node := &Node{Value: 10, Right: &Node{Value: 5}}
	if isBST(node) {
		t.Error("tree with right child less than parent should not be BST")
	}
}

func TestValidBST_IsBST(t *testing.T) {
	node := &Node{
		Value: 10,
		Left:  &Node{Value: 5},
		Right: &Node{Value: 15, Left: &Node{Value: 12}, Right: &Node{Value: 20}},
	}
	if !isBST(node) {
		t.Error("valid BST should return true")
	}
}

func TestInvalidBST_IsBST(t *testing.T) {
	node := &Node{
		Value: 10,
		Left:  &Node{Value: 5, Left: &Node{Value: 3}, Right: &Node{Value: 7}},
		Right: &Node{Value: 15, Left: &Node{Value: 12}, Right: &Node{Value: 8}},
	}
	if isBST(node) {
		t.Error("invalid BST should return false")
	}
}

func TestInvalidDescendant_IsBST(t *testing.T) {
	node := &Node{
		Value: 10,
		Left:  &Node{Value: 5, Right: &Node{Value: 12}},
	}

	if isBST(node) {
		t.Error("left subtree containing a value greater than the root should not be a BST")
	}
}
