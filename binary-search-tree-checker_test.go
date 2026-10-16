package main

import (
	"runtime"
	"testing"
)

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

func isBST(node *Node) bool {
	workers := runtime.GOMAXPROCS(0)
	if workers < 2 {
		return isBSTParallel(node, nil, nil, nil)
	}

	// Slots represent extra goroutines; the caller is the remaining worker.
	slots := make(chan struct{}, workers-1)
	return isBSTParallel(node, nil, nil, slots)
}

func isBSTParallel(node *Node, min, max *int, slots chan struct{}) bool {
	if node == nil {
		return true
	}
	if min != nil && node.Value < *min {
		return false
	}
	if max != nil && node.Value > *max {
		return false
	}

	left, right := node.Left, node.Right
	if left == nil {
		return isBSTParallel(right, &node.Value, max, slots)
	}
	if right == nil {
		return isBSTParallel(left, min, &node.Value, slots)
	}

	if slots == nil {
		return isBSTParallel(left, min, &node.Value, nil) &&
			isBSTParallel(right, &node.Value, max, nil)
	}

	select {
	case slots <- struct{}{}:
		leftResult := make(chan bool, 1)
		go func() {
			defer func() { <-slots }()
			leftResult <- isBSTParallel(left, min, &node.Value, slots)
		}()

		rightOK := isBSTParallel(right, &node.Value, max, slots)
		leftOK := <-leftResult
		return leftOK && rightOK
	default:
		// No worker slot available: check both sides in this goroutine.
		return isBSTParallel(left, min, &node.Value, slots) &&
			isBSTParallel(right, &node.Value, max, slots)
	}
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
