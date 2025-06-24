package main

import (
	"academy-todo/internal/app/web"
	"fmt"
)

func main() {
	err := web.Start()
	if err != nil {
		fmt.Println(err)
	}
}
