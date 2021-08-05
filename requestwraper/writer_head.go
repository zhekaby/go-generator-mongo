package main

var writerHead = `import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
)

var packageValidator = validator.New()
`
