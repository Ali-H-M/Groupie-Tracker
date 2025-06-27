package funcs

import (
	"html/template"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r) // 404
		return
	}

	if r.Method == http.MethodGet {
		RenderTemplate(w, "index.html")
		return
	}
}

func RenderTemplate(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("templates/" + tmpl)
	if err != nil {
		http.Error(w, "Page not found", http.StatusNotFound) // 404
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError) // Template execution fails 500
	}
}
