//go:generate sqlxsdbg -t=User -db=testdb -tb=_user $GOFILE

package example

// User is a database representation of the table _user in the database testdb
type User struct {
	// ID Comment Block
	// key says it should be used as a conditional on save
	// if key is auto it also means that it is not used on new,
	// because its auto incremented for us
	ID       int     `db:"id" key:"auto"`
	Name     *string `db:"name"`
	Email    string  `db:"email"`
	UserName string  `db:"username"`
	Password string  `db:"password"`
}
