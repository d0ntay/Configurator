package main

import "log/slog"

type Application struct {
	logger *slog.Logger
}

type TemplateResponse struct {
	Templates []Template `json:"templates"`
}

type Template struct {
	Name string `json:"name"`
}

type TemplateDetail struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

type Field struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Pattern     string `json:"pattern"`
}

type GenerateRequest struct {
	ConfigType string            `json:"config_type"`
	Data       map[string]string `json:"data"`
}

type GenerateResponse struct {
	Config   string `json:"config"`
	FileName string `json:"filename"`
}
