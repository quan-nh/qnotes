package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"qnotes/util"
)

var notebookTmpl = template.Must(template.New("notebook").ParseFiles("template/base.html", "template/notebook.html"))

func NotebookHandler(w http.ResponseWriter, r *http.Request) {

	page.NoteBookName = mux.Vars(r)["notebook"]
	page.Title = "note/" + page.NoteBookName

	var err error

	if len(page.Notebooks) == 0 {
		err = getNoteBooks(util.Conf.Repo, &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = getNotes(util.Conf.Repo, &page)
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

func getNotes(repo string, page *Page) error {

	fileInfos, err := ioutil.ReadDir(repo + "/" + page.NoteBookName)
	if err != nil {
		return err
	}

	page.Notes = page.Notes[:0]

	for _, fileInfo := range fileInfos {

		if !fileInfo.IsDir() {
			page.Notes = append(page.Notes, fileInfo.Name())
		}

	}

	return nil
}
