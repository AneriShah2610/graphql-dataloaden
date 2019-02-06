package dal

import (
	"database/sql"
	"fmt"
	"log"

	"sync"

	_ "github.com/lib/pq"
)

// DbConnection model
type DbConnection struct {
	Db *sql.DB
}

var once sync.Once
var instance *DbConnection

// DbConnect for database connection
func DbConnect() (*DbConnection, error) {

	once.Do(func() {
		fmt.Println("starting server")
		configuration, err := LoadConfiguration("dal/config.json")
		if err != nil {
			log.Fatal("Error while starting server", err)
		}
		connectionString := fmt.Sprintf("postgresql://%s@%s:%s/%s?sslmode=disable", configuration.Cockroach.User, configuration.Cockroach.Host, configuration.Cockroach.Port, configuration.Cockroach.DbName)
		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatal("error while initializing database", err)
		}
		fmt.Println("Database successfulyy initialized")
		instance = &DbConnection{
			Db: db,
		}

	})
	return instance, nil
}
