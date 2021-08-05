package tests

import (
	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeviceCreateRequestParamsValidator(t *testing.T) {
	Convey(t.Name(), t, func() {
		Convey("InvalidModel", func() {
			d := &item{
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
			var resp requestwarapper_error_model
			err = json.Unmarshal(b, &resp)
			So(err, ShouldBeNil)
			So(resp.Errors, ShouldNotBeNil)
			e := resp.Errors
			So(len(e), ShouldBeGreaterThan, 0)
			So(e, ShouldContainKey, "num")
			So(e, ShouldContainKey, "type")
			So(e, ShouldContainKey, "assn")
			So(e, ShouldContainKey, "assn1")
			So(e, ShouldContainKey, "MyData.N")
		})
	})
}
