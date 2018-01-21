package example

// User is the interface struct between the database
type User struct {
	ID       int    `db:"id"`
	Name     string `db:"name"`
	Email    string `db:"email"`
	UserName string `db:"username"`
	Password string `db:"password"`
}
