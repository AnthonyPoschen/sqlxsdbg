package example
//This Code is generated DO NOT EDIT

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type userField string
type userSearchType string

const(
	UserSearchTypeLIKE userSearchType = "LIKE"
	UserSearchTypeEQUAL userSearchType = "="
	userTableName string = "_user"

	UserFieldID userField = "id"
	UserFieldName userField = "name"
	UserFieldEmail userField = "email"
	UserFieldUserName userField = "username"
	UserFieldPassword userField = "password"
)

func UserGet( db *sqlx.DB, key userField, value string) (User, error) {
	var result User
	statement := fmt.Sprintf("SELECT * from %s.%s where %s=?", "testdb", userTableName, key)
	return result, db.Unsafe().Get(&result,statement,value)
}

func UserGetMulti( db *sqlx.DB, key userField,searchType userSearchType, value string) ([]User,error){
	var result []User
	statement := fmt.Sprintf("SELECT * from %s.%s where %s %s ?","testdb",userTableName,key,searchType)
	return result, db.Unsafe().Select(&result,statement,value)
}

func UserSave(db *sqlx.DB, in User) error {
	statement := fmt.Sprintf("UPDATE %s.%s SET ?=? ?=? ?=? ?=? WHERE ?=?", "testdb", userTableName)
	_,err := db.Exec(statement,
		UserFieldName, in.Name,
		UserFieldEmail, in.Email,
		UserFieldUserName, in.UserName,
		UserFieldPassword, in.Password,
		UserFieldID, in.ID,
	)
	return err
}

func UserNew(db *sqlx.DB, in User) error {
	statement := fmt.Sprintf("INSERT INTO %s.%s (%s,%s,%s,%s) VALUES (?,?,?,?)",
		"testdb",userTableName,
		UserFieldName,UserFieldEmail,UserFieldUserName,UserFieldPassword)
	_, err := db.Exec(statement,
		in.Name,in.Email,in.UserName,in.Password)
	return err
}
Email, in.UserName, in.Password)
	return err
}
