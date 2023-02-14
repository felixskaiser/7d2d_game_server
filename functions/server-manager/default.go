package serverManager

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving default page")

	layoutTmplPath := filepath.Join(templateDir, "base.tmpl")
	defaultTmplPath := filepath.Join(templateDir, "default.tmpl")

	tmpl, err := template.ParseFiles(layoutTmplPath, defaultTmplPath)
	if err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}

	tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Printf("Error serving server status: %v", err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	return
}
