package web

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func RegisterRoutes(r *http.ServeMux) {
	// web
	r.HandleFunc("GET /{$}", HandleStartPage)
	// files
	r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

}

// HandleStartPage handles the request to get the start page.
func HandleStartPage(w http.ResponseWriter, r *http.Request) {
	templates := template.Must(template.ParseFiles(
		filepath.Join("web/templates", "index.html"),
	))
	tmpl := "index.html"
	err := templates.ExecuteTemplate(w, tmpl, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
