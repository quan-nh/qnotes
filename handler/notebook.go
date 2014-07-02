package handler

import (
	"html/template"
	"net/http"
)

var notebookTmpl = template.Must(template.New("notebook").ParseFiles("template/base.html", "template/notebook.html"))

func NotebookHandler(w http.ResponseWriter, r *http.Request) {
	notebookTmpl.ExecuteTemplate(w, "base", &Page{Title: "Notebook"})
}
