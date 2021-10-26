package main

import (
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	bytes := make([]byte, 255)
	for i := 0; i < len(bytes); i++ {
		bytes[i] = byte(i)
	}
	//printBytes(bytes)
	for i := 0; i < len(bytes); i++ {
		tmp := bytes[i]
		if tmp < 200 && tmp > 0 {
			bytes[i] = tmp + 1
		}

	}
	printBytes(bytes)
	for i := 0; i < len(bytes); i++ {
		tmp := bytes[i]
		if tmp < 200+1 && tmp > 0+1 {
			bytes[i] = tmp - 1
		}

	}
	//printBytes(bytes)
}

func printBytes(bytes []byte) {
	for i := 0; i < len(bytes); i++ {
		fmt.Printf("%v\n", bytes[i])
	}
}
