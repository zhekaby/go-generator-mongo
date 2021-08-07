package main

import (
	"flag"
	"fmt"
	"github.com/zhekaby/go-generator-mongo/common"
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

	if *cs == "" && *csVar == "" {
		fmt.Println("--cs or --cs_var required")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Generating db repositories...")
	fmt.Printf("default connection string: %s\n", *cs)
	fmt.Printf("cs_var default: %s\n", *csVar)
	fmt.Printf("env var for db name: %s\n", *dbVar)

	for _, file := range files {
		p := common.NewParser(file)
		if err := p.Parse(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		w := NewWriter(*cs, *csVar, *dbVar, p)

		if err := w.Write(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}
