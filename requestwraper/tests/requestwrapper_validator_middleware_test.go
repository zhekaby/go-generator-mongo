package tests

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceCreateRequestParamsValidator(t *testing.T) {
	Convey(t.Name(), t, func() {
		Convey("InvalidModel", func() {
			d := &device{
				UserID: "userid",
				Locale: "ro-ro",
				MyData: &data{},
			}
			b, _ := d.MarshalJSON()
			req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
			rr := httptest.NewRecorder()
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			deviceCreateRequestParamsValidator(next).ServeHTTP(rr, req)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			b, err := ioutil.ReadAll(rr.Body)
			So(err, ShouldBeNil)
			t.Logf("response: %s", string(b))
		})
	})
}
