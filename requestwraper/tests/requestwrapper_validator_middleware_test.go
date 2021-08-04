package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceCreateRequestParamsValidator(t *testing.T) {
	d := &device{
		UserID: "userid",
		Locale: "ro-ro",
		MyData: &data{},
	}
	b, _ := d.MarshalJSON()
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	rr := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	DeviceCreateRequestParamsValidator(next).ServeHTTP(rr, req)
}
