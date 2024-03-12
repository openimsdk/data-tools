package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/data-tools/chat/internal"
	"os"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "", "path config file")
	flag.Parse()
	if err := internal.Main(path); err != nil {
		fmt.Println("run error", err)
		return
	}
	fmt.Println("run success")
	os.Exit(0)
}
