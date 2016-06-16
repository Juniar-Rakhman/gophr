package main

import (
	"fmt"
	"html/template"
	"net/http"
	"bytes"
)

var layoutFuncs = template.FuncMap{
	"yield": func() (string, error) {
		return "", fmt.Errorf("yield called inappropriately")
	},
}

var layout = template.Must(
	template.
	New("layout.html").
	Funcs(layoutFuncs).
	ParseFiles("templates/layout.html"),
)

//Must: returns the template if it's valid otherwise panic.
var templates = template.Must(template.New("t").ParseGlob("templates/**/*.html"))

var errorTemplate = `
<html>
	<body>
		<h1>Error rendering template %s</h1>
		<p>%s</p>
	</body>
</html>
`

//this will override layout and replace {{yield}} with name
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {

	//TODO: replace funcs with something more meaningful
	funcs := template.FuncMap{
		"yield": func() (template.HTML, error) {
			buf := bytes.NewBuffer(nil)
			err := templates.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String()), err
		},
	}

	//2 request might come in at the same time, pointing to the same layout.
	//so use separate layout to execute
	layoutClone, _ := layout.Clone()
	layoutClone.Funcs(funcs)
	err := layoutClone.Execute(w, data)

	if err != nil {
		http.Error(
			w,
			fmt.Sprintf(errorTemplate, name, err),
			http.StatusInternalServerError,
		)
	}
}
