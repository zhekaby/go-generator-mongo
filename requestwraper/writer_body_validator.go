package main

var writer_body_validator = `
var {{ .Name }}Map = map[string]string{
	{{ range .Fields }}{{ if .JsonPath }}"{{ .Ns }}" :"{{ .JsonPath }}",{{ else }}"{{ .Ns }}" :"{{ .NsShort }}",{{ end }}
	{{ end }}
}
func {{ .Name }}Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body {{ .BodyTypeName }}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("cant read body"))
			return
		}
		if err := body.UnmarshalJSON(bytes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(fmt.Sprintf("invalid json: %s", err)))
			return
		}
		err = packageValidator.Struct(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			verrs := err.(validator.ValidationErrors)
			m := make(map[string][]interface{}, len(verrs))
			for _, e := range verrs {
				n := DeviceCreateRequestParamsMap[e.Namespace()]
				v, ok := m[n]
				if !ok {
					v = make([]interface{}, 0, 5)
				}
				v = append(v, printField(e))
				m[n] = v
			}
			b, _ := json.Marshal(&requestwarapper_error_model{
				Errors: m,
			})
			w.Write(b)
			return
		}
		next.ServeHTTP(w, r)
	})
}
`
