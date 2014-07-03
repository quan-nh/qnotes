package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"qnotes/util"
)

var noteTmpl = template.Must(template.New("note").ParseFiles("template/base.html", "template/note.html"))

func NoteHandler(w http.ResponseWriter, r *http.Request) {

	page.NoteBookName = mux.Vars(r)["notebook"]
	page.NoteName = mux.Vars(r)["note"]
	page.Title = "note/" + page.NoteName
	page.Action = r.URL.Query().Get("a")

	var err error

	if page.Notebooks == nil {
		err = getNoteBooks(util.Conf.Repo, &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if page.Notes == nil {
		err = getNotes(util.Conf.Repo, &page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = loadContent(util.Conf.Repo, &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = noteTmpl.ExecuteTemplate(w, "base", &page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadContent(repo string, page *Page) error {
	filename := repo + "/" + page.NoteBookName + "/" + page.NoteName + ext

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	page.NoteContents = content

	return nil
}
