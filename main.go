// main.go
package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
)

const SessionName = "share-recipe-cookie"

var (
	conn         *pgx.Conn
	tpl          *template.Template
	sessionStore *sessions.CookieStore
)

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

	gob.Register(User{})
	sessionStore = sessions.NewCookieStore(
		[]byte(viper.GetString("SESSION_AUTHENTICATION_KEY")),
		[]byte(viper.GetString("SESSION_ENCRYPTION_KEY")),
	)

	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24 * 7,
		HttpOnly: true,
	}
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

// getUser returns a user from session s
// on error returns an empty user
func getUser(s *sessions.Session) User {
	val := s.Values["user"]
	var user = User{}
	user, ok := val.(User)
	if !ok {
		return User{Authenticated: false}
	}
	return user
}
