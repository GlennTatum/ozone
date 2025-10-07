package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ozone/models"

	"math/rand/v2" // Use math/rand/v2 for newer Go versions

	"github.com/gocql/gocql"
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

type EventResponse struct {
	PortAssignment int
}

/*
Create a Deployment Resource based on an Event that is being currently being accessed
*/
func (app *App) Event(w http.ResponseWriter, r *http.Request) error {

	// the resource identifier for the kubernetes deployment
	resource_uuid, err := gocql.RandomUUID()
	if err != nil {
		return err
	}

	user := app.Session.GetString(r.Context(), "session")

	min := 31000
	max := 32000

	port := rand.IntN(max-min) + min

	// create a new user
	models.CreateAccount(app.Db, user, resource_uuid, port)

	err = app.CreateResourceFromTeplate("deployment.yml", "code-server", resource_uuid.String(), port)
	if err != nil {
		return err
	}

	err = app.CreateResourceFromTeplate("service.yml", "code-server", resource_uuid.String(), port)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(EventResponse{PortAssignment: port})
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
