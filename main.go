package main

import (
	"fmt"

	"github.com/n-samar/fnttw/fnttw"
)

func main() {
	fmt.Println("Hello world")
	fmt.Println(fnttw.NttFriendlyPrimes(1<<20, 1, 28))
}
