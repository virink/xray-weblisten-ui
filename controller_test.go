package main

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestPortInUse(t *testing.T) {
	p := 0
	i := 0
	for {
		i++
		p = 30000 + rand.Intn(10000)
		fmt.Println("Test port", p)
		if !portInUse(p) {
			break
		}
		if i > 10 {
			break
		}
	}
	fmt.Println(p)
}
