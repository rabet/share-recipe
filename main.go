// main.go
package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
)

var conn *pgx.Conn
var tpl *template.Template

func init() {
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	dburl := viper.GetString("DATABASE_URL")
	dbConn(dburl)

	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {

	port := viper.GetString("PORT")
	server := NewPlayerServer(NewRecipeStore())

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}

// Initialises a connection do Postgres db
func dbConn(dbUrl string) {
	var err error
	conn, err = pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
}

// returns a handle to the DB object
func GetDB() *pgx.Conn {
	return conn
}
