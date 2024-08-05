package main

import (
	"github.com/talenesia/router/cmd"
	"github.com/talenesia/router/config"
)

func main() {
	config.Load()
	root := cmd.New()
	root.Execute()
}
