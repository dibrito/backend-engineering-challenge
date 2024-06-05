package main

import (
	"fmt"
	"os"

	"github.com/dibrito/backend-engineering-challenge/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("check ./result.txt")
	fmt.Println("DONE.")
}
