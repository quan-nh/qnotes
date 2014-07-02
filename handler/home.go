package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"qnotes/util"
)

type Page struct {
	Title        string
	Notebooks    []string
	NoteBookName string
	Notes        []string
	NoteName     string
}

var homeTmpl = template.Must(template.New("home").ParseFiles("template/base.html", "template/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	getNoteBooks()
	err := homeTmpl.ExecuteTemplate(w, "base", &Page{Title: "Home"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNoteBooks() {
	fmt.Println(util.Conf.Repo)
}
