package example

//This Code is generated DO NOT EDIT

import (
	"fmt"

	"database/sql"
	"github.com/jmoiron/sqlx"
)

type userField string
type userSearchType string

const(
	UserSearchTypeLIKE userSearchType = "LIKE"
	UserSearchTypeEQUAL userSearchType = "="
	UserSearchTypeLESSTHAN userSearchType = "<"
	UserSearchTypeGREATERTHAN userSearchType = ">"
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

func UserGetAll(db *sqlx.DB) ([]User,error){
	var result []User
	statement := fmt.Sprintf("SELECT * from %s.%s","testdb",userTableName)
	return result,db.Unsafe().Select(&result,statement)
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

func UserSaveMulti(db *sqlx.DB, in []User) (errList []error) {
	for _ , v := range in {
		errList = append(errList,UserSave(db,v))
		 }
	return
}

func UserNew(db *sqlx.DB, in User) (sql.Result,error) {
	statement := fmt.Sprintf("INSERT INTO %s.%s (%s,%s,%s,%s) VALUES (?,?,?,?)",
		"testdb",userTableName,
		UserFieldName,UserFieldEmail,UserFieldUserName,UserFieldPassword)
	return db.Exec(statement,
		in.Name,in.Email,in.UserName,in.Password)
}

func UserNewMulti(db *sqlx.DB, in []User) (errList []error) {
	for _ , v := range in {
		_ , err := UserNew(db,v)
		errList = append(errList,err)
		}
	return
}

func UserDelete(db *sqlx.DB, in User) error {
	statement := fmt.Sprintf("DELETE FROM %s.%s WHERE ?=?","testdb",userTableName)
	_,err := db.Exec(statement,
		UserFieldID, in.ID,
	)
	return err
}

func UserDeleteMulti(db *sqlx.DB, in []User) (errList []error) {
	for _ , v := range in {
		errList = append(errList,UserDelete(db,v))
	}
	return
}
