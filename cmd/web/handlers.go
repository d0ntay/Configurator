package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi/v5"
)

func (app *Application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://127.0.0.1:5000/v1/templates")
	if err != nil {
		app.logger.Error("failed to fetch templates", "error", err)
		http.Error(w, "could not fetch templates", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	app.logger.Info("template service response", "status", resp.StatusCode)

	var data TemplateResponse

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		app.logger.Error("failed to decode template response", "error", err)
		http.Error(w, "failed template request", http.StatusInternalServerError)
		return
	}

	app.logger.Info("templates loaded", "count", len(data.Templates))

	ts, err := template.ParseFiles(
		"./ui/html/base.tmpl",
		"./ui/html/pages/home.tmpl",
	)

	if err != nil {
		app.logger.Error("failed to parse templates", "error", err)
		http.Error(w, "failed template parse", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.logger.Error("failed to execute template", "error", err)
		http.Error(w, "failed template execution", http.StatusInternalServerError)
		return
	}
}

func (app *Application) configFormPageHandler(w http.ResponseWriter, r *http.Request) {
	configType := chi.URLParam(r, "name")
	if configType == "" {
		http.Error(w, "missing config type", http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(map[string]string{"config_type": configType})
	if err != nil {
		http.Error(w, "failed to encode request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://127.0.0.1:5000/v1/getTemplate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, "failed to contact template service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data TemplateDetail
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		http.Error(w, "failed to parse template response", http.StatusInternalServerError)
		return
	}

	data.Name = configType

	ts, err := template.ParseFiles(
		"./ui/html/base.tmpl",
		"./ui/html/pages/form.tmpl",
	)
	if err != nil {
		http.Error(w, "failed to parse templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		fmt.Println("template exec error:", err)
	}
}

func (app *Application) configGeneratorHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.logger.Error("failed to parse form", "error", err)
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	data := make(map[string]string)
	for key := range r.PostForm {
		if key != "config_type" {
			data[key] = r.FormValue(key)
		}
	}

	reqBody := GenerateRequest{
		ConfigType: r.FormValue("config_type"),
		Data:       data,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		app.logger.Error("failed to encode request", "error", err)
		http.Error(w, "failed to encode request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://127.0.0.1:5000/v1/render", "application/json", bytes.NewBuffer(body))
	if err != nil {
		app.logger.Error("failed to contact config generator service", "error", err)
		http.Error(w, "failed to contact config generator service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var genResp GenerateResponse
	err = json.NewDecoder(resp.Body).Decode(&genResp)
	if err != nil {
		app.logger.Error("failed to parse generator response", "error", err)
		http.Error(w, "failed to parse generator response", http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles(
		"./ui/html/base.tmpl",
		"./ui/html/pages/result.tmpl",
	)
	if err != nil {
		app.logger.Error("failed to parse templates", "error", err)
		http.Error(w, "failed to parse templates", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", genResp)
	if err != nil {
		app.logger.Error("failed to execute template", "error", err)
		http.Error(w, "failed to execute template", http.StatusInternalServerError)
		return
	}
}

func (app *Application) downloadHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	config := r.FormValue("config")
	filename := r.FormValue("filename")

	if filename == "" {
		filename = "download.txt"
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Type", "text/plain")

	w.Write([]byte(config))
}
