package handler

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"qnotes/util"
)

type Page struct {
	Title        string
	Notebooks    []string
	NoteBookName string
	Notes        []string
	NoteName     string
	NoteContents []byte
	Action       string
}

var page Page
var homeTmpl = template.Must(template.New("home").ParseFiles("template/base.html", "template/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// set default value
	page.Title = "note/home"
	page.NoteBookName = ""
	page.Notes = nil
	page.NoteName = ""

	// get notebooks
	if err := getNoteBooks(&page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render template
	if err := homeTmpl.ExecuteTemplate(w, "base", &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// get notebooks in repo.
func getNoteBooks(page *Page) error {

	fileInfos, err := ioutil.ReadDir(util.Conf.Repo)
	if err != nil {
		return err
	}

	page.Notebooks = page.Notebooks[:0]
	for _, fileInfo := range fileInfos {

		if fileInfo.IsDir() && fileInfo.Name() != ".git" {
			page.Notebooks = append(page.Notebooks, fileInfo.Name())
		}

	}

	return nil
}
