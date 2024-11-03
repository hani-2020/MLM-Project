package main

import (
	"fmt"
	"sync"
)

// Member represents a node in the tree.
type Member struct {
	ID     int
	Left   *Member
	Right  *Member
	Parent *Member
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
		tree.Root = &Member{ID: 1}
		tree.buildTree()
	}
	return tree
}

// buildTree creates the binary tree concurrently by building subtrees in parallel.
func (tree *Tree) buildTree() {
	var wg sync.WaitGroup
	buildSubtree(tree.Root, 2, tree.NumMembers, &wg)
	wg.Wait() // Wait for all goroutines to finish
}

// buildSubtree recursively builds the left and right subtrees in parallel if possible.
func buildSubtree(parent *Member, currentID, maxID int, wg *sync.WaitGroup) {
	if currentID > maxID {
		return
	}

	// Left child creation
	left := &Member{ID: currentID, Parent: parent}
	parent.Left = left
	currentID++

	// Right child creation if within bounds
	var right *Member
	if currentID <= maxID {
		right = &Member{ID: currentID, Parent: parent}
		parent.Right = right
		currentID++
	}

	// Spawn goroutines to create left and right subtrees if there are remaining nodes
	if currentID <= maxID {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buildSubtree(left, currentID, maxID, wg)
		}()
	}
	if currentID <= maxID && right != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			buildSubtree(right, currentID+1, maxID, wg)
		}()
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
	numMembers := 100 // Replace with desired number of members
	tree := NewTree(numMembers)
	tree.DisplayTree()
}
