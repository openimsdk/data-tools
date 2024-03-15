package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/data-tools/chat/internal"
	"os"
)

func main() {
	var (
		path string
		mode string
	)
	flag.StringVar(&path, "c", "config.yaml", "path config file")
	flag.StringVar(&mode, "m", "auto", "optional auto toc tob")
	flag.Parse()
	if err := internal.Main(path, mode); err != nil {
		fmt.Println("run error", err)
		os.Exit(1)
		return
	}
	fmt.Println("run success")
	os.Exit(0)
}
