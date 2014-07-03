package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type config struct {
	Port string
	Repo string
}

const configFile = "conf.json"

var configTmpl = template.Must(template.New("config").ParseFiles("template/base.html", "template/config.html"))

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	// set default value
	Page.Title = "note/home"
	Page.NoteBookName = ""
	Page.Notes = nil
	Page.NoteName = ""

	// save
	if r.FormValue("save") == "Save" {
		Page.Conf.Port = r.FormValue("port")
		Page.Conf.Repo = r.FormValue("repo")

		if err := saveConfig(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/config", http.StatusFound)
	}

	// render template
	if err := configTmpl.ExecuteTemplate(w, "base", &Page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoadConfig() error {
	file, err := os.Open(configFile)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	return decoder.Decode(&Page.Conf)
}

func saveConfig() error {
	b, err := json.Marshal(Page.Conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, b, 0600)
}
