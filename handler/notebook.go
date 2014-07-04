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

	page.NoteBookName = mux.Vars(r)["notebook"]
	page.Title = "note/" + page.NoteBookName

	if err := getNoteBooks(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new notebook if it doesn't exist
	if !contains(page.Notebooks, page.NoteBookName) {
		if err := page.createNotebook(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := getNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := notebookTmpl.ExecuteTemplate(w, "base", &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNotes() error {

	fileInfos, err := ioutil.ReadDir(page.Conf.Repo + "/" + page.NoteBookName)
	if err != nil {
		return err
	}

	page.Notes = page.Notes[:0]

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			page.Notes = append(page.Notes, strings.TrimSuffix(fileInfo.Name(), ext))
		}
	}

	return nil
}

func (p *Page) createNotebook() error {
	filename := page.Conf.Repo + "/" + p.NoteBookName
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
