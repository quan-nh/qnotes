package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"qnotes/util"
	"strings"
)

const (
	ext = ".txt"
)

var notebookTmpl = template.Must(template.New("notebook").ParseFiles("template/base.html", "template/notebook.html"))

func NotebookHandler(w http.ResponseWriter, r *http.Request) {

	page.NoteBookName = mux.Vars(r)["notebook"]
	page.Title = "note/" + page.NoteBookName

	var err error

	if page.Notebooks == nil {
		err = getNoteBooks(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = getNotes(&page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = notebookTmpl.ExecuteTemplate(w, "base", &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNotes(page *Page) error {

	fileInfos, err := ioutil.ReadDir(util.Conf.Repo + "/" + page.NoteBookName)
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
