package main

import "strings"

var writerFooter = strings.Replace(`
type requestwarapper_error_model struct {Errors map[string][]interface{} ''json:"errors"''}
type requestwarapper_error_key_default struct {	Key string ''json:"key"''}
`, "''", "`", -1)
