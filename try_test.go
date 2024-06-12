package main

import (
	"fmt"
	"testing"
)

func TestTry(t *testing.T) {
	/* err := helpers.ServerInits()
	if err != nil {
		t.Fatal(err)
	} */
	var x float64
	fmt.Sscanf("434000", "%f", &x)

	fmt.Println(x)
}
