package handler

import (
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

var noteTmpl = template.Must(template.New("note").ParseFiles("template/base.html", "template/note.html"))

func NoteHandler(w http.ResponseWriter, r *http.Request) {

	Page.NoteBookName = mux.Vars(r)["notebook"]
	Page.NoteName = mux.Vars(r)["note"]
	Page.Title = "note/" + Page.NoteName
	Page.Action = r.URL.Query().Get("a")

	if Page.Action == "edit" && r.FormValue("save") == "Save" {
		Page.NoteContents = []byte(r.FormValue("note"))

		if err := Page.save(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/n/"+Page.NoteBookName+"/"+Page.NoteName, http.StatusFound)
	}

	if Page.Action == "delete" && r.FormValue("delete") == "DELETE" {
		if err := Page.delete(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/n/"+Page.NoteBookName, http.StatusFound)
	}

	if r.FormValue("cancel") == "Cancel" {
		http.Redirect(w, r, "/n/"+Page.NoteBookName+"/"+Page.NoteName, http.StatusFound)
	}

	if Page.Notebooks == nil {
		if err := getNoteBooks(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := getNotes(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new note if it doesn't exist
	if !contains(Page.Notes, Page.NoteName) {
		if err := Page.createNote(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := loadContent(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := noteTmpl.ExecuteTemplate(w, "base", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadContent() error {
	filename := Page.Conf.Repo + "/" + Page.NoteBookName + "/" + Page.NoteName + ext

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	Page.NoteContents = content

	return nil
}

func (p *page) createNote() error {
	filename := Page.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	_, err := os.Create(filename)
	return err
}

func (p *page) save() error {
	filename := Page.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return ioutil.WriteFile(filename, p.NoteContents, 0600)
}

func (p *page) delete() error {
	filename := Page.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return os.Remove(filename)
}
