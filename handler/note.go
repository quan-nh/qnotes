package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"qnotes/util"
)

var noteTmpl = template.Must(template.New("note").ParseFiles("template/base.html", "template/note.html"))

func NoteHandler(w http.ResponseWriter, r *http.Request) {

	page.NoteBookName = mux.Vars(r)["notebook"]
	page.NoteName = mux.Vars(r)["note"]
	page.Title = "note/" + page.NoteName
	page.Action = r.URL.Query().Get("a")

	if page.Action == "edit" && r.FormValue("save") == "Save" {
		page.NoteContents = []byte(r.FormValue("note"))

		if err := page.save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/n/"+page.NoteBookName+"/"+page.NoteName, http.StatusFound)
	}

	if page.Action == "delete" && r.FormValue("delete") == "DELETE" {
		if err := page.delete(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/n/"+page.NoteBookName, http.StatusFound)
	}

	if r.FormValue("cancel") == "Cancel" {
		http.Redirect(w, r, "/n/"+page.NoteBookName+"/"+page.NoteName, http.StatusFound)
	}

	if page.Notebooks == nil {
		if err := getNoteBooks(&page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := getNotes(&page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new note if it doesn't exist
	if !contains(page.Notes, page.NoteName) {
		if err := page.createNote(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := loadContent(&page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := noteTmpl.ExecuteTemplate(w, "base", &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadContent(page *Page) error {
	filename := util.Conf.Repo + "/" + page.NoteBookName + "/" + page.NoteName + ext

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	page.NoteContents = content

	return nil
}

func (p *Page) createNote() error {
	filename := util.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	_, err := os.Create(filename)
	return err
}

func (p *Page) save() error {
	filename := util.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return ioutil.WriteFile(filename, p.NoteContents, 0600)
}

func (p *Page) delete() error {
	filename := util.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return os.Remove(filename)
}
