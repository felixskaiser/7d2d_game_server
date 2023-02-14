package serverManager

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving '/'")
	defaultHandlerRespond(w)
}

func defaultHandlerRespond(w http.ResponseWriter) {
	layoutTmplPath := filepath.Join(templateDir, "base.tmpl")
	defaultTmplPath := filepath.Join(templateDir, "default.tmpl")

	tmpl, err := template.ParseFiles(layoutTmplPath, defaultTmplPath)
	if err != nil {
		log.Printf("DEFAULT: Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("DEFAULT: Error serving '/': %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}
