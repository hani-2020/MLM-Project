package main

import (
	"fmt"
)

// Member represents a node in the tree.
type Member struct {
	ID       int
	Left     *Member
	Right    *Member
	Parent   *Member
}

// Tree represents the binary tree.
type Tree struct {
	Root       *Member
	NumMembers int
}

// NewTree initializes and constructs a binary tree with the specified number of members.
func NewTree(numMembers int) *Tree {
	tree := &Tree{NumMembers: numMembers}
	if numMembers > 0 {
		tree.buildTree()
	}
	return tree
}

// buildTree creates a binary tree with NumMembers nodes.
func (tree *Tree) buildTree() {
	tree.Root = &Member{ID: 1}
	queue := []*Member{tree.Root}
	currentID := 2

	// Level-order insertion
	for currentID <= tree.NumMembers {
		current := queue[0]
		queue = queue[1:]

		// Create left child if there's room
		if currentID <= tree.NumMembers {
			left := &Member{ID: currentID, Parent: current}
			current.Left = left
			queue = append(queue, left)
			currentID++
		}

		// Create right child if there's room
		if currentID <= tree.NumMembers {
			right := &Member{ID: currentID, Parent: current}
			current.Right = right
			queue = append(queue, right)
			currentID++
		}
	}
}

// DisplayTree outputs the tree in level order, showing each node's ID, left, right, and parent.
func (tree *Tree) DisplayTree() {
	if tree.Root == nil {
		fmt.Println("Tree is empty.")
		return
	}
	queue := []*Member{tree.Root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		leftID := -1
		rightID := -1
		parentID := -1
		if current.Left != nil {
			leftID = current.Left.ID
			queue = append(queue, current.Left)
		}
		if current.Right != nil {
			rightID = current.Right.ID
			queue = append(queue, current.Right)
		}
		if current.Parent != nil {
			parentID = current.Parent.ID
		}

		fmt.Printf("Member ID: %d, Left: %d, Right: %d, Parent: %d\n", current.ID, leftID, rightID, parentID)
	}
}

func main() {
	// Example usage
	numMembers := 1000000 // Replace with desired number of members
	tree := NewTree(numMembers)
	tree.DisplayTree()
}
