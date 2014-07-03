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
	r.HandleFunc("/n/{notebook}", handler.NotebookHandler)
	r.HandleFunc("/n/{notebook}/{note}", handler.NoteHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)

	err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+util.Conf.Port, nil)
}
