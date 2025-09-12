package api

import (
	"net/http"
)

func (app *App) healthz(w http.ResponseWriter, r *http.Request) {
	// app.db.Query()
}
