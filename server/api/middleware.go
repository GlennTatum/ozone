package api

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gocql/gocql"
)

type ValidationResponse struct {
	isUnauthorized bool
	isServe        bool
}

func (app *App) IsValidUser(u string) (string, bool, error) {
	var account string
	err := app.Db.Query(`SELECT session FROM main.account WHERE session = ?`, u).Scan(&account)
	if err == gocql.ErrNotFound {
		return "", false, nil
	}
	return account, true, nil
}

func (app *App) IsValidEvent(uri *url.URL) (string, bool, error) {
	var event string

	spl := strings.Split(uri.Path, "/")

	err := app.Db.Query(`SELECT id FROM main.event WHERE id = ?`, spl[3]).Scan(&event)
	if err == gocql.ErrNotFound {
		return "", false, nil
	}
	return event, true, nil
}

func (app *App) ValidateAuth(r *http.Request) (ValidationResponse, error) {

	event, event_ok, err := app.IsValidEvent(r.URL)
	if err != nil {
		return ValidationResponse{}, err
	}

	is_session := app.Session.Exists(r.Context(), "session")
	if !is_session {
		fmt.Println(4)
		if event_ok {
			id := rand.Text()
			app.Session.Put(r.Context(), "session", id)
			err := app.Db.Query(`INSERT INTO main.account (session) VALUES (?)`, id).Exec()
			if err != nil {
				return ValidationResponse{}, err
			}
			return ValidationResponse{isServe: true, isUnauthorized: false}, nil
		}
		fmt.Println(5)
		if !event_ok {
			return ValidationResponse{isServe: false, isUnauthorized: true}, nil
		}
	}

	s := app.Session.GetString(r.Context(), "session")
	user, user_ok, err := app.IsValidUser(s)
	if err != nil {
		return ValidationResponse{}, err
	}

	fmt.Println(user_ok, event_ok)
	fmt.Println(user, event)

	if user_ok && event_ok {
		fmt.Println(2)
		return ValidationResponse{isServe: true, isUnauthorized: false}, nil
	}
	if user_ok && !event_ok {
		fmt.Println(3)
		return ValidationResponse{isServe: false, isUnauthorized: true}, nil
	}

	return ValidationResponse{isServe: true, isUnauthorized: false}, nil
}

func (app *App) AuthzHandle(w http.ResponseWriter, r *http.Request, next http.Handler) (bool, error) {
	h, err := app.ValidateAuth(r)
	if err != nil {
		return false, err
	}
	if h.isServe {
		next.ServeHTTP(w, r)
		return true, nil
	}
	if h.isUnauthorized {
		w.WriteHeader(http.StatusUnauthorized)
		return true, nil
	}
	return false, err
}

func (app *App) Authz(next http.Handler) http.Handler {
	/*
		auth middleware controller for accessing the api
		auth flow:
		p: user has a valid session
		q: user has a valid event token
		The following checks are done in the order listed below:
		1. (!p) user has an invalid session token and (!q) enters an invalid event
			a. delete the session token
			b. return Unauthorized
		2. (p) user has a valid session token (!q) and enters an invalid event
			the user probably made a mistake on entering the url session
			a. return Unauthorized
		3. (p) user has a session token and (q) enters a valid event
			a. serve the http request
		NOTE: For 4 and 5 we first check if there is a valid event then run 4.1 if there is an invalid session and run 4.2 if there is no session token
		4.1 (!p) user has an invalid session token and (q) enters an valid event
			the user probably has a session from a previous event
			in this case we will delete the database entry of the previous
			session and create a new account
			a. delete the session from the browser with scs
			b. delete the account WHERE the session["session"] = id if it exists
			c. follow procedure for 2.
		4.2 (!p) user has no session token but (q) enters a valid event
			a. create a session for the user
			b. insert the session["session"] id into the databse
				table: account
				columns   | description
				id	      | the session id registered for the specific event
				resource  | the resource id of the object created in kubernete
			c. serve the http request
	*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ok, err := app.AuthzHandle(w, r, next)
		if err != nil {
			panic(err)
		}
		if !ok {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
}
