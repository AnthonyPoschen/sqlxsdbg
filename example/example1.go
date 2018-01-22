//go:generate sqlxsdbg -t=User $GOFILE

package example

/*
	User is a interface struct between the database,
	the generator looks for the below keywords somewhere in the document to help identify
	meta information needed to generate automatice db entries.

	databaseName:"testdb"
	tableName:"_user"
*/
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
