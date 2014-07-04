package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Config struct {
	Port string
	Repo string
}

const configFile = "conf.json"

var configTmpl = template.Must(template.New("config").ParseFiles("template/base.html", "template/config.html"))

func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	// set default value
	page.Title = "note/home"
	page.NoteBookName = ""
	page.Notes = nil
	page.NoteName = ""

	// save
	if r.FormValue("save") == "Save" {
		page.Conf.Port = r.FormValue("port")
		page.Conf.Repo = r.FormValue("repo")

		if err := saveConfig(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/config", http.StatusFound)
	}

	// render template
	if err := configTmpl.ExecuteTemplate(w, "base", &page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func LoadConfig() (Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&page.Conf)
	return page.Conf, err
}

func saveConfig() error {
	b, err := json.Marshal(page.Conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFile, b, 0600)
}
