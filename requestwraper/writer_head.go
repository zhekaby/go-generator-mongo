package main

import (
	"strings"
)

var writerHead = `import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
)

var packageValidator = validator.New()
`

var writerFooter = strings.Replace(`
func printField(e validator.FieldError) interface{} {
	switch e.Tag() {
	case "required":
		return requestwarapper_error_required
	default:
		return &struct {
		}{}
	}
}

type requestwarapper_error_model struct {
	Errors map[string][]interface{} ''json:"errors"''
}

var requestwarapper_error_required = struct {
	Key string ''json:"key"''
}{
	Key: "required",
}
`, "''", "`", -1)
