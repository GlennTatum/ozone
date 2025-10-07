package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ozone/models"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func StatusReply(w http.ResponseWriter, code int) {
	w.Write([]byte(fmt.Sprint(code)))
}

func (app *App) Home(w http.ResponseWriter, r *http.Request) error {
	b, err := json.Marshal([]byte("codelabs"))
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	return err
}

func (app *App) HomeRoute(w http.ResponseWriter, r *http.Request) {
	app.Home(w, r)
}

/*
Create a Deployment Resource based on an Event that is being currently being accessed
*/
func (app *App) Event(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	// the resource identifier for the kubernetes deployment
	resource_uuid, err := gocql.RandomUUID()
	if err != nil {
		return err
	}

	user := app.Session.GetString(r.Context(), "session")

	// the event code (will have an associated tag on docker hub ex: codeserver:ABCDEFG)
	resp, err := json.Marshal(vars["id"])
	if err != nil {
		return err
	}

	// create a new user
	models.CreateAccount(app.Db, user, resource_uuid)

	err = app.CreateResourceFromTeplate("deployment.yml", "lab", resource_uuid.String())
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) EventRoute(w http.ResponseWriter, r *http.Request) {
	err := app.Event(w, r)
	if err != nil {
		panic(err)
	}
}
