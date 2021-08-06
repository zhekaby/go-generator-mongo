package main

import (
	"flag"
	"fmt"
	"os"
)

var cs = flag.String("cs", "", "default connection string")
var csVar = flag.String("cs_var", "", "env var name represents connection string")
var dbVar = flag.String("db_var", "", "env var name represents db to connect, otherwise db is taken from connection string")

func main() {
	flag.Parse()

	files := flag.Args()
	if len(files) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	for _, file := range files {
		p := NewParser(file)
		if err := p.Parse(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		w := NewWriter(*cs, *csVar, *dbVar, p)

		w.Write()
	}

}
