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

	var err error

	if page.Action == "edit" && r.FormValue("save") == "Save" {
		page.NoteContents = []byte(r.FormValue("note"))
		err = page.save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/n/"+page.NoteBookName+"/"+page.NoteName, http.StatusFound)
	}

	if page.Action == "delete" && r.FormValue("delete") == "DELETE" {
		err = page.delete()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/n/"+page.NoteBookName, http.StatusFound)
	}

	if r.FormValue("cancel") == "Cancel" {
		http.Redirect(w, r, "/n/"+page.NoteBookName+"/"+page.NoteName, http.StatusFound)
	}

	if page.Notebooks == nil {
		err = getNoteBooks(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if page.Notes == nil {
		err = getNotes(&page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	err = loadContent(&page)
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

func loadContent(page *Page) error {
	filename := util.Conf.Repo + "/" + page.NoteBookName + "/" + page.NoteName + ext

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	page.NoteContents = content

	return nil
}

func (p *Page) save() error {
	filename := util.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return ioutil.WriteFile(filename, p.NoteContents, 0600)
}

func (p *Page) delete() error {
	filename := util.Conf.Repo + "/" + p.NoteBookName + "/" + p.NoteName + ext
	return os.Remove(filename)
}
