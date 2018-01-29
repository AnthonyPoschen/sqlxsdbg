# sqlxsdbg

sqlxsdbg is a sqlx simple database generator which generates package level functions that will handle getting and saving the struct to and from a mysql db based from a pre written struct. 

## Installation

```
go install github.com/zanven42/sqlxsdbg
```

## Usage

```go
//go:generate sqlxsdbg -t=structName -db=DBName -tb=TableName $GOFILE
```
place the above in a file that requires generation with appropiate variable names and run
```
go generate
```
in a terminal in the same directory to have a file generated with the following package level functions
* FooGet
* FooGetMulti
* FooSave
* FooSaveMulti
* FooNew

### Recognised Tags
* "db" - the column name in your database
* "key" - if this is a key for your table, used in fetching data as the where clause 
	* "auto" - if key has this value instead of nothing, it is treated as an auto incremented variable and ignored for the New function
## Example

The following example is a simple representation of what is required for this package to work. Please check the example folder for a more detailed example.

file: foo.go
```go
package dbfoo
//go:generate sqlxsdbg -t=Bar -db=foo -tb=bar $GOFILE

// Bar description
type Bar struct {
	ID string `db:"id" key:"auto"`
	Name *string `db:"name"`
	OtherTypes int `db:"othertypes"`
	AreSupported *bool `db:"aresupported"`
}
```

file: foo_gen.go
```go
package dbfoo
//This Code is generated DO NOT EDIT

...

BarGet(...)

BarGetMulti(...)

BarSave(...)

BarSaveMulti(...)

BarNew(...)
```