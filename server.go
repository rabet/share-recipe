package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// PlayerStore stores score information about players
type PlayerStore interface {
	GetRecipeScore(name string) int
	RecordWin(name string)
	FetchRecipes() error
	GetRecipes() ([]Recipe, error)
	AddRecipe(userId, categoryId int, title, desc, link string) error
}

// PlayerServer is a HTTP interface for player information
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer creates a PlayerServer with routing configured
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/recipes", http.HandlerFunc(p.recipesHandler))
	router.Handle("/add-recipe", http.HandlerFunc(p.addRecipeHandler))
	// router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	p.Handler = router

	return p
}

func (p *PlayerServer) recipesHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("content-type", jsonContentType)
	// json.NewEncoder(w).Encode(p.store.GetRecipes())
	recipes, err := p.store.GetRecipes()
	if err != nil {
		handleServerError(w, err)
	}
	tpl.ExecuteTemplate(w, "index.html", struct {
		Recipes []Recipe
	}{
		Recipes: recipes,
	})
}

func (p *PlayerServer) addRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("content-type", jsonContentType)
	// json.NewEncoder(w).Encode(p.store.GetRecipes())

	switch r.Method {
	case http.MethodPost:
		title := r.FormValue("title")
		desc := r.FormValue("desc")
		link := r.FormValue("link")
		categoryId, _ := strconv.Atoi(r.FormValue("categoryId"))

		if (title == "") || (link == "") || categoryId == 0 {
			fmt.Fprint(w, "Zadej nazev, odkaz a kategorii")
		} else {
			err := p.store.AddRecipe(1, categoryId, title, desc, link)
			if err != nil {
				handleServerError(w, err)
			}
			fmt.Fprint(w, "Recept pridan, mnamy mnamyy ;)")
		}
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "add-recipe.html", nil)
	}
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
	player := r.URL.Path[len("/players/"):]

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetRecipeScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func handleServerError(w http.ResponseWriter, err error) {
	// log.WithField("err", err).Info("Error handling session.")
	http.Error(w, "Application Error", http.StatusInternalServerError)
}
