package main

var writerUpdater = `

type {{ .Typ }}Updater interface {
	{{range .Fields}}
	Set{{ .Ns }}(v{{ .Prop }} {{ .Type }}) {{ $.Typ }}Updater
	{{end}}
}

type {{ .Name }}_updater struct {
	updates bson.M
}

func New{{ .Typ }}Updater() {{ .Typ }}Updater {
	return &{{ .Name }}_updater{
		updates: bson.M{},
	}
}

func (u *{{ $.Name }}_updater) compile() bson.M {
	return bson.M{"$set": u.updates}
}

{{range .Fields}}
func (u *{{ $.Name }}_updater) Set{{ .Ns }}(v{{ .Prop }} {{ .Type }}) {{ $.Typ }}Updater {
	u.updates["{{ .JsonPath }}"] = v{{ .Prop }}
	return u
}
{{end}}
`
