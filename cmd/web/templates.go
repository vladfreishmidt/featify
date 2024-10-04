package main

import "github.com/vladfreishmidt/featify/internal/models"

type templateData struct {
	Project  *models.Project
	Projects []*models.Project
}
