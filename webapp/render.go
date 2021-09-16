package webapp

import (
	"embed"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//go:embed templates
var templates embed.FS

func render(w http.ResponseWriter, tmplName string, data interface{}) error {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"ToTitle": strings.Title,
		"ToUpper": strings.ToUpper,
	}

	tmpl, err := template.New(fmt.Sprintf("%s.tmpl", tmplName)).
		Funcs(funcMap).
		ParseFS(templates, fmt.Sprintf("templates/%s.tmpl", tmplName))
	if err != nil {
		return err
	}

	tmpl.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func renderError(w http.ResponseWriter, code int, err error) {
	data := struct {
		Code    int
		Message string
	}{
		Code:    code,
		Message: err.Error(),
	}

	render(w, "error", &data)
}

func getHtmlID() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
