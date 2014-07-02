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
	r.HandleFunc("/{notebook}", handler.NotebookHandler)
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", r)

	err := util.LoadConfig()
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":"+util.Conf.Port, nil)
}
