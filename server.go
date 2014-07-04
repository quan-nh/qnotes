package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"qnotes/handler"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", handler.HomeHandler)
	r.HandleFunc("/n/{notebook}", handler.NotebookHandler)
	r.HandleFunc("/n/{notebook}/{note}", handler.NoteHandler)
	r.HandleFunc("/config", handler.ConfigHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)

	conf, err := handler.LoadConfig()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+conf.Port, nil)
}
