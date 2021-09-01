package main

var writer_body_validator = `
var requestwarapper_{{ .Name }}Map = map[string]string{
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
				n := requestwarapper_{{ $.Name }}Map[e.Namespace()]
				v, ok := m[n]
				if !ok {
					v = make([]interface{}, 0, 5)
				}
				if e.Tag() == "required" {
					v = append(v, &requestwarapper_error_key_default{
						Key: e.Tag(),
					})
					m[n] = v
					continue
				}
				switch e.Namespace() {
				{{ range .Fields }}{{ $f := . }}{{if .Validations }}
				case "{{ .Ns }}":
					switch e.Tag() {
					{{ range $key, $value := .Validations }}
						case "{{ $key }}": v = append(v, requestwarapper_error_{{ $.Name }}_{{ $f.NsCompact }}_{{$key}}){{ end}}
					}
				{{ end}}{{ end }}
				}
				m[n] = v
			}
			b, _ := json.Marshal(&requestwarapper_error_model{
				Errors: m,
			})
			w.Write(b)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "body", &body)))
	})
}
`
