package main

import (
	"net/http"

	mgo "gopkg.in/mgo.v2"

	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type router struct {
	*httprouter.Router
}

func (r *router) Get(path string, handler http.Handler) {
	r.GET(path, wrapHandler(handler))
}

func (r *router) Post(path string, handler http.Handler) {
	r.POST(path, wrapHandler(handler))
}

func (r *router) Put(path string, handler http.Handler) {
	r.PUT(path, wrapHandler(handler))
}

func (r *router) Delete(path string, handler http.Handler) {
	r.DELETE(path, wrapHandler(handler))
}

func NewRouter() *router {
	return &router{httprouter.New()}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		context.Set(r, "params", ps)
		h.ServeHTTP(w, r)
	}
}

func main() {
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)

	appC := appContext{session.DB("hammer")}
	commonHandlers := alice.New(context.ClearHandler, loggingHandler, recoverHandler, acceptHandler)
	router := NewRouter()
	router.Get("/user/:id", commonHandlers.ThenFunc(appC.userHandler))
	router.Get("/users", commonHandlers.ThenFunc(appC.usersHandler))
	http.ListenAndServe(":3000", router)
}
