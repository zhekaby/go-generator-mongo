package main

import (
	"fmt"
	"os"
)

type writer struct {
	*Parser
}

func NewWriter(p *Parser) *writer {
	return &writer{p}
}

func (w *writer) Write() {
	for _, collection := range w.collections {
		write(w.Dir, w.PkgName, collection)
	}
}

func write(dir, pkgName string, coll *collection) error {
	f, err := os.Create(fmt.Sprintf("%s/%s_repository.go", dir, coll.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintf(f, "// Code generated by s4bmongo. DO NOT EDIT.\n")
	fmt.Fprintf(f, "package %s", pkgName)
	return nil
}
