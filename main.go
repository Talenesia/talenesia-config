package main

import "github.com/talenesia/router/cmd"

func main() {
	root := cmd.New()
	root.Execute()
}
