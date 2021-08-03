package main

import (
	"flag"
	"fmt"
	"os"
)

var cs = flag.String("cs", "", "default connection string")
var csVar = flag.String("cs_var", "", "env var name represents connection string")

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, file := range files {
		p := NewParser(*cs, *csVar, file)
		if err := p.Parse(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		w := NewWriter(p)
		w.Write()
	}

}
