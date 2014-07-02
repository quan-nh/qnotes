package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"qnotes/handler"
	"qnotes/util"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)
	r.PathPrefix("/static").Handler(http.FileServer(http.Dir("./static/")))
	r.HandleFunc("/{notebook}", handler.NotebookHandler)
	r.HandleFunc("/{notebook}/{note}", handler.NoteHandler)
	http.Handle("/", r)

	err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+util.Conf.Port, nil)
}
