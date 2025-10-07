package models

import (
	"github.com/gocql/gocql"
)

func CreateAccount(db *gocql.Session, session string, resource gocql.UUID) error {
	return db.Query(
		`UPDATE main.account SET resource = ? WHERE session = ?`,
		resource, session,
	).Exec()

}
