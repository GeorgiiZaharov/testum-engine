package user

import "time"

type CreateUserParams struct {
	Login string
	Mail  string
	Name  string
	Group *string
}

type User struct {
	ID           int       `db:"id"`
	Login        string    `db:"login"`
	Mail         string    `db:"mail"`
	Name         string    `db:"name"`
	Group        *string   `db:"group"`
	IsLecturer   bool      `db:"is_lecturer"`
	DateCreated  time.Time `db:"date_created"`
	DateModified time.Time `db:"date_modified"`
}
