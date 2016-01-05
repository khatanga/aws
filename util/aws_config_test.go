package util

import (
	"testing"
	"fmt"
)

func TestInitialize(*testing.T) {
	Initialize()
	fmt.Println("Testing...")
	fmt.Println(Config)
}