package main

var writerUpdater = `

type {{ .Typ }}Updater interface {
	{{range .Fields}}Set{{ .GoPath }}(v{{ .Prop }} {{ .Type }}) {{ $.Typ }}Updater
	{{end}}
	Changes() map[string]interface{}
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

func (u *{{ $.Name }}_updater) Changes() map[string]interface{} {
	return u.updates
}

{{range .Fields}}
func (u *{{ $.Name }}_updater) Set{{ .GoPath }}(v{{ .Prop }} {{ .Type }}) {{ $.Typ }}Updater {
	u.updates["{{ .BsonPath }}"] = v{{ .Prop }}
	return u
}
{{end}}
`
