package main

import (
	"net/http"
    
    "github.com/julienschmidt/httprouter"
    "github.com/justinas/alice" // Package for middleware chaining
)

func (app *application) routes() http.Handler {
    router := httprouter.New()

    router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        app.notFound(w)
    })

    fileServer := http.FileServer(http.Dir("./ui/static"))
    router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

    dynamicMiddleWare := alice.New(app.sessionManager.LoadAndSave)

    router.Handler(http.MethodGet, "/", dynamicMiddleWare.ThenFunc(app.home))
    router.Handler(http.MethodGet, "/snippet/view/:id", dynamicMiddleWare.ThenFunc(app.snippetView))
    router.Handler(http.MethodGet, "/snippet/create", dynamicMiddleWare.ThenFunc(app.snippetCreate))
    router.Handler(http.MethodPost, "/snippet/create", dynamicMiddleWare.ThenFunc(app.snippetCreatePost))

    standardMiddleWare := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

    return standardMiddleWare.Then(router)
}