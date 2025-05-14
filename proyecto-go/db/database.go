// db/database.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // Variable global

func Init() {
	// Conexión con MySQL
	cnx := "****:***********!@tcp(localhost:*****)/*******"
	var err error
	DB, err = sql.Open("mysql", cnx)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	err = DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Conexión exitosa")
}

// Cierra la conexión a la BD
func CloseConnection() {
	if DB != nil {
		DB.Close()
	}
}
