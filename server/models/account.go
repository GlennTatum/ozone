package models

import (
	"github.com/gocql/gocql"
)

func CreateAccount(db *gocql.Session, session string, resource gocql.UUID, port int) error {
	return db.Query(
		`UPDATE main.account SET resource = ?, port = ? WHERE session = ?`,
		resource, port, session,
	).Exec()
}
