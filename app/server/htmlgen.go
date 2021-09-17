package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func writeHtml(w http.ResponseWriter, viewModel interface{}, filenames ...string) (err error) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	err = templates.ExecuteTemplate(w, "layout", viewModel)
	if err != nil {
		log.Printf("execute template failed: %v", err)
		fmt.Fprintf(w, "%v", err)
	}
	return
}
