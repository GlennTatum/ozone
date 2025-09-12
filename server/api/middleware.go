package api

import (
	"net/http"
)

func (app *App) Authz(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// query cassandra for an active session
		session, ok := app.session.Get(r.Context(), "session").(string)
		if !ok {
			return
		}
		app.logger.Debugf(session)
		// create a cassandra anon user with TIME_UUID()
		// time uuid has a 1 hour lifetime until it is deleted

		next.ServeHTTP(w, r)
	})
}
