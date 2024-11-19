package config

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type Config struct {
	Routes map[string]string
}

var config Config

func LoadConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return err
	}

	return nil
}

func RenderConfigPage(writer http.ResponseWriter) {
	tmpl, err := template.New("configPage").ParseFiles("templates/configPage.html")

	if err != nil {
		http.Error(writer, "Failed to load config page template", http.StatusInternalServerError)
		return
	}

	config := GetConfig()

	err = tmpl.ExecuteTemplate(writer, "configPage.html", config)
	if err != nil {
		http.Error(writer, "Failed to render error page", http.StatusInternalServerError)
	}
}

func GetConfig() Config {
	return config
}
