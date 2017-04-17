package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"

	mgo "gopkg.in/mgo.v2"
)

type apiUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int64  `json:"age"`
}

type userResource struct {
	Data apiUser `json:"data"`
}

type userHandler struct {
	repo Users
}

type appContext struct {
	db *mgo.Database
}

func (c *appContext) usersHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewUserRepository()
	session := c.db.Session.Clone()
	defer session.Close()

	users, err := repo.getAllUsers(session.DB("hammer"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(users)
}

func (c *appContext) userHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewUserRepository()
	session := c.db.Session.Clone()
	defer session.Close()

	params := context.Get(r, "params").(httprouter.Params)

	user, err := repo.getUserByID(session.DB("hammer"), params.ByName("id"))
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	json.NewEncoder(w).Encode(user)
}

func (c *appContext) createUserHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewUserRepository()
	session := c.db.Session.Clone()
	defer session.Close()

	body := context.Get(r, "body").(*userResource)

	err := repo.insertUser(session.DB("hammer"), body.Data.Name, body.Data.Age)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(body)
}
