package main

import (
	"fmt"
)

func main() {
	err := Start()
	if err != nil {
		fmt.Println(err)
	}
}
