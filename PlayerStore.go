package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v4"
)

// NewRecipeStore initialises an empty Recipe store
func NewRecipeStore() *RecipeStore {
	return &RecipeStore{
		[]Recipe{},
		sync.RWMutex{},
	}
}

type Recipe struct {
	Title    string
	Descr    string
	Link     string
	Author   User
	Category string
}

type User struct {
	Id            int
	Username      string
	Code          int
	Authenticated bool
}

// RecipeStore collects data about Recipes in memory
type RecipeStore struct {
	store []Recipe
	// A mutex is used to synchronize read/write access to the map
	lock sync.RWMutex
}

// GetRecipes returns a collection of Recipes
func (i *RecipeStore) GetRecipes() ([]Recipe, error) {
	if len(i.store) <= 0 {
		err := i.FetchRecipes()
		if err != nil {
			fmt.Printf("Cannot fetch recipes: %v", err)
			return nil, err
		}
	}
	return i.store, nil
}

// FetchRecipes fetches a collection of Recipes
func (i *RecipeStore) FetchRecipes() error {
	i.store = nil
	rows, err := conn.Query(context.Background(), "SELECT recipe.title, recipe.descr, recipe.link, people.username, category.title FROM recipe INNER JOIN people ON recipe.people_id = people.id INNER JOIN category ON recipe.category_id = category.id")
	if err != nil {
		fmt.Errorf("Select query error: %v", err)
		return err
	}
	for rows.Next() {
		var r Recipe
		err := rows.Scan(&r.Title, &r.Descr, &r.Link, &r.Author.Username, &r.Category)
		if err != nil {
			fmt.Errorf("Rows Scan error: %v", err)
			return err
		}
		i.store = append(i.store, r)
	}

	return nil
}

// AddRecipe adds new recipe and fetch all recipes
func (i *RecipeStore) AddRecipe(userId, categoryId int, title, desc, link string) error {
	_, err := conn.Exec(context.Background(), "INSERT INTO recipe(people_id, category_id, title, descr, link) values($1, $2, $3, $4, $5)", userId, categoryId, title, desc, link)
	if err != nil {
		fmt.Errorf("Add recipe err: %v", err)
		return err
	}
	i.store = nil
	return nil
}

// GetLeague returns a collection of Recipes
func (i *RecipeStore) FetchUserByCode(code int) (*User, error) {
	var u User
	err := conn.QueryRow(context.Background(), "SELECT id, username, code FROM people WHERE code = $1", code).Scan(&u.Id, &u.Username, &u.Code)
	switch err {
	case nil:
		return &u, nil
	case pgx.ErrNoRows:
		return &u, nil
	default:
		return nil, err
	}
}

// RecordWin will record a Recipe's win
func (i *RecipeStore) RecordWin(name string) {
	// i.lock.Lock()
	// defer i.lock.Unlock()
	// i.store[name]++
}

// GetRecipeScore retrieves scores for a given Recipe
func (i *RecipeStore) GetRecipeScore(name string) int {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return 1
}
