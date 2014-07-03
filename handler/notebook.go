package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	ext = ".txt"
)

var notebookTmpl = template.Must(template.New("notebook").ParseFiles("template/base.html", "template/notebook.html"))

func NotebookHandler(w http.ResponseWriter, r *http.Request) {

	Page.NoteBookName = mux.Vars(r)["notebook"]
	Page.Title = "note/" + Page.NoteBookName

	if err := getNoteBooks(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new notebook if it doesn't exist
	if !contains(Page.Notebooks, Page.NoteBookName) {
		if err := Page.createNotebook(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := getNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := notebookTmpl.ExecuteTemplate(w, "base", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNotes() error {

	fileInfos, err := ioutil.ReadDir(Page.Conf.Repo + "/" + Page.NoteBookName)
	if err != nil {
		return err
	}

	Page.Notes = Page.Notes[:0]

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			Page.Notes = append(Page.Notes, strings.TrimSuffix(fileInfo.Name(), ext))
		}
	}

	return nil
}

func (p *page) createNotebook() error {
	filename := Page.Conf.Repo + "/" + p.NoteBookName
	return os.Mkdir(filename, os.ModeDir)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
