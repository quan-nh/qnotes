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
}

var homeTmpl = template.Must(template.New("home").ParseFiles("template/base.html", "template/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	notebooks, err := getNoteBooks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = homeTmpl.ExecuteTemplate(w, "base", &Page{Title: "Home", Notebooks: notebooks})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getNoteBooks() ([]string, error) {

	fileInfos, err := ioutil.ReadDir(util.Conf.Repo)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {

		if fileInfo.IsDir() && fileInfo.Name() != ".git" {
			result = append(result, fileInfo.Name())
		}

	}

	return result, nil
}
