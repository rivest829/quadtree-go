package quadtree

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNewQuadtree(t *testing.T) {
	mouseX := float64(rand.Intn(640))
	mouseY := float64(rand.Intn(480))
	//
	topRect := Rect{0, 0, 640, 480}
	myTree := NewQuadtree(topRect, 0)
	myObj := Rect{200, 100, 35, 70}
	myTree.Insert(myObj)
	myCursor := Rect{mouseX, mouseY, 20, 20}
	candidates := myTree.Retrieve(myCursor)
	fmt.Println(candidates)
	myTree.Clear()
}
