package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    app.config.env,
		"version": app.config.version,
	}

	err := app.jsonResponse(w, http.StatusOK, data)
	if err != nil {
		app.errorJSON(w, http.StatusInternalServerError, err.Error())
	}
}
