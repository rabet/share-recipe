package main

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	LOGIN_ROUTE = "/login"
	ErrNotFound = FetchErr("Uživatel s tímto kódem neexistuje")
)

// PlayerStore stores score information about players
type PlayerStore interface {
	// GetRecipeScore(name string) int
	// RecordWin(name string)
	FetchRecipes() error
	GetRecipes() ([]Recipe, error)
	AddRecipe(userId, categoryId int, title, desc, link string) error
	FetchUserByCode(code int) (*User, error)
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
	router.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/recipes/", http.StatusFound)
	}))
	router.Handle("/recipes/", http.HandlerFunc(p.recipesHandler))
	router.Handle("/add-recipe", http.HandlerFunc(p.addRecipeHandler))
	router.Handle(LOGIN_ROUTE, http.HandlerFunc(p.loginHandler))
	// router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	p.Handler = router

	return p
}

func (p *PlayerServer) recipesHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("content-type", jsonContentType)
	// json.NewEncoder(w).Encode(p.store.GetRecipes())
	myRecipesPrefix := r.URL.Path[len("/recipes/"):]
	var filteredRecipes []Recipe

	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		fmt.Printf("Get session error: %v", err)
		//handleServerError(w, err)
		return
	}
	user := getUser(session)

	recipes, err := p.store.GetRecipes()
	if err != nil {
		handleServerError(w, err)
	}
	if myRecipesPrefix == "moje" {
		for _, r := range recipes {
			if r.Author.Username == user.Username {
				filteredRecipes = append(filteredRecipes, r)
			}
		}
	} else {
		filteredRecipes = recipes
	}
	tpl.ExecuteTemplate(w, "index.html", struct {
		Recipes []Recipe
	}{
		Recipes: filteredRecipes,
	})
}

func (p *PlayerServer) addRecipeHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("content-type", jsonContentType)
	// json.NewEncoder(w).Encode(p.store.GetRecipes())

	// cookies
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		fmt.Printf("Get session error: %v", err)
		//handleServerError(w, err)
		return
	}
	user := getUser(session)

	/* if auth := user.Authenticated; !auth {
		session.AddFlash("You don't have access!")
		err = session.Save(r, w)
		if err != nil {
			fmt.Printf("Auth error: %v", err)
			//handleServerError(w, err)
			return
		}
		//http.Redirect(w, r, LOGIN_ROUTE, http.StatusFound)
		return
	} */

	switch r.Method {
	case http.MethodPost:
		title := r.FormValue("title")
		desc := r.FormValue("desc")
		link := r.FormValue("link")
		categoryId, _ := strconv.Atoi(r.FormValue("categoryId"))

		if link == "" || title == "" {
			fmt.Fprint(w, "Zadej nazev, odkaz a kategorii")
		} else {
			ok := CheckWebsite(link)
			if !ok {
				fmt.Fprint(w, "Zadaný odkaz na recept nelze ověřit. Asi není platný.")
				return
			}
			err := p.store.AddRecipe(user.Id, categoryId, title, desc, link)
			if err != nil {
				handleServerError(w, err)
			}
			http.Redirect(w, r, "/recipes/", http.StatusFound)
		}
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "add-recipe.html", user)
	}
}
func (p *PlayerServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := sessionStore.Get(r, SessionName)
	if err != nil {
		fmt.Fprintf(w, "Session error: %v", err)
		return
	}

	code, _ := strconv.Atoi(r.FormValue("code"))

	user, err := p.store.FetchUserByCode(code)
	if err == ErrNotFound {
		err = session.Save(r, w)
		if err != nil {
			handleServerError(w, err)
			return
		}
		fmt.Fprint(w, ErrNotFound)
		return
	} else if err != nil {
		handleServerError(w, err)
		return
	}

	user.Authenticated = true

	session.Values["user"] = user

	err = session.Save(r, w)
	if err != nil {
		handleServerError(w, err)
		return
	}
	http.Redirect(w, r, "/add-recipe", http.StatusFound)
}

/* func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
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
} */

// CheckWebsite returns true if the URL returns a 200 status code, false otherwise
func CheckWebsite(url string) bool {
	response, err := http.Head(url)
	if err != nil {
		return false
	}

	if response.StatusCode != http.StatusOK {
		return false
	}

	return true
}

func handleServerError(w http.ResponseWriter, err error) {
	// log.WithField("err", err).Info("Error handling session.")
	http.Error(w, "Application error", http.StatusInternalServerError)
}
