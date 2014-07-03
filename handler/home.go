package handler

import (
	"html/template"
	"io/ioutil"
	"net/http"
)

type page struct {
	Title        string
	Notebooks    []string
	NoteBookName string
	Notes        []string
	NoteName     string
	NoteContents []byte
	Action       string
	Conf         config
}

var Page page
var homeTmpl = template.Must(template.New("home").ParseFiles("template/base.html", "template/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// set default value
	Page.Title = "note/home"
	Page.NoteBookName = ""
	Page.Notes = nil
	Page.NoteName = ""

	// get notebooks
	if err := getNoteBooks(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render template
	if err := homeTmpl.ExecuteTemplate(w, "base", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get notebooks in repo.
func getNoteBooks() error {

	fileInfos, err := ioutil.ReadDir(Page.Conf.Repo)
	if err != nil {
		return err
	}

	Page.Notebooks = Page.Notebooks[:0]
	for _, fileInfo := range fileInfos {

		if fileInfo.IsDir() && fileInfo.Name() != ".git" {
			Page.Notebooks = append(Page.Notebooks, fileInfo.Name())
		}

	}

	return nil
}
