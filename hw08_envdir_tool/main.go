package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	if len(os.Args) < 3 {
		fmt.Println("usage: envdir dir child")
		return
	}
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	RunCmd(os.Args[2:], env)
}
