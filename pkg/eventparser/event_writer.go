package eventparser

import (
	"fmt"
	"io"
	"text/template"
	"time"
	"unicode"
)

var NowFunc = time.Now

func WriteEvents(out io.Writer, events []ParsedEvent, pkg, id, aggregateType string) error {
	if len(events) < 1 {
		return fmt.Errorf("no events found")
	}

	containsPersonalData := false
	containsSerializedData := false
	for _, event := range events {
		if event.PersonalData != nil {
			if event.PersonalData.SerializedFields != nil {
				containsSerializedData = true
			}
			containsPersonalData = true
		}
	}

	return fileTemplate.Execute(out, templateData{
		Timestamp:              NowFunc(),
		Events:                 events,
		AggregateType:          aggregateType,
		ID:                     id,
		Package:                pkg,
		ContainsPersonalData:   containsPersonalData,
		ContainsSerializedData: containsSerializedData,
	})
}

type templateData struct {
	Timestamp              time.Time
	Events                 []ParsedEvent
	AggregateType          string
	ID                     string
	Package                string
	ContainsPersonalData   bool
	ContainsSerializedData bool
}

func lcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

var funcMap = template.FuncMap{
	"LcFirst": lcFirst,
}

var fileTemplate = template.Must(template.New("").Funcs(funcMap).Parse(`// Code generated by go generate; DO NOT EDIT.
// This file was generated at
// {{ .Timestamp }}
package {{ $.Package }}

{{- if .ContainsPersonalData }}

import (
{{- if .ContainsSerializedData }}
	"strconv"
{{ end }}
	"github.com/inklabs/rangedb/pkg/crypto"
)
{{- end }}

{{- range .Events }}

func (e {{ .Name }}) AggregateID() string   { return e.{{ $.ID }} }
func (e {{ .Name }}) AggregateType() string { return "{{ $.AggregateType }}" }
func (e {{ .Name }}) EventType() string     { return "{{ .Name }}" }

{{- if .PersonalData }}
{{- $event := . }}
func (e *{{ .Name }}) Encrypt(encryptor crypto.Encryptor) error {
{{- range $event.PersonalData.Fields }}
	{{ . | LcFirst }}, err := encryptor.Encrypt(e.{{ $event.PersonalData.SubjectID }}, e.{{ . }})
	if err != nil {
		if err == crypto.ErrKeyWasDeleted {
			e.redactPersonalData("")
		}
		return err
	}
{{ end }}
{{- range $key, $value := $event.PersonalData.SerializedFields }}
	string{{ $value }} := strconv.Itoa(e.{{ $value }})
	{{ $key | LcFirst }}, err := encryptor.Encrypt(e.{{ $event.PersonalData.SubjectID }}, string{{ $value }})
	if err != nil {
		if err == crypto.ErrKeyWasDeleted {
			e.redactPersonalData("")
		}
		return err
	}
{{ end }}
	{{- range $event.PersonalData.Fields }}
	e.{{ . }} = {{ . | LcFirst }}
	{{- end }}
	{{- range $k, $v := $event.PersonalData.SerializedFields }}
	e.{{ $v }} = 0
	{{- end }}
	{{- range $k, $v := $event.PersonalData.SerializedFields }}
	e.{{ $k }} = {{ $k | LcFirst }}
	{{- end }}
	return nil
}
func (e *{{ .Name }}) Decrypt(encryptor crypto.Encryptor) error {
{{- range $event.PersonalData.Fields }}
	{{ . | LcFirst }}, err := encryptor.Decrypt(e.{{ $event.PersonalData.SubjectID }}, e.{{ . }})
	if err != nil {
		if err == crypto.ErrKeyWasDeleted {
			e.redactPersonalData("")
		}
		return err
	}
{{ end }}
{{- range $key, $value := $event.PersonalData.SerializedFields }}
	decrypted{{ $value }}, err := encryptor.Decrypt(e.{{ $event.PersonalData.SubjectID }}, e.{{ $key }})
	if err != nil {
		if err == crypto.ErrKeyWasDeleted {
			e.redactPersonalData("")
		}
		return err
	}
	{{ $value | LcFirst }}, err := strconv.Atoi(decrypted{{ $value }})
	if err != nil {
		return err
	}
{{ end }}
	{{- range $event.PersonalData.Fields }}
	e.{{ . }} = {{ . | LcFirst }}
	{{- end }}
	{{- range $k, $v := $event.PersonalData.SerializedFields }}
	e.{{ $v }} = {{ $v | LcFirst }}
	{{- end }}
	{{- range $k, $v := $event.PersonalData.SerializedFields }}
	e.{{ $k }} = ""
	{{- end }}
	return nil
}
func (e *{{ .Name }}) redactPersonalData(redactTo string) {
	{{- range $event.PersonalData.Fields }}
	e.{{ . }} = redactTo
	{{- end }}
	{{- range $k, $v := $event.PersonalData.SerializedFields }}
	e.{{ $v }} = 0
	{{- end }}
}
{{- end }}
{{- end }}
`))
