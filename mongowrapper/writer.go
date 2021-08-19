package main

import (
	"fmt"
	"github.com/zhekaby/go-generator-mongo/common"
	"os"
	"text/template"
)

type writer struct {
	Cs, CsVar, DbVar string
	*common.Parser
	*common.DataView
}

func NewWriter(cs, csVar, dbVar string, p *common.Parser) *writer {
	return &writer{Cs: cs, CsVar: csVar, DbVar: dbVar, Parser: p}
}

func (w *writer) Write() error {
	if err := w.writeClient(); err != nil {
		return err
	}
	for _, c := range w.Collections {
		w.DataView = c
		if err := w.writeCollections(); err != nil {
			return err
		}
	}

	if len(w.Aggregations) > 0 {
		if err := w.writeAggregations(); err != nil {
			return err
		}
	}

	return nil
}

func (w *writer) writeAggregations() error {
	f, err := os.Create(fmt.Sprintf("%s/aggregation_funcs.go", w.Dir))
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("generating %s...", f.Name())
	fmt.Fprintf(f, "// Code generated by mongowrapper. DO NOT EDIT.\n")
	fmt.Fprintf(f, "package %s\n", w.PkgName)
	for _, str := range []string{writerAggregation} {
		tpl, err := template.New("queue").Parse(str)
		if err != nil {
			return err
		}
		err = template.Must(tpl, err).Execute(f, w)
		if err != nil {
			return err
		}
	}
	fmt.Println("done")
	return nil
}

func (w *writer) writeClient() error {
	f, err := os.Create(fmt.Sprintf("%s/client.go", w.Dir))
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("generating %s...", f.Name())
	fmt.Fprintf(f, "// Code generated by mongowrapper. DO NOT EDIT.\n")
	fmt.Fprintf(f, "package %s\n", w.PkgName)
	tpl, err := template.New("queue").Parse(writerClient)
	if err != nil {
		return err
	}
	err = template.Must(tpl, err).Execute(f, w)
	if err != nil {
		return err
	}
	fmt.Println("done")
	return nil

}

func (w *writer) writeCollections() error {
	f, err := os.Create(fmt.Sprintf("%s/repository_%s.go", w.Dir, w.DataView.Name))
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Printf("generating %s...", f.Name())
	fmt.Fprintf(f, "// Code generated by mongowrapper. DO NOT EDIT.\n")
	fmt.Fprintf(f, "package %s\n", w.PkgName)

	for _, str := range []string{writerIface, writerUpdater} {
		tpl, err := template.New("queue").Parse(str)
		if err != nil {
			return err
		}
		err = template.Must(tpl, err).Execute(f, w)
		if err != nil {
			return err
		}
	}

	fmt.Println("done")
	return nil
}
