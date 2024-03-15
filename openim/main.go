package main

import (
	"flag"
	"fmt"
	"github.com/openimsdk/data-tools/openim/internal"
	"os"
)

func main() {
	var path string
	flag.StringVar(&path, "c", "config.yaml", "path config file")
	flag.Parse()
	if err := internal.Main(path); err != nil {
		fmt.Println("run error", err)
		os.Exit(1)
		return
	}
	fmt.Println("run success")
	os.Exit(0)
}
